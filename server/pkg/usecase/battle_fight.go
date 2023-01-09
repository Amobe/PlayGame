package usecase

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/utils"
)

type BattleFightIn struct {
	BattleID string
}

type BattleFightOut struct {
	Affects []battle.Affect
}

type BattleFightUsecase struct {
	battleRepo battle.Repository
}

func NewBattleFightUsecase(battleRepo battle.Repository) *BattleFightUsecase {
	return &BattleFightUsecase{
		battleRepo: battleRepo,
	}
}

func (u *BattleFightUsecase) Execute(in BattleFightIn) (out BattleFightOut, err error) {
	b, err := u.battleRepo.Get(in.BattleID)
	if err != nil {
		err = fmt.Errorf("battle repository get %s: %w", in.BattleID, err)
		return
	}

	fakeSkill := []character.Skill{
		character.NewSkillPoisonHit(),
	}
	err = b.Fight(fakeSkill)
	if err != nil {
		err = fmt.Errorf("battle fight: %w", err)
		return
	}

	err = u.battleRepo.Save(b)
	if err != nil {
		err = fmt.Errorf("battle repository save: %w", err)
	}

	events := b.Events()
	for _, ev := range events {
		fmt.Println(utils.ToString(ev))
		if foughtEv, ok := ev.(battle.EventBattleFought); ok {
			out.Affects = foughtEv.Affects
		}
	}
	return
}
