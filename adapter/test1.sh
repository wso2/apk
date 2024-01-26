minikube image load apk-adapter:1.1.0-SNAPSHOT
kubectl scale deployment --replicas=1  -n apk-integration-test apk-test-setup-wso2-apk-adapter-deployment
