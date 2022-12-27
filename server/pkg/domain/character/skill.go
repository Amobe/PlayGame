package character

import (
	"strconv"
)

type AffectFunc func(attacker, defender AttributeTypeMap) (attackerAffect, defenderAffect []Attribute)

type Skill struct {
	SkillID string
	Name    string
	Affect  AffectFunc
}

func HitAffect(attacker, defender AttributeTypeMap) (attackerAffect, defenderAffect []Attribute) {
	var attackerATK, defenderDEF int
	if atk, ok := attacker[AttributeTypeATK]; ok {
		attackerATK = atk.GetInt()
	}
	if def, ok := defender[AttributeTypeDEF]; ok {
		defenderDEF = def.GetInt()
	}
	damage := -1 * (attackerATK * attackerATK) / (attackerATK + 2*defenderDEF)
	defenderAffect = []Attribute{
		{
			Type:  AttributeTypeHP,
			Value: strconv.Itoa(damage),
		},
	}
	return nil, defenderAffect
}

var SkillMap = map[string]Skill{
	"hit": {SkillID: "hit", Name: "hit", Affect: HitAffect},
}
