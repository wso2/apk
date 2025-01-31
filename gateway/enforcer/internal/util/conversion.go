package util

import (
	"fmt"
	"strconv"
)

// ConvertBytesToInt converts a []byte to an int.
// It assumes the []byte contains a valid numeric string (e.g., "123").
func ConvertBytesToInt(data []byte) (int, error) {
	// Convert the []byte to string
	str := string(data)

	// Convert the string to int
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value: %s, error: %w", str, err)
	}

	return num, nil
}

// ConvertStringToInt converts a string to an integer.
// Returns the converted int and an error if the input is invalid.
func ConvertStringToInt(input string) (int, error) {
	num, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid input: %s, error: %w", input, err)
	}
	return num, nil
}