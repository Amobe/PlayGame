package battle

import "github.com/Amobe/PlayGame/server/internal/domain/vo"

type Minions [6]vo.GroundUnit

func NewAllyMinions(characters []vo.Character) *Minions {
	return NewMinions(true, characters)
}

func NewEnemyMinions(characters []vo.Character) *Minions {
	return NewMinions(false, characters)
}

// NewMinions creates a new Minions of Ally or Enemy.
// Assume that the length of characters is 6.
// The last character is the summoner.
func NewMinions(isAlly bool, characters []vo.Character) *Minions {
	startIdx := 1
	if !isAlly {
		startIdx = 7
	}
	m := &Minions{}
	for i, c := range characters {
		m[i] = vo.NewGroundUnit(vo.GroundIdx(startIdx+i), c)
	}
	return m
}

func (m *Minions) Get(idx vo.CampIdx) vo.GroundUnit {
	return m[idx-1]
}

func (m *Minions) Set(idx vo.CampIdx, unit vo.GroundUnit) {
	m[idx-1] = unit
}

func (m *Minions) GetSummoner() vo.GroundUnit {
	return m.Get(vo.SummonerCampIdx)
}

func (m *Minions) SetSummoner(summoner vo.GroundUnit) {
	m.Set(vo.SummonerCampIdx, summoner)
}
