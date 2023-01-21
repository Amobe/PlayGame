package character_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"

	"github.com/stretchr/testify/assert"
)

func TestSkillHit(t *testing.T) {
	s := character.NewSkillPoisonHit()

	am := vo.NewAttributeTypeMap(vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)))
	dm := vo.NewAttributeTypeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))

	aa, da := s.Use(am, dm)

	assert.Nil(t, aa)
	assert.NotNil(t, da)
}
