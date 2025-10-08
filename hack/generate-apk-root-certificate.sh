#!/bin/bash

# Generate Self-Signed Certificate for APK

# Set default values
DEFAULT_DAYS=10950  # 30 years (365*30)
DEFAULT_KEY="certs/apk.key"
DEFAULT_CERT="certs/apk.crt"
DEFAULT_SUBJ="/CN=apk"

# Check if OpenSSL is installed
if ! command -v openssl &> /dev/null; then
    echo "Error: OpenSSL is not installed. Please install OpenSSL first."
    exit 1
fi

echo "Generating APK root certificate..."
echo "This will create:"
echo "  Private key: $DEFAULT_KEY"
echo "  Certificate: $DEFAULT_CERT"
echo "  Valid for: $DEFAULT_DAYS days (30 years)"
echo "  Subject: $DEFAULT_SUBJ"

rm -rf certs
mkdir certs

# Generate the private key and certificate
openssl req -x509 -nodes -newkey rsa:2048 \
    -keyout "$DEFAULT_KEY" \
    -out "$DEFAULT_CERT" \
    -days "$DEFAULT_DAYS" \
    -subj "$DEFAULT_SUBJ"

# Check if generation was successful
if [ $? -eq 0 ]; then
    echo ""
    echo "Successfully generated:"
    echo "  Private key: $DEFAULT_KEY"
    echo "  Certificate: $DEFAULT_CERT"
    echo ""
    echo "You can verify the certificate with:"
    echo "  openssl x509 -in $DEFAULT_CERT -text -noout"
    
    # Create Kubernetes secret
    echo ""
    echo "Creating Kubernetes secret 'apk-root-certificate'..."
    
    # Check if kubectl is installed
    if ! command -v kubectl &> /dev/null; then
        echo "Warning: kubectl is not installed. Please install kubectl to create the Kubernetes secret."
        echo "You can manually create the secret using the following command:"
        echo "kubectl create secret tls apk-root-certificate --cert=$DEFAULT_CERT --key=$DEFAULT_KEY -n apk-egress-gateway"
    else
        # Create the secret using kubectl
        kubectl create secret tls apk-root-certificate \
            --cert="$DEFAULT_CERT" \
            --key="$DEFAULT_KEY" \
            --namespace=apk-egress-gateway \
            --dry-run=client -o yaml > certs/apk-root-certificate-secret.yaml
        
        if [ $? -eq 0 ]; then
            echo "Kubernetes secret manifest created: certs/apk-root-certificate-secret.yaml"
            echo ""
            echo "To apply the secret to your cluster, run:"
            echo "  kubectl apply -f certs/apk-root-certificate-secret.yaml"
            echo ""
            echo "Or apply directly:"
            echo "  kubectl create secret tls apk-root-certificate --cert=$DEFAULT_CERT --key=$DEFAULT_KEY -n apk-egress-gateway"
        else
            echo "Error creating Kubernetes secret manifest."
        fi
    fi
else
    echo ""
    echo "Error generating certificate."
    exit 1
fi
