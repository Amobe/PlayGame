package battle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
)

type fakeSkill struct {
	used bool
}

func (s *fakeSkill) Use(am, dm character.AttributeTypeMap) (aa, ta []character.Attribute) {
	s.used = true
	return nil, nil
}

func (s *fakeSkill) Name() string {
	return "fakeSkill"
}

// Test skill is used in the battle fight. The ally and enemy skill should be used.
func TestBattle_FightUseSkill(t *testing.T) {
	ally := character.NewCharacter()
	allySkill := &fakeSkill{}
	enemy := character.NewCharacter()
	enemySkill := &fakeSkill{}
	enemyMob := battle.NewMob(enemy, enemySkill)

	b := battle.NewBattle("", ally, enemyMob)
	_, err := b.Fight([]character.Skill{allySkill})

	assert.NoError(t, err)
	assert.True(t, allySkill.used)
	assert.True(t, enemySkill.used)
}
