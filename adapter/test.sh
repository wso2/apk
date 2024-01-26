kubectl scale deployment --replicas=0  -n apk-integration-test apk-test-setup-wso2-apk-adapter-deployment
./gradlew build
watch -n 1 minikube image rm apk-adapter:1.1.0-SNAPSHOT