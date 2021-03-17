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

const (
	AutostartConfig     = "autostart"
	EnvConfig           = "env"
	ProfileConfig       = "profile"
	DefaultProfile      = "default"
	SnapEnv             = "SNAP"
	SnapDataEnv         = "SNAP_DATA"
	SnapInstanceNameEnv = "SNAP_INSTANCE_NAME"
)

// ConfToEnv defines mappings from snap config keys to EdgeX environment variable
// names that are used to override individual service configuration values via a
// .env file read by the snap service wrapper.
var ConfToEnv = map[string]string{
	// [Writable] - not yet supported
	// conf_to_env["writable.log-level"]="BootTimeout"
	// [Service]
	"service.boot-timeout":     "SERVICE_BOOTTIMEOUT",
	"service.check-interval":   "SERVICE_CHECKINTERVAL",
	"service.host":             "SERVICE_HOST",
	"service.server-bind-addr": "SERVICE_SERVERBINDADDR",
	"service.port":             "SERVICE_PORT",
	"service.protocol":         "SERVICE_PROTOCOL",
	"service.read-max-limit":   "SERVICE_READMAXLIMIT",
	"service.startup-msg":      "SERVICE_STARTUPMSG",
	"service.timeout":          "SERVICE_TIMEOUT",

	// [Registry]
	"registry.host": "REGISTRY_HOST",
	"registry.port": "REGISTRY_PORT",
	"registry.type": "REGISTRY_TYPE",

	// [Database]
	"database.type":    "DATABASE_TYPE",
	"database.host":    "DATABASE_HOST",
	"database.port":    "DATABASE_PORT",
	"database.timeout": "DATABASE_TIMEOUT",

	// [SecretStore]
	"secretstore.host":                      "SECRETSTORE_HOST",
	"secretstore.port":                      "SECRETSTORE_PORT",
	"secretstore.path":                      "SECRETSTORE_PATH",
	"secretstore.protocol":                  "SECRETSTORE_PROTOCOL",
	"secretstore.root-ca-cert-path":         "SECRETSTORE_ROOTCACERTPATH",
	"secretstore.server-name":               "SECRETSTORE_SERVERNAME",
	"secretstore.token-file":                "SECRETSTORE_TOKENFILE",
	"secretstore.additional-retry-attempts": "SECRETSTORE_ADDITIONALRETRYATTEMPTS",
	"secretstore.retry-wait-period":         "SECRETSTORE_RETRYWAITPERIOD",

	// [SecretStore.Authentication]
	"secretstore.authentication.auth-type": "SECRETSTORE_AUTHENTICATION_AUTHTYPE",

	// [SecretStoreExclusive]
	"secretstore-ex.host":                      "SECRETSTOREEXCLUSIVE_HOST",
	"secretstore-ex.port":                      "SECRETSTOREEXCLUSIVE_PORT",
	"secretstore-ex.path":                      "SECRETSTOREEXCLUSIVE_PATH",
	"secretstore-ex.protocol":                  "SECRETSTOREEXCLUSIVE_PROTOCOL",
	"secretstore-ex.root-ca-cert-path":         "SECRETSTOREEXCLUSIVE_ROOTCACERTPATH",
	"secretstore-ex.server-name":               "SECRETSTOREEXCLUSIVE_SERVERNAME",
	"secretstore-ex.token-file":                "SECRETSTOREEXCLUSIVE_TOKENFILE",
	"secretstore-ex.additional-retry-attempts": "SECRETSTOREEXCLUSIVE_ADDITIONALRETRYATTEMPTS",
	"secretstore-ex.retry-wait-period":         "SECRETSTOREEXCLUSIVE_RETRYWAITPERIOD",
	// [SecretStore.Authentication]
	"secretstore-ex.authentication.auth-type": "SECRETSTOREEXCLUSIVE_AUTHENTICATION_AUTHTYPE",

	// [Clients.CoreData]
	"clients.coredata.host":     "CLIENTS_COREDATA_HOST",
	"clients.coredata.port":     "CLIENTS_COREDATA_PORT",
	"clients.coredata.protocol": "CLIENTS_COREDATA_PROTOCOL",

	// [Binding]
	"binding.type":            "BINDING_TYPE",
	"binding.subscribe-topic": "BINDING_SUBSCRIBE_TOPIC",
	"binding.publish-topic":   "BINDING_PUBLISH_TOPIC",

	// [MessageBus]
	"message-bus.type": "MESSAGEBUS_TYPE",
	// [MessageBus.SubscribeHost]
	"message-bus.subscribe-host.host":     "MESSAGEBUS_SUBSCRIBEHOST_HOST",
	"message-bus.subscribe-host.port":     "MESSAGEBUS_SUBSCRIBEHOST_PORT",
	"message-bus.subscribe-host.protocol": "MESSAGEBUS_SUBSCRIBEHOST_PROTOCOL",

	// [MessageBus.PublishHost]
	"message-bus.publish-host.host":     "MESSAGEBUS_PUBLISHHOST_HOST",
	"message-bus.publish-host.port":     "MESSAGEBUS_PUBLISHHOST_PORT",
	"message-bus.publish-host.protocol": "MESSAGEBUS_PUBLISHHOST_PROTOCOL",
}
