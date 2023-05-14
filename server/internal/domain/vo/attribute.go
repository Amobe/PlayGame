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

func NewBoolAttribute(attrType AttributeType, value bool) Attribute {
	if value {
		return NewAttribute(attrType, decimal.NewFromInt(1))
	}
	return NewAttribute(attrType, decimal.NewFromInt(0))
}

func (a Attribute) Add(value decimal.Decimal) Attribute {
	return NewAttribute(a.Type, a.Value.Add(value))
}

// EqualTo compares two Attribute objects.
func (a Attribute) EqualTo(other Attribute) bool {
	if a.Type != other.Type {
		return false
	}
	if a.Value.Equal(other.Value) {
		return false
	}
	return true
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
	AttributeTypeTarget         AttributeType = "target"
	AttributeTypeCD             AttributeType = "cd"
	AttributeTypeFCD            AttributeType = "firstcd"
	AttributeTypeATKB           AttributeType = "atkbonus"
	AttributeTypeMATKB          AttributeType = "matkbonus"
	AttributeTypePhysicalDamage AttributeType = "pd"
	AttributeTypeMagicalDamage  AttributeType = "md"
	AttributeTypeHeal           AttributeType = "heal"

	AttributeTypeDead     AttributeType = "dead"
	AttributeTypePoisoned AttributeType = "poisoned"
)

var attributeTypeMap = map[string]AttributeType{
	AttributeTypeHP.String():             AttributeTypeHP,
	AttributeTypeAGI.String():            AttributeTypeAGI,
	AttributeTypeATK.String():            AttributeTypeATK,
	AttributeTypeDEF.String():            AttributeTypeDEF,
	AttributeTypeMATK.String():           AttributeTypeMATK,
	AttributeTypeMDEF.String():           AttributeTypeMDEF,
	AttributeTypeCRI.String():            AttributeTypeCRI,
	AttributeTypeCRIR.String():           AttributeTypeCRIR,
	AttributeTypeCRID.String():           AttributeTypeCRID,
	AttributeTypeCRIDR.String():          AttributeTypeCRIDR,
	AttributeTypeAMP.String():            AttributeTypeAMP,
	AttributeTypeAMPR.String():           AttributeTypeAMPR,
	AttributeTypeStatusH.String():        AttributeTypeStatusH,
	AttributeTypeStatusHR.String():       AttributeTypeStatusHR,
	AttributeTypeDI.String():             AttributeTypeDI,
	AttributeTypeDR.String():             AttributeTypeDR,
	AttributeTypeSDR.String():            AttributeTypeSDR,
	AttributeTypeHR.String():             AttributeTypeHR,
	AttributeTypeHit.String():            AttributeTypeHit,
	AttributeTypeDodge.String():          AttributeTypeDodge,
	AttributeTypeTarget.String():         AttributeTypeTarget,
	AttributeTypeCD.String():             AttributeTypeCD,
	AttributeTypeFCD.String():            AttributeTypeFCD,
	AttributeTypeATKB.String():           AttributeTypeATKB,
	AttributeTypeMATKB.String():          AttributeTypeMATKB,
	AttributeTypePhysicalDamage.String(): AttributeTypePhysicalDamage,
	AttributeTypeMagicalDamage.String():  AttributeTypeMagicalDamage,
	AttributeTypeHeal.String():           AttributeTypeHeal,
	AttributeTypeDead.String():           AttributeTypeDead,
	AttributeTypePoisoned.String():       AttributeTypePoisoned,
}

var defaultAttributeValue = map[AttributeType]decimal.Decimal{}

func GetDefaultAttributeValue(attributeType AttributeType) decimal.Decimal {
	value, ok := defaultAttributeValue[attributeType]
	if !ok {
		return decimal.Zero
	}
	return value
}

var (
	DeadAttribute = NewAttribute(AttributeTypeDead, decimal.NewFromInt(1))
)
