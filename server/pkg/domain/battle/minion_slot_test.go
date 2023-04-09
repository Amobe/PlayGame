package battle

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Amobe/PlayGame/server/pkg/domain/vo"
)

type mockUnit struct {
	*MockUnit
	changed bool
}

func newMockUnit() *mockUnit {
	return &mockUnit{MockUnit: &MockUnit{}}
}

func getMockUnitTakeAffect(groundIdx GroundIdx) *mockUnit {
	u := newMockUnit()
	u.On("GetGroundIdx").Return(groundIdx)
	u.On("TakeAffect", mock.Anything).Return(&mockUnit{MockUnit: &MockUnit{}, changed: true})
	return u
}

func getMockMinions() (ally *Minions, enemy *Minions) {
	return &Minions{
			getMockUnitTakeAffect(1),
			getMockUnitTakeAffect(2),
			getMockUnitTakeAffect(3),
			getMockUnitTakeAffect(4),
			getMockUnitTakeAffect(5),
			getMockUnitTakeAffect(6),
		}, &Minions{
			getMockUnitTakeAffect(7),
			getMockUnitTakeAffect(8),
			getMockUnitTakeAffect(9),
			getMockUnitTakeAffect(10),
			getMockUnitTakeAffect(11),
			getMockUnitTakeAffect(12),
		}
}

func Test_MinionSlot_unitTakeAffect(t *testing.T) {
	type args struct {
		groundIdx GroundIdx
		unit      Unit
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "groundIdx is enemy, the enemy should be changed",
			args: args{
				groundIdx: GroundIdx(8),
				unit:      getMockUnitTakeAffect(8),
			},
		},
		{
			name: "groundIdx is ally, the ally should be changed",
			args: args{
				groundIdx: GroundIdx(4),
				unit:      getMockUnitTakeAffect(4),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allyMinions, enemyMinions := getMockMinions()
			s := &MinionSlot{
				AllyMinions:  allyMinions,
				EnemyMinions: enemyMinions,
			}
			s.unitTakeAffect(tt.args.unit, nil)
			// assert the specific minion changed
			assert.Truef(t, s.getUnit(tt.args.groundIdx).(*mockUnit).changed,
				"unit changed expect true=%t", s.getUnit(tt.args.groundIdx).(*mockUnit).changed)
		})
	}
}

func getMockUnitTakeAttack(groundIdx GroundIdx) *mockUnit {
	u := newMockUnit()
	u.On("GetGroundIdx").Return(groundIdx)
	u.On("GetSkill").Return(vo.Skill{Name: "attack"})
	u.On("TakeAffect", mock.Anything).Return(&mockUnit{MockUnit: &MockUnit{}, changed: true})
	return u
}

func getMissCalculateAttackDamageFn() calculateAttackDamageFn {
	return func(attacker, target Unit, skill vo.Skill) (decimal.Decimal, bool) {
		return decimal.Zero, false
	}
}

func getHitCalculateAttackDamageFn() calculateAttackDamageFn {
	return func(attacker, target Unit, skill vo.Skill) (decimal.Decimal, bool) {
		return decimal.NewFromInt(1), true
	}
}

func Test_MinionSlot_attack(t *testing.T) {
	type fields struct {
		calculateAttackDamageFn calculateAttackDamageFn
	}
	type args struct {
		attacker Unit
		target   Unit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Affect
	}{
		{
			name: "get miss affect",
			fields: fields{
				calculateAttackDamageFn: getMissCalculateAttackDamageFn(),
			},
			args: args{
				attacker: getMockUnitTakeAttack(1),
				target:   getMockUnitTakeAttack(2),
			},
			want: Affect{
				ActorIdx:  1,
				TargetIdx: 2,
				Skill:     "miss",
			},
		},
		{
			name: "get hit affect",
			fields: fields{
				calculateAttackDamageFn: getHitCalculateAttackDamageFn(),
			},
			args: args{
				attacker: getMockUnitTakeAttack(1),
				target:   getMockUnitTakeAttack(2),
			},
			want: Affect{
				ActorIdx:  1,
				TargetIdx: 2,
				Skill:     "attack",
				Attributes: []vo.Attribute{
					{Type: vo.AttributeTypeDamage, Value: decimal.NewFromInt(1)},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allyMinions, enemyMinions := getMockMinions()
			s := &MinionSlot{
				AllyMinions:             allyMinions,
				EnemyMinions:            enemyMinions,
				calculateAttackDamageFn: tt.fields.calculateAttackDamageFn,
			}
			assert.Equalf(t, tt.want, s.attack(tt.args.attacker, tt.args.target), "attack(%v, %v)", tt.args.attacker, tt.args.target)
		})
	}
}

func getMockUnitGetSkill(groundIdx GroundIdx, isDead bool) *mockUnit {
	u := newMockUnit()
	u.On("GetGroundIdx").Return(groundIdx)
	u.On("GetSkill").Return(vo.Skill{
		Name: "test",
		AttributeMap: vo.AttributeMap{
			vo.AttributeTypeTarget: vo.Attribute{
				Type:  vo.AttributeTypeTarget,
				Value: decimal.NewFromInt(2),
			},
		},
	})
	u.On("IsDead").Return(isDead)
	return u
}

func Test_MinionSlot_getAttackerAndTargets(t *testing.T) {
	type fields struct {
		isDead bool
	}
	type args struct {
		idx GroundIdx
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantAttackerGroundIdx GroundIdx
		wantTargetsGroundIdx  []GroundIdx
	}{
		{
			name: "attacker is ally, targets are enemy and not dead",
			fields: fields{
				isDead: false,
			},
			args: args{
				idx: GroundIdx(2),
			},
			wantAttackerGroundIdx: GroundIdx(2),
			wantTargetsGroundIdx:  []GroundIdx{7, 8},
		},
		{
			name: " attacker is ally, targets are enemy and dead, targets will be summoner",
			fields: fields{
				isDead: true,
			},
			args: args{
				idx: GroundIdx(2),
			},
			wantAttackerGroundIdx: GroundIdx(2),
			wantTargetsGroundIdx:  []GroundIdx{12, 12},
		},
		{
			name: "attacker is enemy, targets are ally and not dead",
			fields: fields{
				isDead: false,
			},
			args: args{
				idx: GroundIdx(9),
			},
			wantAttackerGroundIdx: GroundIdx(9),
			wantTargetsGroundIdx:  []GroundIdx{1, 2},
		},
		{
			name: "attacker is enemy, targets are ally and dead, targets will be summoner",
			fields: fields{
				isDead: true,
			},
			args: args{
				idx: GroundIdx(9),
			},
			wantAttackerGroundIdx: GroundIdx(9),
			wantTargetsGroundIdx:  []GroundIdx{6, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MinionSlot{
				AllyMinions: &Minions{
					getMockUnitGetSkill(1, tt.fields.isDead),
					getMockUnitGetSkill(2, tt.fields.isDead),
					getMockUnitGetSkill(3, tt.fields.isDead),
					getMockUnitGetSkill(4, tt.fields.isDead),
					getMockUnitGetSkill(5, tt.fields.isDead),
					getMockUnitGetSkill(6, tt.fields.isDead),
				},
				EnemyMinions: &Minions{
					getMockUnitGetSkill(7, tt.fields.isDead),
					getMockUnitGetSkill(8, tt.fields.isDead),
					getMockUnitGetSkill(9, tt.fields.isDead),
					getMockUnitGetSkill(10, tt.fields.isDead),
					getMockUnitGetSkill(11, tt.fields.isDead),
					getMockUnitGetSkill(12, tt.fields.isDead),
				},
				targetPickerFn: defaultTargetPickerFn,
			}
			gotAttacker, gotTargets := s.getAttackerAndTargets(tt.args.idx)
			assert.Equalf(t, tt.wantAttackerGroundIdx, gotAttacker.GetGroundIdx(), "getAttackerAndTargets(%v)", tt.args.idx)
			for i, target := range gotTargets {
				assert.Equalf(t, tt.wantTargetsGroundIdx[i], target.GetGroundIdx(), "getAttackerAndTargets(%v)", tt.args.idx)
			}
		})
	}
}

func getMockUnitGetAgi(groundIdx GroundIdx, agi int) *mockUnit {
	u := newMockUnit()
	u.On("GetGroundIdx").Return(groundIdx)
	u.On("GetAgi").Return(agi)
	return u
}

func Test_MinionSlot_getActionOrder(t *testing.T) {
	type fields struct {
		AllyMinions  *Minions
		EnemyMinions *Minions
	}
	tests := []struct {
		name   string
		fields fields
		want   []GroundIdx
	}{
		{
			name: "enemy summoner is faster than ally summoner",
			fields: fields{
				AllyMinions: &Minions{
					getMockUnitGetAgi(1, 0),
					getMockUnitGetAgi(2, 0),
					getMockUnitGetAgi(3, 0),
					getMockUnitGetAgi(4, 0),
					getMockUnitGetAgi(5, 0),
					// Only check summoner agi
					getMockUnitGetAgi(6, 50),
				},
				EnemyMinions: &Minions{
					getMockUnitGetAgi(7, 0),
					getMockUnitGetAgi(8, 0),
					getMockUnitGetAgi(9, 0),
					getMockUnitGetAgi(10, 0),
					getMockUnitGetAgi(11, 0),
					// Only check summoner agi
					getMockUnitGetAgi(12, 100),
				},
			},
			want: []GroundIdx{7, 1, 8, 2, 9, 3, 10, 4, 11, 5},
		},
		{
			name: "ally summoner is faster than enemy summoner",
			fields: fields{
				AllyMinions: &Minions{
					getMockUnitGetAgi(1, 0),
					getMockUnitGetAgi(2, 0),
					getMockUnitGetAgi(3, 0),
					getMockUnitGetAgi(4, 0),
					getMockUnitGetAgi(5, 0),
					// Only check summoner agi
					getMockUnitGetAgi(6, 100),
				},
				EnemyMinions: &Minions{
					getMockUnitGetAgi(7, 0),
					getMockUnitGetAgi(8, 0),
					getMockUnitGetAgi(9, 0),
					getMockUnitGetAgi(10, 0),
					getMockUnitGetAgi(11, 0),
					// Only check summoner agi
					getMockUnitGetAgi(12, 50),
				},
			},
			want: []GroundIdx{1, 7, 2, 8, 3, 9, 4, 10, 5, 11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MinionSlot{
				AllyMinions:  tt.fields.AllyMinions,
				EnemyMinions: tt.fields.EnemyMinions,
			}
			assert.Equalf(t, tt.want, s.getActionOrder(), "getActionOrder()")
		})
	}
}

func getMockDeadActorUnit(groundIdx GroundIdx) *mockUnit {
	u := newMockUnit()
	u.On("GetGroundIdx").Return(groundIdx)
	u.On("IsDead").Return(true)
	return u
}

func getMockSkillWithTargetNumber(targetNumber int) vo.Skill {
	return vo.Skill{
		Name: "attack",
		AttributeMap: vo.NewAttributeMap(
			[]vo.Attribute{
				vo.NewAttribute(vo.AttributeTypeTarget, decimal.NewFromInt(int64(targetNumber))),
			}...,
		),
	}
}

func getMockAliveActorUnit(groundIdx GroundIdx, skill vo.Skill) *mockUnit {
	u := newMockUnit()
	u.On("GetGroundIdx").Return(groundIdx)
	u.On("IsDead").Return(false)
	u.On("GetSkill").Return(skill)
	return u
}

func getMockTargetUnit(groundIdx GroundIdx) *mockUnit {
	u := newMockUnit()
	u.On("GetGroundIdx").Return(groundIdx)
	u.On("IsDead").Return(false)
	u.On("TakeAffect", mock.Anything).Return(&mockUnit{MockUnit: &MockUnit{}, changed: true})
	return u
}

func getMockDeadUnit(groundIdx GroundIdx, isDead bool) *mockUnit {
	u := newMockUnit()
	u.On("GetGroundIdx").Return(groundIdx)
	u.On("IsDead").Return(isDead)
	return u
}

func Test_MinionSlot_Act(t *testing.T) {
	type fields struct {
		AllyMinions             *Minions
		EnemyMinions            *Minions
		calculateAttackDamageFn calculateAttackDamageFn
		targetPickerFn          targetPickerFn
	}
	type args struct {
		actorIdx GroundIdx
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantAffects []Affect
		wantStatus  MinionSlotStatus
	}{
		{
			name: "produce no affect when actor is dead",
			fields: fields{
				AllyMinions: &Minions{
					getMockDeadActorUnit(1),
				},
				// when actor is dead, the target picker will return nil.
				EnemyMinions:   &Minions{},
				targetPickerFn: defaultTargetPickerFn,
			},
			args:        args{actorIdx: 1},
			wantAffects: nil,
		},
		{
			name: "actor attacks two enemy units and produce affects",
			fields: fields{
				AllyMinions: &Minions{
					getMockAliveActorUnit(1, getMockSkillWithTargetNumber(2)),
					nil,
					nil,
					nil,
					nil,
					getMockDeadUnit(6, false),
				},
				EnemyMinions: &Minions{
					getMockTargetUnit(7),
					getMockTargetUnit(8),
					nil,
					nil,
					nil,
					getMockDeadUnit(12, false),
				},
				calculateAttackDamageFn: getHitCalculateAttackDamageFn(),
				targetPickerFn:          defaultTargetPickerFn,
			},
			args: args{actorIdx: 1},
			wantAffects: []Affect{
				{
					ActorIdx:  1,
					TargetIdx: 7,
					Skill:     "attack",
					Attributes: []vo.Attribute{
						vo.NewAttribute(vo.AttributeTypeDamage, decimal.NewFromInt(1)),
					},
				},
				{
					ActorIdx:  1,
					TargetIdx: 8,
					Skill:     "attack",
					Attributes: []vo.Attribute{
						vo.NewAttribute(vo.AttributeTypeDamage, decimal.NewFromInt(1)),
					},
				},
			},
		},
		{
			name: "ally won when enemy summoner is dead",
			fields: fields{
				AllyMinions: &Minions{
					getMockAliveActorUnit(1, getMockSkillWithTargetNumber(0)),
					nil,
					nil,
					nil,
					nil,
					getMockDeadUnit(6, false),
				},
				EnemyMinions: &Minions{
					getMockTargetUnit(7),
					getMockTargetUnit(8),
					nil,
					nil,
					nil,
					getMockDeadUnit(12, true),
				},
				targetPickerFn: defaultTargetPickerFn,
			},
			args:        args{actorIdx: 1},
			wantAffects: nil,
			wantStatus:  MinionSlotStatusAllyWon,
		},
		{
			name: "enemy won when ally summoner is dead",
			fields: fields{
				AllyMinions: &Minions{
					getMockAliveActorUnit(1, getMockSkillWithTargetNumber(0)),
					nil,
					nil,
					nil,
					nil,
					getMockDeadUnit(6, true),
				},
				EnemyMinions: &Minions{
					getMockTargetUnit(7),
					getMockTargetUnit(8),
					nil,
					nil,
					nil,
					getMockDeadUnit(12, false),
				},
				targetPickerFn: defaultTargetPickerFn,
			},
			args:        args{actorIdx: 1},
			wantAffects: nil,
			wantStatus:  MinionSlotStatusEnemyWon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MinionSlot{
				AllyMinions:             tt.fields.AllyMinions,
				EnemyMinions:            tt.fields.EnemyMinions,
				calculateAttackDamageFn: tt.fields.calculateAttackDamageFn,
				targetPickerFn:          tt.fields.targetPickerFn,
			}
			assert.Equalf(t, tt.wantAffects, s.act(tt.args.actorIdx), "act(%v)", tt.args.actorIdx)
			assert.Equalf(t, tt.wantStatus, s.Status, "Status(%v)", tt.args.actorIdx)
		})
	}
}

func Test_MinionSlot_Act_ChangeStatus(t *testing.T) {

}
