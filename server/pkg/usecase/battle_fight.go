package usecase

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/domain/battle"
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

	affects, err := b.FightToTheEnd()
	if err != nil {
		err = fmt.Errorf("battle fight: %w", err)
		return
	}

	err = u.battleRepo.Save(b)
	if err != nil {
		err = fmt.Errorf("battle repository save: %w", err)
	}

	out.Affects = affects
	return
}
