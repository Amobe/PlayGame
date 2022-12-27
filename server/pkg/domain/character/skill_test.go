package character_test

import (
	"testing"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/stretchr/testify/assert"
)

func TestSkillHit(t *testing.T) {
	s := character.NewSkillPoisonHit()

	am := character.NewAttributeTypeMap()
	am.Insert(character.Attribute{character.AttributeTypeATK, "10"})
	dm := character.NewAttributeTypeMap()
	dm.Insert(character.Attribute{character.AttributeTypeHP, "100"})

	aa, da := s.Use(am, dm)

	assert.Nil(t, aa)
	assert.NotNil(t, da)
}
