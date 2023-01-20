package character_test

import (
	"testing"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/Amobe/PlayGame/server/pkg/domain/valueobject"

	"github.com/stretchr/testify/assert"
)

func TestEquipMajorWeapon(t *testing.T) {
	e := character.NewEquipment()
	w := character.WeaponSet["Knife"]
	e.EquipWeapon(w)

	want := []valueobject.Attribute{
		{valueobject.AttributeTypeATK, "10"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipMinorWeapon(t *testing.T) {
	e := character.NewEquipment()
	w := character.WeaponSet["Shield"]
	e.EquipWeapon(w)

	want := []valueobject.Attribute{
		{valueobject.AttributeTypeDEF, "10"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipBothHandWeapon(t *testing.T) {
	e := character.NewEquipment()
	w := character.WeaponSet["Axe"]
	e.EquipWeapon(w)

	want := []valueobject.Attribute{
		{valueobject.AttributeTypeATK, "25"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipBothHandWeaponRemoveMinorWeapon(t *testing.T) {
	e := character.NewEquipment()
	e.EquipWeapon(character.WeaponSet["Knife"])
	e.EquipWeapon(character.WeaponSet["Shield"])
	e.EquipWeapon(character.WeaponSet["Axe"])

	want := []valueobject.Attribute{
		{valueobject.AttributeTypeATK, "25"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}
