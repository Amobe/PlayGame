package config

import (
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

func weaponDataToDomain(data configData) (vo.Weapon, error) {
	var id, weaponType, weaponName, weaponSkill string
	var attributes []vo.Attribute
	for k, v := range data {
		switch {
		case k == "id":
			id = v.(string)
		case k == "type":
			weaponType = v.(string)
		case k == "name":
			weaponName = v.(string)
		case k == "skill":
			weaponSkill = v.(string)
		case k == "rarity":
		case vo.ToAttributeType(k) != vo.AttributeTypeUnspecified:
			at := vo.ToAttributeType(k)
			av := decimal.NewFromFloat(v.(float64))
			attributes = append(attributes, vo.NewAttribute(at, av))
		default:
			return vo.Weapon{}, fmt.Errorf("unrecognized weapon data field %s", k)
		}
	}
	return vo.NewWeapon(id, weaponType, weaponName, weaponSkill, attributes...)
}

func skillDataToDomain(data configData) (vo.Skill, error) {
	var id, name string
	var attributes []vo.Attribute
	for k, v := range data {
		switch {
		case k == "id":
			id = v.(string)
		case k == "name":
			name = v.(string)
		case vo.ToAttributeType(k) != vo.AttributeTypeUnspecified:
			at := vo.ToAttributeType(k)
			av := decimal.NewFromFloat(v.(float64))
			attributes = append(attributes, vo.NewAttribute(at, av))
		default:
			return vo.Skill{}, fmt.Errorf("unrecognized skill data field %s", k)
		}
	}
	return vo.NewSkill(id, name, attributes...), nil
}
