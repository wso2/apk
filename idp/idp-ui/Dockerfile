FROM node:20-alpine

ARG USER=wso2apk
ARG USER_ID=10001
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=10001
ARG USER_HOME=/home/${USER}
ARG MOTD='printf "\n\
 Welcome to WSO2 Docker Resources \n\
 --------------------------------- \n\
 This Docker container comprises of a WSO2 product, running with its latest GA release \n\
 which is under the Apache License, Version 2.0. \n\
 Read more about Apache License, Version 2.0 here @ http://www.apache.org/licenses/LICENSE-2.0.\n"'
ENV ENV=${USER_HOME}"/.ashrc"
RUN apk update && apk upgrade --no-cache \
    && apk add  --no-cache tzdata && apk upgrade libssl3 libcrypto3
RUN \
    addgroup -S -g ${USER_GROUP_ID} ${USER_GROUP} \
    && adduser -S -u ${USER_ID} -h ${USER_HOME} -G ${USER_GROUP} ${USER} \
    && chgrp -R 0 ${USER_HOME} \
    && chmod -R g=u ${USER_HOME} \
    && echo ${MOTD} > "${ENV}"

WORKDIR /home/wso2apk
COPY package*.json ./
ENV NODE_ENV production
RUN npm ci --only=production
COPY . .
EXPOSE 9443
USER ${USER_ID}
CMD [ "node", "index.js" ]