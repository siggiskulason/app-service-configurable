// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2017-2021 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package hooks

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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
