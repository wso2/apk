# Sample for Custom Filters

Filters is a set of execution points in the request flow that intercept the request before it goes to the
backend service. They are engaged while the request is processed within the enforcer. The defined set of filters
are applied to all the APIs deployed in the APK. And these filters are engaged inline and if the request
fails at a certain filter, the request will not be forwarded to the next filter and the backend.
The inbuilt set of filters are the authentication filter and the throttling filter.

In this sample, it would read a property (`CustomProperty`) set at the values.yaml under the filter and
set the value as a header for each request.


1. Build the project and create the JAR file.

   Gradle and JDK11 is required to build the project.

    ```
    /.gradlew build
    ```

   Let's assume that the output JAR is named - `sample-custom-filter-1.0-SNAPSHOT.jar`.

2. Add the custom filter to the Enforcer.

    1. Open the values.yaml.
    2. Include the custom filter related configurations.

        - The `className` needs to be the fully qualified `className`.
        - The position denotes the final filter position in the chain, when all the filters are added.
        - By default, the first position is taken by the Authentication Filter and the Throttle Filter is placed as the second filter.
        - As the following example configuration contains `1` as the `position`, it would be executed prior to the Authentication Filter.

    ```yaml
    enforcer:
        configs:
            filters:
            - className: org.example.tests.CustomFilter
                position: 1
                properties:
                - name: CustomProperty
                    value: foo
    ```

3. Create a Docker image.

   Use the APK Enforcer image as the base image and include the JAR into the `/home/wso2/lib/dropins` directory.
   You can build the new image with the following sample Docker file named - `Dockerfile`

    ```
    FROM wso2/enforcer:0.0.1-m7 
    COPY sample-custom-filter-1.0-SNAPSHOT.jar /home/wso2/lib/dropins/sample-custom-filter-1.0-SNAPSHOT.jar
    ```

4. Build the new Enforcer image.

    `docker build -t wso2/enforcer-new:latest . `

5. Start WSO2 APK.

    1. Update the `values.yaml` file of the APK Helm chart.
        - `dp.gatewayRuntime.deployment.enforcer.image` - Use this image (`wso2/enforcer-new:latest`) as the value.
    2. Start APK.

       ```tab="Format"
       helm install <helm-chart-name> . -n <namespace>
       ```

       ```tab="Example"
       helm install apk-test . -n apk
       ```
