package vo

import "github.com/shopspring/decimal"

type Attribute struct {
	Type  AttributeType
	Value decimal.Decimal
}

func NewAttribute(attrType AttributeType, value decimal.Decimal) Attribute {
	return Attribute{
		Type:  attrType,
		Value: value,
	}
}

func (a Attribute) Add(value decimal.Decimal) Attribute {
	return NewAttribute(a.Type, a.Value.Add(value))
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

var defaultAttributeValue = map[AttributeType]decimal.Decimal{}

func GetDefaultAttributeValue(attributeType AttributeType) decimal.Decimal {
	value, ok := defaultAttributeValue[attributeType]
	if !ok {
		return decimal.Zero
	}
	return value
}
