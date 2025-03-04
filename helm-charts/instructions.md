# WSO2 APK - Helm Charts

This file contains the instructions to manage the `helm-charts` component of WSO2 APK. Following are the instructions available in this file.

1. [Add Configurations to log-conf.yaml from values.yaml](#add-configurations-to-log-confyaml-from-valuesyaml)
2. [Generate README.md File](#generate-readmemd-file)

## Add Configurations to log-conf.yaml from values.yaml

`log-conf.yaml` is used to set the configurations in APK as a `ConfigMap`. Also, some configurations in `values.yaml` file can be populated to this file during the `helm install`. Following steps describe a sample scenario.

1. Add a sample configuration to `values.yaml` file as follows.

    ```yaml
    wso2:
      apk:
        dp:
          gatewayRuntime:
            customObject:
              property1: true
              property2: customValue
    ```

2. Add the logic to populate the above configurations to `log-conf.yaml` file as follows.

    ```
    {{ if and .Values.wso2.apk.dp.gatewayRuntime.customObject }}
    [customObject]
    property1 = {{ .Values.wso2.apk.dp.gatewayRuntime.customObject.property1 | default true }}
    property2 = {{ .Values.wso2.apk.dp.gatewayRuntime.customObject.property2 | default "customValue" }}
    {{ end }}
    ```

    > **Note**
    >
    > For advanced helm chart template functions, refer [helm documentation](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/).

3. Run following command from `<APK_HOME>/helm-charts` directory level and verify whether `test.yaml` file is populated correctly.

    ```bash
    helm template test . > test.yaml
    ```

## Generate README.md File.

1. Download and install the latest [helm-docs](https://github.com/norwoodj/helm-docs) executable.

2. Run `helm-docs --version` to verify the installation.

3. Add the relevant changes to `values.yaml.template` file.

4. Run the follwoing command from `<APK_HOME>/helm-charts` directory level.

    ```bash
    helm-docs --values-file values.yaml.template  --document-dependency-values --sort-values-order file
    ```