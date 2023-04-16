package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/utils"
)

type BattleGorm struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement"`
	BattleID     string    `gorm:"column:battle_id;type:varchar(64);uniqueIndex"`
	AllyMinions  []byte    `gorm:"column:ally_minions;type:bytes;serializer:json"`
	EnemyMinions []byte    `gorm:"column:enemy_minions;type:bytes;serializer:json"`
	CreateAt     time.Time `gorm:"column:create_at;type:timestamp;DEFAULT:CURRENT_TIMESTAMP"`
}

func (BattleGorm) FromEntity(b *battle.Battle) (BattleGorm, error) {
	allyMinionsJson, err := utils.MarshalToJSON(b.MinionSlot().AllyMinions)
	if err != nil {
		return BattleGorm{}, fmt.Errorf("marshal ally minions: %w", err)
	}
	enemyMinionsJson, err := utils.MarshalToJSON(b.MinionSlot().EnemyMinions)
	if err != nil {
		return BattleGorm{}, fmt.Errorf("marshal enemy minions: %w", err)
	}

	return BattleGorm{
		BattleID:     b.ID(),
		AllyMinions:  allyMinionsJson,
		EnemyMinions: enemyMinionsJson,
	}, nil
}

func (b BattleGorm) ToEntity() (*battle.Battle, error) {
	var allyMinions, enemyMinions battle.Minions
	if err := utils.UnmarshalFromJSON(b.AllyMinions, &allyMinions); err != nil {
		return nil, fmt.Errorf("unmarshal ally minions: %w", err)
	}
	if err := utils.UnmarshalFromJSON(b.EnemyMinions, &enemyMinions); err != nil {
		return nil, fmt.Errorf("unmarshal enemy minions: %w", err)
	}
	return battle.CreateBattle(b.BattleID, battle.NewMinionSlot(&allyMinions, &enemyMinions))
}

func (BattleGorm) TableName() string {
	return "battle"
}

type BattleGormRepository struct {
	client *gorm.DB
}

func NewBattleGormRepository(dsn string) (*BattleGormRepository, error) {
	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect battle gorm: %w", err)
	}
	client.AutoMigrate(&BattleGorm{})
	return &BattleGormRepository{
		client: client,
	}, nil
}

// Get gets battle by battle id
func (r *BattleGormRepository) Get(battleID string) (*battle.Battle, error) {
	var battleGorm BattleGorm
	if err := r.client.Where("battle_id = ?", battleID).First(&battleGorm).Error; err != nil {
		return nil, fmt.Errorf("find battle: %w", err)
	}
	return battleGorm.ToEntity()
}

// Save saves battle
func (r *BattleGormRepository) Save(b *battle.Battle) error {
	battleGorm, err := BattleGorm{}.FromEntity(b)
	if err != nil {
		return fmt.Errorf("convert battle to gorm: %w", err)
	}
	if err := r.client.Save(&battleGorm).Error; err != nil {
		return fmt.Errorf("save battle: %w", err)
	}
	return nil
}
