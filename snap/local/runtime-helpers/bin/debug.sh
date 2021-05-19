#!/bin/bash -ex

debug=$(snapctl get debug)
logger "app-service-config: debug: $debug"

autostart=$(snapctl get autostart)
logger "app-service-config: autostart: $autostart"

exec "$@"
