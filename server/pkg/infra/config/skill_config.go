package config

import (
	_ "embed"
)

type skillConfigData struct {
	Skills []skillDate `json:"skills"`
}

type skillDate struct {
	Type       string         `json:"type"`
	Attributes map[string]any `json:"attributes"`
}

//go:embed resources/skill.json
var skillConfig []byte
