package calculator

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCalculateHitRate(t *testing.T) {
	type args struct {
		hit   decimal.Decimal
		dodge decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "hit and dodge are equal",
			args: args{
				hit:   decimal.NewFromFloat(10),
				dodge: decimal.NewFromFloat(10),
			},
			want: decimal.NewFromFloat(0.95),
		},
		{
			name: "example 1",
			args: args{
				hit:   decimal.NewFromFloat(10),
				dodge: decimal.NewFromFloat(50),
			},
			want: decimal.NewFromFloat(0.3167),
		},
		{
			name: "example 2",
			args: args{
				hit:   decimal.NewFromFloat(400),
				dodge: decimal.NewFromFloat(50),
			},
			want: decimal.NewFromFloat(1),
		},
		{
			name: "example 3",
			args: args{
				hit:   decimal.NewFromFloat(10),
				dodge: decimal.NewFromFloat(500),
			},
			want: decimal.NewFromFloat(0.0373),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateHitRate(tt.args.hit, tt.args.dodge)
			assert.Truef(t, tt.want.Equal(got), "%s != %s", tt.want.String(), got.String())
		})
	}
}
