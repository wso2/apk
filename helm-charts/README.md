# apk-helm

![Version: 1.0.0-beta](https://img.shields.io/badge/Version-1.0.0--beta-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.16.0](https://img.shields.io/badge/AppVersion-1.16.0-informational?style=flat-square)

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
| wso2.apk.auth.enabled | bool | `true` | Enable Service Account Creation |
| wso2.apk.auth.enableServiceAccountCreation | bool | `true` | Enable Service Account Creation |
| wso2.apk.auth.enableClusterRoleCreation | bool | `true` | Enable Cluster Role Creation |
| wso2.apk.auth.serviceAccountName | string | `"wso2apk-platform"` | Service Account name |
| wso2.apk.auth.roleName | string | `"wso2apk-role"` | Cluster Role name |
| wso2.apk.listener.hostname | string | `"api.am.wso2.com"` | System api listener hostname |
| wso2.apk.idp.issuer | string | `"https://idp.am.wso2.com/token"` | IDP issuer value |
| wso2.apk.idp.authorizeEndpoint | string | `"https://idp.am.wso2.com:9095/oauth2/authorize"` | IDP authorization endpoint |
| wso2.apk.idp.tokenEndpoint | string | `"https://idp.am.wso2.com:9095/oauth2/token"` | IDP token endpoint |
| wso2.apk.idp.revokeEndpoint | string | `"https://idp.am.wso2.com:9095/oauth2/revoke"` | IDP revoke endpoint |
| wso2.apk.idp.usernameClaim | string | `"sub"` | Optionally configure username Claim in JWT. |
| wso2.apk.idp.groupClaim | string | `"groups"` | Optionally configure groups Claim in JWT. |
| wso2.apk.idp.scopeClaim | string | `"scope"` | Optionally configure scope Claim in JWT. |
| wso2.apk.idp.organizationClaim | string | `"organization"` | Optionally configure organization Claim in JWT. |
| wso2.apk.idp.organizationResolver | string | `"controlPlane"` | Optionally configure organization Resolution method for APK (controlPlane/none)). |
| wso2.apk.idp.credentials.secretName | string | `""` | IDP credentials secret name to be configured with  |
| wso2.apk.idp.tls.configMapName | string | `""` | IDP public certificate configmap name |
| wso2.apk.idp.tls.secretName | string | `""` | IDP public certificate secret name |
| wso2.apk.idp.tls.fileName | string | `""` | IDP public certificate file name |
| wso2.apk.idp.signing.jwksEndpoint | string | `""` | IDP jwks endpoint (optional) |
| wso2.apk.idp.signing.configMapName | string | `""` | IDP jwt signing certificate configmap name |
| wso2.apk.idp.signing.secretName | string | `""` | IDP jwt signing certificate secret name |
| wso2.apk.idp.signing.fileName | string | `""` | IDP jwt signing certificate file name |
| wso2.apk.dp.enabled | bool | `true` | Enable the deployment of the Data Plane |
| wso2.apk.dp.gateway.listener.hostname | string | `"gw.wso2.com"` | Gateway Listener Hostname |
| wso2.apk.dp.gateway.listener.secretName | string | `""` | Gateway Listener Certificate Secret Name |
| wso2.apk.dp.gateway.autoscaling.enabled | bool | `false` | Enable autoscaling for Gateway |
| wso2.apk.dp.gateway.autoscaling.minReplicas | int | `1` | Minimum number of replicas for Gateway |
| wso2.apk.dp.gateway.autoscaling.maxReplicas | int | `2` | Maximum number of replicas for Gateway |
| wso2.apk.dp.gateway.autoscaling.targetMemory | int | `80` | Target memory utilization percentage for Gateway |
| wso2.apk.dp.gateway.autoscaling.targetCPU | int | `80` | Target CPU utilization percentage for Gateway |
| wso2.apk.dp.partitionServer.enabled | bool | `false` | Enable partition server for Data Plane. |
| wso2.apk.dp.partitionServer.host | string | `""` | Partition Server Service URL |
| wso2.apk.dp.partitionServer.serviceBasePath | string | `"/api/publisher/v1"` | Partition Server Service Base Path. |
| wso2.apk.dp.partitionServer.partitionName | string | `"default"` | Partition Name. |
| wso2.apk.dp.partitionServer.tls.secretName | string | `"managetment-server-cert"` | TLS secret name for Partition Server Public Certificate. |
| wso2.apk.dp.partitionServer.tls.fileName | string | `"certificate.crt"` | TLS certificate file name. |
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
| wso2.apk.dp.configdeployer.deployment.image | string | `"wso2/config-deployer-service:latest"` | Image |
| wso2.apk.dp.configdeployer.deployment.configs.authrorization | bool | `true` | Enable authorization for runtime api. |
| wso2.apk.dp.configdeployer.deployment.configs.baseUrl | string | `"https://api.am.wso2.com:9095/api/runtime"` | Baseurl for runtime api. |
| wso2.apk.dp.configdeployer.deployment.configs.tls.secretName | string | `""` | TLS secret name for runtime public certificate. |
| wso2.apk.dp.configdeployer.deployment.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.configdeployer.deployment.configs.tls.certFilename | string | `""` | TLS certificate file name. |
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
| wso2.apk.dp.adapter.deployment.image | string | `"wso2/adapter:0.0.1-m8"` | Image |
| wso2.apk.dp.adapter.deployment.security.sslHostname | string | `"adapter"` | Enable security for adapter. |
| wso2.apk.dp.adapter.configs.apiNamespaces | string | `nil` | Optionally configure namespaces to watch for apis. |
| wso2.apk.dp.adapter.configs.tls.secretName | string | `""` | TLS secret name for adapter public certificate. |
| wso2.apk.dp.adapter.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.adapter.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.adapter.logging.level | string | `"INFO"` | Optionally configure logging for adapter. LogLevels can be "DEBG", "FATL", "ERRO", "WARN", "INFO", "PANC" |
| wso2.apk.dp.adapter.logging.logFormat | string | `"TEXT"` | Log format can be "JSON", "TEXT" |
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
| wso2.apk.dp.ratelimiter.deployment.image | string | `"wso2/ratelimiter:0.0.1-m8"` | Image |
| wso2.apk.dp.ratelimiter.deployment.security.sslHostname | string | `"ratelimiter"` | hostname for the rate limiter |
| wso2.apk.dp.ratelimiter.deployment.configs.tls.secretName | string | `"ratelimiter-cert"` | TLS secret name for rate limiter public certificate. |
| wso2.apk.dp.ratelimiter.deployment.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.ratelimiter.deployment.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.ratelimiter.deployment.configs.tls.certCAFilename | string | `""` | TLS CA certificate file name. |
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
| wso2.apk.dp.gatewayRuntime.deployment.router.image | string | `"wso2/router:0.0.1-m8"` | Image |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.tls.secretName | string | `"router-cert"` | TLS secret name for router public certificate. |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.deployment.router.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.deployment.router.logging.wireLogs | object | `{"enable":true}` | Optionally configure logging for router. |
| wso2.apk.dp.gatewayRuntime.deployment.router.logging.wireLogs.enable | bool | `true` | Enable wire logs for router. |
| wso2.apk.dp.gatewayRuntime.deployment.router.logging.accessLogs.enable | bool | `true` | Enable access logs for router. |
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
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.image | string | `"wso2/enforcer:latest"` | Image |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.security.sslHostname | string | `"enforcer"` | hostname for the enforcer |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.tls.secretName | string | `""` | TLS secret name for enforcer public certificate. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.tls.certKeyFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.configs.tls.certFilename | string | `""` | TLS certificate file name. |
| wso2.apk.dp.gatewayRuntime.deployment.enforcer.logging.level | string | `"DEBUG"` | Log level can be one of DEBUG, INFO, WARN, ERROR, OFF |
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
| wso2.apk.dp.gatewayRuntime.analytics.authURL | string | `"https://analytics-event-auth.choreo.dev/auth/v1"` | Choreo analytics auth URL. Not required for ELK type. |
| wso2.apk.dp.gatewayRuntime.analytics.authToken | string | `"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"` | On-prem key generated from Choreo console. Not required for ELK type. |
| wso2.apk.dp.gatewayRuntime.analytics.logFileName | string | `"logs/enforcer_analytics.log"` | Optional: File name of the log file. |
| wso2.apk.dp.gatewayRuntime.analytics.logLevel | string | `"INFO"` | Optional: Log level the analytics data. Can be one of DEBUG, INFO, WARN, ERROR, OFF. |
| wso2.apk.migration.enabled | bool | `false` | It is not recommended to run a production deployment with this flag enabled. |
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
| idp.idpds.deployment.image | string | `"wso2/idp-domain-service:latest"` | Image |
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
| idp.idpui.deployment.image | string | `"wso2/idp-ui:0.0.1-m8"` | Image |
| idp.idpui.configs.idpLoginUrl | string | `"https://idp.am.wso2.com:9095/commonauth/login"` | identity server Login URL |
| idp.idpui.configs.idpAuthCallBackUrl | string | `"https://idp.am.wso2.com:9095/oauth2/auth-callback"` | identity server authCallBackUrl |
| gatewaySystem.enabled | bool | `true` | Enable gateway system to install gateway system components |
| gatewaySystem.enableServiceAccountCreation | bool | `true` |  |
| gatewaySystem.enableClusterRoleCreation | bool | `true` |  |
| gatewaySystem.serviceAccountName | string | `"gateway-api-admission"` |  |
| certmanager.enabled | bool | `true` | Enable certificate manager to generate certificates |
| certmanager.enableClusterIssuer | bool | `true` | Enable cluster issuer to generate certificates |
| certmanager.enableRootCa | bool | `true` | Enable root CA to generate certificates |
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
| postgresql.image.debug | bool | `true` | Enable debug mode |
| redis.enabled | bool | `true` | Enable redis |
| redis.architecture | string | `"standalone"` | Redis® architecture. Allowed values: standalone or replication.  |
| redis.fullnameOverride | string | `"redis"` | String to fully override common.names.fullname template |
| redis.primary.service.ports.redis | int | `6379` | Redis service port |
| redis.auth.enabled | bool | `false` | Enable password authentication	 |
| redis.image.debug | bool | `true` | Enable debug mode |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
