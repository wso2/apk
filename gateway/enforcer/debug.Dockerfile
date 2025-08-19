# --------------------------------------------------------------------
# Copyright (c) 2022, WSO2 LLC. (http://wso2.com) All Rights Reserved.
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

# Stage 1: Build dlv
FROM golang:1.21-alpine AS dlv-builder

RUN apk add --no-cache git bash

RUN go install github.com/go-delve/delve/cmd/dlv@latest


FROM alpine:3.20.3
LABEL maintainer="WSO2 Docker Maintainers <wso2.com>"

RUN apk update && apk upgrade --no-cache \
    && apk add  --no-cache tzdata && apk upgrade libssl3 libcrypto3

ENV LANG=C.UTF-8

ARG APK_USER=wso2
ARG APK_USER_ID=10001
ARG CHECKSUM_AMD64="8302d54fc41d4ffbfea6871b6a04584265d5eabbe738aee2966f9a574b6f17d1"
ARG CHECKSUM_ARM64="8ebe934e7cf29a65ca31d671a97949467fba6977d656b5e6d6c40902ae7402f6"
ARG APK_USER_GROUP=wso2
ARG APK_USER_GROUP_ID=10001
ARG APK_USER_HOME=/home/${APK_USER}
ARG GRPC_HEALTH_PROBE_PATH=/bin/grpc_health_probe
ARG TARGETARCH

ARG MOTD="\n\
 Welcome to WSO2 Docker Resources \n\
 --------------------------------- \n\
 This Docker container comprises of a WSO2 product, running with its latest GA release \n\
 which is under the Apache License, Version 2.0. \n\
 Read more about Apache License, Version 2.0 here @ http://www.apache.org/licenses/LICENSE-2.0.\n"

RUN \
    addgroup -S -g ${APK_USER_GROUP_ID} ${APK_USER_GROUP} \
    && adduser -S -u ${APK_USER_ID} -h ${APK_USER_HOME} -G ${APK_USER_GROUP} ${APK_USER} \
    && mkdir ${APK_USER_HOME}/logs && mkdir -p ${APK_USER_HOME}/artifacts/apis \
    && chown -R ${APK_USER}:${APK_USER_GROUP} ${APK_USER_HOME} \
    && echo '[ ! -z "${TERM}" -a -r /etc/motd ] && cat /etc/motd' >> /etc/bash.bashrc; echo "${MOTD}" > /etc/motd

RUN \
    wget -q https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-${TARGETARCH} \
    && mv grpc_health_probe-linux-${TARGETARCH} ${GRPC_HEALTH_PROBE_PATH} \
    && if [ "${TARGETARCH}" = "amd64" ]; then echo "${CHECKSUM_AMD64}  ${GRPC_HEALTH_PROBE_PATH}" | sha256sum -c -; fi
    
RUN \
    chmod +x ${GRPC_HEALTH_PROBE_PATH} \
    && chown ${APK_USER}:${APK_USER_GROUP} ${GRPC_HEALTH_PROBE_PATH} \
    && chgrp -R 0 ${APK_USER_HOME} \
    && chmod -R g=u ${APK_USER_HOME}

WORKDIR ${APK_USER_HOME}
USER ${APK_USER_ID}

COPY resources/conf/config.toml conf/
COPY resources/check_health.sh .
COPY resources/conf/log_config.toml conf/
COPY ./${TARGETARCH}/main enforcer

# copy dlv from builder
COPY --from=dlv-builder /go/bin/dlv /usr/local/bin/dlv

CMD ./enforcer
