package vo

import "github.com/Amobe/PlayGame/server/internal/utils"

type AttributeMap map[AttributeType]Attribute

func newAttributeMap() AttributeMap {
	return make(AttributeMap)
}

func NewAttributeMap(attrs ...Attribute) AttributeMap {
	m := newAttributeMap()
	if len(attrs) > 0 {
		m.Insert(attrs...)
	}
	return m
}

// Merge merges two AttributeMap
func (a AttributeMap) Merge(b AttributeMap) AttributeMap {
	return MergeAttributeMap(a, b)
}

func (a AttributeMap) Insert(attrs ...Attribute) AttributeMap {
	for _, attr := range attrs {
		a.insert(attr)
	}

	newMap := newAttributeMap()
	utils.CopyMap(newMap, a)
	return newMap
}

func (a AttributeMap) insert(attr Attribute) {
	target, ok := a[attr.Type]
	if !ok {
		a[attr.Type] = attr
		return
	}
	a[attr.Type] = target.Add(attr.Value)
}

func MergeAttributeMap(maps ...AttributeMap) AttributeMap {
	newMap := newAttributeMap()
	for _, m := range maps {
		for _, v := range m {
			newMap.insert(v)
		}
	}
	return newMap
}

func (a AttributeMap) Get(attributeType AttributeType) Attribute {
	attr, ok := a[attributeType]
	if !ok {
		return NewAttribute(attributeType, GetDefaultAttributeValue(attributeType))
	}
	return attr
}

// EqualTo compares two AttributeMap
func (a AttributeMap) EqualTo(b AttributeMap) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if !v.EqualTo(b[k]) {
			return false
		}
	}
	return true
}
