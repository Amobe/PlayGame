package battle

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
)

type Affect struct {
	ActorID    string
	TargetID   string
	Skill      string
	Attributes []character.Attribute
}
