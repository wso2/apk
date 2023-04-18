# apk-custom-ratelimit-interceptor-service

This sample demonstrates the use of custom policy ratelimiting using a simple interceptor service written in NodeJS.

### Steps to try out the sample scenarios.

> Note: Ensure that you have a functional APK deployment with the ratelimiter before proceeding with the scenarios mentioned below.

Execute `kubectl apply -f deployment.yaml -n apk` to deploy the sample interceptor service and the sample backend.

#### Applying a custom ratelimit policy for the Gateway

Gateway-level interceptors can be utilized to implement a custom rate limiting policy at the gateway level. The gateway will invoke the external interceptor service for all requests that pass through it, and the external interceptor service will dynamically set the effective rate limiting policies.

```sh
kubectl apply -f gateway-interceptor -n apk
```

#### Applying a custom ratelimit policy for a specific API

API level interceptors can be utilized to implement a custom rate limiting policy for a specific API. The gateway will invoke the external interceptor service for all requests belonging to this particular API, and the external interceptor service will dynamically set the effective rate limiting policies.

```sh
kubectl apply -f api-interceptor -n apk
```

#### Applying a custom ratelimit policy for a specific resource

Resource level interceptors can be utilized to implement a custom rate limiting policy for a specific resource. The gateway will invoke the external interceptor service for all requests belonging to this particular resource, and the external interceptor service will dynamically set the effective rate limiting policies.

```sh
kubectl apply -f resource-interceptor -n apk
```

### Steps to modify the interceptor service

1. Navigate to `index.js` and add your implementation there. The external interceptor service should implement a POST endpoint with the path `/api/v1/handle-request` to handle requests made by the gateway. Additionally, ensure that the `rateLimitKeys` property is returned with the appropriate rate limit keys.

2. Build the Docker image using the command `docker build -t interceptor-service .`"

3. Deploy the changes by running the command `kubectl apply -f deployment.yaml -n apk`"
