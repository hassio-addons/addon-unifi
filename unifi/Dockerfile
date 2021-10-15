ARG BUILD_FROM=ghcr.io/hassio-addons/ubuntu-base/amd64:6.2.1
# hadolint ignore=DL3006
FROM ${BUILD_FROM}

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Set Unifi args
ARG UNIFI_VERSION
# ARG UNIFI_BRANCH="stable"  Change if you need to move to beta, and edit latest version check from 'stable' to '${UNIFI_BRANCH}'
# Setup base system
ARG BUILD_ARCH=amd64

RUN \
    apt-get update \
    && apt-get install -y --no-install-recommends \
        gnupg2=2.2.4-1ubuntu1.4 \
    \
    && curl -fsSL "https://www.mongodb.org/static/pgp/server-3.4.asc" \
        | apt-key add -\
    && echo "deb http://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.4 multiverse" \
        > /etc/apt/sources.list.d/mongodb.list \
    \
    && apt-get update \
    && apt-get install -y \
	binutils \
	jsvc \
	libcap2 \
	logrotate \
	mongodb-server \
	openjdk-8-jre-headless \
	wget \
    \
    && echo "**** check latest version ****" 
    && \
    if [ -z ${UNIFI_VERSION+x} ]; then \
	 UNIFI_VERSION=$(curl -sX GET http://dl-origin.ubnt.com/unifi/debian/dists/stable/ubiquiti/binary-amd64/Packages \
	 |grep -A 7 -m 1 'Package: unifi' \
	 | awk -F ': ' '/Version/{print $2;exit}' \
	 | awk -F '-' '{print $1}'); \
    fi 
    && curl -o \
    /tmp/unifi.deb -L \
	"https://dl.ui.com/unifi/${UNIFI_VERSION}/unifi_sysvinit_all.deb" && \
    dpkg -i /tmp/unifi.deb 
    && echo "**** cleanup ****" && \
    apt-get clean && \
    rm -rf \
	/tmp/* \
	/var/lib/apt/lists/* \
	/var/tmp/*

# Copy root filesystem
COPY rootfs /

# Build arguments
ARG BUILD_ARCH
ARG BUILD_DATE
ARG BUILD_DESCRIPTION
ARG BUILD_NAME
ARG BUILD_REF
ARG BUILD_REPOSITORY
ARG BUILD_VERSION

# Labels
LABEL \
    io.hass.name="${BUILD_NAME}" \
    io.hass.description="${BUILD_DESCRIPTION}" \
    io.hass.arch="${BUILD_ARCH}" \
    io.hass.type="addon" \
    io.hass.version=${BUILD_VERSION} \
    maintainer="Franck Nijhof <frenck@addons.community>" \
    org.opencontainers.image.title="${BUILD_NAME}" \
    org.opencontainers.image.description="${BUILD_DESCRIPTION}" \
    org.opencontainers.image.vendor="Home Assistant Community Add-ons" \
    org.opencontainers.image.authors="Franck Nijhof <frenck@addons.community>" \
    org.opencontainers.image.licenses="MIT" \
    org.opencontainers.image.url="https://addons.community" \
    org.opencontainers.image.source="https://github.com/${BUILD_REPOSITORY}" \
    org.opencontainers.image.documentation="https://github.com/${BUILD_REPOSITORY}/blob/main/README.md" \
    org.opencontainers.image.created=${BUILD_DATE} \
    org.opencontainers.image.revision=${BUILD_REF} \
    org.opencontainers.image.version=${BUILD_VERSION}
