
FROM apkbase:0.0.1-SNAPSHOT
COPY runtime_domain_service.jar ${USER_HOME}
RUN chown ${USER}:${USER_GROUP} /home/wso2apk/runtime_domain_service.jar
EXPOSE 9443
USER wso2apk
WORKDIR ${USER_HOME}
CMD java -jar runtime_domain_service.jar