package main

import (
	"encoding/json"
)

// ProcessJSON is the default entrypoint symbol looked up by the runner.
// It echoes input and allows simple header/body mutations for testing.
func ProcessJSON(in []byte) ([]byte, error) {
	var m map[string]interface{}
	_ = json.Unmarshal(in, &m)
	// Return a minimal valid mediation-like output
	out := map[string]interface{}{
		"addHeaders": map[string]string{"X-Plugin": "ok"},
		"modifyBody": false,
	}
	return json.Marshal(out)
}
