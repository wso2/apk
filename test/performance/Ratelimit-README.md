# Performance Testing Guide for APK with Rate Limiting

Refer to the README.md file in the test/performance directory for general information on performance testing and environment setup.

This guide outlines the steps to configure the environment and conduct performance tests for the APK product with rate limiting enabled.

## Install APK and Setup APIs

### Sample values.yaml for Rate Limiting
When installing APK, refer to the sample helm/values-cpu2-ratelimit.yaml file.

Ensure the following configurations are updated in the values.yaml file:
- Redis URL and Keys: Set the appropriate Redis URL and authentication keys.

### Sample Secret for Redis Authentication

Apply the helm/secrets.yaml file to the cluster after updating it with the necessary Redis credentials.

### Configure Rate Limiter Service (Headless)

In the helm-charts/templates/data-plane/gateway-components/ratelimiter/ratelimiter-service.yaml file, set the clusterIP to None to make the service headless:

```yaml
spec:
  clusterIP: None
```

### Configure Rate Limiter Deployment

In the helm-charts/templates/data-plane/gateway-components/ratelimiter/ratelimiter-deployment.yaml file, set the following environment variables:

```yaml
- name: LOG_LEVEL
  value: "WARN"
- name: REDIS_TLS_SKIP_HOSTNAME_VERIFICATION
  value: "true"
```

### Setup APIs

Apply the API resources from the ratelimitAPIArtifacts/ directory to the cluster.

## Set Up Azure Redis Cache

Create an Azure Redis Cache (Standard C3) instance in the same region as the Kubernetes cluster.