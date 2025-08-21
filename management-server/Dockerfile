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

FROM alpine:3.17.2
LABEL maintainer="WSO2 Docker Maintainers <wso2.com>"

RUN apk update && apk upgrade --no-cache
RUN apk add  --no-cache tzdata

ENV LANG=C.UTF-8

ARG APK_USER=wso2
ARG APK_USER_ID=802
ARG APK_USER_GROUP=wso2
ARG APK_USER_GROUP_ID=802
ARG APK_USER_HOME=/home/${APK_USER}
ARG GRPC_HEALTH_PROBE_PATH=/bin/grpc_health_probe
ARG TARGETARCH
ENV MGT_SERVER_PRIVATE_KEY_PATH=/home/wso2/security/keystore/mg.key
ENV MGT_SERVER_PUBLIC_CERT_PATH=/home/wso2/security/keystore/mg.pem
ARG CHECKSUM_AMD64="36f1e2bade1672ba7f3ddb4a14e5035c06856cc5fd4e9e538dd1d4878871edeb"
ARG CHECKSUM_ARM64="ac7b46065a9faf15a9ddca4cbdb145442f858fd430cf245fc4f31c33bdc61143"

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
    wget -q https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.40/grpc_health_probe-linux-${TARGETARCH} \
    && mv grpc_health_probe-linux-${TARGETARCH} ${GRPC_HEALTH_PROBE_PATH}\
    && if [ "${TARGETARCH}" = "amd64" ]; then echo "${CHECKSUM_AMD64}  ${GRPC_HEALTH_PROBE_PATH}" | sha256sum -c -; fi

RUN \
    chmod +x ${GRPC_HEALTH_PROBE_PATH} \
    && chown ${APK_USER}:${APK_USER_GROUP} ${GRPC_HEALTH_PROBE_PATH}

WORKDIR ${APK_USER_HOME}
USER ${APK_USER}

COPY resources/conf/config.toml conf/
COPY resources/conf/log_config.toml conf/
COPY ./${TARGETARCH}/main mgt_server
COPY resources/check_health.sh .

CMD ./mgt_server
