
FROM apkbase:0.0.1-SNAPSHOT
ARG USER=wso2apk
ARG USER_ID=802
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=802
ARG USER_HOME=/home/${USER}
COPY runtime_domain_service.jar ${USER_HOME}
RUN chown ${USER}:${USER_GROUP} /home/wso2apk/runtime_domain_service.jar
EXPOSE 9443
USER wso2apk
WORKDIR ${USER_HOME}
CMD java -jar runtime_domain_service.jar