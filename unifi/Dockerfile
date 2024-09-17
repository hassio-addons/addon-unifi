FROM ghcr.io/hassio-addons/ubuntu-base:8.2.0

# Set shell
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Setup base system
RUN \
    apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends binutils libcap2 logrotate mongodb-server openjdk-17-jre-headless \
    \
    && curl -sSLf "https://dl.ui.com/unifi/8.4.62/unifi_sysvinit_all.deb" > /tmp/unifi.deb \
    \
    && dpkg --install /tmp/unifi.deb \
    && apt-get clean \
    && rm -fr \
        /tmp/* \
        /var/cache/* \
        /var/lib/apt/lists/* \
        /var/log/*.log \
        /var/log/apt

# Copy root filesystem
COPY rootfs /

# Health check
HEALTHCHECK --start-period=5m \
    CMD curl --insecure --fail https://localhost:8443 || exit 1
