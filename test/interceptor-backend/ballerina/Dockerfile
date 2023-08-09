FROM --platform=linux/amd64 alpine:3.18.2

ENV LANG='en_US.UTF-8' LANGUAGE='en_US:en' LC_ALL='en_US.UTF-8'

# install JDK Dependencies
RUN apk update \
    && apk add --no-cache tzdata musl-locales musl-locales-lang gcompat \
    && rm -rf /var/cache/apk/*

ENV JAVA_VERSION jdk-11.0.20+8

RUN set -eux; \
    ARCH="$(apk --print-arch)"; \
    case "${ARCH}" in \
       amd64|x86_64) \
         ESUM='0b7e894af823484f270031ed6626f78504a7a44f8173d67433e14321d015e669'; \
         BINARY_URL='https://github.com/adoptium/temurin11-binaries/releases/download/jdk-11.0.20%2B8/OpenJDK11U-jdk_x64_alpine-linux_hotspot_11.0.20_8.tar.gz'; \
         ;; \
       *) \
         echo "Unsupported arch: ${ARCH}"; \
         exit 1; \
         ;; \
    esac; \
	  wget -O /tmp/openjdk.tar.gz ${BINARY_URL}; \
	  echo "${ESUM} */tmp/openjdk.tar.gz" | sha256sum -c -; \
	  mkdir -p /opt/java/openjdk; \
	  tar --extract \
	      --file /tmp/openjdk.tar.gz \
	      --directory /opt/java/openjdk \
	      --strip-components 1 \
	      --no-same-owner \
	  ; \
    rm -rf /tmp/openjdk.tar.gz;

ENV JAVA_HOME=/opt/java/openjdk \
    PATH="/opt/java/openjdk/bin:$PATH"

RUN \
    addgroup -S -g 10000 ineterceptor \
    && adduser -S -u 1001 -h /home/ineterceptor -G ineterceptor ineterceptor \
    && chown -R ineterceptor:ineterceptor /home/ineterceptor
COPY /target/bin/interceptor_service.jar /home/ineterceptor/interceptor_service.jar
EXPOSE 8443 8444 8445
CMD java -jar /home/ineterceptor/interceptor_service.jar