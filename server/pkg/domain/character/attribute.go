package character

import (
	"strconv"
)

type Attribute struct {
	Type  AttributeType
	Value string
}

func (a Attribute) GetInt() int {
	value, err := strconv.ParseInt(a.Value, 10, 64)
	if err != nil {
		return 0
	}
	return int(value)
}

func (a Attribute) GetFloat() float64 {
	value, err := strconv.ParseFloat(a.Value, 64)
	if err != nil {
		return 0
	}
	return value
}

func (a *Attribute) Add(value int) {
	res := a.GetInt() + value
	a.Value = strconv.Itoa(res)
}

type AttributeType string

const (
	AttributeTypeHP       AttributeType = "hp"
	AttributeTypeAGI      AttributeType = "agi"
	AttributeTypeATK      AttributeType = "atk"
	AttributeTypeDEF      AttributeType = "def"
	AttributeTypeMATK     AttributeType = "matk"
	AttributeTypeMDEF     AttributeType = "mdef"
	AttributeTypeCRI      AttributeType = "cri"
	AttributeTypeCRIR     AttributeType = "crir"
	AttributeTypeCRID     AttributeType = "crid"
	AttributeTypeCRIDR    AttributeType = "cridr"
	AttributeTypeAMP      AttributeType = "amp"
	AttributeTypeAMPR     AttributeType = "ampr"
	AttributeTypeStatusH  AttributeType = "sh"  // status hit rate
	AttributeTypeStatusHR AttributeType = "shr" // status hit resist
	AttributeTypeDI       AttributeType = "di"  // damage increase
	AttributeTypeDR       AttributeType = "dr"  // damage reduce
	AttributeTypeSDR      AttributeType = "sdr" // skill damage rate

	AttributeTypeDead     AttributeType = "dead"
	AttributeTypePoisoned AttributeType = "poisoned"
)

type AttributeTypeMap map[AttributeType]Attribute

func NewAttributeTypeMap() AttributeTypeMap {
	return make(AttributeTypeMap)
}

func (a AttributeTypeMap) Insert(attr Attribute) {
	target, ok := a[attr.Type]
	if !ok {
		a[attr.Type] = attr
		return
	}
	target.Add(attr.GetInt())
	a[attr.Type] = target
}
