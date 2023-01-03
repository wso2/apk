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

- **Control Plane** - The APK control plane. It provides API Management capabilities, marketplace capabilities along with domain services and web applications. It consists of the following sub-components:

   - **Back Office** - Responsible for configuring the portal aspects of API including descriptions, documents, images, etc. Also, manage API visibility and lifecycle. The backend component was developed using **Ballerina**

   - **Dev Portal** - Responsible for API consumer interaction. API consumers can discover APIs, read documents, try them out and eventually subscribe to and consume APIs. The backend component was developed using **Ballerina**
  
   - **Admin Portal** - Responsible for configuring rate limit policies, key management services, and other administrative tasks. Backend components developed using **Ballerina**
  
   - **Management Server** - Responsible for communication with data planes and pushing updates. Backend components developed using **Go lang** 

- **Data Plane** - The APK data plane. It provides API runtime capabilities such as gateway, rate-limiting services, and runtime management. It consists of the following sub-components:
  
   - **Runtime Manager** - Responsible for configuring the runtime aspects of API including API endpoints, discovering Kubernetes services, and converting them into APIs, etc. The backend component was developed using **Ballerina**
  
   - **Management Client** - Responsible for communication with the management server(control plane) to push/pull updates and maintain connectivity between the data plane and the control plane. The backend component was developed using **Go** 
  
   - **API Gateway - Router** - Router will intercept incoming API traffic and apply quality of services such as authentication, authorization, and rate limiting. The router uses the **Envoy Proxy** as the core component that does the traffic routing. Required additional extensions were developed using **C++**

   - **API Gateway - Enforcer** - The Enforcer is the component that enforces the API management capabilities such as security, Rate Limiting, analytics, validation and etc. When the Router receives a request, it forwards that request to the Enforcer in order to perform the additional QoS. Plugins were developed using **Java** 
 
   - **Identity Platform** - Responsible for authentication and authorization happens in the data plane.

## Test Product APIs
WSO2 APK comes with Postman collections to test product APIs and developers can use collection of API requests and configure them to test different scenarios. For example, they can reuse available requests to verify that the API returns the correct responses for different requests.
These tests will users t identify potential issues or bugs that may need to be addressed before using it. 
Please refer [Postman Tests](https://github.com/wso2/apk/tree/main/test/postman-tests) section of the repo for more information about tests and test artifacts.

## Getting Started

WSO2 API Kubernetes Platform has released following docker images in the WSO2 public docker hub.

Adapter: [wso2/adapter:0.0.1-m1](https://hub.docker.com/r/wso2/adapter)
Gateway Enforcer: [wso2/choreo-connect-enforcer:1.1.0-ubuntu](https://hub.docker.com/r/wso2/choreo-connect-enforcer)
Gatewary Router: [wso2/choreo-connect-router:1.1.0](https://hub.docker.com/r/wso2/choreo-connect-router)
Management Server: [wso2/management-server:0.0.1-m1](https://hub.docker.com/r/wso2/management-server)
Runtime DS: wso2/runtime-domain-service:0.0.1-m1
Admin DS: [wso2/admin-domain-service:0.0.1-m1](https://hub.docker.com/r/wso2/admin-domain-service)
BackOffice DS: [wso2/backoffice-domain-service:0.0.1-m1](https://hub.docker.com/r/wso2/backoffice-domain-service)
BackOffice Internal DS: [wso2/backoffice-internal-domain-service:0.0.1-m1](https://hub.docker.com/r/wso2/backoffice-internal-domain-service)
Devportal DS: [wso2/devportal-domain-service:0.0.1-m1](https://hub.docker.com/r/wso2/devportal-domain-service)

### Before you begin...

* Install [Helm](https://helm.sh/docs/intro/install/)
  and [Kubernetes client](https://kubernetes.io/docs/tasks/tools/install-kubectl/) <br><br>

* An already setup [Kubernetes cluster](https://kubernetes.io/docs/setup). If you want to run it on the local you can use Minikube or Kind or a similar software.<br><br>

* Install [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/). If you are using Minikube you can install ingress by running ```minikube addons enable ingress```<br><br>

### Steps to deploy APK DS servers and CloudNativePG

```HELM-HOME``` = apk/helm

1. Execute ``` helm repo add bitnami https://charts.bitnami.com/bitnami ``` and ```helm repo add jetstack https://charts.jetstack.io```
2. Clone the repo and cd into the ```HELM-HOME``` folder.
3. Execute ``` helm dependency build ``` command to download the dependent charts.
4. Now execute ```helm install apk-test . -n apk``` to install the APK components.
5. Verify the deployment by executing ```kubectl get pods -n apk```

### To Access Deployment through local machine

#### Using ingress
   Execute ``` kubectl get ing -n apk ``` command

#### Minikube Flow
   Execute ``` minikube tunnel ``` command

#### Rancher Flow

## Issue management
We use GitHub to track all of our bugs and feature requests. Each issue we track has a variety of metadata:
- **Epic**. An epic represents a feature area for APK as a whole. Epics are fairly broad in scope and are basically product-level things.Each issue is ultimately part of an epic.
- **Milestone**. Each issue is assigned a milestone. This is 0.1, 0.2, ..., or 'Nebulous Future'. The milestone indicates when we think the issue should get addressed.
- **Priority**. Each issue has a priority which is represented by the column in the [Prioritization]() project. Priority can be one of P1, P2, or >P2. The priority indicates how important it is to address the issue within the milestone. P1 says that themilestone cannot be considered achieved if the issue isn't resolved.

