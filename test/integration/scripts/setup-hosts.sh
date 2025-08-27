#!/usr/bin/env bash

# kubectl wait --timeout=5m -n gateway-system deployment/gateway-api-admission-server --for=condition=Available
# kubectl wait --timeout=5m -n gateway-system job/gateway-api-admission --for=condition=Complete
# kubectl wait --timeout=5m -n apk-integration-test deployment/apk-test-setup-wso2-apk-adapter-deployment --for=condition=Available
kubectl wait --timeout=5m -n apk-integration-test --for=condition=Available -l app.kubernetes.io/component=proxy,app.kubernetes.io/managed-by=envoy-gateway deployment
IP=$(kubectl get svc -n apk-integration-test -l app.kubernetes.io/component=proxy,app.kubernetes.io/managed-by=envoy-gateway -o jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}')
sudo echo "$IP localhost" | sudo tee -a /etc/hosts
sudo echo "$IP all-http-methods-for-wildcard.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP api-policy-with-jwt-generator.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP backend-base-path.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP path-param-api.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP fetch-api-definition.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP fetch-non-existing-api-definition.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP gateway-integration-test-infra.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP no-base-path.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP api.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP idp.am.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP disable-api-security.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP disable-resource-security.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP prod-api.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP sand-api.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP resource-scopes.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP trailing-slash.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP interceptor-api.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP interceptor-resource.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP cors-policy.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default-api-version-ratelimit.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default-api-version.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP disable-api-level-jwt.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP disable-resource-level-jwt.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP disable-api-level-jwt1.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP disable-api-level-jwt2.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP disable-resource-level-jwt1.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP disable-resource-level-jwt2.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP default-api-version-ratelimit-resource-level.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP ratelimit-priority.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP different-endpoint-with-same-route.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP backend-api-key-security.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP custom-auth-header.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP gql.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP grpc.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP api-level-jwt.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP resource-level-jwt.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP ai-guardrail-wordcount-jsonpath.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP ai-guardrail-wordcount-inverted.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "$IP all-http-methods-for-wildcard-sentencecount.test.gw.wso2.com" | sudo tee -a /etc/hosts
sudo echo "255.255.255.255 broadcasthost" | sudo tee -a /etc/hosts
sudo echo "::1 localhost" | sudo tee -a /etc/hosts
