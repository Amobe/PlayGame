package vo_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"

	"github.com/stretchr/testify/assert"
)

func TestEquipMajorWeapon(t *testing.T) {
	e := vo.NewEquipment()
	w := vo.WeaponSet["Knife"]

	want := []vo.Attribute{
		vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(10)),
	}
	got := e.EquipWeapon(w).GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipMinorWeapon(t *testing.T) {
	e := vo.NewEquipment()
	w := vo.WeaponSet["Shield"]

	want := []vo.Attribute{
		vo.NewAttribute(vo.AttributeTypeDEF, decimal.NewFromInt(10)),
	}
	got := e.EquipWeapon(w).GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipBothHandWeapon(t *testing.T) {
	e := vo.NewEquipment()
	w := vo.WeaponSet["Axe"]

	want := []vo.Attribute{
		vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(25)),
	}
	got := e.EquipWeapon(w).GetAttributes()

	assert.ElementsMatch(t, want, got)
}

func TestEquipBothHandWeaponRemoveMinorWeapon(t *testing.T) {
	e := vo.NewEquipment()
	res := e.EquipWeapon(vo.WeaponSet["Knife"]).
		EquipWeapon(vo.WeaponSet["Shield"]).
		EquipWeapon(vo.WeaponSet["Axe"])

	want := []vo.Attribute{
		vo.NewAttribute(vo.AttributeTypeATK, decimal.NewFromInt(25)),
	}
	got := res.GetAttributes()

	assert.ElementsMatch(t, want, got)
}
