package util

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
)

// ToJSONString converts any object to a JSON string
func ToJSONString(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// IsValidJSON checks if a string is a valid JSON
func IsValidJSON(s string) bool {
	return json.Valid([]byte(s))
}

// DecompressIfGzip Decompress GZIP if the response is compressed
func DecompressIfGzip(data []byte) ([]byte, error) {
	if len(data) < 2 {
		return data, nil // Not GZIP
	}

	// GZIP magic number check
	if data[0] == 0x1f && data[1] == 0x8b {
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to create GZIP reader: %w", err)
		}
		defer reader.Close()

		uncompressedData, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress GZIP: %w", err)
		}
		return uncompressedData, nil
	}

	return data, nil // Not compressed, return as is
}
