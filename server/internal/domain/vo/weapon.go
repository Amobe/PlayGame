package vo

import (
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	WeaponEmpty = Weapon{
		WeaponID:     "empty",
		WeaponType:   WeaponTypeEmpty,
		AttributeMap: NewAttributeMap(),
	}
	WeaponKnife = Weapon{
		WeaponID:   "knife",
		WeaponType: WeaponTypeKnife,
		Slot:       WeaponSlotMajorHand,
		Skill: NewSkill("slash", NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(2)),
			NewAttribute(AttributeTypeCD, decimal.NewFromInt(2)),
			NewAttribute(AttributeTypeATKB, decimal.NewFromFloat(1.4)),
			NewBoolAttribute(AttributeTypePhysicalDamage, true),
		)),
		AttributeMap: NewAttributeMap(),
	}
	WeaponDagger = Weapon{
		WeaponID:   "dagger",
		WeaponType: WeaponTypeDagger,
		Slot:       WeaponSlotAny,
		Skill: NewSkill("stab", NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(1)),
			NewAttribute(AttributeTypeCD, decimal.NewFromInt(1)),
			NewAttribute(AttributeTypeATKB, decimal.NewFromFloat(1.8)),
			NewBoolAttribute(AttributeTypePhysicalDamage, true),
		)),
		AttributeMap: NewAttributeMap(),
	}
	WeaponBow = Weapon{
		WeaponID:   "bow",
		WeaponType: WeaponTypeBow,
		Slot:       WeaponSlotBothHand,
		Skill: NewSkill("shoot", NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(5)),
			NewAttribute(AttributeTypeCD, decimal.NewFromInt(1)),
			NewAttribute(AttributeTypeATKB, decimal.NewFromFloat(0.6)),
			NewBoolAttribute(AttributeTypePhysicalDamage, true),
		)),
		AttributeMap: NewAttributeMap(),
	}
	WeaponAxe = Weapon{
		WeaponID:   "axe",
		WeaponType: WeaponTypeAxe,
		Slot:       WeaponSlotBothHand,
		Skill: NewSkill("chop", NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(1)),
			NewAttribute(AttributeTypeCD, decimal.NewFromInt(3)),
			NewAttribute(AttributeTypeATKB, decimal.NewFromInt(5)),
			NewBoolAttribute(AttributeTypePhysicalDamage, true),
		)),
		AttributeMap: NewAttributeMap(),
	}
	WeaponRod = Weapon{
		WeaponID:   "rod",
		WeaponType: WeaponTypeRod,
		Slot:       WeaponSlotMajorHand,
		Skill: NewSkill("magic_ball", NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(4)),
			NewAttribute(AttributeTypeCD, decimal.NewFromInt(3)),
			NewAttribute(AttributeTypeMATKB, decimal.NewFromInt(3)),
			NewBoolAttribute(AttributeTypeMagicalDamage, true),
		)),
		AttributeMap: NewAttributeMap(),
	}
	WeaponWand = Weapon{
		WeaponID:   "wand",
		WeaponType: WeaponTypeWand,
		Slot:       WeaponSlotMajorHand,
		Skill: NewSkill("magic_hit", NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(1)),
			NewAttribute(AttributeTypeCD, decimal.NewFromInt(2)),
			NewAttribute(AttributeTypeMATKB, decimal.NewFromFloat(2.1)),
			NewBoolAttribute(AttributeTypeMagicalDamage, true),
		)),
		AttributeMap: NewAttributeMap(),
	}
	WeaponHolyBook = Weapon{
		WeaponID:   "holy_book",
		WeaponType: WeaponTypeBook,
		Slot:       WeaponSlotMajorHand,
		Skill: NewSkill("cure", NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(-1)),
			NewAttribute(AttributeTypeCD, decimal.NewFromInt(2)),
			NewAttribute(AttributeTypeMATKB, decimal.NewFromFloat(1.2)),
			NewBoolAttribute(AttributeTypeHeal, true),
		)),
		AttributeMap: NewAttributeMap(),
	}
	WeaponCurseBook = Weapon{
		WeaponID:   "curse_book",
		WeaponType: WeaponTypeBook,
		Slot:       WeaponSlotMajorHand,
		Skill: NewSkill("curse", NewAttributeMap(
			NewAttribute(AttributeTypeTarget, decimal.NewFromInt(3)),
			NewAttribute(AttributeTypeCD, decimal.NewFromInt(2)),
			NewAttribute(AttributeTypeMATKB, decimal.NewFromInt(9)),
			NewBoolAttribute(AttributeTypeMagicalDamage, true),
		)),
		AttributeMap: NewAttributeMap(),
	}
	WeaponArrow  = Weapon{}
	WeaponShield = Weapon{}
)

type Weapon struct {
	WeaponID     string
	WeaponType   WeaponType
	Name         string
	Slot         WeaponSlot
	Skill        Skill
	AttributeMap AttributeMap
}

func NewWeapon(id string, weaponTypeStr string, weaponName string, skill Skill, attributes AttributeMap) (Weapon, error) {
	weaponType := ToWeaponType(weaponTypeStr)
	if weaponType == WeaponTypeUnspecified {
		return Weapon{}, fmt.Errorf("weapon type is unspecified: %s", weaponTypeStr)
	}

	return Weapon{
		WeaponID:     id,
		WeaponType:   weaponType,
		Name:         weaponName,
		AttributeMap: attributes,
	}, nil
}

func (w Weapon) ID() string {
	return w.WeaponID
}

type WeaponType string

func (w WeaponType) String() string {
	return string(w)
}

func ToWeaponType(s string) WeaponType {
	t, ok := weaponTypeMap[s]
	if !ok {
		return WeaponTypeUnspecified
	}
	return t
}

const (
	WeaponTypeUnspecified WeaponType = "unspecified"
	WeaponTypeEmpty       WeaponType = "empty"
	WeaponTypeKnife       WeaponType = "knife"
	WeaponTypeDagger      WeaponType = "dagger"
	WeaponTypeBow         WeaponType = "bow"
	WeaponTypeAxe         WeaponType = "axe"
	WeaponTypeRod         WeaponType = "rod"
	WeaponTypeWand        WeaponType = "wand"
	WeaponTypeBook        WeaponType = "book"
	WeaponTypeArrow       WeaponType = "arrow"
	WeaponTypeShield      WeaponType = "shield"
)

type WeaponSlot string

func (w WeaponSlot) String() string {
	return string(w)
}

const (
	WeaponSlotAny       WeaponSlot = "any"
	WeaponSlotMajorHand WeaponSlot = "major"
	WeaponSlotMinorHand WeaponSlot = "minor"
	WeaponSlotBothHand  WeaponSlot = "both"
)

var weaponTypeMap = map[string]WeaponType{
	WeaponTypeEmpty.String():  WeaponTypeEmpty,
	WeaponTypeKnife.String():  WeaponTypeKnife,
	WeaponTypeDagger.String(): WeaponTypeDagger,
	WeaponTypeBow.String():    WeaponTypeBow,
	WeaponTypeAxe.String():    WeaponTypeAxe,
	WeaponTypeRod.String():    WeaponTypeRod,
	WeaponTypeWand.String():   WeaponTypeWand,
	WeaponTypeBook.String():   WeaponTypeBook,
	WeaponTypeArrow.String():  WeaponTypeArrow,
	WeaponTypeShield.String(): WeaponTypeShield,
}

//var weaponTypeSlotMap = map[WeaponType]WeaponSlot{
//	WeaponTypeEmpty:  WeaponSlotAny,
//	WeaponTypeKnife:  WeaponSlotMajorHand,
//	WeaponTypeDagger: WeaponSlotAny,
//	WeaponTypeBow:    WeaponSlotMajorHand,
//	WeaponTypeAxe:    WeaponSlotBothHand,
//	WeaponTypeRod:    WeaponSlotAny,
//	WeaponTypeWand:   WeaponSlotBothHand,
//	WeaponTypeOrb:    WeaponSlotAny,
//	WeaponTypeBook:   WeaponSlotAny,
//	WeaponTypeArrow:  WeaponSlotMinorHand,
//	WeaponTypeShield: WeaponSlotMinorHand,
//}

//var weaponTypePairMap = map[string]bool{
//	getWeaponTypePairKey(WeaponTypeKnife, WeaponTypeShield):  true,
//	getWeaponTypePairKey(WeaponTypeDagger, WeaponTypeDagger): true,
//	getWeaponTypePairKey(WeaponTypeDagger, WeaponTypeShield): true,
//	getWeaponTypePairKey(WeaponTypeBow, WeaponTypeArrow):     true,
//	getWeaponTypePairKey(WeaponTypeRod, WeaponTypeOrb):       true,
//	getWeaponTypePairKey(WeaponTypeRod, WeaponTypeBook):      true,
//	getWeaponTypePairKey(WeaponTypeOrb, WeaponTypeRod):       true,
//	getWeaponTypePairKey(WeaponTypeOrb, WeaponTypeBook):      true,
//	getWeaponTypePairKey(WeaponTypeBook, WeaponTypeRod):      true,
//	getWeaponTypePairKey(WeaponTypeBook, WeaponTypeOrb):      true,
//}

//func getWeaponTypePairKey(major, minor WeaponType) string {
//	return fmt.Sprintf("%s:%s", major.String(), minor.String())
//}
