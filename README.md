# WSO2 Kubernetes Gateway

<a href="https://wso2.com/">  <img src="https://raw.githubusercontent.com/wso2/apk/main/logo/wso2-logo.png" alt="WSO2 logo" title="WSO2" height="100" width="180" /></a>

---
Introducing WSO2 Kubernetes Gateway, a cutting-edge API management solution designed to leverage the power of Kubernetes for seamless and scalable deployments. Kubernetes Gateway harnesses Kubernetes' native features, enabling automatic scaling based on load and configurable parameters, utilizing rich Kubernetes metrics.

At the core of Kubernetes Gateway's robust gateway solution is the meticulously selected Envoy technology, known for exceptional performance, lightweight nature, and perfect compatibility within the Kubernetes Gateway ecosystem. Kubernetes Gateway extends beyond traditional gateways with purpose-built extensions addressing specific API management use cases. Some of these extensions have been contributed back to the Envoy community, reflecting our commitment to collaborative innovation.

WSO2 Kubernetes Gateway adheres to the Kubernetes Gateway API specification, an open-source project managed by the SIG-NETWORK community. This specification introduces vital resources such as GatewayClass, Gateway, HTTPRoute, TCPRoute, and Service, augmenting service networking capabilities in Kubernetes. By adhering to this specification, WSO2 Kubernetes Gateway seamlessly integrates with Kubernetes service networking, leveraging expressive and extensible interfaces to enhance API management functionality within Kubernetes deployments.

Some characteristics of Kubernetes Gateway
- Kubernetes Gateway's microservices architecture offers advantages such as easy scalability and seamless upgrades, harnessing the benefits of the architecture for agility and flexibility.
- The separation of the control plane and data plane in Kubernetes Gateway allows users to integrate any control plane of their choice, providing maximum flexibility and customization.
- Kubernetes Gateway is an evolving open-source solution that delivers advanced API management capabilities and is designed for cloud-native architectures, seamlessly integrating with Kubernetes.
- With seamless CI/CD integration, Kubernetes Gateway supports a streamlined GitOps approach for efficient deployment and management of APIs.
- Kubernetes Gateway aims to provide API marketplace capabilities, enabling sharing, discovery, and reusability of APIs while focusing on efficient governance and administration.
- With its Kubernetes-native approach, exceptional characteristics, microservices architecture, and commitment to collaboration and innovation, Kubernetes Gateway sets a new standard for API management.

For more information about Kubernetes Gateway release planning and project management information, visit [APK Project Dashboard](https://github.com/orgs/wso2/projects/80/)

For in-depth information about WSO2 API Management Platform, visit [WSO2 API Management](https://wso2.com/api-manager/)

To ask questions and get assistance from our community, visit [WSO2 Discord](https://discord.com/invite/Xa5VubmThw?utm_source=wso2-dev&utm_medium=link&utm_campaign=wso2-dev_link_from-dev-homepage_221002)

To learn how to participate in our overall community, visit [our community page](https://wso2.com/community/)

In this README:
- [Introduction](#introduction)
- [Getting Started](#getting-started)
- [Issue management](#issue-management)

You'll find many other useful documents on our [Documentation](https://wso2.com/documentation/).
## Introduction
[Kubernetes Gateway](https://github.com/wso2/apk) is an open-source platform for providing complete API Management capabilities on top of the Kubernetes cluster management platform.

Kubernetes Gateway is composed of these components:

___

<img src="https://raw.githubusercontent.com/wso2/apk/main/logo/architecture.png" alt="API Architecture" title="APKArchitecture" />

___


<!---
   - **Runtime Manager** - Responsible for configuring the runtime aspects of API including API endpoints, discovering Kubernetes services, and converting them into APIs, etc. The backend component was developed using **Ballerina**
  
   - **Management Client** - Responsible for communication with the management server(control plane) to push/pull updates and maintain connectivity between the data plane and the control plane. The backend component was developed using **Go** 
  -->
The Kubernetes Gateway architecture consists of both control plane and data plane components. In the data plane, we have the Config Service, an open API for generating Kubernetes Gateway configurations and Kubernetes API artifacts based on inputs like OpenAPI schema files. The Deployer Service enables the creation of API artifacts within the gateway runtime, requiring a valid access token for secure deployment.

These components efficiently generate configurations and deploy API artifacts within the data plane. The gateway partition comprises the Router, Enforcer, and Management Client. The Router intercepts API traffic, applying QoS policies for optimal performance. The Enforcer handles authentication and authorization, ensuring authorized access. The Management Client configures and synchronizes the Router and Enforcer, ensuring the gateway partition's smooth operation.

The architecture also includes the Rate Limiting Service, which manages rate limits for API calls. The Router communicates with the Rate Limiter to enforce quota compliance. To facilitate distributed counters across gateways, Redis serves as a shared information store for rate limiting.

## Test Product APIs
WSO2 Kubernetes Gateway comes with Postman collections to test product APIs and developers can use collection of API requests and configure them to test different scenarios. For example, they can reuse available requests to verify that the API returns the correct responses for different requests.
These tests will allow users to identify potential issues or bugs that may need to be addressed before using it. 
Please refer [Postman Tests](https://github.com/wso2/apk/tree/main/test/postman-tests) section of the repo for more information about tests and test artifacts.

## Getting Started

To tryout Kubernetes Gateway please refer to this [document](https://apk.docs.wso2.com/en/latest/get-started/quick-start-guide/).


### Before you begin...

* Install [Helm](https://helm.sh/docs/intro/install/) (3.11.x)
  and [Kubernetes client](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

* Setup a [Kubernetes cluster](https://kubernetes.io/docs/setup). If you want to run it on the local you can use Minikube or Kind or a similar software.

* Setup deployment namespace.
    ```bash
    kubectl create namespace <namespace>
    ```

### Steps to deploy Kubernetes Gateway DS servers and CloudNativePG

```HELM-HOME``` = apk/helm-charts

1. Execute `helm repo add bitnami https://charts.bitnami.com/bitnami` and `helm repo add jetstack https://charts.jetstack.io`.
2. Clone the repo and cd into the `HELM-HOME` folder.
3. Execute `helm dependency build` command to download the dependent charts.
4. Now execute `helm install apk-test .` to install the Kubernetes Gateway components.

    > **Optional**
    >
    > To deploy control plane components only use `--set wso2.apk.dp.enabled=false`
    >
    > To deploy data plane components only use `--set wso2.apk.cp.enabled=false`

5. Verify the deployment by executing ```kubectl get pods```

### To Access Deployment through local machine

- You can either, identify the gateway-service external IP address to invoke the API through the APK gateway
    ```bash
    kubectl get svc | grep gateway-service
    ```

- or, port forward to router-service to use localhost.
    ```bash
    kubectl port-forward svc/apk-test-wso2-apk-gateway-service 9095:9095
    ```

## Quick Start Kubernetes Gateway with Kubernetes client
Follow the instruction below to deploy an API using the `kubectl`.

1. Create API CR and create production and/or sandbox HTTPRoute CRs, and service for the API backend. You can find a sample CR set in `developer/tryout/samples/` folder in this repository.

2. Apply CRs to kubernetes API server using the kubectl.
    ```bash
    kubectl apply -f developer/tryout/samples/
    ```
    > **Note**
    >
    > Services should be created in a different namespace than Kubernetes Gateway or Kubernetes System namespaces.
    >
    > APIs should be created in the Kubernetes Gateway deployment namespace.
    >
    > Provide the router service external ip to `{router_service}` in below commands.

3. Get a token to invoke the System API.
    ```bash
    ACCESS_TOKEN=$(curl --location --request POST "https://{router_service}:9095/oauth2/token" \
    --header "Host: idp.am.wso2.com" \
    --header "Authorization: Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==" \
    --header "Content-Type: application/x-www-form-urlencoded" \
    --data-urlencode "grant_type=client_credentials" | jq -r ".access_token")
    ```

4. List the created API and retrieve API's `id`.
    ```bash
    curl --location --request GET "https://{router_service}:9095/api/runtime/apis" \
    --header "Host: api.am.wso2.com" \
    --header "Authorization: Bearer $ACCESS_TOKEN"
    ```

5. Get a token to invoke the created API. Provide the API's `id` to `{api_id}` in below command.
    ```bash
    INTERNAL_KEY=$(curl --location --request POST "https://{router_service}:9095/api/runtime/apis/{api_id}/generate-key" \
    --header "Content-Type: application/json" \
    --header "Accept: application/json" \
    --header "Host: api.am.wso2.com" \
    --header "Authorization: Bearer $ACCESS_TOKEN" | jq -r ".apikey")
    ```

6. Invoke the API.
    ```bash
    curl --location --request GET "https://{router_service}:9095/http-bin-api/1.0.8/get" \
    --header "HOST: gw.wso2.com" \
    --header "Internal-Key: $INTERNAL_KEY"
    ```

## Run domain services APIs in Kubernetes Gateway with postman
[Test Postman collection](#test/postman-tests/README.md)

## Build Kubernetes Gateway Components

### Pre-requisites
1. Install Java JDK 17.
2. Install Gradle(7.6).
3. Install Ballerina Ballerina version: 2201.10.2 (Swan Lake Update 10).
4. Install Go.
5. Install Lua.
6. Docker Runtime Up and Running.

### Build all components

Run `apk/build-apk.sh` file.
```bash
sh build-apk.sh
```

### Build single component

For example: building Runtime Domain Service
```bash
cd runtime/runtime-domain-service
gradle build
```

## Issue management
We use GitHub to track all of our bugs and feature requests. Each issue we track has a variety of metadata:
- **Epic**. An epic represents a feature area for Kubernetes Gateway as a whole. Epics are fairly broad in scope and are basically product-level things.Each issue is ultimately part of an epic.
- **Milestone**. Each issue is assigned a milestone. This is 0.1, 0.2, ..., or 'Nebulous Future'. The milestone indicates when we think the issue should get addressed.
- **Priority**. Each issue has a priority which is represented by the column in the [Prioritization]() project. Priority can be one of P1, P2, or >P2. The priority indicates how important it is to address the issue within the milestone. P1 says that themilestone cannot be considered achieved if the issue isn't resolved.

