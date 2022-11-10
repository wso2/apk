
FROM eclipse-temurin:11.0.17_8-jre-alpine
ARG USER=wso2apk
ARG USER_ID=802
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=802
ARG USER_HOME=/home/${USER}
RUN \
    addgroup -S -g ${USER_GROUP_ID} ${USER_GROUP} \
    && adduser -S -u ${USER_ID} -h ${USER_HOME} -G ${USER_GROUP} ${USER}

COPY runtime_domain_service.jar ${USER_HOME}
RUN chown ${USER}:${USER_GROUP} /home/wso2apk/runtime_domain_service.jar
EXPOSE 9443
USER wso2apk
WORKDIR ${USER_HOME}
CMD java -jar runtime_domain_service.jar