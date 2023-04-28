FROM --platform=linux/x86_64 adoptopenjdk/openjdk11:jre-11.0.11_9-alpine
RUN \
    addgroup -S -g 10000 ineterceptor \
    && adduser -S -u 1001 -h /home/ineterceptor -G ineterceptor ineterceptor \
    && chown -R ineterceptor:ineterceptor /home/ineterceptor
COPY /target/bin/interceptor_service.jar /home/ineterceptor/interceptor_service.jar
EXPOSE 8443 8444
CMD java -jar /home/ineterceptor/interceptor_service.jar