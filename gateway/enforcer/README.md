# WSO2 APK - Enforcer

This guide has information to setup enforcer component for development.

## Prerequisites

The following should be installed in your dev machine.

- [Gradle](https://gradle.org/install/) >= 7.5.1 version
- [Docker](https://docs.docker.com/engine/install/ubuntu/) >= 17.03 version
- [Java](https://adoptium.net/installation/) >= 11-jdk version

## Setting up for debugging

1. Create `enforcer-service.yaml` as follows and copy it to `helm-chart/templates/data-plane/gateway-components/gateway-runtime`.

    ```yaml
    {{ if or .Values.wso2.apk.dp.enabled .Values.wso2.apk.cp.enabled }}
    apiVersion: v1
    kind: Service
    metadata:
        name: {{ template "apk-helm.resource.prefix" . }}-enforcer-service
        namespace : {{ .Release.Namespace }}
    spec:
    # label keys and values that must match in order to receive traffic for this service
        selector:
        {{ include "apk-helm.pod.selectorLabels" (dict "root" . "app" .Values.wso2.apk.dp.gatewayRuntime.appName ) | indent 4}}
        ports:
        - name: endpoint1
          protocol: TCP
          port: 8081
        - name: endpoint2
          protocol: TCP
          port: 9001
        - name: debug
          protocol: TCP
          port: 5006
    {{- end -}}
    ```

2. Make following changes to `helm-chart/templates/data-plane/gateway-components/gateway-runtime/gateway-runtime-deployment.yaml` file.

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

3. Start WSO2 APK.

4. Port forward the port 5006.

    ```bash
    kubectl port-forward <apk-gateway-runtime-deployment-pod-name> -n apk 5006:5006
    ```

5. Start debugging from port 5006 in IntelliJ IDEA.