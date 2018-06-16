#!/usr/bin/with-contenv bash
# ==============================================================================
# Community Hass.io Add-ons: UniFi Controller
# Ensures the data of the UniFi controller is store outside the container
# ==============================================================================
# shellcheck disable=SC1091
source /usr/lib/hassio-addons/base.sh

if ! hass.directory_exists '/data/unifi/data'; then
    mkdir -p /data/unifi/data
fi
rm -fr /usr/lib/unifi/data
ln -s /data/unifi/data /usr/lib/unifi/data

if ! hass.directory_exists '/backup/unifi'; then
    mkdir -p /backup/unifi
fi
rm -fr /usr/lib/unifi/data/backup
ln -s /backup/unifi /usr/lib/unifi/data/backup
