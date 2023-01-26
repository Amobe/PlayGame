package vo

import "fmt"

var WeaponEmpty = Weapon{
	WeaponID:     "empty",
	WeaponType:   WeaponTypeEmpty,
	AttributeMap: NewAttributeMap(),
}

type Weapon struct {
	WeaponID     string
	WeaponType   WeaponType
	Name         string
	Slot         WeaponSlot
	SkillType    SkillType
	AttributeMap AttributeMap
}

func NewWeapon(id string, weaponTypeStr string, weaponName string, skillTypeStr string, attributes ...Attribute) (Weapon, error) {
	weaponType := ToWeaponType(weaponTypeStr)
	if weaponType == WeaponTypeUnspecified {
		return Weapon{}, fmt.Errorf("weapon type is unspecified: %s", weaponTypeStr)
	}

	return Weapon{
		WeaponID:     id,
		WeaponType:   weaponType,
		Name:         weaponName,
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
	WeaponTypeOrb         WeaponType = "orb"
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
	WeaponTypeOrb.String():    WeaponTypeOrb,
	WeaponTypeBook.String():   WeaponTypeBook,
	WeaponTypeArrow.String():  WeaponTypeArrow,
	WeaponTypeShield.String(): WeaponTypeShield,
}

var weaponTypeSlotMap = map[WeaponType]WeaponSlot{
	WeaponTypeEmpty:  WeaponSlotAny,
	WeaponTypeKnife:  WeaponSlotMajorHand,
	WeaponTypeDagger: WeaponSlotAny,
	WeaponTypeBow:    WeaponSlotMajorHand,
	WeaponTypeAxe:    WeaponSlotBothHand,
	WeaponTypeRod:    WeaponSlotAny,
	WeaponTypeWand:   WeaponSlotBothHand,
	WeaponTypeOrb:    WeaponSlotAny,
	WeaponTypeBook:   WeaponSlotAny,
	WeaponTypeArrow:  WeaponSlotMinorHand,
	WeaponTypeShield: WeaponSlotMinorHand,
}

var weaponTypePairMap = map[string]bool{
	getWeaponTypePairKey(WeaponTypeKnife, WeaponTypeShield):  true,
	getWeaponTypePairKey(WeaponTypeDagger, WeaponTypeDagger): true,
	getWeaponTypePairKey(WeaponTypeDagger, WeaponTypeShield): true,
	getWeaponTypePairKey(WeaponTypeBow, WeaponTypeArrow):     true,
	getWeaponTypePairKey(WeaponTypeRod, WeaponTypeOrb):       true,
	getWeaponTypePairKey(WeaponTypeRod, WeaponTypeBook):      true,
	getWeaponTypePairKey(WeaponTypeOrb, WeaponTypeRod):       true,
	getWeaponTypePairKey(WeaponTypeOrb, WeaponTypeBook):      true,
	getWeaponTypePairKey(WeaponTypeBook, WeaponTypeRod):      true,
	getWeaponTypePairKey(WeaponTypeBook, WeaponTypeOrb):      true,
}

func getWeaponTypePairKey(major, minor WeaponType) string {
	return fmt.Sprintf("%s:%s", major.String(), minor.String())
}
