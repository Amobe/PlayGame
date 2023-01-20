package character_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/skill"
	"github.com/Amobe/PlayGame/server/pkg/domain/valueobject"
)

func act(actor, target character.Character, skill skill.Skill) (affectedActor, affectedTarget character.Character) {
	aa, ta := skill.Use(actor.AttributeMap(), target.AttributeMap())
	actor.Affect(aa)
	target.Affect(ta)
	return actor, target
}

type skillTest struct{}

func (s *skillTest) Use(am, dm valueobject.AttributeTypeMap) (aa, ta []valueobject.Attribute) {
	return nil, []valueobject.Attribute{{valueobject.AttributeTypeHP, "-100"}}
}

func (s *skillTest) Name() string {
	return "skillTest"
}

func TestAct(t *testing.T) {
	attacker := character.NewCharacter("attacker",
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "100"},
	)
	defender := character.NewCharacter("defender",
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "100"},
	)
	skill := &skillTest{}

	wantAttacker := character.NewCharacter("attacker",
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "100"},
	)
	wantDefender := character.NewCharacter("defender",
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "0"},
		valueobject.Attribute{Type: valueobject.AttributeTypeDead},
	)

	affectedAttacker, affectedDefender := act(attacker, defender, skill)

	assert.Equal(t, wantAttacker, affectedAttacker)
	assert.Equal(t, wantDefender, affectedDefender)
}
