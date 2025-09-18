package main

import "errors"

var (
	ErrNotFound = errors.New("record not found")
)

type InmemoryRepository struct {
	store map[int]string
}

func NewInmemoryRepository() InmemoryRepository {
	store := map[int]string{
		1: "1st record",
		2: "2nd record",
		3: "3rd record",
	}

	return InmemoryRepository{store}
}

func (r InmemoryRepository) GetUser(id int) (string, error) {
	value, ok := r.store[id]
	if !ok {
		return "", ErrNotFound
	}
	return value, nil
}
