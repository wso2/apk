### Kubernetes YAML files - Local

```sh
helm template apk-eg ./helm-charts -n apk-egress-gateway --create-namespace -f helm-charts/values-local.yaml > local-test-v2.yaml
```

### Install Helm Chart

```sh
cd helm-charts
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add jetstack https://charts.jetstack.io
helm dependency build
helm install cert-manager jetstack/cert-manager \
    --namespace apk-egress-gateway \
    --create-namespace \
    --version v1.17.1 \
    --set crds.enabled=true
helm upgrade -i apk-eg . \
    -n apk-egress-gateway \
    --create-namespace \
    -f helm-charts/values-local.yaml
```

### Run Cucumber Tests

```
cd helm-charts
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add jetstack https://charts.jetstack.io
helm dependency build
helm install cert-manager jetstack/cert-manager \
    --namespace apk-integration-test \
    --create-namespace \
    --version v1.17.1 \
    --set crds.enabled=true

helm install apk-test-setup -n apk-integration-test -f values-local.yaml --create-namespace . \
        --set wso2.apk.dp.ratelimiter.requestTimeoutInMillis=800 \
        --set wso2.apk.dp.gatewayRuntime.analytics.enabled=true \
        --set 'wso2.apk.dp.gatewayRuntime.analytics.publishers[0].enabled=true' \
        --set 'wso2.apk.dp.gatewayRuntime.analytics.publishers[0].type=elk'

cd ../test/cucumber-tests
kubectl apply -f ./CRs/artifacts.yaml
./gradlew runTests
```
