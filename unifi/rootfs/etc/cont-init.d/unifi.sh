#!/usr/bin/with-contenv bashio
# ==============================================================================
# Community Hass.io Add-ons: UniFi Controller
# Configures the UniFi Controller
# ==============================================================================
readonly properties="/data/unifi/data/system.properties"

bashio::config.require.ssl

# Ensures the data of the UniFi controller is store outside the container
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

# Enable small files on MongoDB
if ! bashio::fs.file_exists "${properties}"; then
    touch "${properties}"
    echo "unifi.db.extraargs=--smallfiles" > "${properties}"
fi

#shellcheck disable=SC2016
sed -i \
    '/^unifi.db.extraargs=/{h;s/=.*/=--smallfiles/};${x;/^$/{s//unifi.db.extraargs=--smallfiles/;H};x}' \
    "${properties}"
