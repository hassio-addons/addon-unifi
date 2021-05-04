#!/usr/bin/with-contenv bashio
# ==============================================================================
# Home Assistant Community Add-on: UniFi Controller
# Configures the UniFi Controller
# ==============================================================================
readonly KEYSTORE="/usr/lib/unifi/data/keystore"
readonly properties="/data/unifi/data/system.properties"
declare certfile
declare keyfile
declare root_chain
declare tempcert

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

# Identrust cross-signed CA cert needed by the java keystore for import.
# Can get original here: https://www.identrust.com/certificates/trustid/root-download-x3.html
root_chain=$(cat <<-END
-----BEGIN CERTIFICATE-----
MIIDSjCCAjKgAwIBAgIQRK+wgNajJ7qJMDmGLvhAazANBgkqhkiG9w0BAQUFADA/
MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMT
DkRTVCBSb290IENBIFgzMB4XDTAwMDkzMDIxMTIxOVoXDTIxMDkzMDE0MDExNVow
PzEkMCIGA1UEChMbRGlnaXRhbCBTaWduYXR1cmUgVHJ1c3QgQ28uMRcwFQYDVQQD
Ew5EU1QgUm9vdCBDQSBYMzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AN+v6ZdQCINXtMxiZfaQguzH0yxrMMpb7NnDfcdAwRgUi+DoM3ZJKuM/IUmTrE4O
rz5Iy2Xu/NMhD2XSKtkyj4zl93ewEnu1lcCJo6m67XMuegwGMoOifooUMM0RoOEq
OLl5CjH9UL2AZd+3UWODyOKIYepLYYHsUmu5ouJLGiifSKOeDNoJjj4XLh7dIN9b
xiqKqy69cK3FCxolkHRyxXtqqzTWMIn/5WgTe1QLyNau7Fqckh49ZLOMxt+/yUFw
7BZy1SbsOFU5Q9D8/RhcQPGX69Wam40dutolucbY38EVAjqr2m7xPi71XAicPNaD
aeQQmxkqtilX4+U9m5/wAl0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNV
HQ8BAf8EBAMCAQYwHQYDVR0OBBYEFMSnsaR7LHH62+FLkHX/xBVghYkQMA0GCSqG
SIb3DQEBBQUAA4IBAQCjGiybFwBcqR7uKGY3Or+Dxz9LwwmglSBd49lZRNI+DT69
ikugdB/OEIKcdBodfpga3csTS7MgROSR6cz8faXbauX+5v3gTt23ADq1cEmv8uXr
AvHRAosZy5Q6XkjEGB5YGV8eAlrwDPGxrancWYaLbumR9YbK+rlmM6pZW87ipxZz
R8srzJmwN0jP41ZL9c8PDHIyh8bwRLtTcm1D9SZImlJnt1ir/md2cXjbDaJWFBM5
JDGFoqgCWjBH4d1QB7wCCZAA62RjYJsWvIjJEubSfZGL+T0yjWW06XyxV3bqxbYo
Ob8VZRzI9neWagqNdwvYkQsEjgfbKbYK7p2CNTUQ
-----END CERTIFICATE-----
END
)

# Stop running this script if SSL is disabled
if bashio::config.false 'ssl'; then
 exit 0
fi

# Initialize keystore, in case it does not exist yet
if ! bashio::fs.file_exists "${KEYSTORE}"; then
    bashio::log.debug 'Intializing keystore...'
    keytool \
        -genkey \
        -keyalg RSA \
        -alias unifi \
        -keystore "${KEYSTORE}" \
        -storepass aircontrolenterprise \
        -keypass aircontrolenterprise \
        -validity 1825 \
        -keysize 4096 \
        -dname "cn=UniFi" || \
        bashio::exit.nok "Failed creating UniFi keystore"
fi

bashio::log.debug 'Injecting SSL certificate into the controller...'

certfile="/ssl/$(bashio::config 'certfile')"
keyfile="/ssl/$(bashio::config 'keyfile')"
tempcert=$(mktemp)

# Adds Identrust cross-signed CA cert in case of letsencrypt
if [[ $(openssl x509 -noout -ocsp_uri -in "${certfile}") == *"letsencrypt"* ]]; then
    echo "${root_chain}" > "${tempcert}"
    cat "${certfile}" >> "${tempcert}"
else
    cat "${certfile}" > "${tempcert}"
fi

bashio::log.debug 'Preparing certificate in a format UniFi accepts...'
openssl pkcs12 \
    -export  \
    -passout pass:aircontrolenterprise \
    -in "${tempcert}" \
    -inkey "${keyfile}" \
    -out "${tempcert}" \
    -name unifi

bashio::log.debug 'Removing existing certificate from UniFi protected keystore...'
keytool \
    -delete \
    -alias unifi \
    -keystore "${KEYSTORE}" \
    -deststorepass aircontrolenterprise

bashio::log.debug 'Inserting certificate into UniFi keystore...'
keytool \
    -trustcacerts \
    -importkeystore \
    -deststorepass aircontrolenterprise \
    -destkeypass aircontrolenterprise \
    -destkeystore "${KEYSTORE}" \
    -srckeystore "${tempcert}" \
    -srcstoretype PKCS12 \
    -srcstorepass aircontrolenterprise \
    -alias unifi

# Cleanup
rm -f "${tempcert}"
