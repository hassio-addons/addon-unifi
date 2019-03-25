#!/usr/bin/with-contenv bashio
# ==============================================================================
# Community Hass.io Add-ons: UniFi Controller
# Ensures the data of the UniFi controller is store outside the container
# ==============================================================================
readonly properties="/data/unifi/data/system.properties"

if ! bashio::fs.file_exists "${properties}"; then
    touch "${properties}"
    echo "unifi.db.extraargs=--smallfiles" > "${properties}"
fi

#shellcheck disable=SC2016
sed -i \
    '/^unifi.db.extraargs=/{h;s/=.*/=--smallfiles/};${x;/^$/{s//unifi.db.extraargs=--smallfiles/;H};x}' \
    "${properties}"
