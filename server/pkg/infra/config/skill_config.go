package config

import (
	_ "embed"
)

type skillConfigData struct {
	Skills []struct {
		Type       string         `json:"type"`
		Attributes map[string]any `json:"attributes"`
	} `json:"skills"`
}

//go:embed resources/skill.json
var skillConfig []byte
