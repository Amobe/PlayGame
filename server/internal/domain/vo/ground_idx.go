package vo

// GroundIdx is the index on the battleground.
// The battleground contains the ally and enemy minions.
// The GroundIdx for ally minions is between 1 and 5.
// The GroundIdx for ally summoner is 6.
// The GroundIdx for enemy minions is between 7 and 11.
// The GroundIdx for enemy summoner is 12.
type GroundIdx int

var (
	AllySummonerGroundIdx  = GroundIdx(6)
	EnemySummonerGroundIdx = GroundIdx(12)
)

func (g GroundIdx) ToCampIdx() CampIdx {
	if g > AllySummonerGroundIdx {
		return CampIdx(g - 6)
	}
	return CampIdx(g)
}

func (g GroundIdx) IsEnemy() bool {
	return g > 6
}

func (g GroundIdx) GetOppositeIdx() GroundIdx {
	if g.IsEnemy() {
		return g - 6
	}
	return g + 6
}

func (g GroundIdx) ToInt32() int32 {
	return int32(g)
}

// EqualTo returns true if the GroundIdx is equal to the other GroundIdx.
func (g GroundIdx) EqualTo(other GroundIdx) bool {
	return g == other
}
