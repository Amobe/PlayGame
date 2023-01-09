package character_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
)

func act(actor, target character.Character, skill character.Skill) (affectedActor, affectedTarget character.Character) {
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
	attackerAttr := character.NewAttributeTypeMap()
	attackerAttr.Insert(character.Attribute{Type: character.AttributeTypeHP, Value: "100"})
	attacker := character.Character{
		Basement: attackerAttr,
	}
	defender := character.NewCharacter("defender")
	skill := &skillTest{}

	wantAttackerAttr := character.NewAttributeTypeMap()
	wantAttackerAttr.Insert(character.Attribute{Type: character.AttributeTypeHP, Value: "100"})
	wantAttacker := character.Character{
		Basement: wantAttackerAttr,
	}
	wantDefenderAttr := character.NewAttributeTypeMap()
	wantDefenderAttr.Insert(
		character.Attribute{Type: character.AttributeTypeHP, Value: "-100"},
		character.Attribute{Type: character.AttributeTypeDead},
	)
	wantDefender := character.Character{
		Basement: wantDefenderAttr,
	}

	affectedAttacker, affectedDefender := act(attacker, defender, skill)

	assert.Equal(t, wantAttacker, affectedAttacker)
	assert.Equal(t, wantDefender, affectedDefender)
}
