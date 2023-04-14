package vo

// CampIdx is the index of the camp slots.
// The camp slots only contain the minions with same camp.
// The CampIdx of the minions is between 1 and 5.
// The CampIdx of the summoner is 6.
type CampIdx int

var SummonerCampIdx = CampIdx(6)

func (c CampIdx) IsSummoner() bool {
	return c == SummonerCampIdx
}
