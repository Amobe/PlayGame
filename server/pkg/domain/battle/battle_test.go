package battle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

func TestBattle_Create(t *testing.T) {
	allyMinions := &battle.Minions{}
	enemyMinions := &battle.Minions{}
	minionSlot := battle.NewMinionSlot(allyMinions, enemyMinions)
	b, err := battle.CreateBattle("1", minionSlot)
	assert.NoError(t, err)
	assert.Equal(t, "1", b.ID())
	assert.Equal(t, battle.StatusUnspecified, b.Status())
	assert.Equal(t, minionSlot, b.MinionSlot())
}

func getSummoner(isDead bool) *battle.MockUnit {
	u := &battle.MockUnit{}
	u.On("GetAgi").Return(0)
	u.On("IsDead").Return(isDead)
	return u
}

func getUnit(isDead bool) *battle.MockUnit {
	u := &battle.MockUnit{}
	u.On("IsDead").Return(isDead)
	u.On("GetSkill").Return(vo.Skill{})
	return u
}

func getDeadUnit() *battle.MockUnit {
	return getUnit(true)
}

func TestBattle_FightToTheEnd_Draw(t *testing.T) {
	allyMinions := &battle.Minions{
		getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getSummoner(false),
	}
	enemyMinions := &battle.Minions{
		getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getSummoner(false),
	}
	minionSlot := battle.NewMinionSlot(allyMinions, enemyMinions)
	b, err := battle.CreateBattle("1", minionSlot)
	assert.NoError(t, err)

	affects, err := b.FightToTheEnd()
	assert.NoError(t, err)
	assert.Zero(t, len(affects))
	assert.Equalf(t, battle.StatusDraw, b.Status(), "battle status should be draw, but got %s", b.Status())
}

func TestBattle_FightToTheEnd_AllyWon(t *testing.T) {
	allyMinions := &battle.Minions{
		getUnit(false), getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getSummoner(false),
	}
	enemyMinions := &battle.Minions{
		getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getSummoner(true),
	}
	minionSlot := battle.NewMinionSlot(allyMinions, enemyMinions)
	b, err := battle.CreateBattle("1", minionSlot)
	assert.NoError(t, err)

	affects, err := b.FightToTheEnd()
	assert.NoError(t, err)
	assert.Zero(t, len(affects))
	assert.Equalf(t, battle.StatusWon, b.Status(), "battle status should be won, but got %s", b.Status())
}

func TestBattle_FightToTheEnd_AllyLost(t *testing.T) {
	alleyMinions := &battle.Minions{
		getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getSummoner(true),
	}
	enemyMinions := &battle.Minions{
		getUnit(false), getDeadUnit(), getDeadUnit(), getDeadUnit(), getDeadUnit(), getSummoner(false),
	}
	minionSlot := battle.NewMinionSlot(alleyMinions, enemyMinions)
	b, err := battle.CreateBattle("1", minionSlot)
	assert.NoError(t, err)

	affects, err := b.FightToTheEnd()
	assert.NoError(t, err)
	assert.Zero(t, len(affects))
	assert.Equalf(t, battle.StatusLost, b.Status(), "battle status should be lost, but got %s", b.Status())
}
