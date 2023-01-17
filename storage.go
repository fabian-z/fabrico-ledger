package main

import (
	"errors"
	"sync"

	"golang.org/x/crypto/sha3"
)

// Using SHA3-512
type FabricationDataHash [64]byte

// TODO interface currently assuming data will fit into memory of node, change to streaming interface
type FabricationDataStore interface {
	StoreData(payload []byte) (FabricationDataHash, error)
	GetData(hash FabricationDataHash) ([]byte, error)
	DeleteData(hash FabricationDataHash) error
}

// TODO implement file storage backend
// TODO implement data logistic vendor interfaces

type MemoryStore struct {
	sync.Mutex
	store map[FabricationDataHash][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[FabricationDataHash][]byte),
	}
}

func (mem *MemoryStore) StoreData(payload []byte) (FabricationDataHash, error) {

	hash := sha3.New512()
	hash.Write(payload)

	var address FabricationDataHash
	copy(address[:], hash.Sum(nil))

	mem.Lock()
	mem.store[address] = payload
	mem.Unlock()

	return address, nil
}

func (mem *MemoryStore) GetData(address FabricationDataHash) ([]byte, error) {
	mem.Lock()
	val, ok := mem.store[address]
	mem.Unlock()

	if !ok {
		return nil, errors.New("data does not exist in memory store")
	}

	return val, nil
}

func (mem *MemoryStore) DeleteData(address FabricationDataHash) error {
	mem.Lock()
	defer mem.Unlock()

	delete(mem.store, address)
	return nil
}
