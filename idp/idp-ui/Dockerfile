FROM node:16-alpine

ARG USER=wso2apk
ARG USER_ID=802
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=802
ARG USER_HOME=/home/${USER}

RUN groupadd --system -g ${USER_GROUP_ID} ${USER_GROUP} && useradd --system --create-home --home-dir ${USER_HOME} --no-log-init -g ${USER_GROUP_ID} -u ${USER_ID} ${USER}

WORKDIR /home/wso2apk
COPY package*.json ./
ENV NODE_ENV production
RUN npm ci --only=production
COPY . .
EXPOSE 3000
CMD [ "node", "index.js" ]