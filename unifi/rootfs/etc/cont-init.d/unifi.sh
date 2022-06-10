#!/command/with-contenv bashio
# ==============================================================================
# Home Assistant Community Add-on: UniFi Network Application
# Configures the UniFi Network Application
# ==============================================================================
readonly KEYSTORE="/usr/lib/unifi/data/keystore"
readonly properties="/data/unifi/data/system.properties"

# Ensures the data of the UniFi Network Application is store outside the container
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

# If there is no keystore yet, we are good to go
if ! bashio::fs.file_exists "${KEYSTORE}"; then
    # Prevent migration next time
    touch /data/keystore_reset
    bashio::exit.ok
fi

# If the keystore has never been reset, we are going to do that... once
# This is a migration path, for people that previously had a custom SSL
# certificate.
# It will trigger a new self-signed certificate.
# This logic can be removed and cleaned up at a later moment
if ! bashio::fs.file_exists "/data/keystore_reset"; then
    rm -f -r "${KEYSTORE}"
    touch /data/keystore_reset
    bashio::exit.ok
fi
