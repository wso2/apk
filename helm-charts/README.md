# Helm chart for APK

## Prerequisites

* Install [Helm](https://helm.sh/docs/intro/install/)
  and [Kubernetes client](https://kubernetes.io/docs/tasks/tools/install-kubectl/) <br><br>

* An already setup [Kubernetes cluster](https://kubernetes.io/docs/setup). If you want to run it on the local you can use Minikube or Kind or a similar software.<br><br>

* Install [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/). If you are using Minikube you can install ingress by running ```minikube addons enable ingress```<br><br>

## Steps to deploy APK DS servers and CloudNativePG

```HELM-HOME``` = apk/helm

1. Execute ``` helm repo add bitnami https://charts.bitnami.com/bitnami ``` and ```helm repo add jetstack https://charts.jetstack.io```
2. Clone the repo and cd into the ```HELM-HOME``` folder.
3. Execute ``` helm dependency build ``` command to download the dependent charts.
4. Now execute ```helm install apk-test . -n apk``` to install the APK components.


