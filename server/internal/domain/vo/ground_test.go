package vo

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_Ground_GetMatchOrder(t *testing.T) {
	type fields struct {
		AllyAGI  Attribute
		EnemyAGI Attribute
	}
	tests := []struct {
		name    string
		fields  fields
		want    []GroundIdx
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ally summoner is faster",
			fields: fields{
				AllyAGI:  NewAttribute(AttributeTypeAGI, decimal.NewFromFloat(100)),
				EnemyAGI: NewAttribute(AttributeTypeAGI, decimal.NewFromFloat(50)),
			},
			want:    []GroundIdx{1, 7, 2, 8, 3, 9, 4, 10, 5, 11},
			wantErr: assert.NoError,
		},
		{
			name: "enemy summoner is faster",
			fields: fields{
				AllyAGI:  NewAttribute(AttributeTypeAGI, decimal.NewFromFloat(50)),
				EnemyAGI: NewAttribute(AttributeTypeAGI, decimal.NewFromFloat(100)),
			},
			want:    []GroundIdx{7, 1, 8, 2, 9, 3, 10, 4, 11, 5},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allyCamp, err := NewCamp(
				NewCharacter(1, NewAttributeMap()),
				NewCharacter(2, NewAttributeMap()),
				NewCharacter(3, NewAttributeMap()),
				NewCharacter(4, NewAttributeMap()),
				NewCharacter(5, NewAttributeMap()),
				NewCharacter(6, NewAttributeMap(tt.fields.AllyAGI)),
			)
			assert.NoError(t, err)
			enemyCamp, err := NewCamp(
				NewCharacter(1, NewAttributeMap()),
				NewCharacter(2, NewAttributeMap()),
				NewCharacter(3, NewAttributeMap()),
				NewCharacter(4, NewAttributeMap()),
				NewCharacter(5, NewAttributeMap()),
				NewCharacter(6, NewAttributeMap(tt.fields.EnemyAGI)),
			)
			assert.NoError(t, err)
			g := NewGround(allyCamp, enemyCamp)
			got, err := g.GetMatchOrder()
			if !tt.wantErr(t, err, "GetMatchOrder()") {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMatchOrder()")
		})
	}
}

func Test_Ground_GetMatch(t *testing.T) {
	oneTargetSkill := NewSkill("test", NewAttributeMap(NewAttribute(AttributeTypeTarget, decimal.NewFromInt(1))))

	type fields struct {
		allyCamp  Camp
		enemyCamp Camp
	}
	type args struct {
		idx GroundIdx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Match
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get match with ally ground idx",
			fields: fields{
				allyCamp: Camp{
					NewCharacterWithSkill(1, oneTargetSkill, NewAttributeMap()),
					NewCharacter(2, NewAttributeMap()),
					NewCharacter(3, NewAttributeMap()),
					NewCharacter(4, NewAttributeMap()),
					NewCharacter(5, NewAttributeMap()),
					NewCharacter(6, NewAttributeMap()),
				},
				enemyCamp: Camp{
					NewCharacter(7, NewAttributeMap()),
					NewCharacter(8, NewAttributeMap()),
					NewCharacter(9, NewAttributeMap()),
					NewCharacter(10, NewAttributeMap()),
					NewCharacter(11, NewAttributeMap()),
					NewCharacter(12, NewAttributeMap()),
				},
			},
			args: args{
				idx: 1,
			},
			want: Match{
				Actor:   NewCharacterWithSkill(1, oneTargetSkill, NewAttributeMap()),
				Targets: []Character{NewCharacter(7, NewAttributeMap())},
			},
			wantErr: assert.NoError,
		},
		{
			name: "get match with enemy ground idx",
			fields: fields{
				allyCamp: Camp{
					NewCharacter(1, NewAttributeMap()),
					NewCharacter(2, NewAttributeMap()),
					NewCharacter(3, NewAttributeMap()),
					NewCharacter(4, NewAttributeMap()),
					NewCharacter(5, NewAttributeMap()),
					NewCharacter(6, NewAttributeMap()),
				},
				enemyCamp: Camp{
					NewCharacterWithSkill(7, oneTargetSkill, NewAttributeMap()),
					NewCharacter(8, NewAttributeMap()),
					NewCharacter(9, NewAttributeMap()),
					NewCharacter(10, NewAttributeMap()),
					NewCharacter(11, NewAttributeMap()),
					NewCharacter(12, NewAttributeMap()),
				},
			},
			args: args{
				idx: 7,
			},
			want: Match{
				Actor:   NewCharacterWithSkill(7, oneTargetSkill, NewAttributeMap()),
				Targets: []Character{NewCharacter(1, NewAttributeMap())},
			},
			wantErr: assert.NoError,
		},
		{
			name: "get match when attacker is dead",
			fields: fields{
				allyCamp: Camp{
					NewCharacterWithSkill(1, oneTargetSkill, NewAttributeMap(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))),
					NewCharacter(2, NewAttributeMap()),
					NewCharacter(3, NewAttributeMap()),
					NewCharacter(4, NewAttributeMap()),
					NewCharacter(5, NewAttributeMap()),
					NewCharacter(6, NewAttributeMap()),
				},
				enemyCamp: Camp{
					NewCharacter(7, NewAttributeMap()),
					NewCharacter(8, NewAttributeMap()),
					NewCharacter(9, NewAttributeMap()),
					NewCharacter(10, NewAttributeMap()),
					NewCharacter(11, NewAttributeMap()),
					NewCharacter(12, NewAttributeMap()),
				},
			},
			args: args{
				idx: 1,
			},
			want: Match{
				Actor: NewCharacterWithSkill(1, oneTargetSkill, NewAttributeMap(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGround(tt.fields.allyCamp, tt.fields.enemyCamp)
			got, err := g.GetMatch(tt.args.idx)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMatch(%v)", tt.args.idx)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMatch()")
		})
	}
}

func Test_getTargetsFn(t *testing.T) {
	type fields struct {
		characters []Character
	}
	type args struct {
		number int
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantTargets []Character
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "get 1 target",
			fields: fields{
				characters: []Character{
					NewCharacter(1, NewAttributeMap()),
					NewCharacter(2, NewAttributeMap()),
					NewCharacter(3, NewAttributeMap()),
					NewCharacter(4, NewAttributeMap()),
					NewCharacter(5, NewAttributeMap()),
					NewCharacter(6, NewAttributeMap()),
				},
			},
			args: args{
				number: 1,
			},
			wantTargets: []Character{NewCharacter(1, NewAttributeMap())},
			wantErr:     assert.NoError,
		},
		{
			name: "get 3 target",
			fields: fields{
				characters: []Character{
					NewCharacter(1, NewAttributeMap()),
					NewCharacter(2, NewAttributeMap()),
					NewCharacter(3, NewAttributeMap()),
					NewCharacter(4, NewAttributeMap()),
					NewCharacter(5, NewAttributeMap()),
					NewCharacter(6, NewAttributeMap()),
				},
			},
			args: args{
				number: 3,
			},
			wantTargets: []Character{
				NewCharacter(1, NewAttributeMap()),
				NewCharacter(2, NewAttributeMap()),
				NewCharacter(3, NewAttributeMap()),
			},
			wantErr: assert.NoError,
		},
		{
			name: "get summoner when target dead",
			fields: fields{
				characters: []Character{
					NewCharacter(1, NewAttributeMap(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))),
					NewCharacter(2, NewAttributeMap()),
					NewCharacter(3, NewAttributeMap()),
					NewCharacter(4, NewAttributeMap()),
					NewCharacter(5, NewAttributeMap()),
					NewCharacter(6, NewAttributeMap()),
				},
			},
			args: args{
				number: 1,
			},
			wantTargets: []Character{NewCharacter(6, NewAttributeMap())},
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			camp, err := NewCamp(tt.fields.characters...)
			if !assert.NoErrorf(t, err, "NewCamp()") {
				return
			}
			targetsFn := getTargetsFn(camp)
			got, err := targetsFn(tt.args.number)
			if !tt.wantErr(t, err, fmt.Sprintf("getTargetsFn(%v)", tt.args.number)) {
				return
			}
			assert.Equalf(t, tt.wantTargets, got, "getTargetsFn(%v)", tt.args.number)
		})
	}
}

func Test_Ground_UpdateCharacter(t *testing.T) {
	originGround := NewGround(
		Camp{
			NewCharacter(1, NewAttributeMap()),
			NewCharacter(2, NewAttributeMap()),
			NewCharacter(3, NewAttributeMap()),
			NewCharacter(4, NewAttributeMap()),
			NewCharacter(5, NewAttributeMap()),
			NewCharacter(6, NewAttributeMap()),
		},
		Camp{
			NewCharacter(7, NewAttributeMap()),
			NewCharacter(8, NewAttributeMap()),
			NewCharacter(9, NewAttributeMap()),
			NewCharacter(10, NewAttributeMap()),
			NewCharacter(11, NewAttributeMap()),
			NewCharacter(12, NewAttributeMap()),
		},
	)

	type args struct {
		idx       GroundIdx
		character Character
	}
	tests := []struct {
		name       string
		args       args
		wantGround Ground
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "update ally character",
			args: args{
				idx:       1,
				character: NewCharacter(1, NewAttributeMap(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))),
			},
			wantGround: NewGround(
				Camp{
					NewCharacter(1, NewAttributeMap(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))),
					NewCharacter(2, NewAttributeMap()),
					NewCharacter(3, NewAttributeMap()),
					NewCharacter(4, NewAttributeMap()),
					NewCharacter(5, NewAttributeMap()),
					NewCharacter(6, NewAttributeMap()),
				},
				Camp{
					NewCharacter(7, NewAttributeMap()),
					NewCharacter(8, NewAttributeMap()),
					NewCharacter(9, NewAttributeMap()),
					NewCharacter(10, NewAttributeMap()),
					NewCharacter(11, NewAttributeMap()),
					NewCharacter(12, NewAttributeMap()),
				},
			),
			wantErr: assert.NoError,
		},
		{
			name: "update enemy character",
			args: args{
				idx:       7,
				character: NewCharacter(7, NewAttributeMap(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))),
			},
			wantGround: NewGround(
				Camp{
					NewCharacter(1, NewAttributeMap()),
					NewCharacter(2, NewAttributeMap()),
					NewCharacter(3, NewAttributeMap()),
					NewCharacter(4, NewAttributeMap()),
					NewCharacter(5, NewAttributeMap()),
					NewCharacter(6, NewAttributeMap()),
				},
				Camp{
					NewCharacter(7, NewAttributeMap(NewAttribute(AttributeTypeDead, decimal.NewFromInt(1)))),
					NewCharacter(8, NewAttributeMap()),
					NewCharacter(9, NewAttributeMap()),
					NewCharacter(10, NewAttributeMap()),
					NewCharacter(11, NewAttributeMap()),
					NewCharacter(12, NewAttributeMap()),
				},
			),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := originGround
			got, err := g.UpdateCharacter(tt.args.idx, tt.args.character)
			if !tt.wantErr(t, err, fmt.Sprintf("UpdateCharacter(%v)", tt.args.idx)) {
				return
			}
			assert.Equalf(t, tt.wantGround, got, "UpdateCharacter()")
		})
	}
}
