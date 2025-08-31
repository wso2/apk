package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyID      = "mock-key-id"
)

func init() {
	keyID = os.Getenv("KEY_ID")
	if keyID == "" {
		log.Println("KEY_ID environment variable not set, using default mock-key-id")
		keyID = "mock-key-id"
	}
}

type TokenRequest struct {
	GrantType      string  `json:"grant_type"`
	Code           string  `json:"code,omitempty"`
	RedirectURI    string  `json:"redirect_uri,omitempty"`
	ClientID       *string `json:"client_id,omitempty"`
	ClientSecret   string  `json:"client_secret,omitempty"`
	RefreshToken   string  `json:"refresh_token,omitempty"`
	Scope          string  `json:"scope,omitempty"`
	Username       string  `json:"username,omitempty"`
	Password       string  `json:"password,omitempty"`
	ValidityPeriod int     `json:"validity_period,omitempty"` // You may need to parse this manually if passed as string
	AppUUID        *string `json:"app_uuid,omitempty"`
}

func main() {
	const keyFile = "private_key.pem"
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		// Key file does not exist, generate new key
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		// Save it to file
		err = SavePrivateKeyToFile(privateKey, keyFile)
		if err != nil {
			panic(err)
		}
		fmt.Println("Generated and saved new private key.")
	} else {
		// Load existing key from file
		privateKey, err = LoadPrivateKeyFromFile(keyFile)
		if err != nil {
			panic(err)
		}
		fmt.Println("Loaded private key from file.")
	}

	// if err != nil {
	// 	log.Fatalf("Failed to generate key: %v", err)
	// }
	publicKey = &privateKey.PublicKey

	http.HandleFunc("/oauth2/token", tokenHandler)
	http.HandleFunc("/oauth2/jwks", jwksHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Mock IdP running on :9443")
	log.Fatal(http.ListenAndServeTLS(":9443", "idp.crt", "idp.key", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "healthy"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	var req TokenRequest
	req.GrantType = r.FormValue("grant_type")
	req.Code = r.FormValue("code")
	req.RedirectURI = r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	if clientID != "" {
		req.ClientID = &clientID
	}
	appUUID := r.FormValue("app_uuid")
	if appUUID != "" {
		req.AppUUID = &appUUID
	}
	req.ClientSecret = r.FormValue("client_secret")
	req.RefreshToken = r.FormValue("refresh_token")
	req.Scope = r.FormValue("scope")
	req.Username = r.FormValue("username")
	req.Password = r.FormValue("password")

	if vp := r.FormValue("validity_period"); vp != "" {
		if val, err := strconv.Atoi(vp); err == nil {
			req.ValidityPeriod = val
		}
	}

	// 🔐 You can now switch based on req.GrantType
	switch req.GrantType {
	case "client_credentials":
		handleClientCredentialsGrant(w, req)
	case "password":
		handlePasswordGrant(w, req)
	case "authorization_code":
		handleAuthorizationCodeGrant(w, req)
	case "refresh_token":
		handleRefreshTokenGrant(w, req)
	default:
		http.Error(w, "Unsupported grant type", http.StatusBadRequest)
	}
}

func handleClientCredentialsGrant(w http.ResponseWriter, req TokenRequest) {
	aud := "aud"
	if req.ClientID != nil {
		aud = *req.ClientID
	}
	claims := jwt.MapClaims{
		"sub":   "service-account",
		"aud":   aud,
		"scope": req.Scope,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iss":   "https://your-idp-host",
		"jti":   newJTI(),
	}
	if req.ClientID != nil {
		claims["client_id"] = *req.ClientID
	}
	if req.AppUUID != nil {
		claims["application"] = map[string]interface{}{
			"id":   1,
			"uuid": *req.AppUUID,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"access_token": signedToken,
		"token_type":   "Bearer",
		"expires_in":   3600,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handlePasswordGrant(w http.ResponseWriter, req TokenRequest) {
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	claims := jwt.MapClaims{
		"sub":   req.Username,
		"aud":   req.ClientID,
		"scope": req.Scope,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iss":   "https://your-idp-host",
		"jti":   newJTI(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"access_token": signedToken,
		"token_type":   "Bearer",
		"expires_in":   3600,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleAuthorizationCodeGrant(w http.ResponseWriter, req TokenRequest) {
	if req.Code == "" || req.RedirectURI == "" {
		http.Error(w, "Authorization code and redirect URI are required", http.StatusBadRequest)
		return
	}
	// Here you would typically validate the authorization code and redirect URI
	claims := jwt.MapClaims{
		"sub":   "user-id-from-code",
		"aud":   req.ClientID,
		"scope": req.Scope,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iss":   "https://your-idp-host",
		"jti":   newJTI(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"access_token":  signedToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
		"refresh_token": "some-refresh-token",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleRefreshTokenGrant(w http.ResponseWriter, req TokenRequest) {
	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}
	// Here you would typically validate the refresh token
	claims := jwt.MapClaims{
		"sub":   "user-id-from-refresh-token",
		"aud":   req.ClientID,
		"scope": req.Scope,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iss":   "https://your-idp-host",
		"jti":   newJTI(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyID
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"access_token":  signedToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
		"refresh_token": "new-refresh-token",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func jwksHandler(w http.ResponseWriter, r *http.Request) {
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())

	jwk := map[string]interface{}{
		"kty": "RSA",
		"kid": keyID,
		"use": "sig",
		"alg": "RS256",
		"n":   n,
		"e":   e,
	}

	jwks := map[string]interface{}{
		"keys": []interface{}{jwk},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}

type DeterministicReader struct {
	data []byte
	pos  int
}

func NewDeterministicReader(seed []byte) *DeterministicReader {
	return &DeterministicReader{
		data: seed,
		pos:  0,
	}
}

func (r *DeterministicReader) Read(p []byte) (int, error) {
	fmt.Print(p)
	fmt.Print("Length of p: ", len(p))
	n := len(p)
	for i := 0; i < n; i++ {
		p[i] = r.data[r.pos]
		r.pos = (r.pos + 1) % len(r.data) // cycle through seed repeatedly
	}
	return n, nil
}

// newJTI returns a random, URL-safe identifier for the JWT ID claim
func newJTI() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// SavePrivateKeyToFile saves RSA private key to a file in PEM format
func SavePrivateKeyToFile(key *rsa.PrivateKey, filename string) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}
	pemData := pem.EncodeToMemory(pemBlock)
	return ioutil.WriteFile(filename, pemData, 0600) // permission 600 to keep it private
}

// LoadPrivateKeyFromFile loads RSA private key from PEM file
func LoadPrivateKeyFromFile(filename string) (*rsa.PrivateKey, error) {
	pemData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing RSA private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
