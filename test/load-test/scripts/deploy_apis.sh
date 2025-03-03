#!/bin/bash

# API deployment details
API_URL="https://api.am.wso2.com:9095/api/deployer/1.3.0/apis/deploy"
HOST_HEADER="api.am.wso2.com"

# Token generation details
TOKEN_URL="https://idp.am.wso2.com:9095/oauth2/token"
TOKEN_HOST="idp.am.wso2.com"
CLIENT_CREDENTIALS="Basic NDVmMWM1YzgtYTkyZS0xMWVkLWFmYTEtMDI0MmFjMTIwMDAyOjRmYmQ2MmVjLWE5MmUtMTFlZC1hZmExLTAyNDJhYzEyMDAwMg=="  # Replace with your credentials

# File paths
APK_CONFIG_TEMPLATE="/Users/admin/Documents/1000APIsTest/sample/EmployeeService.apk-conf"
DEFINITION_FILE="/Users/admin/Documents/1000APIsTest/sample/EmployeeServiceDefinition.json"
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

    if [ "$ACCESS_TOKEN" == "null" ] || [ -z "$ACCESS_TOKEN" ]; then
        echo "Failed to retrieve JWT token. Exiting..."
        exit 1
    fi
}

# Deploy 1000 APIs
for i in {1..1000}
do
    # Get a fresh JWT token for each request
    get_jwt_token

    # Modify APK Configuration
    sed "s/name: \"EmployeeServiceAPI[0-9]*\"/name: \"EmployeeServiceAPI$i\"/;
         s|basePath: \"/employee[0-9]*\"|basePath: \"/employee$i\"|" "$APK_CONFIG_TEMPLATE" > "$TEMP_APK_CONFIG"

    echo "Deploying API #$i..."

    # Deploy API using curl
    curl -k --location "$API_URL" \
        --header "Host: $HOST_HEADER" \
        --header "Authorization: Bearer $ACCESS_TOKEN" \
        --form "apkConfiguration=@$TEMP_APK_CONFIG" \
        --form "definitionFile=@$DEFINITION_FILE"

    # Wait for 15 seconds before checking resource utilization
    sleep 15

    # Get resource utilization of pods and containers in the apk namespace
    TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")
    
    kubectl get pods -n apk --no-headers -o custom-columns=":metadata.name" | while read -r POD_NAME; do
        # Get container-wise resource utilization
        kubectl top pod "$POD_NAME" -n apk --containers --no-headers | while read -r line; do
            POD_NAME=$(echo "$line" | awk '{print $1}')
            CONTAINER_NAME=$(echo "$line" | awk '{print $2}')
            CPU_USAGE=$(echo "$line" | awk '{print $3}')
            MEMORY_USAGE=$(echo "$line" | awk '{print $4}')

            # Ensure both CPU and Memory values exist
            if [ -n "$CPU_USAGE" ] && [ -n "$MEMORY_USAGE" ]; then
                echo "$i,$TIMESTAMP,$POD_NAME,$CONTAINER_NAME,$CPU_USAGE,$MEMORY_USAGE" >> "$CSV_FILE"
            else
                echo "Warning: Missing resource values for Pod: $POD_NAME, Container: $CONTAINER_NAME"
            fi
        done
    done

    echo "API #$i deployed and resource utilization recorded."
done

echo "All APIs deployed successfully!"
