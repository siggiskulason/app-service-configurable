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

package main

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"os"
	"strconv"

	"github.com/canonical/app-service-configurable/hooks"
)

var log syslog.Writer

func handleVal(p string, k string, v interface{}, flatConf map[string]interface{}) {
	var mk string

	// top level keys don't include "env", so no separator needed
	if p == "" {
		mk = k
	} else {
		mk = fmt.Sprintf("%s.%s", p, k)
	}

	log.Debug(fmt.Sprintf("handleVal: mk: %s", mk))

	switch t := v.(type) {
	case string:
		log.Debug(fmt.Sprintf("ADDING %s=%s to flatConf", k, t))
		flatConf[mk] = t
	case bool:
		log.Debug(fmt.Sprintf("ADDING %s=%v to flatConf", k, t))
		flatConf[mk] = strconv.FormatBool(t)
	case float64:
		log.Debug(fmt.Sprintf("ADDING %s=%v to flatConf", k, t))
		flatConf[mk] = strconv.FormatFloat(t, 'f', -1, 64)
	case map[string]interface{}:
		log.Debug(fmt.Sprintf("FOUND AN OBJECT"))

		for k, v := range t {
			handleVal(mk, k, v, flatConf)
		}
	default:
		log.Err("I DON'T KNOW!!!!")
	}
}

func handleSvcConf(env string) {
	log.Debug(fmt.Sprintf("edgex-asc:configure:handleSvcConf config is %s", env))

	if env == "" {
		return
	}

	var m map[string]interface{}
	var flatConf = make(map[string]interface{})
	//flatM = make(map[string]interface{})

	err := json.Unmarshal([]byte(env), &m)
	if err != nil {
		log.Err(fmt.Sprintf("edgex-asc:configure:handleSvcConf: failed to unmarshall env; %v", err))
		return
	}

	for k, v := range m {
		handleVal("", k, v, flatConf)
	}

	path := fmt.Sprintf("%s/config/res/service.env", hooks.SnapData)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Err(fmt.Sprintf("edgex-asc:configure:handleSvcConf: can't open %s - %v", path, err))
		os.Exit(1)
	}

	defer f.Close()

	log.Debug(fmt.Sprintf("edgex-asc:configure:handleSvcConf about write %s", path))
	for k, v := range flatConf {
		log.Debug(fmt.Sprintf("%s=%v", k, v))
		_, err := f.WriteString(fmt.Sprintf("export %s=%s\n", hooks.ConfToEnv[k], v))
		if err != nil {
			log.Err(fmt.Sprintf("edgex-asc:configure:handleSvcConf: can't open %s - %v", path, err))
			os.Exit(1)
		}
	}
}

func handleProf(prof string) {
	log.Debug(fmt.Sprintf("edgex-asc:configure:handleProf: profile is %s", prof))

	if prof == "" || prof == "default" {
		return
	}

	path := fmt.Sprintf("%s/config/res/%s/configuration.toml", hooks.SnapData, prof)
	log.Debug(fmt.Sprintf("edgex-asc:configure:handleProf: checking if %s exists", path))

	_, err := os.Stat(path)
	if err != nil {
		log.Err(fmt.Sprintf("edgex-asc:configure:handleProf: invalid setting profile %s specified; no configuration.toml found", prof))
		os.Exit(1)
	}

	log.Debug(fmt.Sprintf("edgex-asc:configure:handleProf: OK!!!"))
}

func main() {
	var env, prof string

	log, err := syslog.New(syslog.LOG_INFO, "edgex-asc:configure")
	if err != nil {
		return
	}

	err = hooks.GetEnvVars()
	if err != nil {
		log.Crit(fmt.Sprintf("edgex-asc:configure: %v", err))
		os.Exit(1)
	}

	log.Debug("edgex-asc:configure hook running")

	cli := hooks.NewSnapCtl()

	prof, err = cli.Config(hooks.ProfileConfig)
	if err != nil {
		log.Crit(fmt.Sprintf("edgex-asc:install: error reading config key 'env' - %v", err))
		os.Exit(1)
	}
	handleProf(prof)

	env, err = cli.Config(hooks.EnvConfig)
	if err != nil {
		log.Crit(fmt.Sprintf("edgex-asc:install: error reading config key 'env' - %v", err))
		os.Exit(1)
	}
	handleSvcConf(env)
}
