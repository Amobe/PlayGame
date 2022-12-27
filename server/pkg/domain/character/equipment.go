package character

type Equipment struct {
	MajorHand Weapon
	MinorHand Weapon
	Suite     Suite
}

func NewEquipment() Equipment {
	return Equipment{
		MajorHand: EmptyWeapon,
		MinorHand: EmptyWeapon,
		Suite:     EmptySuite,
	}
}

func (e *Equipment) EquipWeapon(w Weapon) {
	switch w.Slot {
	case WeaponSlotMajorHand:
		e.MajorHand = w
	case WeaponSlotMinorHand:
		e.MinorHand = w
	case WeaponSlotBothHand:
		e.MajorHand = w
		e.MinorHand = EmptyWeapon
	}
	e.Suite = EmptySuite
	pairSet := e.MajorHand.Name + "+" + e.MinorHand.Name
	if suiteName, ok := SuitePairSet[pairSet]; ok {
		e.Suite = SuiteSet[suiteName]
	}
}

func (e *Equipment) GetAttributes() []Attribute {
	var sum []Attribute
	sum = append(sum, e.MajorHand.Attributes...)
	sum = append(sum, e.MinorHand.Attributes...)
	sum = append(sum, e.Suite.Attributes...)
	return sum
}

type Weapon struct {
	ID         string
	Name       WeaponName
	Slot       WeaponSlot
	Attributes []Attribute
}

type WeaponName string
type WeaponSlot string

const (
	WeaponSlotAny       WeaponSlot = "any"
	WeaponSlotMajorHand WeaponSlot = "major_hand"
	WeaponSlotMinorHand WeaponSlot = "minor_hand"
	WeaponSlotBothHand  WeaponSlot = "both_hand"
)

type Suite struct {
	ID         string
	Name       string
	Pair       []string
	Attributes []Attribute
}

var EmptyWeapon = Weapon{"60dad481-527d-4132-bf2f-7e8eab8ce136", "MajorEmpty", WeaponSlotAny, nil}

var WeaponSet = map[string]Weapon{
	"Empty": EmptyWeapon,
	"Knife": {"4b3867ab-e54d-4f34-a014-c2f87e1906f5", "Knife", WeaponSlotMajorHand, []Attribute{
		{AttributeTypeATK, "10"},
	}},
	"Shield": {"de977351-4a5f-4559-8b0c-ff09337a979d", "Shield", WeaponSlotMinorHand, []Attribute{
		{AttributeTypeDEF, "10"},
	}},
	"Book": {"c578dbe6-9758-4c4b-a9f4-22e577a2b9bb", "Book", WeaponSlotMajorHand, []Attribute{
		{AttributeTypeMATK, "10"},
	}},
	"Ball": {"d8a7a657-fb7a-4ce9-83b9-f2a833032dc0", "Ball", WeaponSlotMinorHand, []Attribute{
		{AttributeTypeMDEF, "10"},
	}},
	"Axe": {"be52ce42-2a8e-4324-b84b-fbad5761f586", "Axe", WeaponSlotBothHand, []Attribute{
		{AttributeTypeATK, "25"},
	}},
}

var EmptySuite = Suite{"2d462197-e311-4d29-8e2c-6df9a2f76582", "Empty", nil, nil}

var SuiteSet = map[string]Suite{
	"Empty": EmptySuite,
	"Physical": {"87d1454d-c5b9-48c3-9928-1a9e003ee9c6", "Physical",
		[]string{"Knife", "Shield"},
		[]Attribute{
			{AttributeTypeATK, "10"},
			{AttributeTypeDEF, "10"},
		},
	},
	"Magical": {"09a2b7e8-2943-493b-bc48-3d413e969bca", "Magical",
		[]string{"Book", "Ball"},
		[]Attribute{
			{AttributeTypeMATK, "10"},
			{AttributeTypeMDEF, "10"},
		},
	},
}

var SuitePairSet = map[WeaponName]string{
	"Knife+Shield": "Physical",
	"Book+Ball":    "Magical",
}
