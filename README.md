# WSO2 APK - API Platform for Kubernetes

<a href="https://wso2.com/">  <img src="https://raw.githubusercontent.com/wso2/apk/main/logo/wso2-logo.png" alt="WSO2 logo" title="WSO2" height="100" width="180" /></a>

---
A complete API Management platform for modern cloud-native architectures. This is an open platform to design, develop and manage APIs in Kubernetes environments. APK is a cloud-native API Management platform in K8s utilizing native K8s capabilities.
- For more information about APK release planning and project management information, visit [APK Project Dashboard](https://github.com/orgs/wso2/projects/80/)
- For in-depth information about WSO2 API Management Platform, visit [WSO2 API Management](https://wso2.com/api-manager/)- To ask questions and get assistance from our community, visit [WSO2 Discord](https://discord.com/invite/Xa5VubmThw?utm_source=wso2-dev&utm_medium=link&utm_campaign=wso2-dev_link_from-dev-homepage_221002)- To learn how to participate in our overall community, visit [our community page](https://wso2.com/community/)

In this README:
- [Introduction](#introduction)- [Issue management](#issue-management)
You'll find many other useful documents on our [Documentation](https://wso2.com/documentation/).
## Introduction
[APK](https://github.com/wso2/apk) is an open-source platform for providing complete API Management capabilities on top of the Kubernetes cluster management platform.
APK is composed of these components:

___

<img src="https://raw.githubusercontent.com/wso2/apk/main/logo/architecture.png" alt="API Architecture" title="APKArchitecture" />

___

- **Control Plane** - The APK control plane. It provides API Management capabilities, marketplace capabilities along with domain services and web applications. It consists of the following sub-components:

   - **Back Office** - Responsible for configuring the portal aspects of API including descriptions, documents, images, etc. Also, manage API visibility and lifecycle. The backend component was developed using **Ballerina/Java** while the frontend component was developed using **ReactJS**

   - **Dev Portal** - Responsible for API consumer interaction. API consumers can discover APIs, read documents, try them out and eventually subscribe to and consume APIs. The backend component was developed using **Ballerina/Java** and the frontend component was developed using **ReactJS**
  
   - **Admin Portal** - Responsible for configuring rate limit policies, key management services, and other administrative tasks. Backend components developed using **Ballerina/Java** and frontend components developed using **ReactJS**
  
   - **Management Server** - Responsible for communication with data planes and pushing updates. Backend components developed using **Ballerina** and frontend components developed using **ReactJS**

- **Data Plane** - The APK data plane. It provides API runtime capabilities such as gateway, rate-limiting services, and runtime management. It consists of the following sub-components:
  
   - **Runtime Manager** - Responsible for configuring the runtime aspects of API including API endpoints, discovering Kubernetes services, and converting them into APIs, etc. The backend component was developed using **Ballerina** and the frontend component was developed using **ReactJS**
  
   - **Management Client** - Responsible for communication with the management server(control plane) to push/pull updates and maintain connectivity between the data plane and the control plane. The backend component was developed using **Go** 
  
   - **API Gateway - Router** - Router will intercept incoming API traffic and apply quality of services such as authentication, authorization, and rate limiting. The router uses the **Envoy Proxy** as the core component that does the traffic routing. Required additional extensions were developed using **C++**

   - **API Gateway - Enforcer** - The Enforcer is the component that enforces the API management capabilities such as security, Rate Limiting, analytics, validation and etc. When the Router receives a request, it forwards that request to the Enforcer in order to perform the additional QoS. Plugins were developed using **Java** 
 
   - **Identity Platform** - Responsible for authentication and authorization happens in the data plane.

- **Operator** - The component provides user-friendly options to operate the APK platform. Operator implementation is done using **Go**


## Issue management
We use GitHub to track all of our bugs and feature requests. Each issue we track has a variety of metadata:
- **Epic**. An epic represents a feature area for APK as a whole. Epics are fairly broad in scope and are basically product-level things.Each issue is ultimately part of an epic.
- **Milestone**. Each issue is assigned a milestone. This is 0.1, 0.2, ..., or 'Nebulous Future'. The milestone indicates when we think the issue should get addressed.
- **Priority**. Each issue has a priority which is represented by the column in the [Prioritization]() project. Priority can be one of P1, P2, or >P2. The priority indicates how important it is to address the issue within the milestone. P1 says that themilestone cannot be considered achieved if the issue isn't resolved.

