package vo

import "fmt"

var WeaponEmpty = Weapon{
	WeaponID:     "empty",
	WeaponType:   WeaponTypeEmpty,
	Slot:         WeaponSlotAny,
	AttributeMap: nil,
}

type Weapon struct {
	WeaponID     string
	WeaponType   WeaponType
	Name         string
	Slot         WeaponSlot
	Attributes   []Attribute // TODO: deprecated
	AttributeMap AttributeMap
}

func NewWeapon(id string, weaponTypeStr string, weaponName string, slotStr string, attributes ...Attribute) (Weapon, error) {
	weaponType := ToWeaponType(weaponTypeStr)
	if weaponType == WeaponTypeUnspecified {
		return Weapon{}, fmt.Errorf("weapon type is unspecified: %s", weaponTypeStr)
	}
	slot := ToWeaponSlot(slotStr)
	if slot == WeaponSlotUnspecified {
		return Weapon{}, fmt.Errorf("weapon slot is unspecified: %s", slotStr)

	}
	return Weapon{
		WeaponID:     id,
		WeaponType:   weaponType,
		Name:         weaponName,
		Slot:         slot,
		AttributeMap: NewAttributeMap(attributes...),
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
	switch s {
	case WeaponTypeKnife.String():
		return WeaponTypeKnife
	case WeaponTypeShield.String():
		return WeaponTypeShield
	case WeaponTypeDagger.String():
		return WeaponTypeDagger
	case WeaponTypeAxe.String():
		return WeaponTypeAxe
	default:
		return WeaponTypeUnspecified
	}
}

const (
	WeaponTypeUnspecified WeaponType = "unspecified"
	WeaponTypeEmpty       WeaponType = "empty"
	WeaponTypeKnife       WeaponType = "knife"
	WeaponTypeShield      WeaponType = "shield"
	WeaponTypeDagger      WeaponType = "dagger"
	WeaponTypeAxe         WeaponType = "axe"
)

type WeaponSlot string

func (w WeaponSlot) String() string {
	return string(w)
}

func ToWeaponSlot(s string) WeaponSlot {
	switch s {
	case WeaponSlotAny.String():
		return WeaponSlotAny
	case WeaponSlotMajorHand.String():
		return WeaponSlotMajorHand
	case WeaponSlotMinorHand.String():
		return WeaponSlotMinorHand
	case WeaponSlotBothHand.String():
		return WeaponSlotBothHand
	default:
		return WeaponSlotUnspecified
	}
}

const (
	WeaponSlotUnspecified WeaponSlot = "unspecified"
	WeaponSlotAny         WeaponSlot = "any"
	WeaponSlotMajorHand   WeaponSlot = "major"
	WeaponSlotMinorHand   WeaponSlot = "minor"
	WeaponSlotBothHand    WeaponSlot = "both"
)
