package dto

// LLMRequest defines the OpenAI request structure
type LLMRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a single message in the OpenAI request
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMResponse defines the OpenAI response structure
type LLMResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

// Usage represents token usage details
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Choice represents a completion choice from OpenAI
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Delta        []any   `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}

// GetKey extracts the last message's content (key) from the request
func (r *LLMRequest) GetKey() (string, bool) {
	if len(r.Messages) == 0 || r.Messages[len(r.Messages)-1].Content == "" {
		return "", false
	}
	return r.Messages[len(r.Messages)-1].Content, true
}

// GetValue extracts the assistant's response content (value) from the response
func (r *LLMResponse) GetValue() (string, bool) {
	if len(r.Choices) == 0 || r.Choices[0].Message.Content == "" {
		return "", false
	}
	return r.Choices[0].Message.Content, true
}
