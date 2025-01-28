package util

import (
	"encoding/json"
)
// ToJSONString converts any object to a JSON string
func ToJSONString(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
