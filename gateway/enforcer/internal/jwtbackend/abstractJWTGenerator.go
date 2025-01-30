package jwtbackend

import (
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	dto "github.com/wso2/apk/gateway/enforcer/internal/dto"
)

const (
	none          = "NONE"
	sha256withRSA = "SHA256withRSA"
)

// AbstractAPIMgtGatewayJWTGenerator is an interface for generating JWT tokens.
type AbstractAPIMgtGatewayJWTGenerator struct {
	JwtConfigurationDto *dto.JWTConfiguration
	DialectURI          string
	SignatureAlgorithm  string
	mutex               sync.Mutex
}

// SetJWTConfigurationDto sets the JWT configuration DTO.
func (g *AbstractAPIMgtGatewayJWTGenerator) SetJWTConfigurationDto(jwtConfigurationDto *dto.JWTConfiguration) {
	g.JwtConfigurationDto = jwtConfigurationDto
	g.DialectURI = jwtConfigurationDto.ConsumerDialectURI
	if g.DialectURI == "" {
		g.DialectURI = "http://wso2.org/claims"
	}
	g.SignatureAlgorithm = jwtConfigurationDto.SignatureAlgorithm
	if g.SignatureAlgorithm != none && g.SignatureAlgorithm != sha256withRSA {
		g.SignatureAlgorithm = sha256withRSA
	}
}

// GetJWTConfigurationDto gets the JWT configuration DTO.
func (g *AbstractAPIMgtGatewayJWTGenerator) GetJWTConfigurationDto() *dto.JWTConfiguration {
	return g.JwtConfigurationDto
}

// GenerateToken generates a JWT token.
func (g *AbstractAPIMgtGatewayJWTGenerator) GenerateToken(jwtInfoDto *JWTInfoDto, signatureAlgorithm string, signJWT func(string) ([]byte, error)) (string, error) {
	jwtHeader, err := g.BuildHeader(nil, signatureAlgorithm)
	if err != nil {
		return "", fmt.Errorf("error building JWT header: %w", err)
	}

	jwtBody, err := g.BuildBody(jwtInfoDto)
	if err != nil {
		return "", fmt.Errorf("error building JWT body: %w", err)
	}

	base64UrlEncodedHeader := Encode([]byte(jwtHeader))
	base64UrlEncodedBody := Encode([]byte(jwtBody))

	if signatureAlgorithm == "SHA256withRSA" {
		assertion := base64UrlEncodedHeader + "." + base64UrlEncodedBody

		signedAssertion, err := signJWT(assertion)
		if err != nil {
			return "", fmt.Errorf("error signing JWT: %w", err)
		}

		base64UrlEncodedAssertion := Encode(signedAssertion)
		return base64UrlEncodedHeader + "." + base64UrlEncodedBody + "." + base64UrlEncodedAssertion, nil
	}

	return base64UrlEncodedHeader + "." + base64UrlEncodedBody + ".", nil
}

func (g *AbstractAPIMgtGatewayJWTGenerator) populateStandardClaims(jwtInfoDto *JWTInfoDto) map[string]ClaimValueDTO {
	claims := make(map[string]ClaimValueDTO)
	for key, value := range jwtInfoDto.Claims {
		claims[key] = ClaimValueDTO{
			Value: value.Value,
			Type:  value.Type, // Ensure `Type` is also assigned if it exists in `ClaimValueDTO`.
		}
	}
	return claims
}

func (g *AbstractAPIMgtGatewayJWTGenerator) populateCustomClaims(jwtInfoDto *JWTInfoDto) map[string]ClaimValueDTO {
	claims := make(map[string]ClaimValueDTO)
	for key, value := range jwtInfoDto.Claims {
		claims[key] = ClaimValueDTO{
			Value: value.Value,
			Type:  value.Type, // Ensure `Type` is also assigned if it exists in `ClaimValueDTO`.
		}
	}
	return claims
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

// GenerateHeader generates the JWT header.
func GenerateHeader(jwtConfigurationDto *dto.JWTConfiguration, signatureAlgorithm string) (string, error) {
	if signatureAlgorithm == "NONE" {
		return `{"typ":"JWT","alg":"NONE"}`, nil
	}

	header := fmt.Sprintf(`{"typ":"JWT","alg":"RS256"`)

	if jwtConfigurationDto.UseKid {
		header += fmt.Sprintf(`,"kid":"%v"`, jwtConfigurationDto.UseKid)
	} else {
		thumbprint, err := GenerateThumbprint("SHA-1", jwtConfigurationDto.PublicCert, true)
		if err != nil {
			return "", fmt.Errorf("error in generating public certificate thumbprint: %w", err)
		}
		header += fmt.Sprintf(`,"x5t":"%s"`, thumbprint)
	}

	header += "}"
	return header, nil
}

// AddCertToHeader adds the certificate to the JWT header.
func AddCertToHeader(jwtConfigurationDto *dto.JWTConfiguration, signatureAlgorithm string) (string, error) {
	header, err := GenerateHeader(jwtConfigurationDto, signatureAlgorithm)
	if err != nil {
		return "", fmt.Errorf("error in obtaining keystore: %w", err)
	}
	return header, nil
}

// BuildHeader builds the JWT header.
func (g *AbstractAPIMgtGatewayJWTGenerator) BuildHeader(jwtConfigurationDto *dto.JWTConfiguration, signatureAlgorithm string) (string, error) {
	var jwtHeader string
	if signatureAlgorithm == "NONE" {
		jwtHeader = `{"typ":"JWT","alg":"NONE"}`
	} else if signatureAlgorithm == "SHA256withRSA" {
		header, err := AddCertToHeader(jwtConfigurationDto, signatureAlgorithm)
		if err != nil {
			return "", err
		}
		jwtHeader = header
	}
	return jwtHeader, nil
}

// BuildBody builds the JWT body.
func (g *AbstractAPIMgtGatewayJWTGenerator) BuildBody(jwtInfoDto *JWTInfoDto) (string, error) {
	claims := make(map[string]interface{})

	// Populate standard claims
	for key, value := range g.populateStandardClaims(jwtInfoDto) {
		claims[key] = value
	}

	// Populate custom claims
	for key, claim := range g.populateCustomClaims(jwtInfoDto) {
		var finalValue interface{} = claim.Value
		if strVal, ok := claim.Value.(string); ok {
			switch strings.ToLower(claim.Type) {
			case "bool":
				finalValue = strVal == "true"
			case "int":
				finalValue = parseToInt(strVal)
			case "long":
				finalValue = parseToInt(strVal)
			case "float":
				finalValue = parseToFloat(strVal)
			case "date":
				parsedDate, err := time.Parse("2006-01-02", strVal)
				if err == nil {
					finalValue = parsedDate
				}
			}
		}
		claims[key] = finalValue
	}

	// Convert claims to JSON
	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("error marshaling claims to JSON: %w", err)
	}

	return string(jsonClaims), nil
}

func parseToInt(value string) int64 {
	parsed, _ := time.ParseDuration(value + "s")
	return int64(parsed.Seconds())
}

func parseToFloat(value string) float64 {
	parsed, _ := time.ParseDuration(value + "s")
	return float64(parsed.Seconds())
}

// Encode encodes a byte slice to a base64 URL encoded string.
func Encode(stringToBeEncoded []byte) string {
	return base64.RawURLEncoding.EncodeToString(stringToBeEncoded)
}
