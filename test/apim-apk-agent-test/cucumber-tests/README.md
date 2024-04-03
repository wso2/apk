# WSO2 APK Cucumber Based Integration Tests

This folder contains APK cucumber integration tests that is used to test APK product capabilities

## Pre-requisites

1. [Helm](https://helm.sh/docs/intro/install/)
2. [Kubernetes Client](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
3. [Kubernetes Cluster](https://kubernetes.io/docs/setup)

## Overview

This test module use the [Cucumber](https://cucumber.io/) framework and [Gherkin](https://cucumber.io/docs/gherkin/) Syntax to write the integration tests.

## Writing Tests

The tests are written using the Cucumber syntax in feature files located under the `src/test/resources/tests` directory. Each feature file represents a set of related test scenarios written in the Gherkin language.

To create a new feature, follow these steps:

1. Create a new feature file in the `src/test/resources/tests` directory with a `.feature` extension.

2. Write your test scenarios in Gherkin syntax.

3. Step definitions are written in Java and can be found under the `src/test/java/org/wso2/apk/integration` directory. You may need to add new step definitions to support your new scenarios.

4. Add the new feature file to the `src/resources/testng.xml` file to run the tests.

## Run or Debug the Integration Tests Locally

1. Setup the deployment namespace.

    ```bash
    kubectl create namespace apk
    ```

2. Go to the `helm-charts/` directory and run following commands to install the APK components.

    ```bash
    helm repo add bitnami https://charts.bitnami.com/bitnami
    helm repo add jetstack https://charts.jetstack.io
    helm dependency build
    helm install apk . -n apk
    ```

3. Port forward router-service to use localhost.

    ```bash
    kubectl port-forward svc/apk-wso2-apk-router-service -n apk 9095:9095
    ```

4. Add the following DNS mappings to `/etc/hosts` file.

    ```bash
    IP=127.0.0.1
    sudo echo "$IP idp.am.wso2.com" | sudo tee -a /etc/hosts
    sudo echo "$IP api.am.wso2.com" | sudo tee -a /etc/hosts
    sudo echo "$IP default.gw.wso2.com" | sudo tee -a /etc/hosts
    sudo echo "$IP default.sandbox.gw.wso2.com" | sudo tee -a /etc/hosts
    ```
   
5. Go to the `test/apim-apk-agent-test/apim-cp-helm-chart` directory and run following commands to install the APIM CP component.

(Optional)If you are on a Mac, need to add following changes to values.yaml file in apim-cp-helm-chart directory.
```bash
   runAsUser: 802
   registry: "docker.io"
   repository: "sampathrajapakse/wso2am"
   digest: "sha256:9b269e0f3b9b23f48320d16235624acbbfa455aeda6592aa4d51631c85652c2a"
```
Then do the following command to deploy APIM CP component.

```bash
    helm dependency build
    helm install apim . -n apk
```

6. Port forward router-service to use localhost.

    ```bash
    kubectl port-forward svc/apim-wso2am-cp-1-service -n apk 9443:9443
    ```
   
7. Go to the `test/apim-apk-agent-test/agent-helm-chart` directory and run following commands to install the APIM APK Agent components.

    ```bash
    helm dependency build
    helm install apim-apk-agent . -n apk
    ```

8. Run or debug integration tests.

   - Run following command from `test/apim-apk-agent-test/cucumber-tests` directory to run the integration tests.

       ```bash
       gradle runTests
       ```
   - To run a single test, update the `test/apim-apk-agent-test/cucumber-tests/src/resources/testng.xml` file with the required feature file and run the above command.

   - Run the `gradle runTests --debug-jvm` command and attach the debugger in the IDE to debug the integration tests.
