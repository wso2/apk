
FROM ubuntu:24.10

ENV LANG='en_US.UTF-8' LANGUAGE='en_US:en' LC_ALL='en_US.UTF-8'

# Upgrade Ubuntu Dependencies
RUN apt-get update \
    && apt-get upgrade -y

# install JDK Dependencies
RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends tzdata curl wget ca-certificates fontconfig locales \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && locale-gen en_US.UTF-8 \
    && rm -rf /var/lib/apt/lists/*

ENV JAVA_VERSION jdk-17.0.8.1+1

RUN set -eux; \
    ARCH="$(dpkg --print-architecture)"; \
    case "${ARCH}" in \
        amd64|i386:x86-64) \
            ESUM='ab68857594792474a3049ede09ea1178e42df29803a6a41be771794f571b2d4e'; \
            BINARY_URL='https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.8.1%2B1/OpenJDK17U-jre_x64_linux_hotspot_17.0.8.1_1.tar.gz'; \
            ;; \
        aarch64|arm64) \
            ESUM='0a1c5c9ee9d20832c87bd1e99a4c4a96947b59bb35c72683fe895d705f202737'; \
            BINARY_URL='https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.8.1%2B1/OpenJDK17U-jre_aarch64_linux_hotspot_17.0.8.1_1.tar.gz'; \
            ;; \
        *) \
            echo "Unsupported arch: ${ARCH}"; \
            exit 1; \
            ;; \
    esac; \
    curl -LfsSo /tmp/openjdk.tar.gz ${BINARY_URL}; \
    echo "${ESUM} */tmp/openjdk.tar.gz" | sha256sum -c -; \
    mkdir -p /opt/java/openjdk; \
    cd /opt/java/openjdk; \
    tar -xf /tmp/openjdk.tar.gz --strip-components=1; \
    rm -rf /tmp/openjdk.tar.gz;

ENV JAVA_HOME=/opt/java/openjdk \
    PATH="/opt/java/openjdk/bin:$PATH"
ENV METRICS_ENABLED=false

ARG USER=wso2apk
ARG USER_ID=10001
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=10001
ARG USER_HOME=/home/${USER}

RUN groupadd --system -g ${USER_GROUP_ID} ${USER_GROUP} \
    && useradd --system --create-home --home-dir ${USER_HOME} --no-log-init -g ${USER_GROUP_ID} -u ${USER_ID} ${USER}

ADD config-deployer ${USER_HOME}/config-deployer
COPY docker-entrypoint.sh ${USER_HOME}
COPY /lib/jmx_prometheus_javaagent-0.20.0.jar /home/${USER}/lib/

RUN chown -R ${USER} ${USER_HOME}/config-deployer \
    && chown ${USER} /home/${USER}/docker-entrypoint.sh \
    && chgrp -R 0 ${USER_HOME} \
    && chmod -R g=u ${USER_HOME}
EXPOSE 9443 9444
USER ${USER_ID}
WORKDIR ${USER_HOME}
ENTRYPOINT ["sh", "/home/wso2apk/docker-entrypoint.sh"]