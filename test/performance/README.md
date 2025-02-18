# Performance Testing Guide for APK Product

This guide will assist you in performing a performance test for the APK product. To successfully conduct the test, you'll need to set up two JMeter slave servers, one JMeter master server, and a Kubernetes cluster in a private network. The Kubernetes platform can be any of the following: EKS, AKS, GKE, or others.

## Table of Contents

- [Performance Testing Guide for APK Product](#performance-testing-guide-for-apk-product)
  - [Table of Contents](#table-of-contents)
  - [Install APK and Setup APIs](#install-apk-and-setup-apis)
    - [Prerequisites](#prerequisites)
    - [Install APK to the Cluster](#install-apk-to-the-cluster)
    - [Configure APK Components](#configure-apk-components)
    - [Extract Gateway LoadBalancer IP](#extract-gateway-loadbalancer-ip)
  - [Setup JMeter Servers](#setup-jmeter-servers)
    - [Server Hardware Requirements](#server-hardware-requirements)
    - [Prerequisites](#prerequisites-1)
    - [Generate Custom JWT Tokens](#generate-custom-jwt-tokens)
    - [Add Payloads to JMeter](#add-payloads-to-jmeter)
    - [Update Hosts File](#update-hosts-file)
    - [Test Connectivity to APK](#test-connectivity-to-apk)
    - [Start JMeter Slave Servers](#start-jmeter-slave-servers)
    - [Run Performance Test](#run-performance-test)
    - [Generate plots and summary](#generate-plots-and-summary)

---

## Install APK and Setup APIs

### Prerequisites

Ensure you have a Kubernetes cluster (version >= 1.26) set up in the same subnet as the JMeter servers.

### Install APK to the Cluster

1. Create a namespace for the performance test:

```bash
kubectl create ns apk-perf-test
```

2. Navigate to the `helm-charts/templates/data-plane/gateway-components/gateway-runtime/gateway-runtime-deployment.yaml` file and update `JAVA_OPTS` Xms and Xmx values to `-Xms1500m -Xmx1500m`.

3. Similarly, update the `TRAILING_ARGS` log levels to 'warn' in `helm-charts/templates/data-plane/gateway-components/gateway-runtime/gateway-runtime-deployment.yaml`.

4. Go to `helm-charts/templates/data-plane/gateway-components/log-conf.yaml` and add the provided configuration to the `config.toml` value.
5. Go to `helm-charts/templates/data-plane/gateway-components/gateway-runtime/default-jwt-issuer.yaml` and update the content with 
    spec.name: "custom-jwt-issuer"
    spec.issuer: "https://localhost:9443/oauth2/token"
    spec.certificate:
      configMapRef:
        name: custom-jwt-cm
        key: "wso2carboncustom.pem"
  

6. Update resource requests and limits in `helm-charts/values.yaml` as per the table provided in your original document.

7. Install APK to the cluster using Helm:

```bash
cd <APK_HOME>/helm-charts
helm3 install apk-perftest -n apk-perf-test .
```

### Configure APK Components

1. Navigate to `<APK_HOME>test/performance/artifacts/` and create the API:

```bash
kubectl -n apk-perf-test apply -f .
```
Note: If you want to use a custom keystore.jks file other than `APK_HOME/test/performance/keystore.jks` then update the `APK_HOME/test/performance/artifacts/custom-jwt-cm.yaml` with the correct public key.

### Extract Gateway LoadBalancer IP

Execute the following command to determine the private IP of the gateway load balancer:

```bash
IP=$(kubectl get svc apk-perftest-wso2-apk-gateway-service -n apk-perf-test --output jsonpath='{.status.loadBalancer.ingress[0].ip}')
# Note: IP should be within the range of your subnet addresses.
```

## Setup JMeter Servers

### Server Hardware Requirements

Each JMeter server should have at least the following hardware configuration:

- CPU: 8 cores
- Memory: 16 GB

### Prerequisites

Before setting up the JMeter servers, ensure the following prerequisites are met:

- Java 11
- JMeter version 5.5 - Install JMeter on all three servers by downloading and extracting the JMeter package in the {$HOME} directory.
- kubectl 1.27 or higher
- jq
- login to az cluster

Run following commands in JMeter VMs;

```bash
wget "https://github.com/adoptium/temurin11-binaries/releases/download/jdk-11.0.26%2B4/OpenJDK11U-jdk_x64_linux_hotspot_11.0.26_4.tar.gz"
sudo mkdir -p /opt/temurin
sudo tar -xzf OpenJDK11U-jdk_x64_linux_hotspot_11.0.26_4.tar.gz -C /opt/temurin
sudo tee /etc/profile.d/temurin.sh <<EOF
export JAVA_HOME=/opt/temurin/jdk-11.0.26+4
export PATH=\$JAVA_HOME/bin:\$PATH
EOF
sudo chmod +x /etc/profile.d/temurin.sh
source /etc/profile.d/temurin.sh
java -version
sudo snap install kubectl --classic
wget https://archive.apache.org/dist/jmeter/binaries/apache-jmeter-5.5.tgz
tar -xvzf  apache-jmeter-5.5.tgz apache-jmeter-5.5
sudo snap install jq

cd apache-jmeter-5.5/bin/
mkdir payloads
chmod 755 payloads
cd 
git clone https://github.com/wso2/apk
cd apk/test/performance/jwt-tokens/
./generate-jwt-tokens.sh -t 10000 -c 1234 
cd 
cp apk/test/performance/payloads/ apache-jmeter-5.5/bin/ -r
sudo nano /etc/hosts
# add <IP> default.gw.wso2.com 
sudo apt install azure-cli
az login
#connect to cluster
```


### Generate Custom JWT Tokens

Follow these steps in all three jmeter servers

1. Clone the APK repository and navigate to the `jwt-tokens` directory:

```bash
git clone https://github.com/wso2/apk
cd apk/test/performance/jwt-tokens/
```

2. [Optional] If needed, replace `<APK_HOME>test/performance/jwt-tokens/keystore.jks` with your keystore file.

3. Generate JWT tokens using the following command:

```bash
./generate-jwt-tokens.sh -t <number-of-tokens> -c <consumer-key>
```

### Add Payloads to JMeter

Follow these steps in all three jmeter servers

1. Create a 'payloads' folder inside the JMeter bin folder (`<JMETER_HOME>/bin`).

2. Copy and paste the contents of `<APK_HOME>/test/performance/payloads` into the 'payloads' folder.

### Update Hosts File

Follow these steps in all three jmeter servers

Add the extracted IP to `/etc/hosts` and map it to `default.gw.wso2.com` hostname.

### Test Connectivity to APK

Execute the following command to make a test request to APK using a previously generated access token:

```bash
curl -k "https://default.gw.wso2.com:9095/test-definition-default/3.14/employee" --header "Authorization: Bearer $access_token" -d "{"sds":"dsdsd"}" -X POST
```

### Start JMeter Slave Servers

On both JMeter slave servers, execute the following commands from {$HOME} to start the servers.

```sh
heap_size="1g"

echo "Start Server"
echo "Heap: ${heap_size}"

cd ./apache-jmeter-5.5/bin
export HEAP="-Xms${heap_size} -Xmx${heap_size}"
nohup ./jmeter-server >> ~/perf_test.out 2>&1 &
echo $! > nohupid.txt
tail ~/perf_test.out -f
cd -
```

### Run Performance Test

On the JMeter master node, navigate to `<APK_HOME>test/performance` and run the command to initiate the performance test.

```sh
# Replace tokens_path according to your need.
nohup ./run-test-jmeter-client.sh   -n 'cpu-2' -r '<ip-address-of-slave-1>,<ip-address-of-slave-2>' -d '1200' -t "${HOME}/apk/test/performance/jwt-tokens/target/jwt-tokens-10.csv" >> ~/perf_test.out  2>&1 &
echo $! > nohupid.txt
tail -f ~/perf_test.out
```

### Generate plots and summary

to generate plots and summary use scripts inside `<APK_HOME>/test/performance/generate-results` 

These steps guide you through the process of setting up the environment and conducting a performance test for the APK product. If you encounter any issues or need further assistance, refer to the APK documentation or seek support from the APK team.
