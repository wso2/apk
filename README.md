# WSO2 APK - API Platform for Kubernetes

<a href="https://wso2.com/">  <img src="https://raw.githubusercontent.com/wso2/apk/main/logo/wso2-logo.png" alt="WSO2 logo" title="WSO2" height="100" width="180" /></a>

---
Cloud native API management refers to the use of API management tools and practices that are designed to be natively integrated with cloud computing environments, such as Kubernetes. These tools and practices help organizations to manage and govern their APIs in a more effective and efficient way, leveraging the benefits of cloud computing such as scalability, reliability, and cost-effectiveness.  Cloud native API management is also a way of designing, building, and deploying APIs in a cloud environment. There are several characteristics that are commonly associated with cloud native API management:

WSO2 APK, designed to help you build, deploy, and manage APIs in a cloud environment. Our platform is built on top of a microservices architecture and uses containerization technologies to ensure scalability and flexibility. With features like automatic failover and load balancing, our APK platform is designed to be highly available and able to handle large numbers of API requests without performance degradation. We've also added support for continuous delivery and deployment, so you can quickly and easily push updates to your API services.


Some characteristics of APK
- Scalability: Designed to scale up and down based on demand, allowing them to handle large numbers of API requests without performance degradation.
- High availability: Designed to be highly available, with features like automatic failover and load balancing to ensure that API services are always available to clients.
- Elasticity: Designed to be elastic, meaning that they can quickly and easily adapt to changes in demand and workload.
- Microservices architecture: Built using a microservices architecture, which allows for more flexible and scalable deployment of API services.
- Containerization: Use containerization technologies to package and deploy product API services/ implementations in a cloud environment.
- Continuous delivery and deployment: Support continuous delivery and deployment, allowing developers to quickly and easily push updates to product.



For more information about APK release planning and project management information, visit [APK Project Dashboard](https://github.com/orgs/wso2/projects/80/)

For in-depth information about WSO2 API Management Platform, visit [WSO2 API Management](https://wso2.com/api-manager/)

To ask questions and get assistance from our community, visit [WSO2 Discord](https://discord.com/invite/Xa5VubmThw?utm_source=wso2-dev&utm_medium=link&utm_campaign=wso2-dev_link_from-dev-homepage_221002)

To learn how to participate in our overall community, visit [our community page](https://wso2.com/community/)

In this README:
- [Introduction](#introduction)
- [Getting Started](#getting-started)
- [Issue management](#issue-management)

You'll find many other useful documents on our [Documentation](https://wso2.com/documentation/).
## Introduction
[APK](https://github.com/wso2/apk) is an open-source platform for providing complete API Management capabilities on top of the Kubernetes cluster management platform.

APK is composed of these components:

___

<img src="https://raw.githubusercontent.com/wso2/apk/main/logo/architecture.png" alt="API Architecture" title="APKArchitecture" />

___


<!---
   - **Runtime Manager** - Responsible for configuring the runtime aspects of API including API endpoints, discovering Kubernetes services, and converting them into APIs, etc. The backend component was developed using **Ballerina**
  
   - **Management Client** - Responsible for communication with the management server(control plane) to push/pull updates and maintain connectivity between the data plane and the control plane. The backend component was developed using **Go** 
  -->
- **Data Plane** - The APK data plane. It provides API runtime capabilities such as gateway, rate-limiting services, and runtime management. It consists of the following sub-components:

   - **Config and Deploy APIs** - Responsible for configuring the runtime aspects of API including API endpoints, rate limiting policies, and converting API schemas into API configurations, etc. This API implementation done using **Ballerina**
 
   - **API Gateway - Router** - Router will intercept incoming API traffic and apply quality of services such as authentication, authorization, and rate limiting. The router uses the **Envoy Proxy** as the core component that does the traffic routing. Required additional extensions were developed using **C++**

   - **API Gateway - Enforcer** - The Enforcer is the component that enforces the API management capabilities such as security, Rate Limiting, analytics, validation and etc. When the Router receives a request, it forwards that request to the Enforcer in order to perform the additional QoS. Plugins were developed using **Java** 
 
   - **Identity Platform** - Responsible for authentication and authorization happens in the data plane.

## Test Product APIs
WSO2 APK comes with Postman collections to test product APIs and developers can use collection of API requests and configure them to test different scenarios. For example, they can reuse available requests to verify that the API returns the correct responses for different requests.
These tests will allow users to identify potential issues or bugs that may need to be addressed before using it. 
Please refer [Postman Tests](https://github.com/wso2/apk/tree/main/test/postman-tests) section of the repo for more information about tests and test artifacts.

## Getting Started

WSO2 API Kubernetes Platform has released following docker images in the WSO2 public docker hub.

* Adapter: [wso2/adapter:0.0.1-m12](https://hub.docker.com/r/wso2/adapter)
* Gateway Enforcer: [wso2/enforcer:0.0.1-m12](https://hub.docker.com/r/wso2/enforcer/tags)
* Gatewary Router: [wso2/router:0.0.1-m12](https://hub.docker.com/r/wso2/router)
* IDP DS: [wso2/idp-domain-service:0.0.1-m12](https://hub.docker.com/r/wso2/devportal-domain-service)
* IDP UI: [wso2/idp-ui:0.0.1-m12](https://hub.docker.com/r/wso2/devportal-domain-service)
* Ratelimiter: [wso2/ratelimiter:0.0.1-m12](https://hub.docker.com/r/wso2/ratelimiter)
* Config Deployer: [wso2/config-deployer:0.0.1-m12](https://hub.docker.com/r/wso2/config-deployer-service/)

### Before you begin...

* Install [Helm](https://helm.sh/docs/intro/install/) (3.11.x)
  and [Kubernetes client](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

* Setup a [Kubernetes cluster](https://kubernetes.io/docs/setup). If you want to run it on the local you can use Minikube or Kind or a similar software.

* Setup deployment namespace.
    ```bash
    kubectl create namespace <namespace>
    ```

### Steps to deploy APK DS servers and CloudNativePG

```HELM-HOME``` = apk/helm-charts

1. Execute `helm repo add bitnami https://charts.bitnami.com/bitnami` and `helm repo add jetstack https://charts.jetstack.io`.
2. Clone the repo and cd into the `HELM-HOME` folder.
3. Execute `helm dependency build` command to download the dependent charts.
4. Now execute `helm install apk-test . -n apk` to install the APK components.

    > **Optional**
    >
    > To deploy control plane components only use `--set wso2.apk.dp.enabled=false`
    >
    > To deploy data plane components only use `--set wso2.apk.cp.enabled=false`

5. Verify the deployment by executing ```kubectl get pods -n apk```

### To Access Deployment through local machine

- You can either, identify the router-service external IP address to invoke the API through the APK gateway
    ```bash
    kubectl get svc -n apk | grep router-service
    ```

- or, port forward to router-service to use localhost.
    ```bash
    kubectl port-forward svc/apk-test-wso2-apk-router-service -n apk 9095:9095
    ```

## Quick Start APK with Kubernetes client
Follow the instruction below to deploy an API using the `kubectl`.

1. Create API CR and create production and/or sandbox HTTPRoute CRs, and service for the API backend. You can find a sample CR set in `developer/tryout/samples/` folder in this repository.

2. Apply CRs to kubernetes API server using the kubectl.
    ```bash
    kubectl apply -f developer/tryout/samples/ -n apk
    ```
    > **Note**
    >
    > Services should be created in a different namespace than APK or Kubernetes System namespaces.
    >
    > APIs should be created in the APK deployment namespace.
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

## Run domain services APIs in APK with postman
[Test Postman collection](#test/postman-tests/README.md)

## Build APK Components

### Pre-requisites
1. Install Java JDK 11.
2. Install Gradle(7.5.1).
3. Install Ballerina Ballerina version: 2201.3.1 (Swan Lake Update 3).
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
- **Epic**. An epic represents a feature area for APK as a whole. Epics are fairly broad in scope and are basically product-level things.Each issue is ultimately part of an epic.
- **Milestone**. Each issue is assigned a milestone. This is 0.1, 0.2, ..., or 'Nebulous Future'. The milestone indicates when we think the issue should get addressed.
- **Priority**. Each issue has a priority which is represented by the column in the [Prioritization]() project. Priority can be one of P1, P2, or >P2. The priority indicates how important it is to address the issue within the milestone. P1 says that themilestone cannot be considered achieved if the issue isn't resolved.

