#!/usr/bin/env bash

kubectl wait --timeout=5m -n apk-integration-test deployment/apk-test-setup-wso2-apk-adapter-deployment --for=condition=Available
kubectl wait --timeout=15m -n apk-integration-test deployment/apk-test-setup-wso2-apk-gateway-runtime-deployment --for=condition=Available
IP=$(kubectl get svc apk-test-setup-wso2-apk-router-service -n apk-integration-test --output jsonpath='{.status.loadBalancer.ingress[0].ip}')
sudo echo "$IP localhost" | sudo tee -a /etc/hosts
sudo echo "$IP idp.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP api.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "255.255.255.255 broadcasthost" | sudo tee -a /etc/hosts
sudo echo "::1 localhost" | sudo tee -a /etc/hosts
