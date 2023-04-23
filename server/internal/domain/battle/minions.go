package battle

import "github.com/Amobe/PlayGame/server/internal/domain/vo"

type Minions [6]vo.Character

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
	//startIdx := 1
	//if !isAlly {
	//	startIdx = 7
	//}
	m := &Minions{}
	copy(m[:], characters)
	return m
}

func (m *Minions) Get(idx vo.CampIdx) vo.Character {
	return m[idx-1]
}

func (m *Minions) Set(idx vo.CampIdx, unit vo.Character) {
	m[idx-1] = unit
}

func (m *Minions) GetSummoner() vo.Character {
	return m.Get(vo.SummonerCampIdx)
}

func (m *Minions) SetSummoner(summoner vo.Character) {
	m.Set(vo.SummonerCampIdx, summoner)
}
