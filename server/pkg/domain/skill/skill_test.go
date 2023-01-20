package skill_test

import (
	"testing"

	"github.com/Amobe/PlayGame/server/pkg/domain/skill"
	"github.com/Amobe/PlayGame/server/pkg/domain/valueobject"

	"github.com/stretchr/testify/assert"
)

func TestSkillHit(t *testing.T) {
	s := skill.NewSkillPoisonHit()

	am := valueobject.NewAttributeTypeMap()
	am.Insert(valueobject.Attribute{valueobject.AttributeTypeATK, "10"})
	dm := valueobject.NewAttributeTypeMap()
	dm.Insert(valueobject.Attribute{valueobject.AttributeTypeHP, "100"})

	aa, da := s.Use(am, dm)

	assert.Nil(t, aa)
	assert.NotNil(t, da)
}
