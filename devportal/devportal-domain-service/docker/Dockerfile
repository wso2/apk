FROM apkbase:0.0.1-SNAPSHOT

ARG USER=wso2apk
ARG USER_HOME=/home/${USER}

COPY docker-entrypoint.sh ${USER_HOME}
ADD devportal ${USER_HOME}/devportal
RUN chown -R ${USER} ${USER_HOME}/devportal
RUN chown ${USER} /home/${USER}/docker-entrypoint.sh

EXPOSE 9443
USER wso2apk
WORKDIR ${USER_HOME}

ENTRYPOINT ["sh", "/home/wso2apk/docker-entrypoint.sh"]