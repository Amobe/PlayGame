package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
	"github.com/Amobe/PlayGame/server/pkg/infra/inmem"
)

const weaponJsonPath = "files/weapon.json"

type weaponConfigData struct {
	Weapons []weaponData `json:"weapons"`
}

type weaponData struct {
	Id         string         `json:"id"`
	Type       string         `json:"type"`
	Name       string         `json:"name"`
	Slot       string         `json:"slot"`
	Attributes map[string]any `json:"attributes"`
}

func initWeaponRepository() *inmem.WeaponRepository {
	r := inmem.NewInmemWeaponRepository()

	d := weaponConfigData{}
	weaponConfigFile, err := configFiles.ReadFile(weaponJsonPath)
	if err != nil {
		panic(fmt.Errorf("read %s from config files: %w", weaponJsonPath, err))
	}
	if err := json.Unmarshal(weaponConfigFile, &d); err != nil {
		panic(fmt.Errorf("json unmarshal weapon config file: %w", err))
	}

	for _, data := range d.Weapons {
		job, err := weaponDataToDomain(data)
		if err != nil {
			panic(fmt.Errorf("weapon data to domain: %w", err))
		}
		if _, err := r.Get(job.ID()); !errors.Is(err, inmem.ErrorRecordNotFound) {
			panic(fmt.Sprintf("weapon id is duplicated: %s", job.ID()))
		}
		_ = r.Create(job)
	}

	return r
}

func weaponDataToDomain(data weaponData) (vo.Weapon, error) {
	var attributes []vo.Attribute
	for k, v := range data.Attributes {
		switch a := v.(type) {
		case int:
			attributes = append(attributes, vo.NewAttribute(vo.AttributeType(k), decimal.NewFromInt(int64(a))))
		case float64:
			attributes = append(attributes, vo.NewAttribute(vo.AttributeType(k), decimal.NewFromFloat(a)))
		case string:
			// TODO: add string attribute
			continue
		default:
			panic(fmt.Sprintf("unknown type: %T", a))
		}
	}
	return vo.NewWeapon(data.Id, data.Type, data.Name, data.Slot, attributes...)

}
