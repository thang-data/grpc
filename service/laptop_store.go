package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/copier"
	"gitlab.com/dem1/dem1/pb"
)

// ErrAlreadyExists is returned when a record with the same ID already exists in the store
var ErrAlreadyExists = errors.New("record already exists")

// LaptopStore in an interface to store laptop
type LaptopStore interface {
	// Save saves the laptop to the store
	Save(laptop *pb.Laptop) error
}

type InMemoryLaptopStore struct{
	mutex sync.RWMutex
	data map[string]*pb.Laptop
}

// NewMemoryLaptopStore return a new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore{
	return  &InMemoryLaptopStore {
		data: make(map[string]*pb.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error{
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil{
		return ErrAlreadyExists
	}

	// deep copy

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil{
		return fmt.Errorf("cannot copy laptop data: %v", err)
	}

	store.data[laptop.Id] = other
	return nil
}