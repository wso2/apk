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

# # Stage 1: Build the native image
FROM ghcr.io/graalvm/native-image-community:17-ol8 as build

# # Copy the JAR file into the build stage
COPY idp_domain_service.jar .

# # Build the native image
RUN native-image -jar idp_domain_service.jar --no-fallback -H:Name="idp_domain_service" -H:+StaticExecutableWithDynamicLibC

FROM ubuntu:22.04

# Update the package list and install build-essential
RUN apt-get update && apt-get install -y build-essential


ARG USER=wso2apk
ARG USER_ID=10001
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=10001
ARG USER_HOME=/home/${USER}

RUN groupadd --system -g ${USER_GROUP_ID} ${USER_GROUP} \
    && useradd --system --create-home --home-dir ${USER_HOME} --no-log-init -g ${USER_GROUP_ID} -u ${USER_ID} ${USER}

ADD idp ${USER_HOME}/idp
COPY --from=build /app/idp_domain_service ${USER_HOME}/idp
COPY docker-entrypoint.sh ${USER_HOME}

RUN chown -R ${USER} ${USER_HOME}/idp \
    && chown ${USER} /home/${USER}/docker-entrypoint.sh \
    && chgrp -R 0 ${USER_HOME} \
    && chmod -R g=u ${USER_HOME}

EXPOSE 9443
USER ${USER_ID}
WORKDIR ${USER_HOME}
ENTRYPOINT ["sh", "/home/wso2apk/docker-entrypoint.sh"]