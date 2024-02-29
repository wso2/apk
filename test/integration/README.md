# WSO2 APK - Integration Tests

This module contains the integration tests for the WSO2 APK. Following instructions will guide you to run or debug the integration tests.

## Pre-requisites

1. [Helm](https://helm.sh/docs/intro/install/)
2. [Kubernetes Client](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
3. [Kubernetes Cluster](https://kubernetes.io/docs/setup)
4. [Golang](https://go.dev/doc/install)

## Run Integration Tests using Gradle

If you have setup `Kind` and wish to run the integration tests using Gradle, then execute following command to run all the integration tests.

```bash
./gradlew integration_test
```

## Run or Debug the Integration Tests Locally

1. Setup the deployment namespace.
    
    ```bash
    kubectl create namespace apk-integration-test
    ```

2. Go to the `apk/helm-charts` directory and run following commands to install the APK components.

    ```bash
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo add jetstack https://charts.jetstack.io
    helm dependency build
    helm install apk-test-setup . -n apk-integration-test
    ```

3. Port forward router-service to use localhost.

    ```bash
    kubectl port-forward svc/apk-test-setup-wso2-apk-gateway-service -n apk-integration-test 9095:9095
    ```

4. Add all DNS mappings to `/etc/hosts` file. Refer to `scripts/run-tests.sh` file for the domain names.

    ```bash
    IP=127.0.0.1
    sudo echo "$IP <DomainName>" | sudo tee -a /etc/hosts
    ```

5. Run or debug integration tests.

    > **Note**
    >
    > If you need to run only one test case, then change `integration_test.go` file with the test name you want to run.
    >
    > ```diff
    > - cSuite.Run(t, tests.IntegrationTests)
    > + cSuite.Run(t, []suite.IntegrationTest{tests.<TestName>})
    > ```

    - Run following command from `apk/test/integration` directory to run the integration tests.

        ```bash
        go test -v integration_test.go
        ```

    - Click on `debug test` option in the IDE to debug the integration tests.
