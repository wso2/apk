### Kubernetes YAML files - Local

```sh
helm template apk-eg ./helm-charts -n apk-egress-gateway --create-namespace -f helm-charts/values-local.yaml > local-test-v2.yaml
```
### Kubernetes YAML files - Choreo

```sh
helm template apk-eg ./helm-charts -n apk-egress-gateway --create-namespace -f helm-charts/values-choreo.yaml > choreo-test.yaml
```

### Install Helm Chart

```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add jetstack https://charts.jetstack.io
helm dependency build
helm install cert-manager jetstack/cert-manager \
    --namespace apk-egress-gateway \
    --create-namespace \
    --version v1.17.1 \
    --set crds.enabled=true
helm upgrade -i apk-eg ./helm-charts \
    -n apk-egress-gateway \
    --create-namespace \
    -f helm-charts/values-local.yaml
```

### Run Cucumber Tests

```
helm repo add bitnami https://charts.bitnami.com/bitnami
        helm repo add jetstack https://charts.jetstack.io
        helm dependency build
        helm install cert-manager jetstack/cert-manager \
          --namespace apk-integration-test \
          --create-namespace \
          --version v1.17.1 \
          --set crds.enabled=true

helm install apk-test-setup -n apk-integration-test -f values-local.yaml --create-namespace . \
        --set wso2.apk.dp.gatewayRuntime.analytics.enabled=true \
        --set 'wso2.apk.dp.gatewayRuntime.analytics.publishers[0].enabled=true' \
        --set 'wso2.apk.dp.gatewayRuntime.analytics.publishers[0].type=elk'
```
