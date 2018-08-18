ARG BUILD_FROM=hassioaddons/ubuntu-base:2.1.0
# hadolint ignore=DL3006
FROM ${BUILD_FROM}

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Setup base system
ARG BUILD_ARCH=amd64
RUN \
    URL="http://archive.ubuntu.com/ubuntu/" \
    && if [ "${BUILD_ARCH}" = "armhf" ] || [ "${BUILD_ARCH}" = "aarch64" ]; then \
        URL="http://ports.ubuntu.com/ubuntu-ports/"; \
    fi \
    && echo "deb ${URL} xenial universe" \
        > /etc/apt/sources.list.d/xenial-universe.list \
    \
    && echo "deb ${URL} xenial main" \
        > /etc/apt/sources.list.d/xenial-main.list \
    \
    && apt-get update \
    \
    && apt-get install -y --no-install-recommends \
        binutils=2.30-20ubuntu2~18.04 \
        jsvc=1.0.15-8 \
        libcap2=1:2.25-1.2 \
        mongodb-server=1:2.6.10-0ubuntu1 \
        openjdk-8-jdk-headless=8u171-b11-0ubuntu0.18.04.1 \
    \
    && curl -J -L -o /tmp/unifi.deb \
        "https://dl.ubnt.com/unifi/5.8.28/unifi_sysvinit_all.deb" \
    \
    && dpkg --install /tmp/unifi.deb \
    \
    && rm -rf /tmp/* \
    && rm -rf /root/.gnupg \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/* \
    && rm -rf /var/{cache,log}/*

# Copy root filesystem
COPY rootfs /

# Build arugments
ARG BUILD_DATE
ARG BUILD_REF
ARG BUILD_VERSION

# Labels
LABEL \
    io.hass.name="UniFi Controller" \
    io.hass.description="Manage your UniFi network using a web browser" \
    io.hass.arch="${BUILD_ARCH}" \
    io.hass.type="addon" \
    io.hass.version=${BUILD_VERSION} \
    maintainer="Franck Nijhof <frenck@addons.community>" \
    org.label-schema.description="Manage your UniFi network using a web browser" \
    org.label-schema.build-date=${BUILD_DATE} \
    org.label-schema.name="UniFi Controller" \
    org.label-schema.schema-version="1.0" \
    org.label-schema.url="https://community.home-assistant.io/t/community-hass-io-add-on-unifi-controller/56297?u=frenck" \
    org.label-schema.usage="https://github.com/hassio-addons/addon-unifi/tree/master/README.md" \
    org.label-schema.vcs-ref=${BUILD_REF} \
    org.label-schema.vcs-url="https://github.com/hassio-addons/addon-unifi" \
    org.label-schema.vendor="Community Hass.io Add-ons"
