#!/usr/bin/env bash

kubectl apply -f ./CRs/artifacts.yaml

kubectl wait deployment/httpbin -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait deployment/backend-retry-deployment -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait deployment/dynamic-backend -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait deployment/llm-deployment -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait deployment/llm-deployment-subs -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait deployment/llm-deployment-header -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait deployment/interceptor-service-deployment -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait deployment/graphql-faker -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait deployment/grpc-backend -n apk-integration-test --for=condition=available --timeout=600s
kubectl wait --timeout=5m -n apk-integration-test deployment/apk-test-setup-wso2-apk-adapter-deployment --for=condition=Available
kubectl wait --timeout=15m -n apk-integration-test deployment/apk-test-setup-wso2-apk-gateway-runtime-deployment --for=condition=Available
IP=$(kubectl get svc apk-test-setup-wso2-apk-gateway-service -n apk-integration-test --output jsonpath='{.status.loadBalancer.ingress[0].ip}')
CC_IP=$(kubectl get svc apk-test-setup-wso2-apk-common-controller-web-server-service -n apk-integration-test --output jsonpath='{.status.loadBalancer.ingress[0].ip}')
sudo echo "$IP localhost" | sudo tee -a /etc/hosts
sudo echo "$IP idp.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$CC_IP apk-test-setup-wso2-apk-common-controller-service.apk-integration-test.svc" | sudo tee -a /etc/hosts
sudo echo "$IP api.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP org1.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP org2.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP org3.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP org4.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default.sandbox.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default-dev.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default-qa.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP org3-qa.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP org4-qa.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP org4-dev.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "255.255.255.255 broadcasthost" | sudo tee -a /etc/hosts
sudo echo "::1 localhost" | sudo tee -a /etc/hosts
