minikube image load apk-adapter:1.1.0-SNAPSHOT
kubectl scale deployment --replicas=1 -n apk apk-test-wso2-apk-adapter-deployment
