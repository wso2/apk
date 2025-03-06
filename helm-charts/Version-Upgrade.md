# Updating APK Version

This guide outlines the process of upgrading from APK v1.2.0 installation to APK v1.3.0 installation.

## In-Place Upgrade

The in-place upgrade process transitions your existing APK v1.2.0 installation to APK v1.3.0. Prior to implementing these steps in a production environment, it is advised to apply and validate them in lower environments.

- Ensure APK v1.2.0 is currently installed in the cluster.

    **Note:** The steps provided below assume that APK v1.2.0 is installed in the `default` namespace under the release name `apk`. Replace the dot (.) with the appropriate APK v1.3.0 Helm chart name and version, which is `wso2apk/apk-helm --version 1.3.0`.

- Install/Update CRDs for APK v1.3.0.

    ```bash
    (helm template apk . -f crds-upgrade-values.yaml -n default && helm show crds .) > apk-v1.3.0-crds.yaml

    kubectl apply -f apk-v1.3.0-crds.yaml
    ```

- Upgrade the existing APK v1.2.0 installation to APK v1.3.0.

    ```bash
    helm upgrade --reuse-values apk . -f ./in-place-upgrade-values.yaml --set skipCrds=true
    ```

These steps will seamlessly transition your APK installation to the latest version, ensuring continued functionality and compatibility.