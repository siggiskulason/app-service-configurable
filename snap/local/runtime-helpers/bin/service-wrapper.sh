#!/bin/bash -e

PROFILE_OPT=""

profile=$(snapctl get "profile")

if [ ! -z "$profile" ]; then
    if [ "$profile" != "default" ]; then
        PROFILE_OPT="-profile $profile"
    fi
fi

# disable TLS usage
export SECRETSTORE_PROTOCOL='http'
export SECRETSTORE_ROOTCACERTPATH=""
export SECRETSTOREEXCLUSIVE_PROTOCOL='http'
export SECRETSTOREEXCLUSIVE_ROOTCACERTPATH=""

$SNAP/bin/app-service-configurable -confdir $SNAP_DATA/config/res $PROFILE_OPT -cp -r

