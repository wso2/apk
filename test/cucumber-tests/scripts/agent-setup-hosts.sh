#!/usr/bin/env bash
kubectl apply -f ./CRs/agent-artifacts.yaml
kubectl wait deployment/apim-wso2am-acp-deployment-1 -n apk --for=condition=available --timeout=600s
kubectl wait --timeout=5m -n apk deployment/apk-wso2-kgw-common-controller-deployment --for=condition=Available
kubectl wait --timeout=5m -n apk deployment/envoy-ratelimit --for=condition=Available
kubectl wait --timeout=5m -n apk deployment -l app.kubernetes.io/component=proxy --for=condition=Available
kubectl wait --timeout=5m -n apk deployment -l app.kubernetes.io/app=apim-common-agent --for=condition=Available
IP=$(kubectl get svc -n apk -l app.kubernetes.io/component=proxy -o jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')
ING_IP=$(kubectl get ing -n apk apim-wso2am-acp-ingress --output=jsonpath='{.status.loadBalancer.ingress[0].ip}')
CC_IP=$(kubectl get svc apk-wso2-kgw-common-controller-service -n apk --output jsonpath='{.spec.clusterIP}')
echo "IP: $IP"
echo "ING_IP: $ING_IP"
echo "CC_IP: $CC_IP"
sudo echo "$IP localhost" | sudo tee -a /etc/hosts
sudo echo "$IP idp.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$CC_IP apk-wso2-kgw-common-controller-service.apk.svc" | sudo tee -a /etc/hosts
sudo echo "$ING_IP am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP api.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default.sandbox.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP sandbox.default.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "255.255.255.255 broadcasthost" | sudo tee -a /etc/hosts
sudo echo "::1 localhost" | sudo tee -a /etc/hosts
