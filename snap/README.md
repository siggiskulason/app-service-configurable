# EdgeX Foundry App Service Configurable Snap
[![snap store badge](https://raw.githubusercontent.com/snapcore/snap-store-badges/master/EN/%5BEN%5D-snap-store-black-uneditable.png)](https://snapcraft.io/edgex-app-service-configurable)

This folder contains snap packaging for the EdgeX Foundry's App Service Configurable application service.

## Installation

### Installing snapd
The snap can be installed on any system that supports snaps. You can see how to install 
snaps on your system [here](https://snapcraft.io/docs/installing-snapd/6735).

However for full security confinement, the snap should be installed on an 
Ubuntu 16.04 LTS or later Desktop or Server, or a system running Ubuntu Core 16 or later.

### Installing EdgeX App Service Configurable as a snap
The snap is published in the snap store at https://snapcraft.io/edgex-app-service-configurable.
You can see the current revisions available for your machine's architecture by running the command:

```bash
$ snap info edgex-app-service-configurable
```

The latest stable version of the snap can be installed using:

```bash
$ sudo snap install edgex-app-service-configurable
```

**Note** - this snap has only been tested on Ubuntu Desktop/Server 16.04 LTS (or greater) and Ubuntu Core 16 (or greater).

## Using the EdgeX App Service Configurable snap
The App Service Configurable application service allows a variety of use cases to be met by simply providing configuration (vs. writing code). For more information about this service, please refer to the top-level [README.md](https://github.com/edgexfoundry/app-service-configurable/blob/hanoi/README.md) fil
e. As with some device service snaps, this service is disabled by default when first installed (although this can be overridden). This allows a profile to be specified (required), and any required configuration changes to be made prior to the service being enabled/started.

### Profiles
On install, all of the services profile-specific ```configuration.toml``` files are copied to sub-directories found in the snap's writable area:

```
/var/snap/edgex-app-service-configurable/current/config/res/
```

To configure a specific profile use the ```snap set``` command:

```
$ sudo snap set edgex-app-service-configurable profile=push-to-core
```

In addition to instructing the service to read a different configuration file, the profile will also be used to name the service when it registers itself to the system.

### Autostart
By default, the edgex-app-service-configurable snap disables the service on install, as the expectation is that the default profile configuration files will be customized, and thus this behavior allows the profile ```configuration.toml``` files in $SNAP_DATA to be modified before the service is first started.

This behavior can be overridden by setting the ```autostart``` is set to "true". This is useful when one or more configuration profiles are being provided via configuration or gadget snap content interface. If specified, care also needs to be taken to also ensure that profile has been set to a valid profile.

**Note** - this option is typically set from a gadget snap.

### Rich Configuration
While it's possible on Ubuntu Core to provide additional profiles via gadget snap content interface, quite often only minor changes to
existing profiles are required. These changes can be accomplished via support for EdgeX environment variable configuration overrides via
the snap's configure and install hooks. If the service has already been started, setting one of these overrides currently requires the
service to be restarted via the command-line or snapd's REST API. If the overrides are provided via the snap configuration defaults
capability of a gadget snap, the overrides will be picked up when the services are first started.

The following syntax is used to specify service-specific configuration overrides:


```env.<stanza>.<config option>```

For instance, to setup an override of the service's Port use:

```$ sudo snap set env.service.port=2112```

And restart the service:

```$ sudo snap restart edgex-app-service-configurable```

**Note** - at this time changes to configuration values in the [Writable] section are not supported.

For details on the mapping of configuration options to Config options, please refer to "Appendix A - edgex-app-service-configurable Configuration options".

### Startup environment variables

EdgeX services by default wait 60s for dependencies (e.g. Core Data) to become available, and will exit after this time if the dependencies aren't met. The following options can be used to override this startup behavior on systems where it takes longer than expected for the dependent services provided by the edgexfoundry snap to start. Note, both options below are specified as a number of seconds.
    
To change the default startup duration (60 seconds), for a service to complete the startup, aka bootstrap, phase of execution by using the following command:

```bash
$ sudo snap set edgex-app-service-configurable startup-duration=60
```

The following environment variable overrides the retry startup interval or sleep time before a failure is retried during the start-up, aka bootstrap, phase of execution by using the following command:

```bash
$ sudo snap set edgex-app-service-configurable startup-interval=1
```

**Note** - Should the environment variables be modified after the service has started, the service must be restarted.

### Multiple Instances
Multiple instances of edgex-app-service-configurable can be installed by using snap [Parallel Installs](https://snapcraft.io/docs/parallel-installs). This is an experimental snap feature and must be first be enabled by running this command:
```
sudo snap set system experimental.parallel-instances=true
```
Now you can install multiple instances of the edgex-app-service-configurable snap by specifying a unique instance name when you install the snap. The instance name is a unique suffix which is appended to the snap name following the “_” character used as a delimeter. This name only needs to be specified for the second and susbequent instances of the snap.
```
sudo snap install edgex-app-service-configurable edgex-app-service-configurable_http
```
or
```
sudo snap install edgex-app-service-configurable edgex-app-service-configurable_mqtt
```
**Note** – you must ensure that any configuration values that might cause conflict between the multiple instances (e.g. port, log file path, …) must be modified before enabling the snap’s service.

### Secret Store Usage
Some profile configuration.toml files specify configuration which requires Secret Store (aka Vault) support, however this snap doesn't fully support secure secrets without manual intervention (see below). The snap can also be configured to use insecure secrets as can be done via docker-compose by setting the option ```security-secret-store=off```. Ex.

```
sudo snap set edgex-app-service-configurable security-secret-store=off
```

## Configuration Options
This section documents the edgex-app-service-configurable's configuration options.

### Autostart
```autostart								// true | yes```

### API Gateway
```env.security-proxy.add-proxy-routes				// ADD_PROXY_ROUTES```

### Secret Store
```env.security-secret-store.add-secretstore-tokens		// ADD_SECRETSTORE_TOKENS```

### Service Environment Configuration Overrides
**Note** - all of the configuration options below must be specified with the prefix: 'env.'

```
[Service]
service.boot-timeout            // Service.BootTimeout
service.check-interval          // Service.CheckInterval
service.host                    // Service.Host
service.server-bind-addr        // Service.ServerBindAddr
service.port                    // Service.Port
service.protocol                // Service.Protocol
service.read-max-limit          // Service_ReadMaxLimit
service.startup-msg             // Service.StartupMsg
service.timeout                 // Service.Timeout

[Registry]
registry.host                   // Registry.Host
registry.port                   // Registry.Port
registry.type                   // Registry.Type

[Database]
database.primary.type           // Database.Type
database.host                   // Database.Host
database.primary.port           // Database.Port
database.primary.timeout        // Database.Timeout

[SecretStore]
secretstore.host                       // SecretStore.Host
secretstore.port                       // SecretStore.Port
secretstore.path                       // SecretStore.Path
secretstore.protocol                   // SecretStore.Protocol
secretstore.root-ca-cert-path          // SecretStore.RootCaCertPath
secretstore.server-name                // SecretStore.ServerName
secretstore.token-file                 // SecretStore.TokenFile
secretstore.additional-retry-attempts  // SecretStore.AdditionalRetryAttempts
secretstore.retry-wait-period          // SecretStore.RetryWaitPeriod

[SecretStore.Authentication]
secretstore.authentication.auth-type   // SecretStore.Authentication.AuthType

[SecretStoreExclusive]
secretstore-ex.host                    // SecretStoreExclusive.Host
secretstore-ex.port                    // SecretStoreExclusive.Port
secretstore-ex.path                    // SecretStoreExclusive.Path
secretstore-ex.protocol                // SecretStoreExclusive.Protocol
secretstore-ex.root-ca-cert-path       // SecretStoreExclusive.RootCaCertPath
secretstore-ex.server-name             // SecretStoreExclusive.ServerName
secretstore-ex.token-file              // SecretStoreExclusive.TokenFile

secretstore-ex.additional-retry-attempts // SecretStoreExclusive.AdditionalRetryAttempts
secretstore-ex.retry-wait-period         // SecretStoreExclusive.RetryWaitPeriod

[SecretStoreExclusive.Authentication]
secretstore-ex.authentication.auth-type  // SecretStoreExclusive.Authentication.AuthType

[Clients.CoreData]
clients.data.host                 // Clients.Data.Host
clients.data.port                 // Clients.Data.Port
clients.data.protocol             // Clients.Data.Protocol

**Note** - the key used vs. the key used for the edgexfoundry snap (e.g. Data).

[Binding]
binding.subscribe-topic           // Binding.Subscribe_Topic
binding.publish-topic             // Binding.PublishTopic

[MessageBus]
message-bus.type                  // MessageBus.Type

[MessageBus.SubscribeHost]
message-bus.subscribe-host.host      // MessageBus.SubscribeHost.Host
message-bus.subscribe-host.port      // MessageBus.SubscribeHost.Port
message-bus.subscribe-host.protocol  // MessageBus.SubscribeHost.Protocol

[MessageBus.PublishHost]
message-bus.publish-host.host        // MessageBus.PublishHost.Host
message-bus.publish-host.port        // MessageBus.PublishHost.Port
message-bus.publish-host.protocol    // MessageBus.PublishHost.Protocol
```
