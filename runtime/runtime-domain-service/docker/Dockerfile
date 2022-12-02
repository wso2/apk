
FROM apkbase:0.0.1-SNAPSHOT
ARG USER=wso2apk
ARG USER_ID=802
ARG USER_GROUP=wso2
ARG USER_GROUP_ID=802
ARG USER_HOME=/home/${USER}
ADD runtime ${USER_HOME}/runtime
COPY docker-entrypoint.sh ${USER_HOME}
RUN chown -R ${USER} ${USER_HOME}/runtime
RUN chown ${USER} /home/${USER}/docker-entrypoint.sh
EXPOSE 9443
USER wso2apk
WORKDIR ${USER_HOME}
ENTRYPOINT ["sh", "/home/wso2apk/docker-entrypoint.sh"]