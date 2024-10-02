# --------------------------------------------------------------------
# Copyright (c) 2023, WSO2 LLC. (http://wso2.com) All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# -----------------------------------------------------------------------

FROM ubuntu:22.04

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

LABEL maintainer="WSO2 Docker Maintainers <wso2.com>"

ENV LANG=C.UTF-8

ARG APK_USER=wso2
ARG APK_USER_ID=10001
ARG APK_USER_GROUP=wso2
ARG APK_USER_GROUP_ID=10001
ARG APK_USER_HOME=/home/${APK_USER}
ARG GRPC_HEALTH_PROBE_PATH=/bin/grpc_health_probe
ARG TARGETARCH
ARG APK_VERSION=1.0-SNAPSHOT
ENV VERSION=${APK_VERSION}
ENV JAVA_OPTS=""
ENV ENFORCER_HOME=${APK_USER_HOME}
ARG CHECKSUM_AMD64="3ddaf85583613c97693e9b8aaa251dac07e73e366e159a7ccadbcf553117fcef"
ARG CHECKSUM_ARM64="5e17ff4c055f075b58a1cd7ec37843d989cd0072340222a4fd0730773382027e"
ENV ENFORCER_PRIVATE_KEY_PATH=/home/wso2/security/keystore/mg.key
ENV ENFORCER_PUBLIC_CERT_PATH=/home/wso2/security/keystore/mg.pem
ENV TRUSTED_CA_CERTS_PATH=/home/wso2/security/truststore
ENV ADAPTER_HOST_NAME=adapter
ENV ADAPTER_HOST=adapter
ENV ADAPTER_XDS_PORT=18000
ENV COMMON_CONTROLLER_HOST_NAME=common-controller
ENV COMMON_CONTROLLER_HOST=common-controller
ENV COMMON_CONTROLLER_XDS_PORT=18002
ENV ENFORCER_LABEL="default"
ENV XDS_MAX_MSG_SIZE=4194304
ENV XDS_MAX_RETRIES=3
#todo update the connection string
ENV APPLICATIONINSIGHTS_CONNECTION_STRING=InstrumentationKey=;IngestionEndpoint=https://westus2-2.in.applicationinsights.azure.com/

ARG MOTD="\n\
 Welcome to WSO2 Docker Resources \n\
 --------------------------------- \n\
 This Docker container comprises of a WSO2 product, which is under the Apache License, Version 2.0. \n\
 Read more about Apache License, Version 2.0 here @ http://www.apache.org/licenses/LICENSE-2.0.\n"

RUN \
    apt-get update && apt-get dist-upgrade -y && apt-get autoclean -y && apt-get autoremove -y
RUN \
    groupadd --system -g ${APK_USER_GROUP_ID} ${APK_USER_GROUP} \
    && useradd --system --create-home --home-dir ${APK_USER_HOME} --no-log-init -g ${APK_USER_GROUP} -u ${APK_USER_ID} ${APK_USER} \
    && mkdir ${APK_USER_HOME}/logs && mkdir -p ${APK_USER_HOME}/lib/dropins \
    && chown -R ${APK_USER}:${APK_USER_GROUP} ${APK_USER_HOME} \
    && echo '[ ! -z "${TERM}" -a -r /etc/motd ] && cat /etc/motd' >> /etc/bash.bashrc; echo "${MOTD}" > /etc/motd

RUN \
    wget -q https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.34/grpc_health_probe-linux-${TARGETARCH} \
    && mv grpc_health_probe-linux-${TARGETARCH} ${GRPC_HEALTH_PROBE_PATH}\
    && if [ "${TARGETARCH}" = "amd64" ]; then echo "${CHECKSUM_AMD64}  ${GRPC_HEALTH_PROBE_PATH}" | sha256sum -c -; fi

RUN \
    chmod +x ${GRPC_HEALTH_PROBE_PATH} \
    && chown ${APK_USER}:${APK_USER_GROUP} ${GRPC_HEALTH_PROBE_PATH} \
    && chgrp -R 0 ${APK_USER_HOME} \
    && chmod -R g=u ${APK_USER_HOME}

WORKDIR ${APK_USER_HOME}
USER ${APK_USER_ID}

COPY resources/lib lib/
COPY resources/conf/log4j2.properties conf/log4j2.properties
COPY resources/security security
COPY resources/check_health.sh .

#todo add applicationinsights.json file and point it to the appInsightsAgent jar
#Add the agent using JVM arg -javaagent:/home/wso2/conf/applicationinsights-agent-3.1.1.jar
#Add the config file using System property -Dapplicationinsights.configuration.file=/home/wso2/conf/applicationinsights.json

CMD java -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath="${ENFORCER_HOME}/logs/heap-dump.hprof" $JAVA_OPTS -Dlog4j.configurationFile="${ENFORCER_HOME}/conf/log4j2.properties" -DtracingEnabled="true" -cp "lib/*:lib/dropins/*" org.wso2.apk.enforcer.server.AuthServer
