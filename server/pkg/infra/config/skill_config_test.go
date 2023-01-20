package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestXXX(t *testing.T) {
	s := skillConfigData{}
	_ = json.Unmarshal(skillConfig, &s)
	fmt.Printf("%+v\n", s)
	t.Fail()
}
