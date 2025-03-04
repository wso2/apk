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
ADD config-deployer ${USER_HOME}/config-deployer
COPY docker-entrypoint.sh ${USER_HOME}
COPY /lib/jmx_prometheus_javaagent-0.20.0.jar ${USER_HOME}/lib/

# Set permissions
RUN chown -R ${USER}:${USER_GROUP} ${USER_HOME}/config-deployer \
    && chown ${USER}:${USER_GROUP} ${USER_HOME}/docker-entrypoint.sh \
    && chmod +x ${USER_HOME}/docker-entrypoint.sh

# Expose ports
EXPOSE 9443 9444

# Set user and working directory
USER ${USER_ID}
WORKDIR ${USER_HOME}

# Set entrypoint
ENTRYPOINT ["sh", "/home/wso2apk/docker-entrypoint.sh"]
