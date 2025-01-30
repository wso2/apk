package jwtbackend

// ClaimValueDTO represents a claim with a value and a type.
type ClaimValueDTO struct {
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

// NewClaimValueDTO creates a new instance of ClaimValueDTO.
func NewClaimValueDTO(value interface{}, valueType string) *ClaimValueDTO {
	return &ClaimValueDTO{
		Value: value,
		Type:  valueType,
	}
}

// GetValue returns the value of the claim.
func (c *ClaimValueDTO) GetValue() interface{} {
	return c.Value
}

// GetType returns the type of the claim.
func (c *ClaimValueDTO) GetType() string {
	return c.Type
}
