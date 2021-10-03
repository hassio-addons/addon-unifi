ARG BUILD_FROM=ghcr.io/hassio-addons/ubuntu-base/amd64:6.2.1
# hadolint ignore=DL3006
FROM ${BUILD_FROM}

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

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
    && apt-get install -y --no-install-recommends \
        binutils=2.30-21ubuntu1~18.04.5 \
        jsvc=1.0.15-8 \
        libcap2=1:2.25-1.2 \
        libssl1.0.0=1.0.2n-1ubuntu5.6 \
        logrotate=3.11.0-0.1ubuntu1 \
        mongodb-org-server=3.4.24 \
        openjdk-8-jdk-headless=8u292-b10-0ubuntu1~18.04 \
    \
    && curl -J -L -o /tmp/unifi.deb \
        "https://dl.ui.com/unifi/6.4.54/unifi_sysvinit_all.deb" \
    \
    && dpkg --install /tmp/unifi.deb \
    \
    && apt-get remove -y --purge \
        gnupg2 \
    \
    && apt-get clean \
    && rm -fr \
        /tmp/* \
        /root/.gnupg \
        /var/{cache,log}/* \
        /var/lib/apt/lists/*

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
