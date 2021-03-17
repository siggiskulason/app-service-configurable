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

package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var SnapData string
var Snap string
var SnapInst string

// CtlCli is the test obj for overridding functions
type CtlCli struct{}

// SnapCtl interface provides abstration for unit testing
type SnapCtl interface {
	Config(key string) (string, error)
	SetConfig(key string, val string) error
	Stop(svc string, disable bool) error
}

func GetEnvVars() error {
	SnapData = os.Getenv(SnapDataEnv)
	if SnapData == "" {
		return errors.New("SNAP_DATA is not set")
	}

	Snap = os.Getenv(SnapEnv)
	if Snap == "" {
		return errors.New("SNAP is not set")
	}

	SnapInst = os.Getenv(SnapInstanceNameEnv)
	if SnapInst == "" {
		return errors.New("SNAP_INSTANCE_NAME is not set")
	}

	return nil
}

// NewSnapCtl returns a normal runtime client
func NewSnapCtl() *CtlCli {
	return &CtlCli{}
}

// ModClient returns a testing client
//func NewTestCtl(g Getter) *Client {
//	return &Client{getter: g}
//}

// Get uses snapctl to get a value from a key, or returns error
func (cc *CtlCli) Config(key string) (string, error) {
	out, err := exec.Command("snapctl", "get", key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Get uses snapctl to set a config value from a key, or returns error
func (cc *CtlCli) SetConfig(key string, val string) error {

	err := exec.Command("snapctl", "set", fmt.Sprintf("%s=%s", key, val)).Run()
	if err != nil {
		return errors.New(fmt.Sprintf("snapctl SET failed for %s - %v", key, err))
	}
	return nil
}

// Stop uses snapctrl to stop a service and optional disable it
func (cc *CtlCli) Stop(svc string, disable bool) error {
	var cmd *exec.Cmd

	if disable {
		cmd = exec.Command("snapctl", "stop", "--disable", svc)
	} else {
		cmd = exec.Command("snapctl", "stop", svc)

	}

	err := cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("snapctl stop failed - %v", err))
	}

	return nil
}

// p is the current prefix of the config key being processed (e.g. "service", "security.auth")
// k is the key name of the current JSON object being processed
// vJSON is the current object
// flatConf is a map containing the configuration keys/values processed thus far
func flattenConfigJSON(p string, k string, vJSON interface{}, flatConf map[string]string) {
	var mk string

	// top level keys don't include "env", so no separator needed
	if p == "" {
		mk = k
	} else {
		mk = fmt.Sprintf("%s.%s", p, k)
	}

	switch t := vJSON.(type) {
	case string:
		flatConf[mk] = t
	case bool:
		flatConf[mk] = strconv.FormatBool(t)
	case float64:
		flatConf[mk] = strconv.FormatFloat(t, 'f', -1, 64)
	case map[string]interface{}:

		for k, v := range t {
			flattenConfigJSON(mk, k, v, flatConf)
		}
	default:
		panic(fmt.Sprintf("internal error: invalid JSON configuration from snapd - prefix: %s key: %s obj: %v", p, k, t))
	}
}

// HandleEdgeXConfig processes snap configuration which can be used to override
// app-service-configurable configuration via environment variables sourced by
// the snap service wrapper script.
func HandleEdgeXConfig(envJSON string) error {

	if envJSON == "" {
		return nil
	}

	var m map[string]interface{}
	var flatConf = make(map[string]string)

	err := json.Unmarshal([]byte(envJSON), &m)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to unmarshall EdgeX config - %v", err))
	}

	for k, v := range m {
		flattenConfigJSON("", k, v, flatConf)
	}

	b := bytes.Buffer{}
	for k, v := range flatConf {
		env, ok := ConfToEnv[k]
		if !ok {
			return errors.New(fmt.Sprintf("invalid EdgeX config option - %s", k))
		}

		_, err := fmt.Fprintf(&b, "export %s=%s\n", env, v)
		if err != nil {
			return err
		}
	}

	path := fmt.Sprintf("%s/config/res/service.env", SnapData)
	err = ioutil.WriteFile(path, b.Bytes(), 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to write service.env file - %v", err))
	}

	return nil
}
