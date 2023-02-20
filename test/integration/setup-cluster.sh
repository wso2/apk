#!/usr/bin/env bash

# In git action install kind and kubectl
# go install sigs.k8s.io/kind@v0.17.0
# curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

# Create a kind cluster with k8s version 1.25.3
kind create cluster --image "kindest/node:v1.25.3" --name "apk-dp-tests"

# Install metalLB
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.7/config/manifests/metallb-native.yaml

# Wait for MetalLB to become available.
kubectl rollout status -n metallb-system deployment/controller --timeout 5m
kubectl rollout status -n metallb-system daemonset/speaker --timeout 5m

# Setup address pool used by loadbalancers
subnet=$(docker network inspect kind | jq -r '.[].IPAM.Config[].Subnet | select(contains(":") | not)')
address_first_octets=$(echo "${subnet}" | awk -F. '{printf "%s.%s",$1,$2}')
address_range="${address_first_octets}.255.200-${address_first_octets}.255.250"
kubectl apply -f - <<EOF
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  namespace: metallb-system
  name: kube-services
spec:
  addresses:
  - ${address_range}
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: kube-services
  namespace: metallb-system
spec:
  ipAddressPools:
  - kube-services
EOF

# Create new namespace to install chart
kubectl create ns apk

# Install wso2 apk chart with cp diasabled
helm install apk-test ../../helm-charts -n apk --set wso2.apk.cp.enabled=false

# Wait gateway resources to be available.
kubectl wait --timeout=5m -n gateway-system deployment/gateway-api-admission-server --for=condition=Available
kubectl wait --timeout=5m -n gateway-system job/gateway-api-admission --for=condition=Complete
kubectl wait --timeout=5m -n apk deployment/apk-test-wso2-apk-adapter-deployment --for=condition=Available

# Run tests
# go test -v integration_test.go

# kind delete cluster --name "apk-dp-tests"
