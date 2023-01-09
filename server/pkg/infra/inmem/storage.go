package inmem

import (
	"fmt"

	"github.com/Amobe/PlayGame/server/pkg/utils/domain"
)

type inmemStorage[T domain.Aggregator] struct {
	storage map[string]T
}

func newInmemStorage[T domain.Aggregator]() *inmemStorage[T] {
	return &inmemStorage[T]{
		storage: make(map[string]T),
	}
}

func (i *inmemStorage[T]) Create(v T) error {
	i.storage[v.ID()] = v
	return nil
}

func (i *inmemStorage[T]) Get(id string) (v T, err error) {
	v, ok := i.storage[id]
	if !ok {
		err = fmt.Errorf("not found")
		return
	}
	return v, nil
}

func (i *inmemStorage[T]) Save(v T) error {
	i.storage[v.ID()] = v
	return nil
}

func (i *inmemStorage[T]) Delete(id string) error {
	delete(i.storage, id)
	return nil
}
