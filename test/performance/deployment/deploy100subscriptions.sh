#!/bin/bash

for i in {101..1000}
do
result=$(echo "scale=0; $i/100" | bc)
    echo "apiVersion: cp.wso2.com/v1alpha2
kind: Subscription
metadata:
  name: perf-test-subscription-$i
spec:
  api:
    name: "APIDefinitionEndpointDefault"
    version: \""3.14"\"
  organization: "default"
  subscriptionStatus: "UNBLOCKED"
---">>perf-test-subsciption-$result.yaml
done

for i in {101..1000}
do
result=$(echo "scale=0; $i/100" | bc)
    echo "apiVersion: cp.wso2.com/v1alpha2
kind: Application
metadata:
  name: perf-test-application-$i
spec:
  name: "application-$i"
  organization: "default"
  owner: "admin"
  securitySchemes:
    oauth2:
      environments:
        - appId: "$(uuidgen)"
          envId: "Default"
          keyType: "PRODUCTION"
---">>perf-test-subsciption-$result.yaml
done

for i in {101..1000}
do
result=$(echo "scale=0; $i/100" | bc)
    echo "apiVersion: cp.wso2.com/v1alpha2
kind: ApplicationMapping
metadata:
  name: perf-test-application-mapping-$i
spec:
  applicationRef: "perf-test-application-$i"
  subscriptionRef: "perf-test-subscription-$i"
---">>perf-test-subsciption-$result.yaml
done