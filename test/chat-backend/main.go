package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// ChatRequest represents the input payload for the /chat/completions endpoint
type ChatRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// ChatResponse represents the response for the /chat/completions endpoint
type ChatResponse struct {
	Model  string `json:"model"`
	Output string `json:"output"`
}

// SupportedModelsResponse represents the response for the /get endpoint
type SupportedModelsResponse struct {
	Models []string `json:"models"`
}

// Supported models
var supportedModels = []string{"gpt-4o", "gpt-3.5", "gpt-4.5"}

func main() {
	rand.Seed(time.Now().UnixNano())

	// POST /chat/completions
	http.HandleFunc("/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req ChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check if the model is supported
		isValidModel := false
		for _, model := range supportedModels {
			if model == req.Model {
				isValidModel = true
				break
			}
		}
		if !isValidModel {
			http.Error(w, fmt.Sprintf("Model %s is not supported", req.Model), http.StatusBadRequest)
			return
		}

		// Simulate token count and response
		remainingTokens := rand.Intn(100) // Randomly assign remaining tokens
		remainingRequests := rand.Intn(100)
		response := ChatResponse{
			Model:  req.Model,
			Output: fmt.Sprintf("Processed input for model %s: %s", req.Model, req.Input),
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-ratelimit-remaining-tokens", fmt.Sprintf("%d", remainingTokens))
		w.Header().Set("x-ratelimit-remaining-requests", fmt.Sprintf("%d", remainingRequests))
		json.NewEncoder(w).Encode(response)
	})

	// GET /get
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		response := SupportedModelsResponse{Models: supportedModels}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Start the server
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
