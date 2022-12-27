package character_test

import (
	"testing"

	"github.com/Amobe/PlayGame/server/pkg/domain/character"
	"github.com/stretchr/testify/assert"
)

func TestEquipMajorWeapon(t *testing.T) {
	e := character.NewEquipment()
	w := character.WeaponSet["Knife"]
	e.EquipWeapon(w)

	want := []character.Attribute{
		{character.AttributeTypeATK, "10"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipMinorWeapon(t *testing.T) {
	e := character.NewEquipment()
	w := character.WeaponSet["Shield"]
	e.EquipWeapon(w)

	want := []character.Attribute{
		{character.AttributeTypeDEF, "10"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipBothHandWeapon(t *testing.T) {
	e := character.NewEquipment()
	w := character.WeaponSet["Axe"]
	e.EquipWeapon(w)

	want := []character.Attribute{
		{character.AttributeTypeATK, "25"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipBothHandWeaponRemoveMinorWeapon(t *testing.T) {
	e := character.NewEquipment()
	e.EquipWeapon(character.WeaponSet["Knife"])
	e.EquipWeapon(character.WeaponSet["Shield"])
	e.EquipWeapon(character.WeaponSet["Axe"])

	want := []character.Attribute{
		{character.AttributeTypeATK, "25"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}
