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

4. Add the new feature file to the `apk/test/cucumber-tests/src/test/resources/testng.xml` file to run the tests.

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

4. Add the following DNS mappings to `/etc/hosts` file.

    ```bash
    IP=127.0.0.1
    sudo echo "$IP idp.am.wso2.com" | sudo tee -a /etc/hosts
    sudo echo "$IP api.am.wso2.com" | sudo tee -a /etc/hosts
    sudo echo "$IP default.gw.wso2.com" | sudo tee -a /etc/hosts
    sudo echo "$IP default.sandbox.gw.wso2.com" | sudo tee -a /etc/hosts
    ```

5. Run or debug integration tests.

   - Run following command from `apk/test/cucumber-tests` directory to run the integration tests.

       ```bash
             ./gradlew runTests
       ```
    - To run a single test, update the `apk/test/cucumber-tests/src/test/resources/testng.xml` file with the required feature file and run the above command.

     - Run the `./gradlew runTests --debug-jvm` command and attach the debugger in the IDE to debug the integration tests.

### Advanced run options

- Run a single feature file:

    ```bash
    ./gradlew runTests -Dcucumber.features=src/test/resources/tests/api/BackendAPIKeyAuthNew.feature
    ```

- Run multiple feature paths (comma-separated):

    ```bash
    ./gradlew runTests -Dcucumber.features=src/test/resources/tests/api/EndpointNew.feature,src/test/resources/tests/api/HeaderModifierNew.feature
    ```

- Filter by scenario name (regex):

    ```bash
    ./gradlew runTests -Dcucumber.filter.name='^Testing API level and resource.*'
    ```

- Filter by tags (Cucumber v7 AND/OR syntax):

    ```bash
    ./gradlew runTests -Dcucumber.filter.tags='@smoke'
    ./gradlew runTests -Dcucumber.filter.tags='@smoke and not @wip'
    ./gradlew runTests -Dcucumber.filter.tags='@smoke or @regression'
    ```

- Rerun only failed scenarios from the last run (failures are written to `build/cucumber-rerun.txt`):

    ```bash
    ./gradlew runTests -Dcucumber.features=@build/cucumber-rerun.txt
    ```