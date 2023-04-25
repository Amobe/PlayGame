package database

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Amobe/PlayGame/server/internal/domain/battle"
	"github.com/Amobe/PlayGame/server/internal/utils"
	"github.com/Amobe/PlayGame/server/internal/utils/domain"
)

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

func (b BattleEventGorm) ToDomainEvent() (domain.Event, error) {
	// TODO: use generic to replace this
	switch b.EventType {
	case battle.EventBattleCreated{}.Name():
		var event battle.EventBattleCreated
		if err := utils.UnmarshalFromJSON(b.EventData, &event); err != nil {
			return nil, fmt.Errorf("unmarshal battle created event: %w", err)
		}
		return event, nil
	case battle.EventBattleFought{}.Name():
		var event battle.EventBattleFought
		if err := utils.UnmarshalFromJSON(b.EventData, &event); err != nil {
			return nil, fmt.Errorf("unmarshal battle fought event: %w", err)
		}
		return event, nil
	case battle.EventBattleWon{}.Name():
		var event battle.EventBattleWon
		if err := utils.UnmarshalFromJSON(b.EventData, &event); err != nil {
			return nil, fmt.Errorf("unmarshal battle won event: %w", err)
		}
		return event, nil
	case battle.EventBattleLost{}.Name():
		var event battle.EventBattleLost
		if err := utils.UnmarshalFromJSON(b.EventData, &event); err != nil {
			return nil, fmt.Errorf("unmarshal battle lost event: %w", err)
		}
		return event, nil
	case battle.EventBattleDraw{}.Name():
		var event battle.EventBattleDraw
		if err := utils.UnmarshalFromJSON(b.EventData, &event); err != nil {
			return nil, fmt.Errorf("unmarshal battle draw event: %w", err)
		}
		return event, nil
	}
	return nil, fmt.Errorf("unknown event type: %s", b.EventType)
}

func (BattleEventGorm) TableName() string {
	return "battle_event"
}

type BattleGormRepository struct {
	client *gorm.DB
}

func NewBattleGormRepository(client *gorm.DB) (*BattleGormRepository, error) {
	if err := client.AutoMigrate(&BattleEventGorm{}); err != nil {
		return nil, fmt.Errorf("migrate battle event gorm: %w", err)
	}
	return &BattleGormRepository{
		client: client,
	}, nil
}

// Get gets battle by battle id
func (r *BattleGormRepository) Get(battleID string) (*battle.Battle, error) {
	var battleEventGormList []BattleEventGorm
	if err := r.client.
		Order("created_at asc, id asc").
		Where("battle_id = ?", battleID).
		Find(&battleEventGormList).
		Error; err != nil {
		return nil, fmt.Errorf("find battle event: %w", err)
	}
	battleEventList := make([]domain.Event, 0, len(battleEventGormList))
	for _, battleEventGorm := range battleEventGormList {
		event, err := battleEventGorm.ToDomainEvent()
		if err != nil {
			return nil, fmt.Errorf("battle event gorm to domain event: %w", err)
		}
		battleEventList = append(battleEventList, event)
	}
	battleEntity, err := battle.AggregatorLoader(battleEventList)
	if err != nil {
		return nil, fmt.Errorf("aggregate battle: %w", err)
	}
	return battleEntity, nil
}

// Save saves battle
func (r *BattleGormRepository) Save(b *battle.Battle) error {
	battleEventGormList, err := BattleEventGormListFromEntity(b.ID(), b.Events())
	if err != nil {
		return fmt.Errorf("battle event from entity: %w", err)
	}
	if err := r.client.CreateInBatches(battleEventGormList, 100).Error; err != nil {
		return fmt.Errorf("save battle event: %w", err)
	}
	return nil
}
