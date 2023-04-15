package utils

import (
	"encoding/json"
	"fmt"
)

// MarshalToJSON is a helper function that marshal a object to json.
func MarshalToJSON(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}
	return string(b), nil
}

func UnmarshalFromJSON(jsonStr string, obj interface{}) error {
	if err := json.Unmarshal([]byte(jsonStr), obj); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	return nil
}
