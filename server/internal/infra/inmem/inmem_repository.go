package inmem

import (
	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/domain/stage"
	"github.com/Amobe/PlayGame/server/internal/domain/vo"
)

var _ stage.Repository = &StageRepository{}

type StageRepository = InmemStorage[*stage.Stage]

func NewInmemStageRepository() *StageRepository {
	s := NewInmemStorage[*stage.Stage]()
	s.Create(stage1)
	s.Create(stage2)
	return s
}

var _ battle.Repository = &BattleRepository{}

type BattleRepository = inmemEventStorage[*battle.Battle]

func NewInmemBattleRepository() *BattleRepository {
	return newInmemEventStorage[*battle.Battle](battle.AggregatorLoader)
}

var (
	stage1 = &stage.Stage{
		StageID: "stage1",
		Characters: []vo.Character{
			vo.NewCharacter(7, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
			vo.NewCharacter(8, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
			vo.NewCharacter(9, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
			vo.NewCharacter(10, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
			vo.NewCharacter(11, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
			vo.NewCharacter(12, vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)))),
		},
	}

	stage2 = &stage.Stage{
		StageID: "stage2",
		Characters: []vo.Character{
			vo.NewCharacterWithSkill(
				7,
				vo.SkillSlash,
				vo.NewAttributeMap(
					vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
					vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
				)),
			vo.NewCharacterWithSkill(
				8,
				vo.SkillSlash,
				vo.NewAttributeMap(
					vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
					vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
				)),
			vo.NewCharacterWithSkill(
				9,
				vo.SkillSlash,
				vo.NewAttributeMap(
					vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
					vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
				)),
			vo.NewCharacterWithSkill(
				10,
				vo.SkillSlash,
				vo.NewAttributeMap(
					vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
					vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
				)),
			vo.NewCharacterWithSkill(
				11,
				vo.SkillSlash,
				vo.NewAttributeMap(
					vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)),
					vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
				)),
			vo.NewCharacter(
				12,
				vo.NewAttributeMap(
					vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(100)),
				)),
		},
	}
)
