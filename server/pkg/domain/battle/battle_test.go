package battle_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/valueobject"
)

type fakeSkill struct {
	usedFunc func() (aa, ta []valueobject.Attribute)
}

func (s fakeSkill) Use(am, dm valueobject.AttributeTypeMap) (aa, ta []valueobject.Attribute) {
	if s.usedFunc != nil {
		return s.usedFunc()
	}
	return nil, nil
}

func (s fakeSkill) Name() string {
	return "fakeSkill"
}

// Test skill is used in the battle fight. The ally and enemy skill should be used.
func TestBattle_FightUseSkill(t *testing.T) {
	ally := character.NewCharacter(
		"ally",
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "10"},
	)
	isAllySkillUsed := false
	allySkill := fakeSkill{
		usedFunc: func() (aa, ta []valueobject.Attribute) {
			isAllySkillUsed = true
			return
		},
	}
	allySlot := battle.NewSlot(allySkill)
	enemy := character.NewCharacter(
		"enemy",
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "10"},
	)
	isEnemySkillUsed := false
	enemySkill := fakeSkill{
		usedFunc: func() (aa, ta []valueobject.Attribute) {
			isEnemySkillUsed = true
			return
		},
	}
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
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "10"},
	)
	allySkill := fakeSkill{
		usedFunc: func() (aa, ta []valueobject.Attribute) {
			usedOrder = append(usedOrder, "ally")
			return
		},
	}
	allySlot := battle.NewSlot(allySkill)
	enemy := character.NewCharacter(
		"enemy",
		valueobject.Attribute{Type: valueobject.AttributeTypeAGI, Value: "10"},
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "10"},
	)
	enemySkill := fakeSkill{
		usedFunc: func() (aa, ta []valueobject.Attribute) {
			usedOrder = append(usedOrder, "enemy")
			return
		},
	}
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
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "10"},
	)
	allySlot := battle.NewSlot(fakeSkill{})
	enemy := character.NewCharacter("enemy")
	enemySlot := battle.NewSlot(fakeSkill{})

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
	allySlot := battle.NewSlot(fakeSkill{})
	enemy := character.NewCharacter(
		"enemy",
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "10"},
	)
	enemySlot := battle.NewSlot(fakeSkill{})

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
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "51"},
	)
	enemy := character.NewCharacter("enemy",
		valueobject.Attribute{Type: valueobject.AttributeTypeHP, Value: "1"},
	)
	enemySkill := fakeSkill{
		usedFunc: func() (aa, ta []valueobject.Attribute) {
			return nil, []valueobject.Attribute{
				{Type: valueobject.AttributeTypeHP, Value: "-1"},
			}
		},
	}
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
	assert.Equal(t, 1, b.Fighter("ally").AttributeMap().Get(valueobject.AttributeTypeHP))
	assert.Equal(t, 1, len(drawEvents))
}
