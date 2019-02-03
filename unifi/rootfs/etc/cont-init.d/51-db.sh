#!/usr/bin/with-contenv bash
# ==============================================================================
# Community Hass.io Add-ons: UniFi Controller
# Ensures the data of the UniFi controller is store outside the container
# ==============================================================================
# shellcheck disable=SC1091
source /usr/lib/hassio-addons/base.sh

readonly properties="/data/unifi/data/system.properties"

if ! hass.file_exists "${properties}"; then
    touch "${properties}"
    echo "unifi.db.extraargs=--smallfiles" > "${properties}"
fi

#shellcheck disable=SC2016
sed -i \
    '/^unifi.db.extraargs=/{h;s/=.*/=--smallfiles/};${x;/^$/{s//unifi.db.extraargs=--smallfiles/;H};x}' \
    "${properties}"
