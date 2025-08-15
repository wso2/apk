## Interceptor service sample

The sample interceptor service code is generated using [swagger-codegen](https://github.com/swagger-api/swagger-codegen) using the following command:
```
swagger-codegen generate -i https://raw.githubusercontent.com/wso2/apk/main/developer/resources/interceptor-service-open-api-v1.yaml -l nodejs-server -o ./interceptors
```
You can use the same command with the supported languages you prefer and generate stub code to implement your interceptor service.

You can go through following to build, deploy, and test the sample interceptor service. Follow first two points to better understand the sample.

- [Build the interceptor sample image and test request and response path interceptor invokations](#build-the-interceptor-sample-image-and-test-request-and-response-path-interceptor-invokations)
- [Build the legacy backend sample image and test](#build-the-interceptor-sample-image-and-test-request-and-response-path-interceptor-invokations)
- [Try out the API with interceptor APIPolicy in APK](#build-the-interceptor-sample-image-and-test-request-and-response-path-interceptor-invokations)

### Build the interceptor sample image and test request and response path interceptor invocations

1. Build the sample interceptor service.
   ```sh
   sh ./build-interceptor-service.sh
   ```
   Here, the requestBody is Base64 encoded.

2. Start the interceptor service.
   ```sh
   docker run --name interceptor-nodejs -p 9081:9081 wso2am/apk-sample-xml-interceptor-nodejs:v1.0.0
   ```

3. Test the interceptor service.
   Invoke request path 
   ```sh
   curl https://localhost:9081/api/v1/handle-request \
      -H "content-type: application/json" \
      -H "accept: application/json" \
      -d '{"requestBody": "eyJuYW1lIjoiVGhlIFByaXNvbmVyIn0K"}' 
      -kv
   ```
   Sample request path response
   ```json
   {
     "headersToAdd": {
       "x-user": "admin"
     },
     "headersToReplace": {
       "content-type": "application/xml"
     },
     "body": "PG5hbWU+VGhlIFByaXNvbmVyPC9uYW1lPg=="
   }
   ```

   Invoke response path 
   ```sh
   curl https://localhost:9081/api/v1/handle-response \
      -H "content-type: application/json" \
      -H "accept: application/json" \
      -d '{"responseCode": 200}' \
      -kv
   ```
   Sample response path response
   ```json
   {
     "responseCode": 201
   }
   ```

4. Remove container
   ```sh
   docker rm -f interceptor-nodejs
   ```

### Build the legacy backend sample image and test

1. Build the backend service.
    ```sh
    sh ./build-legacy-backend-service.sh
    ```

2. Test the backend service.
   ```sh
   docker run --name lagacy-backend -p 9082:9082 wso2am/apk-sample-legacy-backend-nodejs:v1.0.0
   ```
   
   In another shell
    ```sh
    curl -X POST http://localhost:9082/books \
      -d '<name>The Prisoner</name>' \
      -H 'x-user: admin' -v
    ```
   
   Remove the container
   ```shell
   docker rm -f lagacy-backend
   ```

### Try out the API with interceptor APIPolicy in APK 

1. Build the interceptor service and backend service container images:
    ```
    sh ./build-interceptor-service.sh
    sh ./build-legacy-backend-service.sh
    ```

3. Create namepace called `interceptor` to deploy all the interceptor resources:
    ```
    kubectl create ns interceptor
    ```

4. Change directory to `/k8s-resources` and run following to deploy all the API related resources:
    ```sh
    kubectl apply -f .
    ```
    When you run above, you apply following resources: 
    | File  | Description |
    | ------------- | ------------- |
    | api.yaml  | `API` resource for the API |
    | api-httproute.yaml  | `HTTPRoute` resource for the API |
    | api-backend.yaml | Contains `Backend`, `Service` and `Deployment` resources represeting the API backend  |
    | api-policy.yaml  | `APIPolicy` resource to wire the interceptor service with the API  |
    | api-req-interceptor-svc.yaml  | `InterceptorService` resource to define the interceptor service in the request path  |
    | api-res-interceptor-svc.yaml  | `InterceptorService` resource to define the interceptor service in the response path  |
    | interceptor-backend.yaml  | Contains `Backend`, `Service` and `Deployment` resources represeting the interceptor service backend  |
    | interceptor-certificate.yaml | `Certificate` resource to generate cert/key pair and CA cert for the interceptor service |


5. Add loopback ip to `interceptor.gw.wso2.com` mapping to your `/etc/hosts` file:
    ```
    127.0.0.1       interceptor.gw.wso2.com
    ```

6. Port forward the gateway runtime service:
    ```sh
    kubectl port-forward svc/apk-test-router  9095:9095
    ```

7. Get an access token to invoke the API:
    ```sh
    curl -k --location 'https://idp.am.wso2.com:9095/oauth2/token' \
    --header 'Host: idp.am.wso2.com' \
    --header 'Authorization: Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg==' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'grant_type=client_credentials'
    ```

8. Invoke the API:
    ```sh
    curl "https://interceptor.gw.wso2.com:9095/interceptor/1.0.0/books" \
    -H "Authorization:Bearer $TOKEN" \
    -H "content-type: application/json" \
    -d '{"name":"The Prisoner"}' -v -k 
    ```
