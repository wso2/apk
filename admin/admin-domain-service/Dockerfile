FROM apkbase:0.0.1-SNAPSHOT

ARG USER=wso2apk
ARG USER_HOME=/home/${USER}

COPY target/bin/admin_service.jar ${USER_HOME}
COPY Config.toml ${USER_HOME}

RUN chown wso2apk /home/wso2apk/admin_service.jar

EXPOSE 9443
USER wso2apk
WORKDIR ${USER_HOME}

CMD java -jar admin_service.jar
