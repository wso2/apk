# WSO2 APK - API Platform for Kubernetes


<a href="https://wso2.com/">
   <img src="https://raw.githubusercontent.com/wso2/apk/main/logo/wso2-logo.png"
        alt="WSO2 logo" title="Istio" height="100" width="180" />
</a>

---
A complete API Management platform for modern cloud native architectures. This is an open platform to design, develop and manage APIs in Kubernetes environments. APK is a cloud native API Management platform in K8s utilizing native K8s capabilities.


- For in-depth information about WSO2 API Management Platform, visit [WSO2 API Management](https://wso2.com/api-manager/)
- To ask questions and get assistance from our community, visit [WSO2 Discord](https://discord.com/invite/Xa5VubmThw?utm_source=wso2-dev&utm_medium=link&utm_campaign=wso2-dev_link_from-dev-homepage_221002)
- To learn how to participate in our overall community, visit [our community page](https://wso2.com/community/)

In this README:

- [Introduction](#introduction)
- [Issue management](#issue-management)

You'll find many other useful documents on our [Documentation](https://wso2.com/documentation/).

## Introduction

[APK](https://github.com/wso2/apk) is an open source platform for providing complete API Management capabilities on top of the Kubernetes cluster management platform.

APK is composed of these components:

- **Control Plane** - The APK control plane. It provides API Management capabilities, marketplace capabilities along with domain services and web applications. It consists of the following sub-components:

   - **Back Office** - Responsible for configuring the portal aspects of API including descriptions, documents, images etc. Also manage API visibility and lifecycle

   - **Dev Portal** - Responsible for API consumer interaction. API consumers can discover APIs, read documents, try out and eventually subscribe and consume APIs.

   - **Admin Portal** - Responsible for configuring rate limit policies, key management services and other administrative tasks.

   - **Management Service** - Responsible for communication with dataplanes and pushing updates.

- **Data Plane** - The APK data plane. It provides API runtime capabilities such as gateway, rate limiting services, runtime management. It consists of the following sub-components:

   - **Runtime Manager** - Responsible for configuring the runtime aspects of API including API endpoints, discovering Kubernetes services and converting them into APIs etc.

   - **API Gateway(Router)** - Router will intercept incoming API traffic and apply quality of services such as authentication, authorization and rate limiting.

   - **API Gateway(Enforcer)** - Communicate with the router and enforce quality of services for APIs.

   - **Identity Platform** - Responsible for authentication and authorization happens in dataplane.

- **Operator** - The component provides user-friendly options to operate the APK platform.



## Issue management

We use GitHub to track all of our bugs and feature requests. Each issue we track has a variety of metadata:

- **Epic**. An epic represents a feature area for APK as a whole. Epics are fairly broad in scope and are basically product-level things.
Each issue is ultimately part of an epic.

- **Milestone**. Each issue is assigned a milestone. This is 0.1, 0.2, ..., or 'Nebulous Future'. The milestone indicates when we think the issue should get addressed.

- **Priority**. Each issue has a priority which is represented by the column in the [Prioritization]() project. Priority can be one of P1, P2, or >P2. The priority indicates how important it is to address the issue within the milestone. P1 says that the
milestone cannot be considered achieved if the issue isn't resolved.

