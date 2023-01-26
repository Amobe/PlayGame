package inmem

import (
	"fmt"
)

var (
	ErrorRecordNotFound = fmt.Errorf("record not found")
)

type Indexer interface {
	ID() string
}

type InmemStorage[T Indexer] struct {
	storage map[string]T
}

func NewInmemStorage[T Indexer]() *InmemStorage[T] {
	return &InmemStorage[T]{
		storage: make(map[string]T),
	}
}

func (i *InmemStorage[T]) Create(v T) error {
	i.storage[v.ID()] = v
	return nil
}

func (i *InmemStorage[T]) Get(id string) (v T, err error) {
	v, ok := i.storage[id]
	if !ok {
		err = ErrorRecordNotFound
		return
	}
	return v, nil
}

func (i *InmemStorage[T]) Save(v T) error {
	i.storage[v.ID()] = v
	return nil
}

func (i *InmemStorage[T]) Delete(id string) error {
	delete(i.storage, id)
	return nil
}
