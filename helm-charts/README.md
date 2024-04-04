# apk-helm

![Version: 1.1.0-alpha](https://img.shields.io/badge/Version-1.1.0--alpha-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.16.0](https://img.shields.io/badge/AppVersion-1.16.0-informational?style=flat-square)

A Helm chart for APK components

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://charts.bitnami.com/bitnami | postgresql | 11.9.6 |
| https://charts.bitnami.com/bitnami | redis | 17.8.0 |
| https://charts.jetstack.io | cert-manager | v1.10.1 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| wso2.subscription.imagePullSecrets | string | `""` | Optionally specify image pull secrets. |
| wso2.apk.webhooks.validatingwebhookconfigurations | bool | `true` |  |
| wso2.apk.webhooks.mutatingwebhookconfigurations | bool | `true` |  |
| wso2.apk.auth.enabled | bool | `true` | Enable Service Account Creation |
| wso2.apk.auth.enableServiceAccountCreation | bool | `true` | Enable Service Account Creation |
| wso2.apk.auth.enableClusterRoleCreation | bool | `true` | Enable Cluster Role Creation |
| wso2.apk.auth.serviceAccountName | string | `"wso2apk-platform"` | Service Account name |
| wso2.apk.auth.roleName | string | `"wso2apk-role"` | Cluster Role name |
| wso2.apk.listener.hostname | string | `"api.am.wso2.com"` | System api listener hostname |
| wso2.apk.listener.port | int | `9095` | Gatewaylistener port |
| wso2.apk.listener.secretName | string | `"system-api-listener-cert"` | System api listener certificates. If you are using a custom certificate. |
| wso2.apk.idp.issuer | string | `"https://idp.am.wso2.com/token"` | IDP issuer value |
| wso2.apk.idp.usernameClaim | string | `"sub"` |  |
| wso2.apk.idp.scopeClaim | string | `"scope"` | Optionally configure scope Claim in JWT. |
| wso2.apk.idp.organizationClaim | string | `"organization"` | Optionally configure organization Claim in JWT. |
| wso2.apk.idp.organizationResolver | string | `"none"` | Optionally configure organization Resolution method for APK (none)). |
| wso2.apk.idp.tls.configMapName | string | `""` | IDP public certificate configmap name |
| wso2.apk.idp.tls.secretName | string | `""` | IDP public certificate secret name |
| wso2.apk.idp.tls.fileName | string | `""` | IDP public certificate file name |
| wso2.apk.idp.signing.jwksEndpoint | string | `""` | IDP jwks endpoint (optional) |
| wso2.apk.idp.signing.configMapName | string | `""` | IDP jwt signing certificate configmap name |
| wso2.apk.idp.signing.secretName | string | `""` | IDP jwt signing certificate secret name |
| wso2.apk.idp.signing.fileName | string | `""` | IDP jwt signing certificate file name |
| wso2.apk.cp.enableApiPropagation | bool | `false` | Enable controlplane connection |
| wso2.apk.cp.enabledSubscription | bool | `false` | Enable controlplane connection |
| wso2.apk.cp.host | string | `"apim-apk-agent-service.apk.svc.cluster.local"` | Hostname of the APK agent service |
| wso2.apk.cp.skipSSLVerification | bool | `false` | Skip SSL verification |
| wso2.apk.cp.persistence | object | `{"type":"K8s"}` | Provide persistence mode DB/K8s |
| wso2.apk.dp.enabled | bool | `true` | Enable the deployment of the Data Plane |
| wso2.apk.dp.environment.name | string | `"Development"` | Environment Name of the Data Plane |
| wso2.apk.dp.gatewayClass | object | `{"name":"wso2-apk-default"}` | GatewayClass custom resource name |
| wso2.apk.dp.gateway.name | string | `"wso2-apk-default"` | Gateway custom resource name |
| wso2.apk.dp.gateway.listener.hostname | string | `"gw.wso2.com"` | Gateway Listener Hostname |
| wso2.apk.dp.gateway.listener.secretName | string | `""` | Gateway Listener Certificate Secret Name |
| wso2.apk.dp.gateway.listener.dns | list | `["*.gw.wso2.com","*.sandbox.gw.wso2.com","prod.gw.wso2.com"]` | DNS entries for gateway listener certificate |
| wso2.apk.dp.gateway.httpListener.enabled | bool | `false` | HTTP listener enabled or not |
| wso2.apk.dp.gateway.httpListener.hostname | string | `"api.am.wso2.com"` | HTTP listener hostname |
| wso2.apk.dp.gateway.httpListener.port | int | `9080` | HTTP listener port |
| wso2.apk.dp.gateway.autoscaling.enabled | bool | `false` | Enable autoscaling for Gateway |
| wso2.apk.dp.gateway.autoscaling.minReplicas | int | `1` | Minimum number of replicas for Gateway |
| wso2.apk.dp.gateway.autoscaling.maxReplicas | int | `2` | Maximum number of replicas for Gateway |
| wso2.apk.dp.gateway.autoscaling.targetMemory | int | `80` | Target memory utilization percentage for Gateway |
| wso2.apk.dp.gateway.autoscaling.targetCPU | int | `80` | Target CPU utilization percentage for Gateway |
| wso2.apk.dp.redis.type | string | `"single"` | Redis type |
| wso2.apk.dp.redis.url | string | `"redis-master:6379"` | Redis URL |
| wso2.apk.dp.redis.tls | bool | `false` | TLS enabled  |
| wso2.apk.dp.redis.auth.certificatesSecret | string | `nil` | Redis ceritificate secret |
| wso2.apk.dp.redis.auth.secretKey | string | `nil` | Redis secret key |
| wso2.apk.dp.redis.poolSize | string | `nil` | Redis pool size |
| wso2.apk.dp.partitionServer.enabled | bool | `false` | Enable partition server for Data Plane. |
| wso2.apk.dp.partitionServer.host | string | `""` | Partition Server Service URL |
| wso2.apk.dp.partitionServer.serviceBasePath | string | `"/api/publisher/v1"` | Partition Server Service Base Path. |
| wso2.apk.dp.partitionServer.partitionName | string | `"default"` | Partition Name. |
| wso2.apk.dp.partitionServer.tls.secretName | string | `"managetment-server-cert"` | TLS secret name for Partition Server Public Certificate. |
| wso2.apk.dp.partitionServer.tls.fileName | string | `"certificate.crt"` | TLS certificate file name. |
| wso2.apk.dp.configdeployer.enabled | bool | `true` |  |
| wso2.apk.dp.configdeployer.deployment.resources.requests.memory | string | `"128Mi"` | CPU request for the container |
| wso2.apk.dp.configdeployer.deployment.resources.requests.cpu | string | `"100m"` | Memory request for the container |
| wso2.apk.dp.configdeployer.deployment.resources.limits.memory | string | `"1028Mi"` | CPU limit for the container |
| wso2.apk.dp.configdeployer.deployment.resources.limits.cpu | string | `"1000m"` | Memory limit for the container |
| wso2.apk.dp.configdeployer.deployment.readinessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.configdeployer.deployment.readinessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.configdeployer.deployment.readinessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.configdeployer.deployment.livenessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.configdeployer.deployment.livenessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.configdeployer.deployment.livenessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.configdeployer.deployment.strategy | string | `"RollingUpdate"` | Deployment strategy |
| wso2.apk.dp.configdeployer.deployment.replicas | int | `1` | Number of replicas |
| wso2.apk.dp.configdeployer.deployment.imagePullPolicy | string | `"Always"` | Image pull policy |
| wso2.apk.dp.configdeployer.deployment.image | string | `"wso2/apk-config-deployer-service:1.1.0-alpha"` | Image |
| wso2.apk.dp.configdeployer.deployment.configs.authorization | bool | `true` | Enable authorization for runtime api. |
| wso2.apk.dp.configdeployer.deployment.configs.baseUrl | string | `"https://api.am.wso2.com:9095/api/runtime"` | Baseurl for runtime api. |
| wso2.apk.dp.configdeployer.deployment.configs.tls.secretName | string | `""` | TLS secret name for runtime public certificate. |
| wso2.apk.dp.configdeployer.deployment.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.configdeployer.deployment.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.configdeployer.vhosts | list | `[{"hosts":["gw.wso2.com"],"name":"default","type":"production"},{"hosts":["sandbox.gw.wso2.com"],"name":"default","type":"sandbox"}]` | List of vhost |
| wso2.apk.dp.adapter.deployment.resources.requests.memory | string | `"128Mi"` | CPU request for the container |
| wso2.apk.dp.adapter.deployment.resources.requests.cpu | string | `"100m"` | Memory request for the container |
| wso2.apk.dp.adapter.deployment.resources.limits.memory | string | `"1028Mi"` | CPU limit for the container |
| wso2.apk.dp.adapter.deployment.resources.limits.cpu | string | `"1000m"` | Memory limit for the container |
| wso2.apk.dp.adapter.deployment.readinessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.adapter.deployment.readinessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.adapter.deployment.readinessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.adapter.deployment.livenessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.adapter.deployment.livenessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.adapter.deployment.livenessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.adapter.deployment.strategy | string | `"RollingUpdate"` | Deployment strategy |
| wso2.apk.dp.adapter.deployment.replicas | int | `1` | Number of replicas |
| wso2.apk.dp.adapter.deployment.imagePullPolicy | string | `"Always"` | Image pull policy |
| wso2.apk.dp.adapter.deployment.image | string | `"wso2/apk-adapter:1.1.0-alpha"` | Image |
| wso2.apk.dp.adapter.deployment.security.sslHostname | string | `"adapter"` | Enable security for adapter. |
| wso2.apk.dp.adapter.configs.apiNamespaces | string | `nil` | Optionally configure namespaces to watch for apis. |
| wso2.apk.dp.adapter.configs.tls.secretName | string | `""` | TLS secret name for adapter public certificate. |
| wso2.apk.dp.adapter.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.adapter.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.adapter.logging.level | string | `"INFO"` | Optionally configure logging for adapter. LogLevels can be "DEBG", "FATL", "ERRO", "WARN", "INFO", "PANC" |
| wso2.apk.dp.adapter.logging.logFile | string | `"logs/adapter.log"` | Log file name |
| wso2.apk.dp.adapter.logging.logFormat | string | `"TEXT"` | Log format can be "JSON", "TEXT" |
| wso2.apk.dp.commonController.deployment.resources.requests.memory | string | `"128Mi"` | Memory request for the container |
| wso2.apk.dp.commonController.deployment.resources.requests.cpu | string | `"100m"` | CPU request for the container |
| wso2.apk.dp.commonController.deployment.resources.limits.memory | string | `"1028Mi"` | Memory limit for the container |
| wso2.apk.dp.commonController.deployment.resources.limits.cpu | string | `"1000m"` | CPU limit for the container |
| wso2.apk.dp.commonController.deployment.readinessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before readinessProbe probes are initiated. |
| wso2.apk.dp.commonController.deployment.readinessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.commonController.deployment.readinessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.commonController.deployment.livenessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.commonController.deployment.livenessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.commonController.deployment.livenessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.commonController.deployment.strategy | string | `"RollingUpdate"` | Deployment strategy |
| wso2.apk.dp.commonController.deployment.replicas | int | `1` | Number of replicas |
| wso2.apk.dp.commonController.deployment.imagePullPolicy | string | `"Always"` | Image pull policy |
| wso2.apk.dp.commonController.deployment.image | string | `"wso2/apk-common-controller:1.1.0-alpha"` | Image |
| wso2.apk.dp.commonController.deployment.security.sslHostname | string | `"commoncontroller"` | hostname for the common controller |
| wso2.apk.dp.commonController.deployment.configs.apiNamespaces | list | `["apk-v12"]` | Optionally configure namespaces to watch for apis,ratelimitpolicies,etc. |
| wso2.apk.dp.commonController.deployment.redis.host | string | `"redis-master"` | Redis host |
| wso2.apk.dp.commonController.deployment.redis.port | string | `"6379"` | Redis port |
| wso2.apk.dp.commonController.deployment.redis.username | string | `"default"` | Redis user name |
| wso2.apk.dp.commonController.deployment.redis.password | string | `""` | Redis password |
| wso2.apk.dp.commonController.deployment.redis.tlsEnabled | bool | `false` | Redis TLS enabled or not |
| wso2.apk.dp.commonController.deployment.redis.userCertPath | string | `"/home/wso2/security/keystore/commoncontroller.crt"` | Redis user cert to use for redis connections |
| wso2.apk.dp.commonController.deployment.redis.userKeyPath | string | `"/home/wso2/security/keystore/commoncontroller.key"` | Redis user key to use for redis connections |
| wso2.apk.dp.commonController.deployment.redis.cACertPath | string | `"/home/wso2/security/keystore/commoncontroller.crt"` | Redis CA cert to use for redis connections |
| wso2.apk.dp.commonController.deployment.redis.channelName | string | `"wso2-apk-revoked-tokens-channel"` | Token revocation subscription channel name |
| wso2.apk.dp.commonController.deployment.database.host | string | `"wso2apk-db-service.apk"` |  |
| wso2.apk.dp.commonController.deployment.database.port | int | `5432` |  |
| wso2.apk.dp.commonController.deployment.database.username | string | `"wso2carbon"` |  |
| wso2.apk.dp.commonController.deployment.database.password | string | `"wso2carbon"` |  |
| wso2.apk.dp.commonController.deployment.database.poolOptions.poolMaxConns | int | `4` |  |
| wso2.apk.dp.commonController.deployment.database.poolOptions.poolMinConns | int | `0` |  |
| wso2.apk.dp.commonController.deployment.database.poolOptions.poolMaxConnLifetime | string | `"1h"` |  |
| wso2.apk.dp.commonController.deployment.database.poolOptions.poolMaxConnIdleTime | string | `"1h"` |  |
| wso2.apk.dp.commonController.deployment.database.poolOptions.poolHealthCheckPeriod | string | `"1m"` |  |
| wso2.apk.dp.commonController.deployment.database.poolOptions.poolMaxConnLifetimeJitter | string | `"1s"` |  |
| wso2.apk.dp.ratelimiter.enabled | bool | `true` | Enable the deployment of the Rate Limiter |
| wso2.apk.dp.ratelimiter.deployment.resources.requests.memory | string | `"128Mi"` | CPU request for the container |
| wso2.apk.dp.ratelimiter.deployment.resources.requests.cpu | string | `"100m"` | Memory request for the container |
| wso2.apk.dp.ratelimiter.deployment.resources.limits.memory | string | `"1028Mi"` | CPU limit for the container |
| wso2.apk.dp.ratelimiter.deployment.resources.limits.cpu | string | `"1000m"` | Memory limit for the container |
| wso2.apk.dp.ratelimiter.deployment.readinessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.ratelimiter.deployment.readinessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.ratelimiter.deployment.readinessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.ratelimiter.deployment.livenessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.ratelimiter.deployment.livenessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.ratelimiter.deployment.livenessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.ratelimiter.deployment.strategy | string | `"RollingUpdate"` | Deployment strategy |
| wso2.apk.dp.ratelimiter.deployment.replicas | int | `1` | Number of replicas |
| wso2.apk.dp.ratelimiter.deployment.imagePullPolicy | string | `"Always"` | Image pull policy |
| wso2.apk.dp.ratelimiter.deployment.image | string | `"wso2/apk-ratelimiter:1.1.0-alpha"` | Image |
| wso2.apk.dp.ratelimiter.deployment.security.sslHostname | string | `"ratelimiter"` | hostname for the rate limiter |
| wso2.apk.dp.ratelimiter.deployment.configs.tls.secretName | string | `"ratelimiter-cert"` | TLS secret name for rate limiter public certificate. |
| wso2.apk.dp.ratelimiter.deployment.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.ratelimiter.deployment.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.ratelimiter.deployment.configs.tls.certCAFilename | string | `""` | TLS CA certificate file name. |
| wso2.apk.dp.gatewayRuntime.service.annotations | string | `nil` | Gateway service related annotations. |
| wso2.apk.dp.gatewayRuntime.deployment.replicas | int | `1` | Number of replicas |
| wso2.apk.dp.gatewayRuntime.deployment.router.resources.requests.memory | string | `"128Mi"` | CPU request for the container |
| wso2.apk.dp.gatewayRuntime.deployment.router.resources.requests.cpu | string | `"100m"` | Memory request for the container |
| wso2.apk.dp.gatewayRuntime.deployment.router.resources.limits.memory | string | `"1028Mi"` | CPU limit for the container |
| wso2.apk.dp.gatewayRuntime.deployment.router.resources.limits.cpu | string | `"1000m"` | Memory limit for the container |
| wso2.apk.dp.gatewayRuntime.deployment.router.readinessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.gatewayRuntime.deployment.router.readinessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.gatewayRuntime.deployment.router.readinessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.gatewayRuntime.deployment.router.livenessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.gatewayRuntime.deployment.router.livenessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.gatewayRuntime.deployment.router.livenessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.gatewayRuntime.deployment.router.strategy | string | `"RollingUpdate"` | Deployment strategy |
| wso2.apk.dp.gatewayRuntime.deployment.router.imagePullPolicy | string | `"Always"` | Image pull policy |
| wso2.apk.dp.gatewayRuntime.deployment.router.image | string | `"wso2/apk-router:1.1.0-alpha"` | Image |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.enforcerResponseTimeoutInSeconds | int | `20` | The timeout for response coming from enforcer to route per API request |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.useRemoteAddress | bool | `false` | If configured true, router appends the immediate downstream ip address to the x-forward-for header |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.systemHost | string | `"localhost"` | System hostname for system API resources (eg: /testkey and /health) |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.enableIntelligentRouting | bool | `false` | Enable Semantic Versioning based Intelligent Routing for Gateway |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.tls.secretName | string | `"router-cert"` | TLS secret name for router public certificate. |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.upstream.tls.verifyHostName | bool | `true` | Enable/Disable Verifying host name |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.upstream.tls.disableSslVerification | bool | `false` | Disable SSL verification |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.upstream.dns.dnsRefreshRate | int | `5000` | DNS refresh rate in miliseconds |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.upstream.dns.respectDNSTtl | bool | `false` | set cluster’s DNS refresh rate to resource record’s TTL which comes from DNS resolution |
| wso2.apk.dp.gatewayRuntime.deployment.router.logging.wireLogs | object | `{"enable":true}` | Optionally configure logging for router. |
| wso2.apk.dp.gatewayRuntime.deployment.router.logging.wireLogs.enable | bool | `true` | Enable wire logs for router. |
| wso2.apk.dp.gatewayRuntime.deployment.router.logging.accessLogs.enable | bool | `true` | Enable access logs for router. |
| wso2.apk.dp.gatewayRuntime.deployment.router.logging.accessLogs.logfile | string | `"/tmp/envoy.access.log"` | Log file name |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.resources.requests.memory | string | `"128Mi"` | CPU request for the container |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.resources.requests.cpu | string | `"100m"` | Memory request for the container |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.resources.limits.memory | string | `"1028Mi"` | CPU limit for the container |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.resources.limits.cpu | string | `"1000m"` | Memory limit for the container |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.readinessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.readinessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.readinessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.livenessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.livenessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.livenessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.strategy | string | `"RollingUpdate"` | Deployment strategy |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.imagePullPolicy | string | `"Always"` | Image pull policy |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.image | string | `"wso2/apk-enforcer:1.1.0-alpha"` | Image |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.security.sslHostname | string | `"enforcer"` | hostname for the enforcer |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.tls.secretName | string | `""` | TLS secret name for enforcer public certificate. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.authService | object | `{"keepAliveTime":600,"maxHeaderLimit":8192,"maxMessageSize":1000000000,"threadPool":{"coreSize":400,"keepAliveTime":600,"maxSize":1000,"queueSize":2000}}` | The configurations of gRPC netty based server in Enforcer that handles the incoming requests from ext_authz |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.mandateSubscriptionValidation | bool | `false` | Specifies whether subscription validation is mandated for all APIs. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.logging.level | string | `"DEBUG"` | Log level can be one of DEBUG, INFO, WARN, ERROR, OFF |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.logging.logFile | string | `"logs/enforcer.log"` | Log file name |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.host | string | `"redis-master"` | Redis host |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.port | string | `"6379"` | Redis port |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.username | string | `"default"` | Redis user name |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.password | string | `""` | Redis password |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.tlsEnabled | bool | `false` | Redis TLS enabled or not |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.userCertPath | string | `"/home/wso2/security/keystore/commoncontroller.crt"` |  |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.userKeyPath | string | `"/home/wso2/security/keystore/commoncontroller.key"` | Redis user key to use for redis connections |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.cACertPath | string | `"/home/wso2/security/keystore/commoncontroller.crt"` | Redis CA cert to use for redis connections |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.redis.channelName | string | `"wso2-apk-revoked-tokens-channel"` | Token revocation subscription channel name |
| wso2.apk.dp.gatewayRuntime.tracing.enabled | bool | `true` | Enable/Disable tracing in gateway runtime. |
| wso2.apk.dp.gatewayRuntime.tracing.type | string | `"zipkin"` | Type of tracer exporter (e.g: azure, zipkin). Use zipkin type for Jaeger as well. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.host | string | `"jaeger"` | Jaeger/Zipkin host. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.port | string | `"9411"` | Jaeger/Zipkin port. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.endpoint | string | `"/api/v2/spans"` | Jaeger/Zipkin collector endpoint path. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.instrumentationName | string | `"APK"` | Library Name to be tagged in traces (`otel.library.name`). |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.maximumTracesPerSecond | string | `"2"` | Maximum number of sampled traces per second string. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.maxPathLength | string | `"256"` | Maximum length of the request path to extract and include in the HttpUrl tag. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.connectionString | string | `"https://otlp.nr-data.net"` | New Relic OTLP gRPC collector endpoint. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.authHeaderName | string | `"api-key"` | Auth header name. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.authHeaderValue | string | `"<INGEST_LICENSE_KEY>"` | Auth header value. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.connectionTimeout | string | `"20"` | Connection timeout for the otlp service. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.tls.enabled | bool | `true` | Enable/Disable TLS for the otlp service. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.tls.secretName | string | `"ratelimiter-cert"` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.tracing.configProperties.tls.certCAFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.analytics.enabled | bool | `true` | Enable/Disable analytics in gateway runtime. |
| wso2.apk.dp.gatewayRuntime.analytics.type | string | `"Choreo"` | Type of analytics data publisher. Can be "Choreo" or "ELK". |
| wso2.apk.dp.gatewayRuntime.analytics.secretName | string | `"choreo-analytics-secret"` | Choreo analytics secret. |
| wso2.apk.dp.gatewayRuntime.analytics.properties | object | `{"property_name":"property_value"}` | Property values for the analytics. |
| wso2.apk.dp.gatewayRuntime.analytics.publishers | list | `[{"configProperties":{"auth.api.token":"$env{analytics_authToken}","auth.api.url":"$env{analytics_authURL}"},"enabled":true,"type":"default"},{"enabled":true,"type":"elk"}]` | Analytics Publishers |
| wso2.apk.dp.gatewayRuntime.analytics.logFileName | string | `"logs/enforcer_analytics.log"` | Optional: File name of the log file. |
| wso2.apk.dp.gatewayRuntime.analytics.logLevel | string | `"INFO"` | Optional: Log level the analytics data. Can be one of DEBUG, INFO, WARN, ERROR, OFF. |
| wso2.apk.dp.gatewayRuntime.analytics.receiver | object | `{"keepAliveTime":600,"maxHeaderLimit":8192,"maxMessageSize":1000000000,"threadPool":{"coreSize":10,"keepAliveTime":600,"maxSize":100,"queueSize":1000}}` | gRPC access log service within Enforcer |
| wso2.apk.dp.gatewayRuntime.analytics.receiver.maxMessageSize | int | `1000000000` | Maximum message size in bytes |
| wso2.apk.dp.gatewayRuntime.analytics.receiver.maxHeaderLimit | int | `8192` | Maximum header size in bytes |
| wso2.apk.dp.gatewayRuntime.analytics.receiver.keepAliveTime | int | `600` | Keep alive time of gRPC access log connection |
| wso2.apk.dp.gatewayRuntime.analytics.receiver.threadPool | object | `{"coreSize":10,"keepAliveTime":600,"maxSize":100,"queueSize":1000}` | Thread pool configuration for gRPC access log server |
| wso2.apk.dp.gatewayRuntime.analytics.receiver.threadPool.coreSize | int | `10` | Minimum number of workers to keep alive |
| wso2.apk.dp.gatewayRuntime.analytics.receiver.threadPool.maxSize | int | `100` | Maximum pool size |
| wso2.apk.dp.gatewayRuntime.analytics.receiver.threadPool.keepAliveTime | int | `600` | Timeout in seconds for idle threads waiting for work |
| wso2.apk.dp.gatewayRuntime.analytics.receiver.threadPool.queueSize | int | `1000` | Queue size of the worker threads |
| wso2.apk.metrics.enabled | bool | `false` | Enable Prometheus metrics |
| idp.enabled | bool | `true` | Enable Non production identity server |
| idp.listener.hostname | string | `"idp.am.wso2.com"` | identity server hostname |
| idp.listener.secretName | string | `"idp-tls"` | identity server certificate |
| idp.database.driver | string | `"org.postgresql.Driver"` | identity server database driver |
| idp.database.url | string | `"jdbc:postgresql://wso2apk-db-service:5432/WSO2AM_DB"` | identity server database url |
| idp.database.host | string | `"wso2apk-db-service"` | identity server database host |
| idp.database.port | int | `5432` | identity server database port |
| idp.database.databaseName | string | `"WSO2AM_DB"` | identity server database name |
| idp.database.username | string | `"wso2carbon"` | identity server database username |
| idp.database.secretName | string | `"apk-db-secret"` | identity server database password secret name |
| idp.database.secretKey | string | `"DB_PASSWORD"` | identity server database password secret key |
| idp.database.validationQuery | string | `"SELECT 1"` | identity server database validation query |
| idp.database.validationTimeout | int | `250` | identity server database validation timeout |
| idp.idpds.config.issuer | string | `"https://idp.am.wso2.com/token"` | identity server issuer url |
| idp.idpds.config.keyId | string | `"gateway_certificate_alias"` | identity server keyId |
| idp.idpds.config.hostname | string | `"idp.am.wso2.com"` | identity server hostname. |
| idp.idpds.config.loginPageURl | string | `"https://idp.am.wso2.com:9095/authenticationEndpoint/login"` | identity server login page url |
| idp.idpds.config.loginErrorPageUrl | string | `"https://idp.am.wso2.com:9095/authenticationEndpoint/error"` | identity server login error page url |
| idp.idpds.config.loginCallBackURl | string | `"https://idp.am.wso2.com:9095/authenticationEndpoint/login-callback"` | identity server login callback page url |
| idp.idpds.deployment.resources.requests.memory | string | `"128Mi"` | CPU request for the container |
| idp.idpds.deployment.resources.requests.cpu | string | `"100m"` | Memory request for the container |
| idp.idpds.deployment.resources.limits.memory | string | `"1028Mi"` | CPU limit for the container |
| idp.idpds.deployment.resources.limits.cpu | string | `"1000m"` | Memory limit for the container |
| idp.idpds.deployment.readinessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| idp.idpds.deployment.readinessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| idp.idpds.deployment.readinessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| idp.idpds.deployment.livenessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| idp.idpds.deployment.livenessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| idp.idpds.deployment.livenessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| idp.idpds.deployment.strategy | string | `"RollingUpdate"` | Deployment strategy |
| idp.idpds.deployment.replicas | int | `1` | Number of replicas |
| idp.idpds.deployment.imagePullPolicy | string | `"Always"` | Image pull policy |
| idp.idpds.deployment.image | string | `"wso2/apk-idp-domain-service:1.1.0-alpha"` | Image |
| idp.idpui.deployment.resources.requests.memory | string | `"128Mi"` | CPU request for the container |
| idp.idpui.deployment.resources.requests.cpu | string | `"100m"` | Memory request for the container |
| idp.idpui.deployment.resources.limits.memory | string | `"1028Mi"` | CPU limit for the container |
| idp.idpui.deployment.resources.limits.cpu | string | `"1000m"` | Memory limit for the container |
| idp.idpui.deployment.readinessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| idp.idpui.deployment.readinessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| idp.idpui.deployment.readinessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| idp.idpui.deployment.livenessProbe.initialDelaySeconds | int | `20` | Number of seconds after the container has started before liveness probes are initiated. |
| idp.idpui.deployment.livenessProbe.periodSeconds | int | `20` | How often (in seconds) to perform the probe. |
| idp.idpui.deployment.livenessProbe.failureThreshold | int | `5` | Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| idp.idpui.deployment.strategy | string | `"RollingUpdate"` | Deployment strategy |
| idp.idpui.deployment.replicas | int | `1` | Number of replicas |
| idp.idpui.deployment.imagePullPolicy | string | `"Always"` | Image pull policy |
| idp.idpui.deployment.image | string | `"wso2/apk-idp-ui:1.1.0-alpha"` | Image |
| idp.idpui.configs.idpLoginUrl | string | `"https://idp.am.wso2.com:9095/commonauth/login"` | identity server Login URL |
| idp.idpui.configs.idpAuthCallBackUrl | string | `"https://idp.am.wso2.com:9095/oauth2/auth-callback"` | identity server authCallBackUrl |
| gatewaySystem.enabled | bool | `true` | Enable gateway system to install gateway system components |
| gatewaySystem.enableServiceAccountCreation | bool | `true` |  |
| gatewaySystem.enableClusterRoleCreation | bool | `true` |  |
| gatewaySystem.serviceAccountName | string | `"gateway-api-admission"` |  |
| gatewaySystem.applyGatewayWehbhookJobs | bool | `true` |  |
| certmanager.enabled | bool | `true` | Enable certificate manager to generate certificates |
| certmanager.enableClusterIssuer | bool | `true` | Enable cluster issuer to generate certificates |
| certmanager.enableRootCa | bool | `true` | Enable root CA to generate certificates |
| certmanager.rootCaSecretName | string | `"apk-root-certificate"` | Enable CA certificate secret name. |
| certmanager.listeners.issuerName | string | `"selfsigned-issuer"` | Issuer name |
| certmanager.listeners.issuerKind | string | `"ClusterIssuer"` | Issuer kind |
| certmanager.servers.issuerName | string | `"selfsigned-issuer"` | Issuer name |
| certmanager.servers.issuerKind | string | `"ClusterIssuer"` | Issuer kind |
| postgresql.enabled | bool | `true` | Enable postgresql database |
| postgresql.fullnameOverride | string | `"wso2apk-db-service"` | String to fully override common.names.fullname template |
| postgresql.auth.database | string | `"WSO2AM_DB"` | Name for a custom database to create |
| postgresql.auth.postgresPassword | string | `"wso2carbon"` | Password for the "postgres" admin user. Ignored if auth.existingSecret is provided |
| postgresql.auth.username | string | `"wso2carbon"` | Name for a custom user to create |
| postgresql.auth.password | string | `"wso2carbon"` | Password for the custom user to create. Ignored if auth.existingSecret is provided |
| postgresql.primary.extendedConfiguration | string | `"max_connections = 400\n"` | Extended PostgreSQL Primary configuration (appended to main or default configuration) |
| postgresql.primary.initdb.scriptsConfigMap | string | `"postgres-initdb-scripts-configmap"` | ConfigMap with PostgreSQL initialization scripts |
| postgresql.primary.initdb.user | string | `"wso2carbon"` | Specify the PostgreSQL username to execute the initdb scripts |
| postgresql.primary.initdb.password | string | `"wso2carbon"` | Specify the PostgreSQL password to execute the initdb scripts |
| postgresql.primary.service.ports.postgresql | int | `5432` | PostgreSQL service port |
| postgresql.primary.podSecurityContext.enabled | bool | `true` | Enable pod security context |
| postgresql.primary.podSecurityContext.fsGroup | string | `nil` | Pod security context fsGroup |
| postgresql.primary.podSecurityContext.runAsNonRoot | bool | `true` | Pod security context runAsNonRoot |
| postgresql.primary.podSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` | Pod security context seccomp profile type |
| postgresql.primary.containerSecurityContext.enabled | bool | `true` | Enable container security context |
| postgresql.primary.containerSecurityContext.allowPrivilegeEscalation | bool | `false` | Container security context allow privilege escalation |
| postgresql.primary.containerSecurityContext.capabilities.drop | list | `["ALL"]` | Container security context capabilities drop |
| postgresql.primary.containerSecurityContext.runAsUser | string | `nil` | Container security context runAsUser |
| redis.enabled | bool | `true` | Enable redis |
| redis.architecture | string | `"standalone"` | Redis® architecture. Allowed values: standalone or replication.  |
| redis.fullnameOverride | string | `"redis"` | String to fully override common.names.fullname template |
| redis.primary.service.ports.redis | int | `6379` | Redis service port |
| redis.master.podSecurityContext.enabled | bool | `true` | Enable pod security context |
| redis.master.podSecurityContext.fsGroup | string | `nil` | Pod security context fsGroup |
| redis.master.podSecurityContext.runAsNonRoot | bool | `true` | Pod security context runAsNonRoot |
| redis.master.podSecurityContext.seccompProfile.type | string | `"RuntimeDefault"` | Pod security context seccomp profile type |
| redis.master.containerSecurityContext.enabled | bool | `true` | Enable container security context |
| redis.master.containerSecurityContext.allowPrivilegeEscalation | bool | `false` | Container security context allow privilege escalation |
| redis.master.containerSecurityContext.capabilities.drop | list | `["ALL"]` | Container security context capabilities drop |
| redis.master.containerSecurityContext.runAsUser | string | `nil` | Container security context runAsUser |
| redis.auth.enabled | bool | `false` | Enable password authentication	 |
| skipCrds | bool | `false` | Skip generate of CRD templates |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.13.1](https://github.com/norwoodj/helm-docs/releases/v1.13.1)
