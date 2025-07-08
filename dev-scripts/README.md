### Kubernetes YAML files

```sh
helm template apk-eg ./helm-charts -n apk-egress-gateway --create-namespace -f helm-charts/values-dev.yaml > local-test.yaml
```

### Install Helm Chart

```sh
helm upgrade -i apk-eg ./helm-charts -n apk-egress-gateway --create-namespace -f helm-charts/values-dev.yaml
```
