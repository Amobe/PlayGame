package vo

type AttributeMap map[AttributeType]Attribute

func NewAttributeTypeMap(attrs ...Attribute) AttributeMap {
	m := make(AttributeMap)
	if len(attrs) > 0 {
		m.Insert(attrs...)
	}
	return m
}

func (a AttributeMap) Insert(attrs ...Attribute) {
	for _, attr := range attrs {
		target, ok := a[attr.Type]
		if !ok {
			a[attr.Type] = attr
			continue
		}
		a[attr.Type] = target.Add(attr.Value)
	}
}

func (a AttributeMap) Get(attributeType AttributeType) Attribute {
	attr, ok := a[attributeType]
	if !ok {
		return NewAttribute(attributeType, GetDefaultAttributeValue(attributeType))
	}
	return attr
}
