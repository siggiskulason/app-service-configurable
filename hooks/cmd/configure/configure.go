// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2021 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 * SPDX-License-Identifier: Apache-2.0'
 */

package main

import (
	"errors"
	"fmt"
	"log/syslog"
	"os"

	"github.com/canonical/app-service-configurable/hooks"
)

var log *syslog.Writer

// validateProfile processes the snap 'profile' configure option, ensuring that the directory
// and associated configuration.toml file in $SNAP_DATA both exist.
//
func validateProfile(prof string) error {
	log.Debug(fmt.Sprintf("edgex-asc:configure:validateProfile: profile is %s", prof))

	if prof == "" || prof == hooks.DefaultProfile {
		return nil
	}

	path := fmt.Sprintf("%s/config/res/%s/configuration.toml", hooks.SnapData, prof)
	log.Debug(fmt.Sprintf("edgex-asc:configure:validateProfile: checking if %s exists", path))

	_, err := os.Stat(path)
	if err != nil {
		return errors.New(fmt.Sprintf("profile %s has no configuration.toml", prof))
	}

	return nil
}

func main() {
	var err error
	var envJSON, prof string

	log, err = syslog.New(syslog.LOG_INFO, "edgex-asc:configure")
	if err != nil {
		log.Crit(fmt.Sprintf("Creating new syslog instance failed: %v", err))
		os.Exit(1)
	}

	err = hooks.GetEnvVars()
	if err != nil {
		log.Crit(fmt.Sprintf("Error reading SNAP environment variables: %v", err))
		os.Exit(1)
	}

	log.Debug("edgex-asc:configure hook running")

	cli := hooks.NewSnapCtl()

	prof, err = cli.Config(hooks.ProfileConfig)
	if err != nil {
		log.Crit(fmt.Sprintf("Error reading config 'profile': %v", err))
		os.Exit(1)
	}

	validateProfile(prof)
	if err != nil {
		log.Crit(fmt.Sprintf("Error validating profile: %v", err))
		os.Exit(1)
	}

	envJSON, err = cli.Config(hooks.EnvConfig)
	if err != nil {
		log.Crit(fmt.Sprintf("Reading config 'env' failed: %v", err))
		os.Exit(1)
	}

	err = hooks.HandleEdgeXConfig(envJSON)
	if err != nil {
		log.Crit(fmt.Sprintf("HandleEdgeXConfig failed: %v", err))
		os.Exit(1)
	}
}
