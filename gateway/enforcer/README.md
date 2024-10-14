# WSO2 APK - Enforcer

This guide has information to setup enforcer component for development.

## Prerequisites

The following should be installed in your development machine.

- [Gradle](https://gradle.org/install/) >= 7.5.1 version
- [Docker](https://docs.docker.com/engine/install/ubuntu/) >= 17.03 version
- [Java](https://adoptium.net/installation/) >= 17-jdk version

## Setting up for debugging

1. Make following changes to `helm-chart/templates/data-plane/gateway-components/gateway-runtime/gateway-runtime-deployment.yaml` file.

    ```diff
    containers:
      - name: enforcer	
        image: {{ .Values.wso2.apk.dp.gatewayRuntime.deployment.enforcer.image }}	
        imagePullPolicy: {{ .Values.wso2.apk.dp.gatewayRuntime.deployment.enforcer.imagePullPolicy }}	
        ports:	
          - containerPort: 8081	
            protocol: "TCP"	
          - containerPort: 9001	
            protocol: "TCP"	
    +     - containerPort: 5006	
    +       protocol: "TCP"
    ...
    ...
          - name: JAVA_OPTS	
    -       value: -Dhttpclient.hostnameVerifier=AllowAll -Xms512m -Xmx512m -XX:MaxRAMFraction=2
    +       value: -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5006 -Dhttpclient.hostnameVerifier=AllowAll -Xms512m -Xmx512m -XX:MaxRAMFraction=2
    ```

2. Start WSO2 API Platform for K8s in you local k8s cluster.

3. Port forward the port 5006.

    ```bash
    kubectl port-forward <apk-gateway-runtime-deployment-pod-name> -n apk 5006:5006
    ```

4. Start debugging from port 5006 in IntelliJ IDEA.
