package equipment_test

import (
	"testing"

	"github.com/Amobe/PlayGame/server/pkg/domain/equipment"
	"github.com/stretchr/testify/assert"
)

func TestEquipMajorWeapon(t *testing.T) {
	e := equipment.NewEquipment()
	w := equipment.WeaponSet["Knife"]
	e.EquipWeapon(w)

	want := []equipment.Attribute{
		{equipment.AttributeValueTypeATK, "10"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipMinorWeapon(t *testing.T) {
	e := equipment.NewEquipment()
	w := equipment.WeaponSet["Shield"]
	e.EquipWeapon(w)

	want := []equipment.Attribute{
		{equipment.AttributeValueTypeDEF, "10"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipBothHandWeapon(t *testing.T) {
	e := equipment.NewEquipment()
	w := equipment.WeaponSet["Axe"]
	e.EquipWeapon(w)

	want := []equipment.Attribute{
		{equipment.AttributeValueTypeATK, "25"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipBothHandWeaponRemoveMinorWeapon(t *testing.T) {
	e := equipment.NewEquipment()
	e.EquipWeapon(equipment.WeaponSet["Knife"])
	e.EquipWeapon(equipment.WeaponSet["Shield"])
	e.EquipWeapon(equipment.WeaponSet["Axe"])

	want := []equipment.Attribute{
		{equipment.AttributeValueTypeATK, "25"},
	}
	got := e.GetAttributes()

	assert.ElementsMatch(t, want, got)
}
