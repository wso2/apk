FROM alpine:3.20.3
LABEL maintainer="WSO2 Docker Maintainers <wso2.com>"
RUN apk update && apk upgrade --no-cache \
    && apk add  --no-cache tzdata && apk upgrade libssl3 libcrypto3

ENV LANG=C.UTF-8

ARG APK_USER=wso2
ARG APK_USER_ID=10001
ARG APK_USER_GROUP=wso2
ARG APK_USER_GROUP_ID=10001
ARG APK_USER_HOME=/home/${APK_USER}
ARG TARGETARCH
ARG TAG_NAME
ARG SHORT_SHA
ENV TAG_NAME=$TAG_NAME
ENV SHORT_SHA=$SHORT_SHA


ARG MOTD="\n\
 Welcome to WSO2 Docker Resources \n\
 --------------------------------- \n\
 This Docker container comprises of a WSO2 product, running with its latest GA release \n\
 which is under the Apache License, Version 2.0. \n\
 Read more about Apache License, Version 2.0 here @ http://www.apache.org/licenses/LICENSE-2.0.\n"

# Add CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN \
    addgroup -S -g ${APK_USER_GROUP_ID} ${APK_USER_GROUP} \
    && adduser -S -u ${APK_USER_ID} -h ${APK_USER_HOME} -G ${APK_USER_GROUP} ${APK_USER} \
    && mkdir ${APK_USER_HOME}/logs && mkdir -p ${APK_USER_HOME}/artifacts/apis \
    && chown -R ${APK_USER}:${APK_USER_GROUP} ${APK_USER_HOME} \
    && echo '[ ! -z "${TERM}" -a -r /etc/motd ] && cat /etc/motd' >> /etc/bash.bashrc; echo "${MOTD}" > /etc/motd
RUN \
    chgrp -R 0 ${APK_USER_HOME} \
    && chmod -R g=u ${APK_USER_HOME}

WORKDIR ${APK_USER_HOME}


# Use non-root user
USER ${APK_USER_ID}

# Expose port
EXPOSE 8080

# Copy binary from builder
COPY ./${TARGETARCH}/main idp-domain-service

# Run the binary
ENTRYPOINT ["./idp-domain-service"]

