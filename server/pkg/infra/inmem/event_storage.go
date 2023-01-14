package inmem

import (
	"fmt"
	"time"

	"github.com/Amobe/PlayGame/server/pkg/utils"
	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

type aggregatorLoader[T domain.Aggregator] func([]domain.Event) (T, error)

type inmemEventStorage[T domain.Aggregator] struct {
	storage map[string][]eventRecord
	loader  aggregatorLoader[T]
}

func newInmemEventStorage[T domain.Aggregator](loader aggregatorLoader[T]) *inmemEventStorage[T] {
	return &inmemEventStorage[T]{
		storage: make(map[string][]eventRecord),
		loader:  loader,
	}
}

func (s *inmemEventStorage[T]) Save(t T) error {
	var records []eventRecord
	appliedEvents := t.Events()
	for i := range appliedEvents {
		record := eventRecord{
			Event:     appliedEvents[i],
			Version:   i + t.Version(),
			ApplyTime: time.Now(),
		}
		records = append(records, record)
	}

	events := s.storage[t.ID()]
	events = append(events, records...)
	s.storage[t.ID()] = events
	fmt.Println(utils.ToString(s.storage))
	return nil
}

func (s *inmemEventStorage[T]) Get(id string) (t T, err error) {
	records, ok := s.storage[id]
	if !ok {
		err = fmt.Errorf("events not found")
		return
	}
	var events []domain.Event
	for _, r := range records {
		events = append(events, r.Event)
	}
	return s.loader(events)
}

type eventRecord struct {
	Event     domain.Event
	Version   int
	ApplyTime time.Time
}
