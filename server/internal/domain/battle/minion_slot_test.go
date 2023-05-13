package battle

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/Amobe/PlayGame/server/internal/domain/vo"
)

func getMockGround() vo.Ground {
	return vo.NewGround(
		vo.Camp{
			vo.NewCharacter(1, vo.NewAttributeMap()),
			vo.NewCharacter(2, vo.NewAttributeMap()),
			vo.NewCharacter(3, vo.NewAttributeMap()),
			vo.NewCharacter(4, vo.NewAttributeMap()),
			vo.NewCharacter(5, vo.NewAttributeMap()),
			vo.NewCharacter(6, vo.NewAttributeMap()),
		}, vo.Camp{
			vo.NewCharacter(7, vo.NewAttributeMap()),
			vo.NewCharacter(8, vo.NewAttributeMap()),
			vo.NewCharacter(9, vo.NewAttributeMap()),
			vo.NewCharacter(10, vo.NewAttributeMap()),
			vo.NewCharacter(11, vo.NewAttributeMap()),
			vo.NewCharacter(12, vo.NewAttributeMap()),
		},
	)
}

func Test_MinionSlot_unitTakeAffect(t *testing.T) {
	type args struct {
		groundUnit vo.Character
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "groundIdx is enemy, the enemy should be changed",
			args: args{
				groundUnit: vo.NewCharacter(8, vo.NewAttributeMap()),
			},
		},
		{
			name: "groundIdx is ally, the ally should be changed",
			args: args{
				groundUnit: vo.NewCharacter(4, vo.NewAttributeMap()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMinionSlot(getMockGround())
			err := s.unitTakeAffect(tt.args.groundUnit, nil)
			assert.NoError(t, err)
		})
	}
}

func getMissCalculateAttackDamageFn() calculateAttackDamageFn {
	return func(attacker, target vo.Character, skill vo.Skill) (decimal.Decimal, bool) {
		return decimal.Zero, false
	}
}

func getHitCalculateAttackDamageFn() calculateAttackDamageFn {
	return func(attacker, target vo.Character, skill vo.Skill) (decimal.Decimal, bool) {
		return decimal.NewFromInt(1), true
	}
}

func Test_MinionSlot_attack(t *testing.T) {
	type fields struct {
		calculateAttackDamageFn calculateAttackDamageFn
	}
	type args struct {
		attacker vo.Character
		target   vo.Character
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
				attacker: vo.NewCharacter(1, vo.NewAttributeMap()),
				target:   vo.NewCharacter(2, vo.NewAttributeMap()),
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
				attacker: vo.NewCharacterWithSkill(1, vo.SkillSlash, vo.NewAttributeMap()),
				target:   vo.NewCharacter(2, vo.NewAttributeMap()),
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
			s := &MinionSlot{
				Ground:                  getMockGround(),
				calculateAttackDamageFn: tt.fields.calculateAttackDamageFn,
			}
			got, err := s.attack(tt.args.attacker, tt.args.target)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, got, "attack(%v, %v)", tt.args.attacker, tt.args.target)
		})
	}
}

func Test_MinionSlot_Act(t *testing.T) {
	getAttackSkillWithTarget := func(targetNumber int) vo.Skill {
		targets := vo.NewAttribute(vo.AttributeTypeTarget, decimal.NewFromInt(int64(targetNumber)))
		return vo.NewSkill("attack", vo.NewAttributeMap(targets))
	}

	type fields struct {
		AllyMinions             vo.Camp
		EnemyMinions            vo.Camp
		calculateAttackDamageFn calculateAttackDamageFn
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
				AllyMinions: vo.Camp{
					vo.NewCharacter(1, vo.NewAttributeMap(vo.DeadAttribute)),
				},
				// when actor is dead, the target picker will return nil.
				EnemyMinions: vo.Camp{},
			},
			args:        args{actorIdx: 1},
			wantAffects: nil,
		},
		{
			name: "actor attacks two enemy units and produce affects",
			fields: fields{
				AllyMinions: vo.Camp{
					vo.NewCharacterWithSkill(1, getAttackSkillWithTarget(2), vo.NewAttributeMap()),
					vo.NewCharacter(2, vo.NewAttributeMap()),
					vo.NewCharacter(3, vo.NewAttributeMap()),
					vo.NewCharacter(4, vo.NewAttributeMap()),
					vo.NewCharacter(5, vo.NewAttributeMap()),
					vo.NewCharacter(6, vo.NewAttributeMap()),
				},
				EnemyMinions: vo.Camp{
					vo.NewCharacter(7, vo.NewAttributeMap()),
					vo.NewCharacter(8, vo.NewAttributeMap()),
					vo.NewCharacter(9, vo.NewAttributeMap()),
					vo.NewCharacter(10, vo.NewAttributeMap()),
					vo.NewCharacter(11, vo.NewAttributeMap()),
					vo.NewCharacter(12, vo.NewAttributeMap()),
				},
				calculateAttackDamageFn: getHitCalculateAttackDamageFn(),
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
				AllyMinions: vo.Camp{
					vo.NewCharacterWithSkill(1, getAttackSkillWithTarget(0), vo.NewAttributeMap()),
					vo.NewCharacter(2, vo.NewAttributeMap()),
					vo.NewCharacter(3, vo.NewAttributeMap()),
					vo.NewCharacter(4, vo.NewAttributeMap()),
					vo.NewCharacter(5, vo.NewAttributeMap()),
					vo.NewCharacter(6, vo.NewAttributeMap()),
				},
				EnemyMinions: vo.Camp{
					vo.NewCharacter(7, vo.NewAttributeMap()),
					vo.NewCharacter(8, vo.NewAttributeMap()),
					vo.NewCharacter(9, vo.NewAttributeMap()),
					vo.NewCharacter(10, vo.NewAttributeMap()),
					vo.NewCharacter(11, vo.NewAttributeMap()),
					vo.NewCharacter(12, vo.NewAttributeMap(vo.DeadAttribute)),
				},
			},
			args:        args{actorIdx: 1},
			wantAffects: nil,
			wantStatus:  MinionSlotStatusAllyWon,
		},
		{
			name: "enemy won when ally summoner is dead",
			fields: fields{
				AllyMinions: vo.Camp{
					vo.NewCharacterWithSkill(1, getAttackSkillWithTarget(0), vo.NewAttributeMap()),
					vo.NewCharacter(2, vo.NewAttributeMap()),
					vo.NewCharacter(3, vo.NewAttributeMap()),
					vo.NewCharacter(4, vo.NewAttributeMap()),
					vo.NewCharacter(5, vo.NewAttributeMap()),
					vo.NewCharacter(6, vo.NewAttributeMap(vo.DeadAttribute)),
				},
				EnemyMinions: vo.Camp{
					vo.NewCharacter(7, vo.NewAttributeMap()),
					vo.NewCharacter(8, vo.NewAttributeMap()),
					vo.NewCharacter(9, vo.NewAttributeMap()),
					vo.NewCharacter(10, vo.NewAttributeMap()),
					vo.NewCharacter(11, vo.NewAttributeMap()),
					vo.NewCharacter(12, vo.NewAttributeMap()),
				},
			},
			args:        args{actorIdx: 1},
			wantAffects: nil,
			wantStatus:  MinionSlotStatusEnemyWon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MinionSlot{
				Ground:                  vo.NewGround(tt.fields.AllyMinions, tt.fields.EnemyMinions),
				calculateAttackDamageFn: tt.fields.calculateAttackDamageFn,
			}
			gotAffects, err := s.act(tt.args.actorIdx)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantAffects, gotAffects, "act(%v)", tt.args.actorIdx)
			assert.Equalf(t, tt.wantStatus, s.Status, "Status(%v)", tt.args.actorIdx)
		})
	}
}

func Test_MinionSlot_Act_ChangeStatus(t *testing.T) {

}
