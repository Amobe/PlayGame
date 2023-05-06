package vo

import "fmt"

const CampSize = 6

type Camp [CampSize]Character

func NewCamp(characters []Character) (Camp, error) {
	if len(characters) > CampSize {
		return Camp{}, fmt.Errorf("invalid characters length: %d", len(characters))
	}
	c := Camp{}
	copy(c[:], characters)
	return c, nil
}

func (c Camp) Get(idx CampIdx) (Character, error) {
	if idx > CampSize {
		return Character{}, fmt.Errorf("invalid camp idx: %d", idx)
	}
	return c[idx-1], nil
}

func (c Camp) Set(idx CampIdx, unit Character) (Camp, error) {
	if idx > CampSize {
		return Camp{}, fmt.Errorf("invalid camp idx: %d", idx)
	}
	newCamp := Camp{}
	copy(newCamp[:], c[:])
	newCamp[idx-1] = unit
	return newCamp, nil
}

func (c Camp) GetSummoner() (Character, error) {
	return c.Get(SummonerCampIdx)
}

func (c Camp) SetSummoner(summoner Character) (Camp, error) {
	return c.Set(SummonerCampIdx, summoner)
}
