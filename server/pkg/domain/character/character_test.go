package character_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/skill"
)

func act(actor, target character.Character, skill skill.Skill) (affectedActor, affectedTarget character.Character) {
	aa, ta := skill.Use(actor.AttributeMap(), target.AttributeMap())
	actor.Affect(aa)
	target.Affect(ta)
	return actor, target
}

type skillTest struct{}

func (s *skillTest) Use(am, dm character.AttributeTypeMap) (aa, ta []character.Attribute) {
	return nil, []character.Attribute{{character.AttributeTypeHP, "-100"}}
}

func (s *skillTest) Name() string {
	return "skillTest"
}

func TestAct(t *testing.T) {
	attacker := character.NewCharacter("attacker",
		character.Attribute{Type: character.AttributeTypeHP, Value: "100"},
	)
	defender := character.NewCharacter("defender",
		character.Attribute{Type: character.AttributeTypeHP, Value: "100"},
	)
	skill := &skillTest{}

	wantAttacker := character.NewCharacter("attacker",
		character.Attribute{Type: character.AttributeTypeHP, Value: "100"},
	)
	wantDefender := character.NewCharacter("defender",
		character.Attribute{Type: character.AttributeTypeHP, Value: "0"},
		character.Attribute{Type: character.AttributeTypeDead},
	)

	affectedAttacker, affectedDefender := act(attacker, defender, skill)

	assert.Equal(t, wantAttacker, affectedAttacker)
	assert.Equal(t, wantDefender, affectedDefender)
}
