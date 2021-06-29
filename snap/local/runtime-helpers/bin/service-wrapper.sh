#!/bin/bash -e

PROFILE_OPT=""

profile=$(snapctl get "profile")

if [ ! -z "$profile" ]; then
    if [ "$profile" != "default" ]; then
        PROFILE_OPT="-profile $profile"
    fi
fi

SecretStore_TokenFile="$SNAP_DATA/$profile/secrets-token.json"

$SNAP/bin/app-service-configurable -confdir $SNAP_DATA/config/res $PROFILE_OPT -cp -r

