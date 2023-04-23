package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/utils"
	"github.com/Amobe/PlayGame/server/internal/utils/domain"
)

type BattleGorm struct {
	gorm.Model
	BattleID     string `gorm:"column:battle_id;type:varchar(64);uniqueIndex"`
	AllyMinions  []byte `gorm:"column:ally_minions;type:bytes;serializer:json"`
	EnemyMinions []byte `gorm:"column:enemy_minions;type:bytes;serializer:json"`
}

func BattleGormFromEntity(b *battle.Battle) (BattleGorm, error) {
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

type BattleEventGorm struct {
	gorm.Model
	BattleID  string `gorm:"column:battle_id;type:varchar(64);index"`
	EventType string `gorm:"column:event_type;type:varchar(64)"`
	EventData []byte `gorm:"column:event_data;type:bytes;serializer:json"`
}

func BattleEventGormListFromEntity(battleID string, events []domain.Event) ([]BattleEventGorm, error) {
	battleEventGormList := make([]BattleEventGorm, 0, len(events))
	for _, event := range events {
		eventData, err := utils.MarshalToJSON(event)
		if err != nil {
			return nil, fmt.Errorf("battle event from entity: %w", err)
		}
		eventGorm := BattleEventGorm{
			BattleID:  battleID,
			EventType: event.Name(),
			EventData: eventData,
		}
		battleEventGormList = append(battleEventGormList, eventGorm)
	}
	return battleEventGormList, nil
}

func (BattleEventGorm) TableName() string {
	return "battle_event"
}

type BattleGormRepository struct {
	client *gorm.DB
}

func NewBattleGormRepository(config Config) (*BattleGormRepository, error) {
	client, err := gorm.Open(postgres.Open(GetDSN(config)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect battle gorm: %w", err)
	}
	if err := client.AutoMigrate(&BattleGorm{}); err != nil {
		return nil, fmt.Errorf("migrate battle gorm: %w", err)
	}
	if err := client.AutoMigrate(&BattleEventGorm{}); err != nil {
		return nil, fmt.Errorf("migrate battle event gorm: %w", err)
	}
	if err != nil {
		return nil, err
	}
	return &BattleGormRepository{
		client: client,
	}, nil
}

// Get gets battle by battle id
func (r *BattleGormRepository) Get(battleID string) (*battle.Battle, error) {
	battleGorm, err := r.get(battleID)
	if err != nil {
		return nil, fmt.Errorf("get battle gorm: %w", err)
	}
	return battleGorm.ToEntity()
}

// Save saves battle
func (r *BattleGormRepository) Save(b *battle.Battle) error {
	gormBattle, err := r.get(b.ID())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.create(b)
	} else if err != nil {
		return fmt.Errorf("get battle gorm: %w", err)
	}
	return r.update(gormBattle.Model, b)
}

// get gets battle by battle id
func (r *BattleGormRepository) get(battleID string) (*BattleGorm, error) {
	var battleGorm BattleGorm
	if err := r.client.Where("battle_id = ?", battleID).First(&battleGorm).Error; err != nil {
		return nil, fmt.Errorf("find battle: %w", err)
	}
	return &battleGorm, nil
}

// create creates battle
func (r *BattleGormRepository) create(b *battle.Battle) error {
	battleGorm, err := BattleGormFromEntity(b)
	if err != nil {
		return fmt.Errorf("from entity: %w", err)
	}
	if err := r.client.Create(&battleGorm).Error; err != nil {
		return fmt.Errorf("create battle: %w", err)
	}
	return nil
}

// update updates battle
func (r *BattleGormRepository) update(gormModel gorm.Model, b *battle.Battle) error {
	battleGorm, err := BattleGormFromEntity(b)
	if err != nil {
		return fmt.Errorf("battle from entity: %w", err)
	}
	battleGorm.Model = gormModel

	battleEventGormList, err := BattleEventGormListFromEntity(b.ID(), b.Events())
	if err != nil {
		return fmt.Errorf("battle event from entity: %w", err)
	}

	if err := r.client.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&battleGorm).Where("battle_id", b.ID()).Error; err != nil {
			return fmt.Errorf("save battle: %w", err)
		}
		if err := tx.CreateInBatches(battleEventGormList, 100).Error; err != nil {
			return fmt.Errorf("create battle event: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
