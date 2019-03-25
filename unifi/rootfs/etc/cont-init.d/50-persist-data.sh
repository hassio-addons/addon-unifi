#!/usr/bin/with-contenv bashio
# ==============================================================================
# Community Hass.io Add-ons: UniFi Controller
# Ensures the data of the UniFi controller is store outside the container
# ==============================================================================
if ! bashio::fs.directory_exists '/data/unifi/data'; then
    mkdir -p /data/unifi/data
fi
rm -fr /usr/lib/unifi/data
ln -s /data/unifi/data /usr/lib/unifi/data

if ! bashio::fs.directory_exists '/backup/unifi'; then
    mkdir -p /backup/unifi
fi
rm -fr /usr/lib/unifi/data/backup
ln -s /backup/unifi /usr/lib/unifi/data/backup
