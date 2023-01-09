package battle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
)

type fakeSkill struct {
	usedFunc func()
}

func (s fakeSkill) Use(am, dm character.AttributeTypeMap) (aa, ta []character.Attribute) {
	s.usedFunc()
	return nil, nil
}

func (s fakeSkill) Name() string {
	return "fakeSkill"
}

// Test skill is used in the battle fight. The ally and enemy skill should be used.
func TestBattle_FightUseSkill(t *testing.T) {
	ally := character.NewCharacter("ally")
	isAllySkillUsed := false
	allySkill := fakeSkill{
		usedFunc: func() {
			isAllySkillUsed = true
		},
	}
	enemy := character.NewCharacter("enemy")
	isEnemySkillUsed := false
	enemySkill := fakeSkill{
		usedFunc: func() {
			isEnemySkillUsed = true
		},
	}
	enemyMobs := []battle.Fighter{battle.NewMob(enemy, enemySkill)}

	b, _ := battle.CreateBattle("", ally, enemyMobs)

	err := b.Fight([]character.Skill{allySkill})
	assert.NoError(t, err)
	assert.True(t, isAllySkillUsed)
	assert.True(t, isEnemySkillUsed)
}
