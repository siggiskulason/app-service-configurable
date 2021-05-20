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
	"os"
	"path/filepath"
	"strings"

	hooks "github.com/canonical/edgex-snap-hooks"
)

var cli *hooks.CtlCli = hooks.NewSnapCtl()

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
	var debug = false
	var disable = false
	var err error

	status, err := cli.Config("debug")
	if err != nil {
		fmt.Println(fmt.Sprintf("edgex-asc:install: can't read value of 'debug': %v", err))
		os.Exit(1)
	}
	if status == "true" {
		debug = true
	}

	if err = hooks.Init(debug, "edgex-app-service-configurable"); err != nil {
		fmt.Println(fmt.Sprintf("edgex-asc:install: initialization failure: %v", err))
		os.Exit(1)

	}

	err = installProfiles()
	if err != nil {
		hooks.Error(fmt.Sprintf("edgex-asc:install: %v", err))
		os.Exit(1)
	}

	profile, err := cli.Config(hooks.ProfileConfig)
	if err != nil {
		hooks.Error(fmt.Sprintf("Reading config 'profile' failed: %v", err))
		os.Exit(1)
	}

	if profile == "" {
		// set default profile
		err = cli.SetConfig(hooks.ProfileConfig, "default")
		if err != nil {
			hooks.Error(fmt.Sprintf("Can't SET DEFAULT PROFILE - %v", err))
			os.Exit(1)
		}
	}

	// If autostart is not explicitly set, default to "no"
	// as only example service configuration and profiles
	// are provided by default.
	autostart, err := cli.Config(hooks.AutostartConfig)
	if err != nil {
		hooks.Error(fmt.Sprintf("Reading config 'autostart' failed: %v", err))
		os.Exit(1)
	}
	if autostart == "" {
		hooks.Debug("edgex-asc: autostart is NOT set, initializing to 'no'")
		autostart = "no"
	}

	autostart = strings.ToLower(autostart)
	if autostart == "true" || autostart == "yes" {
		if profile == "default" {
			hooks.Warn(fmt.Sprintf("autostart is %s, but no profile set", autostart))
			disable = true
		}
	} else if autostart == "false" || autostart == "no" {
		disable = true
	} else {
		hooks.Error(fmt.Sprintf("Invalid value for 'autostart' : %s", autostart))
		os.Exit(1)
	}

	// disable because there's no initial configuration or autostart has
	// been explicitly set to false|no.
	if disable {
		err = cli.Stop("app-service-configurable", true)
		if err != nil {
			hooks.Error(fmt.Sprintf("Can't stop service - %v", err))
			os.Exit(1)
		}
	}

	envJSON, err := cli.Config(hooks.EnvConfig)
	if err != nil {
		hooks.Error(fmt.Sprintf("Reading config 'env' failed: %v", err))
		os.Exit(1)
	}

	err = hooks.HandleEdgeXConfig("app-service-configurable", envJSON, nil)
	if err != nil {
		hooks.Error(fmt.Sprintf("HandleEdgeXConfig failed: %v", err))
		os.Exit(1)
	}
}
