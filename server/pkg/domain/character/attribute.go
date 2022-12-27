package character

import "strconv"

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

func (a *Attribute) Add(value int) {
	res := a.GetInt() + value
	a.Value = strconv.Itoa(res)
}

type AttributeType string

const (
	AttributeTypeHP   AttributeType = "hp"
	AttributeTypeAGI  AttributeType = "agi"
	AttributeTypeATK  AttributeType = "atk"
	AttributeTypeDEF  AttributeType = "def"
	AttributeTypeMATK AttributeType = "matk"
	AttributeTypeMDEF AttributeType = "mdef"

	AttributeTypeDead AttributeType = "dead"
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
