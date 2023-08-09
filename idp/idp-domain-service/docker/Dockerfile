#---------------------------------------------------------------
#
# Copyright (c) 2022, WSO2 LLC. (http://www.wso2.com).
#
# WSO2 LLC. licenses this file to you under the Apache License,
# Version 2.0 (the "License"); you may not use this file except
# in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#
#---------------------------------------------------------------

FROM --platform=linux/amd64 alpine:3.18.2

ENV LANG='en_US.UTF-8' LANGUAGE='en_US:en' LC_ALL='en_US.UTF-8'

# install JDK Dependencies
RUN apk update \
    && apk add --no-cache tzdata musl-locales musl-locales-lang gcompat \
    && rm -rf /var/cache/apk/*

ENV JAVA_VERSION jdk-11.0.20+8

RUN set -eux; \
    ARCH="$(apk --print-arch)"; \
    case "${ARCH}" in \
       amd64|x86_64) \
         ESUM='0b7e894af823484f270031ed6626f78504a7a44f8173d67433e14321d015e669'; \
         BINARY_URL='https://github.com/adoptium/temurin11-binaries/releases/download/jdk-11.0.20%2B8/OpenJDK11U-jdk_x64_alpine-linux_hotspot_11.0.20_8.tar.gz'; \
         ;; \
       *) \
         echo "Unsupported arch: ${ARCH}"; \
         exit 1; \
         ;; \
    esac; \
	  wget -O /tmp/openjdk.tar.gz ${BINARY_URL}; \
	  echo "${ESUM} */tmp/openjdk.tar.gz" | sha256sum -c -; \
	  mkdir -p /opt/java/openjdk; \
	  tar --extract \
	      --file /tmp/openjdk.tar.gz \
	      --directory /opt/java/openjdk \
	      --strip-components 1 \
	      --no-same-owner \
	  ; \
    rm -rf /tmp/openjdk.tar.gz;

ENV JAVA_HOME=/opt/java/openjdk \
    PATH="/opt/java/openjdk/bin:$PATH"

ARG USER=wso2apk
ARG USER_ID=802
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=802
ARG USER_HOME=/home/${USER}

RUN addgroup --system -g ${USER_GROUP_ID} ${USER_GROUP} \
    && adduser --system --home ${USER_HOME} -g ${USER_GROUP_ID} -u ${USER_ID} ${USER}

ADD idp ${USER_HOME}/idp
COPY docker-entrypoint.sh ${USER_HOME}

RUN chown -R ${USER} ${USER_HOME}/idp \
    && chown ${USER} /home/${USER}/docker-entrypoint.sh \
    && chgrp -R 0 ${USER_HOME} \
    && chmod -R g=u ${USER_HOME}

EXPOSE 9443
USER ${USER_ID}
WORKDIR ${USER_HOME}
ENTRYPOINT ["sh", "/home/wso2apk/docker-entrypoint.sh"]
