package character_test

import (
	"testing"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/stretchr/testify/assert"
)

func act(actor, target character.Character, skill character.Skill) (affectedActor, affectedTarget character.Character) {
	aa, ta := skill.Affect(actor.GetAttributeMap(), target.GetAttributeMap())
	actor.Affect(aa)
	target.Affect(ta)
	return actor, target
}

func TestAct(t *testing.T) {
	attacker := character.NewCharacter()
	defender := character.NewCharacter()
	skill := character.Skill{
		Affect: func(attacker, defender character.AttributeTypeMap) ([]character.Attribute, []character.Attribute) {
			da := []character.Attribute{
				{
					Type:  character.AttributeTypeHP,
					Value: "-100",
				},
			}
			return nil, da
		},
	}

	wantAttacker := character.NewCharacter()
	wantDefender := character.Character{
		Basement: []character.Attribute{
			{Type: character.AttributeTypeHP, Value: "100"},
			{Type: character.AttributeTypeHP, Value: "-100"},
			{Type: character.AttributeTypeDead},
		},
	}

	affectedAttacker, affectedDefender := act(attacker, defender, skill)

	assert.Equal(t, wantAttacker, affectedAttacker)
	assert.Equal(t, wantDefender, affectedDefender)
}
