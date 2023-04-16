package battle

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/internal/domain/vo"
)

func getMockMinions() (ally *Minions, enemy *Minions) {
	return &Minions{
			vo.NewGroundUnit(1, vo.NewCharacter("1", vo.NewAttributeMap())),
			vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap())),
			vo.NewGroundUnit(3, vo.NewCharacter("3", vo.NewAttributeMap())),
			vo.NewGroundUnit(4, vo.NewCharacter("4", vo.NewAttributeMap())),
			vo.NewGroundUnit(5, vo.NewCharacter("5", vo.NewAttributeMap())),
			vo.NewGroundUnit(6, vo.NewCharacter("6", vo.NewAttributeMap())),
		}, &Minions{
			vo.NewGroundUnit(7, vo.NewCharacter("7", vo.NewAttributeMap())),
			vo.NewGroundUnit(8, vo.NewCharacter("8", vo.NewAttributeMap())),
			vo.NewGroundUnit(9, vo.NewCharacter("9", vo.NewAttributeMap())),
			vo.NewGroundUnit(10, vo.NewCharacter("10", vo.NewAttributeMap())),
			vo.NewGroundUnit(11, vo.NewCharacter("11", vo.NewAttributeMap())),
			vo.NewGroundUnit(12, vo.NewCharacter("12", vo.NewAttributeMap())),
		}
}

func Test_MinionSlot_unitTakeAffect(t *testing.T) {
	type args struct {
		groundUnit vo.GroundUnit
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "groundIdx is enemy, the enemy should be changed",
			args: args{
				groundUnit: vo.NewGroundUnit(vo.GroundIdx(8), vo.NewCharacter("8", vo.NewAttributeMap())),
			},
		},
		{
			name: "groundIdx is ally, the ally should be changed",
			args: args{
				groundUnit: vo.NewGroundUnit(vo.GroundIdx(4), vo.NewCharacter("4", vo.NewAttributeMap())),
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
			s.unitTakeAffect(tt.args.groundUnit, nil)
			//// assert the specific minion changed
			//assert.Truef(t, s.getUnit(tt.args.groundIdx).(*mockUnit).changed,
			//	"unit changed expect true=%t", s.getUnit(tt.args.groundIdx).(*mockUnit).changed)
		})
	}
}

func getMissCalculateAttackDamageFn() calculateAttackDamageFn {
	return func(attacker, target vo.GroundUnit, skill vo.Skill) (decimal.Decimal, bool) {
		return decimal.Zero, false
	}
}

func getHitCalculateAttackDamageFn() calculateAttackDamageFn {
	return func(attacker, target vo.GroundUnit, skill vo.Skill) (decimal.Decimal, bool) {
		return decimal.NewFromInt(1), true
	}
}

func Test_MinionSlot_attack(t *testing.T) {
	type fields struct {
		calculateAttackDamageFn calculateAttackDamageFn
	}
	type args struct {
		attacker vo.GroundUnit
		target   vo.GroundUnit
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   vo.Affect
	}{
		{
			name: "get miss affect",
			fields: fields{
				calculateAttackDamageFn: getMissCalculateAttackDamageFn(),
			},
			args: args{
				attacker: vo.NewGroundUnit(1, vo.NewCharacter("1", vo.NewAttributeMap())),
				target:   vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap())),
			},
			want: vo.Affect{
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
				attacker: vo.NewGroundUnit(1, vo.NewCharacterWithSkill("1", vo.SkillSlash, vo.NewAttributeMap())),
				target:   vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap())),
			},
			want: vo.Affect{
				ActorIdx:  1,
				TargetIdx: 2,
				Skill:     "slash",
				Attributes: vo.NewAttributeMap(
					vo.NewAttribute(vo.AttributeTypeDamage, decimal.NewFromInt(1)),
				),
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

func Test_MinionSlot_getAttackerAndTargets(t *testing.T) {
	type fields struct {
		dead vo.Attribute
	}
	type args struct {
		idx vo.GroundIdx
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantAttackerGroundIdx vo.GroundIdx
		wantTargetsGroundIdx  []vo.GroundIdx
	}{
		{
			name: "attacker is ally, targets are enemy and not dead",
			fields: fields{
				dead: vo.Attribute{},
			},
			args: args{
				idx: vo.GroundIdx(2),
			},
			wantAttackerGroundIdx: vo.GroundIdx(2),
			wantTargetsGroundIdx:  []vo.GroundIdx{7, 8},
		},
		{
			name: " attacker is ally, targets are enemy and dead, targets will be summoner",
			fields: fields{
				dead: vo.DeadAttribute,
			},
			args: args{
				idx: vo.GroundIdx(2),
			},
			wantAttackerGroundIdx: vo.GroundIdx(2),
			wantTargetsGroundIdx:  []vo.GroundIdx{12, 12},
		},
		{
			name: "attacker is enemy, targets are ally and not dead",
			fields: fields{
				dead: vo.Attribute{},
			},
			args: args{
				idx: vo.GroundIdx(9),
			},
			wantAttackerGroundIdx: vo.GroundIdx(9),
			wantTargetsGroundIdx:  []vo.GroundIdx{1, 2},
		},
		{
			name: "attacker is enemy, targets are ally and dead, targets will be summoner",
			fields: fields{
				dead: vo.DeadAttribute,
			},
			args: args{
				idx: vo.GroundIdx(9),
			},
			wantAttackerGroundIdx: vo.GroundIdx(9),
			wantTargetsGroundIdx:  []vo.GroundIdx{6, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MinionSlot{
				AllyMinions: &Minions{
					vo.NewGroundUnit(1, vo.NewCharacter("1", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(3, vo.NewCharacter("3", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(4, vo.NewCharacter("4", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(5, vo.NewCharacter("5", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(6, vo.NewCharacter("6", vo.NewAttributeMap(tt.fields.dead))),
				},
				EnemyMinions: &Minions{
					vo.NewGroundUnit(7, vo.NewCharacter("7", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(8, vo.NewCharacter("8", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(9, vo.NewCharacter("9", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(10, vo.NewCharacter("10", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(11, vo.NewCharacter("11", vo.NewAttributeMap(tt.fields.dead))),
					vo.NewGroundUnit(12, vo.NewCharacter("12", vo.NewAttributeMap(tt.fields.dead))),
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

func Test_MinionSlot_getActionOrder(t *testing.T) {
	type fields struct {
		AllyMinions  *Minions
		EnemyMinions *Minions
	}
	tests := []struct {
		name   string
		fields fields
		want   []vo.GroundIdx
	}{
		{
			name: "enemy summoner is faster than ally summoner",
			fields: fields{
				AllyMinions: &Minions{
					vo.NewGroundUnit(1, vo.NewCharacter("1", vo.NewAttributeMap())),
					vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap())),
					vo.NewGroundUnit(3, vo.NewCharacter("3", vo.NewAttributeMap())),
					vo.NewGroundUnit(4, vo.NewCharacter("4", vo.NewAttributeMap())),
					vo.NewGroundUnit(5, vo.NewCharacter("5", vo.NewAttributeMap())),
					// Only check summoner agi
					vo.NewGroundUnit(6, vo.NewCharacter("6", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(50))))),
				},
				EnemyMinions: &Minions{
					vo.NewGroundUnit(7, vo.NewCharacter("7", vo.NewAttributeMap())),
					vo.NewGroundUnit(8, vo.NewCharacter("8", vo.NewAttributeMap())),
					vo.NewGroundUnit(9, vo.NewCharacter("9", vo.NewAttributeMap())),
					vo.NewGroundUnit(10, vo.NewCharacter("10", vo.NewAttributeMap())),
					vo.NewGroundUnit(11, vo.NewCharacter("11", vo.NewAttributeMap())),
					// Only check summoner agi
					vo.NewGroundUnit(12, vo.NewCharacter("12", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(100))))),
				},
			},
			want: []vo.GroundIdx{7, 1, 8, 2, 9, 3, 10, 4, 11, 5},
		},
		{
			name: "ally summoner is faster than enemy summoner",
			fields: fields{
				AllyMinions: &Minions{
					vo.NewGroundUnit(1, vo.NewCharacter("1", vo.NewAttributeMap())),
					vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap())),
					vo.NewGroundUnit(3, vo.NewCharacter("3", vo.NewAttributeMap())),
					vo.NewGroundUnit(4, vo.NewCharacter("4", vo.NewAttributeMap())),
					vo.NewGroundUnit(5, vo.NewCharacter("5", vo.NewAttributeMap())),
					// Only check summoner agi
					vo.NewGroundUnit(6, vo.NewCharacter("6", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(100))))),
				},
				EnemyMinions: &Minions{
					vo.NewGroundUnit(7, vo.NewCharacter("7", vo.NewAttributeMap())),
					vo.NewGroundUnit(8, vo.NewCharacter("8", vo.NewAttributeMap())),
					vo.NewGroundUnit(9, vo.NewCharacter("9", vo.NewAttributeMap())),
					vo.NewGroundUnit(10, vo.NewCharacter("10", vo.NewAttributeMap())),
					vo.NewGroundUnit(11, vo.NewCharacter("11", vo.NewAttributeMap())),
					// Only check summoner agi
					vo.NewGroundUnit(12, vo.NewCharacter("12", vo.NewAttributeMap(vo.NewAttribute(vo.AttributeTypeAGI, decimal.NewFromInt(50))))),
				},
			},
			want: []vo.GroundIdx{1, 7, 2, 8, 3, 9, 4, 10, 5, 11},
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

func Test_MinionSlot_Act(t *testing.T) {
	getAttackSkillWithTarget := func(targetNumber int) vo.Skill {
		targets := vo.NewAttribute(vo.AttributeTypeTarget, decimal.NewFromInt(int64(targetNumber)))
		return vo.NewSkill("attack", vo.NewAttributeMap(targets))
	}

	type fields struct {
		AllyMinions             *Minions
		EnemyMinions            *Minions
		calculateAttackDamageFn calculateAttackDamageFn
		targetPickerFn          targetPickerFn
	}
	type args struct {
		actorIdx vo.GroundIdx
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantAffects []vo.Affect
		wantStatus  MinionSlotStatus
	}{
		{
			name: "produce no affect when actor is dead",
			fields: fields{
				AllyMinions: &Minions{
					vo.NewGroundUnit(1, vo.NewCharacter("1", vo.NewAttributeMap(vo.DeadAttribute))),
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
					vo.NewGroundUnit(1, vo.NewCharacterWithSkill("1", getAttackSkillWithTarget(2), vo.NewAttributeMap())),
					vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap())),
					vo.NewGroundUnit(3, vo.NewCharacter("3", vo.NewAttributeMap())),
					vo.NewGroundUnit(4, vo.NewCharacter("4", vo.NewAttributeMap())),
					vo.NewGroundUnit(5, vo.NewCharacter("5", vo.NewAttributeMap())),
					vo.NewGroundUnit(6, vo.NewCharacter("6", vo.NewAttributeMap())),
				},
				EnemyMinions: &Minions{
					vo.NewGroundUnit(7, vo.NewCharacter("7", vo.NewAttributeMap())),
					vo.NewGroundUnit(8, vo.NewCharacter("8", vo.NewAttributeMap())),
					vo.NewGroundUnit(9, vo.NewCharacter("9", vo.NewAttributeMap())),
					vo.NewGroundUnit(10, vo.NewCharacter("10", vo.NewAttributeMap())),
					vo.NewGroundUnit(11, vo.NewCharacter("11", vo.NewAttributeMap())),
					vo.NewGroundUnit(12, vo.NewCharacter("12", vo.NewAttributeMap())),
				},
				calculateAttackDamageFn: getHitCalculateAttackDamageFn(),
				targetPickerFn:          defaultTargetPickerFn,
			},
			args: args{actorIdx: 1},
			wantAffects: []vo.Affect{
				{
					ActorIdx:  1,
					TargetIdx: 7,
					Skill:     "attack",
					Attributes: vo.NewAttributeMap(
						vo.NewAttribute(vo.AttributeTypeDamage, decimal.NewFromInt(1)),
					),
				},
				{
					ActorIdx:  1,
					TargetIdx: 8,
					Skill:     "attack",
					Attributes: vo.NewAttributeMap(
						vo.NewAttribute(vo.AttributeTypeDamage, decimal.NewFromInt(1)),
					),
				},
			},
		},
		{
			name: "ally won when enemy summoner is dead",
			fields: fields{
				AllyMinions: &Minions{
					vo.NewGroundUnit(1, vo.NewCharacterWithSkill("1", getAttackSkillWithTarget(0), vo.NewAttributeMap())),
					vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap())),
					vo.NewGroundUnit(3, vo.NewCharacter("3", vo.NewAttributeMap())),
					vo.NewGroundUnit(4, vo.NewCharacter("4", vo.NewAttributeMap())),
					vo.NewGroundUnit(5, vo.NewCharacter("5", vo.NewAttributeMap())),
					vo.NewGroundUnit(6, vo.NewCharacter("6", vo.NewAttributeMap())),
				},
				EnemyMinions: &Minions{
					vo.NewGroundUnit(7, vo.NewCharacter("7", vo.NewAttributeMap())),
					vo.NewGroundUnit(8, vo.NewCharacter("8", vo.NewAttributeMap())),
					vo.NewGroundUnit(9, vo.NewCharacter("9", vo.NewAttributeMap())),
					vo.NewGroundUnit(10, vo.NewCharacter("10", vo.NewAttributeMap())),
					vo.NewGroundUnit(11, vo.NewCharacter("11", vo.NewAttributeMap())),
					vo.NewGroundUnit(12, vo.NewCharacter("12", vo.NewAttributeMap(vo.DeadAttribute))),
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
					vo.NewGroundUnit(1, vo.NewCharacterWithSkill("1", getAttackSkillWithTarget(0), vo.NewAttributeMap())),
					vo.NewGroundUnit(2, vo.NewCharacter("2", vo.NewAttributeMap())),
					vo.NewGroundUnit(3, vo.NewCharacter("3", vo.NewAttributeMap())),
					vo.NewGroundUnit(4, vo.NewCharacter("4", vo.NewAttributeMap())),
					vo.NewGroundUnit(5, vo.NewCharacter("5", vo.NewAttributeMap())),
					vo.NewGroundUnit(6, vo.NewCharacter("6", vo.NewAttributeMap(vo.DeadAttribute))),
				},
				EnemyMinions: &Minions{
					vo.NewGroundUnit(7, vo.NewCharacter("7", vo.NewAttributeMap())),
					vo.NewGroundUnit(8, vo.NewCharacter("8", vo.NewAttributeMap())),
					vo.NewGroundUnit(9, vo.NewCharacter("9", vo.NewAttributeMap())),
					vo.NewGroundUnit(10, vo.NewCharacter("10", vo.NewAttributeMap())),
					vo.NewGroundUnit(11, vo.NewCharacter("11", vo.NewAttributeMap())),
					vo.NewGroundUnit(12, vo.NewCharacter("12", vo.NewAttributeMap())),
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
