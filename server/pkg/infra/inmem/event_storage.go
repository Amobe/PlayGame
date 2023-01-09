package inmem

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

type aggregatorLoader[T domain.Aggregator] func([]domain.Event) (T, error)

type inmemEventStorage[T domain.Aggregator] struct {
	storage map[string][]domain.Event
	loader  aggregatorLoader[T]
}

func newInmemEventStorage[T domain.Aggregator](loader aggregatorLoader[T]) *inmemEventStorage[T] {
	return &inmemEventStorage[T]{
		storage: make(map[string][]domain.Event),
		loader:  loader,
	}
}

func (s *inmemEventStorage[T]) Save(t T) error {
	events := s.storage[t.ID()]
	events = append(events, t.Events()...)
	s.storage[t.ID()] = events
	return nil
}

func (s *inmemEventStorage[T]) Get(id string) (t T, err error) {
	events, ok := s.storage[id]
	if !ok {
		err = fmt.Errorf("events not found")
		return
	}
	return s.loader(events)
}
