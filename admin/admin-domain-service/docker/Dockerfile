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

FROM apkbase:0.0.1-SNAPSHOT

ARG USER=wso2apk
ARG USER_HOME=/home/${USER}

COPY docker-entrypoint.sh ${USER_HOME}
ADD admin ${USER_HOME}/admin
RUN chown -R ${USER} ${USER_HOME}/admin
RUN chown ${USER} /home/${USER}/docker-entrypoint.sh

EXPOSE 9443
USER wso2apk
WORKDIR ${USER_HOME}

ENTRYPOINT ["sh", "/home/wso2apk/docker-entrypoint.sh"]
