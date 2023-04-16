package battle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/vo"
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

func getUndeadUnit(groundIdx vo.GroundIdx) vo.Character {
	return vo.NewCharacterWithSkill(groundIdx, vo.EmptySkill, vo.NewAttributeMap())
}

func getDeadUnit(groundIdx vo.GroundIdx) vo.Character {
	return vo.NewCharacterWithSkill(groundIdx, vo.EmptySkill, vo.NewAttributeMap(vo.DeadAttribute))
}

func TestBattle_FightToTheEnd_Draw(t *testing.T) {
	allyMinions := &battle.Minions{
		getDeadUnit(1),
		getDeadUnit(2),
		getDeadUnit(3),
		getDeadUnit(4),
		getDeadUnit(5),
		getUndeadUnit(6),
	}
	enemyMinions := &battle.Minions{
		getDeadUnit(7),
		getDeadUnit(8),
		getDeadUnit(9),
		getDeadUnit(10),
		getDeadUnit(11),
		getUndeadUnit(12),
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
		getUndeadUnit(1),
		getDeadUnit(2),
		getDeadUnit(3),
		getDeadUnit(4),
		getDeadUnit(5),
		getUndeadUnit(6),
	}
	enemyMinions := &battle.Minions{
		getDeadUnit(7),
		getDeadUnit(8),
		getDeadUnit(9),
		getDeadUnit(10),
		getDeadUnit(11),
		getDeadUnit(12),
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
		getDeadUnit(1),
		getDeadUnit(2),
		getDeadUnit(3),
		getDeadUnit(4),
		getDeadUnit(5),
		getDeadUnit(6),
	}
	enemyMinions := &battle.Minions{
		getUndeadUnit(7),
		getDeadUnit(8),
		getDeadUnit(9),
		getDeadUnit(10),
		getDeadUnit(11),
		getUndeadUnit(12),
	}
	minionSlot := battle.NewMinionSlot(alleyMinions, enemyMinions)
	b, err := battle.CreateBattle("1", minionSlot)
	assert.NoError(t, err)

	affects, err := b.FightToTheEnd()
	assert.NoError(t, err)
	assert.Zero(t, len(affects))
	assert.Equalf(t, battle.StatusLost, b.Status(), "battle status should be lost, but got %s", b.Status())
}
