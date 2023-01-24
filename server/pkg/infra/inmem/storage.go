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

type inmemStorage[T Indexer] struct {
	storage map[string]T
}

func newInmemStorage[T Indexer]() *inmemStorage[T] {
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
		err = ErrorRecordNotFound
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
