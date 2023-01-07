package inmem

import (
	"fmt"
)

type Entity interface {
	ID() string
}

type inmemStorage[T Entity] struct {
	storage map[string]T
}

func newInmemStorage[T Entity]() *inmemStorage[T] {
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
