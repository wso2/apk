# --------------------------------------------------------------------
# Copyright (c) 2023, WSO2 Inc. (http://wso2.com) All Rights Reserved.
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

FROM --platform=linux/x86_64 adoptopenjdk/openjdk11:jre-11.0.11_9-alpine
LABEL maintainer="WSO2 Docker Maintainers <wso2.com>"

RUN apk update && apk upgrade --no-cache
RUN apk add  --no-cache tzdata

LABEL maintainer="WSO2 Docker Maintainers <dev@wso2.org>"

ENV LANG=C.UTF-8

ARG APK_USER=wso2
ARG APK_USER_ID=802
ARG APK_USER_GROUP=wso2
ARG APK_USER_GROUP_ID=802
ARG APK_USER_HOME=/home/${APK_USER}
ARG GRPC_HEALTH_PROBE_PATH=/bin/grpc_health_probe
ARG TARGETARCH
ARG APK_VERSION=1.0-SNAPSHOT

ENV VERSION=${APK_VERSION}
ENV JAVA_OPTS=""
ENV ENFORCER_HOME=${APK_USER_HOME}
ARG CHECKSUM_AMD64="c72704c9cd49fb18f8df28de26e29c2563b6515e38efee65a6c3a29ec7368a91"
ARG CHECKSUM_ARM64="a7b42093b702ea83f4dcbd1dd74e4085c118fe93b9cfd3f205e48506a9ac94a0"
ENV ENFORCER_PRIVATE_KEY_PATH=/home/wso2/security/keystore/mg.key
ENV ENFORCER_PUBLIC_CERT_PATH=/home/wso2/security/keystore/mg.pem
ENV TRUSTED_CA_CERTS_PATH=/home/wso2/security/truststore
ENV ADAPTER_HOST_NAME=adapter
ENV ADAPTER_HOST=adapter
ENV ADAPTER_XDS_PORT=18000
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
    addgroup -S -g ${APK_USER_GROUP_ID} ${APK_USER_GROUP} \
    && adduser -S -u ${APK_USER_ID} -h ${APK_USER_HOME} -G ${APK_USER_GROUP} ${APK_USER} \
    && mkdir ${APK_USER_HOME}/logs && mkdir -p ${APK_USER_HOME}/lib/dropins \
    && chown -R ${APK_USER}:${APK_USER_GROUP} ${APK_USER_HOME} \
    && echo '[ ! -z "${TERM}" -a -r /etc/motd ] && cat /etc/motd' >> /etc/bash.bashrc; echo "${MOTD}" > /etc/motd

RUN \
    wget -q https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.16/grpc_health_probe-linux-${TARGETARCH} \
    && mv grpc_health_probe-linux-${TARGETARCH} ${GRPC_HEALTH_PROBE_PATH}\
    && if [ "${TARGETARCH}" = "amd64" ]; then echo "${CHECKSUM_AMD64}  ${GRPC_HEALTH_PROBE_PATH}" | sha256sum -c -; fi

RUN \
    chmod +x ${GRPC_HEALTH_PROBE_PATH} \
    && chown ${APK_USER}:${APK_USER_GROUP} ${GRPC_HEALTH_PROBE_PATH}



WORKDIR ${APK_USER_HOME}
USER ${APK_USER}

COPY resources/lib lib/
COPY resources/conf/log4j2.properties conf/log4j2.properties
COPY resources/security security
COPY resources/check_health.sh .

#todo add applicationinsights.json file and point it to the appInsightsAgent jar
#Add the agent using JVM arg -javaagent:/home/wso2/conf/applicationinsights-agent-3.1.1.jar
#Add the config file using System property -Dapplicationinsights.configuration.file=/home/wso2/conf/applicationinsights.json

CMD java -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath="${ENFORCER_HOME}/logs/heap-dump.hprof" $JAVA_OPTS -Dlog4j.configurationFile="${ENFORCER_HOME}/conf/log4j2.properties" -DtracingEnabled="true" -cp "lib/*:lib/dropins/*" org.wso2.choreo.connect.enforcer.server.AuthServer
