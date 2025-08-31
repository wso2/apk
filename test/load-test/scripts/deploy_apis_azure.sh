#!/usr/bin/env bash

# API deployment details
API_URL="https://api.am.wso2.com:9095/api/deployer/2.0.0/apis/deploy"
HOST_HEADER="api.am.wso2.com"

# Token generation details
TOKEN_URL="https://idp.am.wso2.com:9095/oauth2/token"
TOKEN_HOST="idp.am.wso2.com"
CLIENT_CREDENTIALS="Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg=="  # Replace with your credentials

# File paths
APK_CONFIG_TEMPLATE="/home/azureuser/apk/test/load-test/sample/EmployeeService.apk-conf"
DEFINITION_FILE="/home/azureuser/apk/test/load-test/sample/EmployeeServiceDefinition.json"
TEMP_APK_CONFIG="/tmp/temp_apk_conf.yaml"
CSV_FILE="resource_utilization.csv"

# Initialize CSV file with headers
echo "API_Number,Timestamp,Pod,Container,CPU(m),Memory(Mi)" > "$CSV_FILE"

# Function to get a new JWT token
get_jwt_token() {
    echo "Fetching new JWT token..."
    RESPONSE=$(curl -k --silent --location "$TOKEN_URL" \
        --header "Host: $TOKEN_HOST" \
        --header "Authorization: $CLIENT_CREDENTIALS" \
        --header "Content-Type: application/x-www-form-urlencoded" \
        --data-urlencode "grant_type=client_credentials" \
        --data-urlencode "scope=apk:api_create")

    # Extract access_token using jq
    ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.access_token')

    if [[ -z "$ACCESS_TOKEN" || "$ACCESS_TOKEN" == "null" ]]; then
        echo "Failed to retrieve JWT token. Exiting..."
        exit 1
    fi
}

# Deploy 1000 APIs
for ((i = 1; i <= 1000; i++)); do
    # Get a fresh JWT token for each request
    get_jwt_token

    # Modify APK Configuration
    sed -E "s|name: \"EmployeeServiceAPI\"|name: \"EmployeeServiceAPI$i\"|;
            s|basePath: \"/employee\"|basePath: \"/employee$i\"|" \
        "$APK_CONFIG_TEMPLATE" > "$TEMP_APK_CONFIG"

    echo "Deploying API #$i..."

    # Deploy API using curl
    DEPLOY_RESPONSE=$(curl -k --silent --location "$API_URL" \
        --header "Host: $HOST_HEADER" \
        --header "Authorization: Bearer $ACCESS_TOKEN" \
        --form "apkConfiguration=@$TEMP_APK_CONFIG" \
        --form "definitionFile=@$DEFINITION_FILE")

    # Check if deployment succeeded
    if echo "$DEPLOY_RESPONSE" | jq -e '.code' &>/dev/null; then
        ERROR_MESSAGE=$(echo "$DEPLOY_RESPONSE" | jq -r '.message')
        echo "Error deploying API #$i: $ERROR_MESSAGE"
        continue
    fi

    echo "API #$i deployed successfully."

    # Wait for 15 seconds before checking resource utilization
    sleep 15

    # Get resource utilization of pods and containers in the apk namespace
    TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")

    kubectl get pods -n apk --no-headers -o custom-columns=":metadata.name" | while read -r POD_NAME; do
        kubectl top pod "$POD_NAME" -n apk --containers --no-headers | while read -r line; do
            POD=$(echo "$line" | awk '{print $1}')
            CONTAINER=$(echo "$line" | awk '{print $2}')
            CPU=$(echo "$line" | awk '{print $3}')
            MEMORY=$(echo "$line" | awk '{print $4}')

            if [[ -n "$CPU" && -n "$MEMORY" ]]; then
                echo "$i,$TIMESTAMP,$POD,$CONTAINER,$CPU,$MEMORY" >> "$CSV_FILE"
            else
                echo "Warning: Missing resource values for Pod: $POD, Container: $CONTAINER"
            fi
        done
    done
done

echo "All APIs deployed successfully!"