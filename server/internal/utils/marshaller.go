package utils

import (
	"encoding/json"
	"fmt"
)

// MarshalToJSON is a helper function that marshal a object to json.
func MarshalToJSON(obj interface{}) ([]byte, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}
	return b, nil
}

func UnmarshalFromJSON(jsonStr []byte, obj interface{}) error {
	if err := json.Unmarshal(jsonStr, obj); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	return nil
}
