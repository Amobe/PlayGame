package vo

import "github.com/shopspring/decimal"

type Equipment struct {
	MajorHand Weapon
	MinorHand Weapon
	Suite     Suite
}

func NewEquipment() Equipment {
	return Equipment{
		MajorHand: WeaponEmpty,
		MinorHand: WeaponEmpty,
		Suite:     EmptySuite,
	}
}

func (e Equipment) EquipWeapon(w Weapon) Equipment {
	switch w.Slot {
	case WeaponSlotMajorHand:
		e.MajorHand = w
	case WeaponSlotMinorHand:
		e.MinorHand = w
	case WeaponSlotBothHand:
		e.MajorHand = w
		e.MinorHand = WeaponEmpty
	}
	e.Suite = EmptySuite
	pairSet := e.MajorHand.Name + "+" + e.MinorHand.Name
	if suiteName, ok := SuitePairSet[pairSet]; ok {
		e.Suite = SuiteSet[suiteName]
	}
	return Equipment{
		MajorHand: e.MajorHand,
		MinorHand: e.MinorHand,
		Suite:     e.Suite,
	}
}

func (e Equipment) GetAttributes() []Attribute {
	var sum []Attribute
	sum = append(sum, e.MajorHand.Attributes...)
	sum = append(sum, e.MinorHand.Attributes...)
	sum = append(sum, e.Suite.Attributes...)
	return sum
}

type Suite struct {
	ID         string
	Name       string
	Pair       []string
	Attributes []Attribute
}

var EmptySuite = Suite{"2d462197-e311-4d29-8e2c-6df9a2f76582", "Empty", nil, nil}

var SuiteSet = map[string]Suite{
	"Empty": EmptySuite,
	"Physical": {"87d1454d-c5b9-48c3-9928-1a9e003ee9c6", "Physical",
		[]string{"Knife", "Shield"},
		[]Attribute{
			NewAttribute(AttributeTypeATK, decimal.NewFromInt(10)),
			NewAttribute(AttributeTypeDEF, decimal.NewFromInt(10)),
		},
	},
	"Magical": {"09a2b7e8-2943-493b-bc48-3d413e969bca", "Magical",
		[]string{"Book", "Ball"},
		[]Attribute{
			NewAttribute(AttributeTypeMATK, decimal.NewFromInt(10)),
			NewAttribute(AttributeTypeMDEF, decimal.NewFromInt(10)),
		},
	},
}

var SuitePairSet = map[string]string{
	"Knife+Shield": "Physical",
	"Book+Ball":    "Magical",
}
