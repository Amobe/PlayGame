package battle_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

// Test skill is used in the battle fight. The ally and enemy skill should be used.
func TestBattle_FightUseSkill(t *testing.T) {
	ally := character.NewCharacter(
		"ally",
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
	)
	isAllySkillUsed := false
	allySkill := vo.NewCustomSkill(func(actorAttr, targetAttr vo.AttributeMap) (aa, ta []vo.Attribute) {
		isAllySkillUsed = true
		return
	})
	allySlot := battle.NewSlot(allySkill)
	enemy := character.NewCharacter(
		"enemy",
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
	)
	isEnemySkillUsed := false
	enemySkill := vo.NewCustomSkill(func(actorAttr, targetAttr vo.AttributeMap) (aa, ta []vo.Attribute) {
		isEnemySkillUsed = true
		return
	})
	enemySlot := battle.NewSlot(enemySkill)

	b, _ := battle.CreateBattle("", ally, enemy, enemySlot)
	b.SetAllySlot(allySlot)

	err := b.Fight()
	assert.NoError(t, err)
	assert.True(t, isAllySkillUsed)
	assert.True(t, isEnemySkillUsed)
}

// Test skill is used in agi order.
func TestBattle_FightInOrder(t *testing.T) {
	var usedOrder []string
	ally := character.NewCharacter(
		"ally",
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
	)
	allySkill := vo.NewCustomSkill(func(actorAttr, targetAttr vo.AttributeMap) (aa, ta []vo.Attribute) {
		usedOrder = append(usedOrder, "ally")
		return
	})

	allySlot := battle.NewSlot(allySkill)
	enemy := character.NewCharacter(
		"enemy",
		vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(10)),
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
	)
	enemySkill := vo.NewCustomSkill(func(actorAttr, targetAttr vo.AttributeMap) (aa, ta []vo.Attribute) {
		usedOrder = append(usedOrder, "enemy")
		return
	})
	enemySlot := battle.NewSlot(enemySkill)

	b, _ := battle.CreateBattle("", ally, enemy, enemySlot)
	_ = b.SetAllySlot(allySlot)

	err := b.Fight()

	want := []string{"enemy", "ally"}
	assert.NoError(t, err)
	assert.Equal(t, want, usedOrder)
}

func TestBattle_FightToWin(t *testing.T) {
	ally := character.NewCharacter(
		"ally",
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
	)
	allySlot := battle.NewSlot()
	enemy := character.NewCharacter("enemy")
	enemySlot := battle.NewSlot()

	b, _ := battle.CreateBattle("", ally, enemy, enemySlot)
	_ = b.SetAllySlot(allySlot)
	_ = b.Fight()

	assert.Equal(t, battle.StatusWon, b.Status())

	events := b.Events()
	lastEvent := events[len(events)-1]
	assert.IsType(t, battle.EventBattleWon{}, lastEvent)
}

func TestBattle_FightToLose(t *testing.T) {
	ally := character.NewCharacter("ally")
	allySlot := battle.NewSlot()
	enemy := character.NewCharacter(
		"enemy",
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
	)
	enemySlot := battle.NewSlot()

	b, _ := battle.CreateBattle("", ally, enemy, enemySlot)
	_ = b.SetAllySlot(allySlot)
	_ = b.Fight()

	assert.Equal(t, battle.StatusLost, b.Status())

	events := b.Events()
	lastEvent := events[len(events)-1]
	assert.IsType(t, battle.EventBattleLost{}, lastEvent)
}

// give an ally has 51 hp, and an enemy attack ally take 1 hp a time.
// when fight number reach to 50 round limitation
// then the battle should be draw
// and the enemy should attack 50 times
// and the ally has 1 hp left
func TestBattle_FightToTheEnd(t *testing.T) {
	ally := character.NewCharacter("ally",
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(51)),
	)
	enemy := character.NewCharacter("enemy",
		vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(1)),
	)
	enemySkill := vo.NewCustomSkill(func(actorAttr, targetAttr vo.AttributeMap) (aa, ta []vo.Attribute) {
		return nil, []vo.Attribute{
			vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(-1)),
		}
	})
	enemySlot := battle.NewSlot(enemySkill)

	b, _ := battle.CreateBattle("", ally, enemy, enemySlot)
	_ = b.FightToTheEnd()

	var foughtEvents []battle.EventBattleFought
	var drawEvents []battle.EventBattleDraw
	for _, battleEvent := range b.Events() {
		switch event := battleEvent.(type) {
		case battle.EventBattleFought:
			foughtEvents = append(foughtEvents, event)
		case battle.EventBattleDraw:
			drawEvents = append(drawEvents, event)
		}
	}

	assert.Equal(t, 50, len(foughtEvents))
	expectedHP := decimal.NewFromInt(1)
	actualHP := b.Fighter("ally").AttributeMap().Get(vo.AttributeTypeHP).Value
	utils.AssertDecimal(t, expectedHP, actualHP)
	assert.Equal(t, 1, len(drawEvents))
}
