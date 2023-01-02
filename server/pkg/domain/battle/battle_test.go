package battle_test

import (
	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTwentyHPCharacter() character.Character {
	return character.Character{
		Basement: []character.Attribute{
			{Type: character.AttributeTypeHP, Value: "20"},
		},
	}
}

type fakeSkill struct {
	used bool
}

func (s *fakeSkill) Use(am, dm character.AttributeTypeMap) (aa, ta []character.Attribute) {
	s.used = true
	return nil, nil
}

// Test skill is used in the battle fight. The ally and enemy skill should be used.
func TestBattle_FightUseSkill(t *testing.T) {
	ally := getTwentyHPCharacter()
	allySkill := &fakeSkill{}
	enemy := getTwentyHPCharacter()
	enemySkill := &fakeSkill{}
	enemyMob := battle.NewMob(enemy, enemySkill)

	b := battle.NewBattle(ally, enemyMob)
	err := b.Fight([]character.Skill{allySkill})

	assert.NoError(t, err)
	assert.True(t, allySkill.used)
	assert.True(t, enemySkill.used)
}
