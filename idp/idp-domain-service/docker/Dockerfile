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

    FROM alpine:3.21.3

    # Set environment variables
    ENV LANG='en_US.UTF-8' LANGUAGE='en_US:en' LC_ALL='en_US.UTF-8'
    
    # Install required dependencies
    RUN apk add --no-cache \
        bash \
        curl \
        tzdata \
        ca-certificates \
        fontconfig \
        libstdc++ \
        openjdk21-jre \
        openssl \
        libc6-compat
    
    # Ensure /tmp directory exists and is writable
    RUN mkdir -p /tmp && chmod 1777 /tmp
    
    # Set Java environment variables
    ENV JAVA_HOME=/usr/lib/jvm/java-21-openjdk \
        PATH="/usr/lib/jvm/java-21-openjdk/bin:$PATH"
    ENV METRICS_ENABLED=false
    
    # Define user and group variables
    ARG USER=wso2apk
    ARG USER_ID=10001
    ARG USER_GROUP=wso2
    ARG USER_GROUP_ID=10001
    ARG USER_HOME=/home/${USER}
    
    # Create user and group
    RUN addgroup -g ${USER_GROUP_ID} ${USER_GROUP} \
        && adduser -D -h ${USER_HOME} -u ${USER_ID} -G ${USER_GROUP} ${USER}
    
    # Add application files
    ADD idp ${USER_HOME}/idp
    COPY docker-entrypoint.sh ${USER_HOME}
    COPY lib/jmx_prometheus_javaagent-0.20.0.jar ${USER_HOME}/lib/
    
    # Set permissions
    RUN chown -R ${USER}:${USER_GROUP} ${USER_HOME}/idp \
        && chown ${USER}:${USER_GROUP} ${USER_HOME}/docker-entrypoint.sh \
        && chmod +x ${USER_HOME}/docker-entrypoint.sh
    
    # Expose port
    EXPOSE 9443
    
    # Set user and working directory
    USER ${USER_ID}
    WORKDIR ${USER_HOME}
    
    # Set entrypoint
    ENTRYPOINT ["sh", "/home/wso2apk/docker-entrypoint.sh"]
    