package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wso2/apk/gateway/enforcer/internal/dto"
)

// GenerateHeader generates the JWT header.
func GenerateHeader(useKid bool, pubCert *x509.Certificate, signatureAlgorithm string) (string, error) {
	if signatureAlgorithm == "NONE" {
		return `{"typ":"JWT","alg":"NONE"}`, nil
	}

	header := fmt.Sprintf(`{"typ":"JWT","alg":"RS256"`)

	if useKid {
		header += fmt.Sprintf(`,"kid":"%v"`, useKid)
	} else {
		thumbprint, err := GenerateThumbprint("SHA-1", pubCert, true)
		if err != nil {
			return "", fmt.Errorf("error in generating public certificate thumbprint: %w", err)
		}
		header += fmt.Sprintf(`,"x5t":"%s"`, thumbprint)
	}

	header += "}"
	return header, nil
}

// GenerateThumbprint generates a thumbprint of the public certificate.
func GenerateThumbprint(hashType string, publicCert *x509.Certificate, usePadding bool) (string, error) {
	if publicCert == nil {
		return "", errors.New("public certificate is nil")
	}

	hash := sha1.New()
	if hashType != "SHA-1" {
		return "", errors.New("unsupported hash type")
	}

	hash.Write(publicCert.Raw)
	digestInBytes := hash.Sum(nil)
	publicCertThumbprint := Hexify(digestInBytes)

	var base64UrlEncodedThumbPrint string
	if usePadding {
		base64UrlEncodedThumbPrint = base64.URLEncoding.EncodeToString([]byte(publicCertThumbprint))
	} else {
		base64UrlEncodedThumbPrint = base64.RawURLEncoding.EncodeToString([]byte(publicCertThumbprint))
	}

	return base64UrlEncodedThumbPrint, nil
}

// Hexify converts a byte slice to a hex string.
func Hexify(bytes []byte) string {
	hexDigits := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	var builder strings.Builder
	builder.Grow(len(bytes) * 2)

	for _, b := range bytes {
		builder.WriteRune(hexDigits[(b&0xf0)>>4])
		builder.WriteRune(hexDigits[b&0x0f])
	}

	return builder.String()
}

// BuildBody builds the JWT body.
func BuildBody(claims map[string]*dto.ClaimValue) (string, error) {
	customClaims := make(map[string]interface{})

	// Populate custom claims
	for key, claim := range claims {
		var finalValue interface{} = claim.Value
		switch strings.ToLower(claim.Type) {
		case "string":
			finalValue = claim.Value
		case "bool":
			finalValue = claim.Value == "true"
		case "int":
			finalValue = parseToInt(claim.Value)
		case "long":
			finalValue = parseToInt(claim.Value)
		case "float":
			finalValue = parseToFloat(claim.Value)
		case "date":
			parsedDate, err := time.Parse("2006-01-02", claim.Value)
			if err == nil {
				finalValue = parsedDate
			}
		}
		customClaims[key] = finalValue
	}

	// Convert claims to JSON
	jsonClaims, err := json.Marshal(customClaims)
	if err != nil {
		return "", fmt.Errorf("error marshaling claims to JSON: %w", err)
	}

	return string(jsonClaims), nil
}

func parseToInt(value string) int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func parseToFloat(value string) float64 {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return parsed
}

// BuildHeader builds the JWT header.
func BuildHeader(publicCert *x509.Certificate, useKid bool, signatureAlgorithm string) (string, error) {
	if signatureAlgorithm == "NONE" {
		return `{"typ":"JWT","alg":"NONE"}`, nil
	}

	header := `{"typ":"JWT","alg":"RS256"`

	if useKid {
		header += `,"kid":"true"`
	} else {
		thumbprint, err := GenerateThumbprint("SHA-1", publicCert, true)
		if err != nil {
			return "", fmt.Errorf("error generating public certificate thumbprint: %w", err)
		}
		header += fmt.Sprintf(`,"x5t":"%s"`, thumbprint)
	}

	header += "}"
	return header, nil
}

// GenerateJWTToken generates a JWT token.
func GenerateJWTToken(signatureAlgo string, useKid bool, publicCert *x509.Certificate, claims map[string]*dto.ClaimValue, privateKey *rsa.PrivateKey) string {
	header, err := BuildHeader(publicCert, true, signatureAlgo)
	if err != nil {
		return ""
	}
	body, err := BuildBody(claims)
	if err != nil {
		return ""
	}
	// Base64 encode header and body

	if signatureAlgo != "SHA256withRSA" {
		base64Header := base64.RawURLEncoding.EncodeToString([]byte(header))
		base64Body := base64.RawURLEncoding.EncodeToString([]byte(body))

		// Concatenate header and body with a period
		unsignedToken := fmt.Sprintf("%s.%s", base64Header, base64Body)
		return fmt.Sprintf("%s.", unsignedToken)
	}
	// Sign the token
	jwtToken, err := signJWT(header, body, privateKey)
	if err != nil {
		return ""
	}

	return jwtToken
	// Generate JWT token
}

// LoadPrivateKey Read Private Key from a PEM file
func LoadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Try parsing as PKCS#1
	if privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return privateKey, nil
	}

	// Try parsing as PKCS#8
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// Ensure it's an RSA key
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an RSA key")
	}

	return rsaKey, nil
}

// Sign JWT using SHA256withRSA
func signJWT(header, payload string, privateKey *rsa.PrivateKey) (string, error) {
	// Create signing string
	signingInput := base64URLEncode([]byte(header)) + "." + base64URLEncode([]byte(payload))

	// Hash the signing input
	hashed := sha256.Sum256([]byte(signingInput))

	// Sign with RSA private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}

	// Encode the signature in Base64URL
	signedJWT := signingInput + "." + base64URLEncode(signature)
	return signedJWT, nil
}

// base64URLEncode Base64 URL Encoding (JWT-safe)
func base64URLEncode(input []byte) string {
	return base64.RawURLEncoding.EncodeToString(input)
}

// Base64Encode encodes a byte slice to a base64 string.
func Base64Encode(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}
