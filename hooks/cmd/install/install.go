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
	"fmt"
	"io/ioutil"
	"log/syslog"
	"os"
	"path/filepath"

	"github.com/canonical/app-service-configurable/hooks"
)

var log *syslog.Writer

// installProfiles copies the profile configuration.toml files from $SNAP to $SNAP_DATA.
func installProfiles() error {
	dataConfP := fmt.Sprintf("%s/config/res", hooks.SnapData)
	snapConfP := fmt.Sprintf("%s/config/res", hooks.Snap)

	configFiles, err := filepath.Glob(filepath.Join(snapConfP, "*", "configuration.toml"))
	if err != nil {
		panic(fmt.Sprintf("internal error: bad glob pattern: %v", err))
	}

	for _, snapConfFile := range configFiles {
		// build the destination SNAP_DATA file by getting the directory name that the glob matched
		dirMatch := filepath.Base(filepath.Dir(snapConfFile))
		if dirMatch == "sample" {
			// TODO: what about sample config dirs ?
			continue
		}

		dataDestFile := filepath.Join(dataConfP, dirMatch, "configuration.toml")
		b, err := ioutil.ReadFile(snapConfFile)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(dataDestFile), 0755)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(dataDestFile, b, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var err error
	log, err = syslog.New(syslog.LOG_INFO, "edgex-asc:install")
	if err != nil {
		return
	}

	err = hooks.GetEnvVars()
	if err != nil {
		log.Crit(fmt.Sprintf("edgex-asc:install: %v", err))
		os.Exit(1)
	}

	err = installProfiles()
	if err != nil {
		log.Crit(fmt.Sprintf("edgex-asc:install: %v", err))
		os.Exit(1)
	}

	cli := hooks.NewSnapCtl()
	svc := fmt.Sprintf("%s.app-service-configurable", hooks.SnapInst)

	// disable app-service-configurable initially because it specific requires configuration
	// with a device profile that will be specific to each installation
	err = cli.Stop(svc, true)
	if err != nil {
		log.Crit(fmt.Sprintf("Can't stop service - %v", err))
		os.Exit(1)
	}

	// set default profile
	err = cli.SetConfig(hooks.ProfileConfig, hooks.DefaultProfile)
	if err != nil {
		log.Crit(fmt.Sprintf("Can't SET DEFAULT PROFILE - %v", err))
		os.Exit(1)
	}

	envJSON, err := cli.Config(hooks.EnvConfig)
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
