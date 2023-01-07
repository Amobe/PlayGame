package utils

import "encoding/json"

func ToString(v any) string {
	str, _ := json.MarshalIndent(v, "", "  ")
	return string(str)
}
