package vo

import "github.com/Amobe/PlayGame/server/pkg/utils"

type AttributeMap map[AttributeType]Attribute

func NewAttributeMap(attrs ...Attribute) AttributeMap {
	m := make(AttributeMap)
	if len(attrs) > 0 {
		m.Insert(attrs...)
	}
	return m
}

func (a AttributeMap) Insert(attrs ...Attribute) AttributeMap {
	for _, attr := range attrs {
		target, ok := a[attr.Type]
		if !ok {
			a[attr.Type] = attr
			continue
		}
		a[attr.Type] = target.Add(attr.Value)
	}

	newMap := NewAttributeMap()
	utils.CopyMap(newMap, a)
	return newMap
}

func (a AttributeMap) Get(attributeType AttributeType) Attribute {
	attr, ok := a[attributeType]
	if !ok {
		return NewAttribute(attributeType, GetDefaultAttributeValue(attributeType))
	}
	return attr
}
