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

func (a AttributeType) String() string {
	return string(a)
}

func ToAttributeType(s string) AttributeType {
	t, ok := attributeTypeMap[s]
	if !ok {
		return AttributeTypeUnspecified
	}
	return t
}

const (
	AttributeTypeUnspecified AttributeType = "unspecified"
	// basement attributes
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
	AttributeTypeHR       AttributeType = "hr"  // heal rate
	AttributeTypeHit      AttributeType = "hit"
	AttributeTypeDodge    AttributeType = "dodge"
	AttributeTypeDamage   AttributeType = "damage"
	// skill attributes
	AttributeTypeTarget AttributeType = "target"
	AttributeTypeCD     AttributeType = "cd"
	AttributeTypeFCD    AttributeType = "firstcd"
	AttributeTypeATKB   AttributeType = "atkbonus"
	AttributeTypeMATKB  AttributeType = "matkbonus"

	AttributeTypeDead     AttributeType = "dead"
	AttributeTypePoisoned AttributeType = "poisoned"
)

var attributeTypeMap = map[string]AttributeType{
	AttributeTypeHP.String():       AttributeTypeHP,
	AttributeTypeAGI.String():      AttributeTypeAGI,
	AttributeTypeATK.String():      AttributeTypeATK,
	AttributeTypeDEF.String():      AttributeTypeDEF,
	AttributeTypeMATK.String():     AttributeTypeMATK,
	AttributeTypeMDEF.String():     AttributeTypeMDEF,
	AttributeTypeCRI.String():      AttributeTypeCRI,
	AttributeTypeCRIR.String():     AttributeTypeCRIR,
	AttributeTypeCRID.String():     AttributeTypeCRID,
	AttributeTypeCRIDR.String():    AttributeTypeCRIDR,
	AttributeTypeAMP.String():      AttributeTypeAMP,
	AttributeTypeAMPR.String():     AttributeTypeAMPR,
	AttributeTypeStatusH.String():  AttributeTypeStatusH,
	AttributeTypeStatusHR.String(): AttributeTypeStatusHR,
	AttributeTypeDI.String():       AttributeTypeDI,
	AttributeTypeDR.String():       AttributeTypeDR,
	AttributeTypeSDR.String():      AttributeTypeSDR,
	AttributeTypeHR.String():       AttributeTypeHR,
	AttributeTypeHit.String():      AttributeTypeHit,
	AttributeTypeDodge.String():    AttributeTypeDodge,
	AttributeTypeTarget.String():   AttributeTypeTarget,
	AttributeTypeCD.String():       AttributeTypeCD,
	AttributeTypeFCD.String():      AttributeTypeFCD,
	AttributeTypeATKB.String():     AttributeTypeATKB,
	AttributeTypeMATKB.String():    AttributeTypeMATKB,
	AttributeTypeDead.String():     AttributeTypeDead,
	AttributeTypePoisoned.String(): AttributeTypePoisoned,
}

var defaultAttributeValue = map[AttributeType]decimal.Decimal{}

func GetDefaultAttributeValue(attributeType AttributeType) decimal.Decimal {
	value, ok := defaultAttributeValue[attributeType]
	if !ok {
		return decimal.Zero
	}
	return value
}
