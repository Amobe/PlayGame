package vo_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/internal/domain/vo"
)

func TestCharacter_Affect_DeadCondition(t *testing.T) {
	type args struct {
		attr vo.AttributeMap
	}
	tests := []struct {
		name  string
		args  args
		wants []vo.Attribute
	}{
		{
			name: "affect attributes should be added to the new character instance",
			args: args{
				attr: vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10))),
			},
			wants: []vo.Attribute{
				vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(20)),
			},
		},
		{
			name: "dead should be added when character's HP equal to zero",
			args: args{
				attr: vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(-10))),
			},
			wants: []vo.Attribute{
				vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(0)),
				vo.NewAttribute(vo.AttributeTypeDead, decimal.NewFromInt(1)),
			},
		},
		{
			name: "dead should be added when character's HP less than zero",
			args: args{
				attr: vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(-20))),
			},
			wants: []vo.Attribute{
				vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(-10)),
				vo.NewAttribute(vo.AttributeTypeDead, decimal.NewFromInt(1)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attributes := vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeHP, decimal.NewFromInt(10)))
			c := vo.NewCharacter("", attributes)
			got := c.TakeAffect(tt.args.attr).GetAttributeMap()
			for _, want := range tt.wants {
				actual := got.Get(want.Type).Value
				assert.Truef(t, want.Value.Equal(actual), "TakeAffect(%v), %s != %s",
					tt.args.attr, want.Value.String(), actual.String())
			}
		})
	}
}
