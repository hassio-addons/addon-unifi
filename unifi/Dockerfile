ARG BUILD_FROM=hassioaddons/ubuntu-base:3.1.3
# hadolint ignore=DL3006
FROM ${BUILD_FROM}

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Setup base system
ARG BUILD_ARCH=amd64
RUN \
    URL="http://archive.ubuntu.com/ubuntu/" \
    && if [ "${BUILD_ARCH}" = "armv7" ] \
        || [ "${BUILD_ARCH}" = "aarch64" ]; then \
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
        binutils=2.30-21ubuntu1~18.04.2 \
        jsvc=1.0.15-8 \
        libcap2=1:2.25-1.2 \
        mongodb-server=1:2.6.10-0ubuntu1 \
        openjdk-8-jdk-headless=8u222-b10-1ubuntu1~18.04.1 \
    \
    && curl -J -L -o /tmp/unifi.deb \
        "https://dl.ubnt.com/unifi/5.10.25/unifi_sysvinit_all.deb" \
    \
    && dpkg --install /tmp/unifi.deb \
    \
    && rm -fr \
        /tmp/* \
        /root/.gnupg \
        /var/{cache,log}/* \
        /var/lib/apt/lists/*

# Copy root filesystem
COPY rootfs /

# Build arguments
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
