# --------------------------------------------------------------------
# 
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.
# -----------------------------------------------------------------------

FROM alpine:3.20.3
LABEL maintainer="WSO2 Docker Maintainers <wso2.com>"

# Set environment variables
ENV LANG='en_US.UTF-8' LANGUAGE='en_US:en' LC_ALL='en_US.UTF-8'

RUN apk update && apk upgrade --no-cache \
    && apk add  --no-cache tzdata && apk upgrade libssl3 libcrypto3

ARG KGW_USER=wso2kgw
ARG KGW_USER_ID=10001
ARG CHECKSUM_AMD64="8302d54fc41d4ffbfea6871b6a04584265d5eabbe738aee2966f9a574b6f17d1"
ARG CHECKSUM_ARM64="8ebe934e7cf29a65ca31d671a97949467fba6977d656b5e6d6c40902ae7402f6"
ARG KGW_USER_GROUP=wso2
ARG KGW_USER_GROUP_ID=10001
ARG KGW_USER_HOME=/home/${KGW_USER}
ARG TARGETARCH=arm64

ARG MOTD="\n\
 Welcome to WSO2 Docker Resources \n\
 --------------------------------- \n\
 This Docker container comprises of a WSO2 product, running with its latest GA release \n\
 which is under the Apache License, Version 2.0. \n\
 Read more about Apache License, Version 2.0 here @ http://www.apache.org/licenses/LICENSE-2.0.\n"

RUN \
    addgroup -S -g ${KGW_USER_GROUP_ID} ${KGW_USER_GROUP} \
    && adduser -S -u ${KGW_USER_ID} -h ${KGW_USER_HOME} -G ${KGW_USER_GROUP} ${KGW_USER} \
    && chown -R ${KGW_USER}:${KGW_USER_GROUP} ${KGW_USER_HOME} \
    && echo '[ ! -z "${TERM}" -a -r /etc/motd ] && cat /etc/motd' >> /etc/bash.bashrc; echo "${MOTD}" > /etc/motd \
    chgrp -R 0 ${KGW_USER_HOME} \
    && chmod -R g=u ${KGW_USER_HOME}

WORKDIR ${KGW_USER_HOME}
USER ${KGW_USER_ID}

ADD ./resources/conf conf
ADD ./resources/security security
COPY ./${TARGETARCH}/main config-deployer

# Expose ports
EXPOSE 9443 9444

CMD ./config-deployer
