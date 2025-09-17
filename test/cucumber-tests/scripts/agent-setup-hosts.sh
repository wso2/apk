#!/usr/bin/env bash
kubectl apply -f ./CRs/agent-artifacts.yaml
kubectl wait deployment/apim-wso2am-cp-deployment-1 -n apk --for=condition=available --timeout=600s
kubectl wait --timeout=5m -n apk deployment/apk-wso2-apk-adapter-deployment --for=condition=Available
kubectl wait --timeout=15m -n apk deployment/apk-wso2-apk-gateway-runtime-deployment --for=condition=Available
kubectl wait --timeout=5m -n apk deployment/apim-apk-agent --for=condition=Available
IP=$(kubectl get svc apk-wso2-apk-gateway-service -n apk --output jsonpath='{.status.loadBalancer.ingress[0].ip}')
ING_IP=$(kubectl get ing -n apk apim-wso2am-acp-ingress --output=jsonpath='{.status.loadBalancer.ingress[0].ip}')
CC_IP=$(kubectl get svc apk-wso2-apk-common-controller-web-server-service -n apk --output jsonpath='{.status.loadBalancer.ingress[0].ip}')
sudo echo "$IP localhost" | sudo tee -a /etc/hosts
sudo echo "$IP idp.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$CC_IP apk-wso2-apk-common-controller-service.apk.svc" | sudo tee -a /etc/hosts
sudo echo "$ING_IP am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP api.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default.sandbox.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP sandbox.default.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "255.255.255.255 broadcasthost" | sudo tee -a /etc/hosts
sudo echo "::1 localhost" | sudo tee -a /etc/hosts
