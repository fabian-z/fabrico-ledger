package main

import (
	"encoding/hex"
	"errors"
	"io"
	"sync"

	"github.com/c2fo/vfs/v6"
	"github.com/c2fo/vfs/v6/vfssimple"
	"github.com/hashicorp/go-uuid"
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

// see https://github.com/C2FO/vfs/blob/master/docs/vfssimple.md
type VFSStore struct {
	baseURL  string
	location vfs.Location
}

// NewVFSStore takes a URL qualified with the storage backend as scheme:
//Local OS: file:///some/path
//Memory: mem:///some/path
//Amazon S3: s3://mybucket/path
//Google Cloud Storage: gs://mybucket/path

func NewVFSStore(baseURL string) *VFSStore {
	store := &VFSStore{
		baseURL: baseURL,
	}

	loc, err := vfssimple.NewLocation(baseURL)
	if err != nil {
		panic(err)
	}

	store.location = loc
	return store
}

func (vfs *VFSStore) StoreData(payload []byte) (FabricationDataHash, error) {
	hash := sha3.New512()
	hash.Write(payload)

	var address FabricationDataHash
	copy(address[:], hash.Sum(nil))

	tmpName, err := uuid.GenerateUUID()
	if err != nil {
		return FabricationDataHash{}, err
	}

	file, err := vfs.location.NewFile(tmpName)
	if err != nil {
		return FabricationDataHash{}, err
	}

	_, err = file.Write(payload)
	if err != nil {
		return FabricationDataHash{}, err
	}
	err = file.Close()
	if err != nil {
		return FabricationDataHash{}, err
	}

	tmpFile, err := vfs.location.NewFile(tmpName)
	if err != nil {
		return FabricationDataHash{}, err
	}
	file, err = vfs.location.NewFile(hex.EncodeToString(address[:]))
	if err != nil {
		return FabricationDataHash{}, err
	}

	err = tmpFile.MoveToFile(file)
	if err != nil {
		return FabricationDataHash{}, err
	}

	return address, nil
}

func (vfs *VFSStore) GetData(address FabricationDataHash) ([]byte, error) {

	file, err := vfs.location.NewFile(hex.EncodeToString(address[:]))
	if err != nil {
		return nil, err
	}

	exist, err := file.Exists()
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("file does not exist")
	}

	res, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (vfs *VFSStore) DeleteData(address FabricationDataHash) error {

	file, err := vfs.location.NewFile(hex.EncodeToString(address[:]))
	if err != nil {
		return err
	}

	exist, err := file.Exists()
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("file does not exist")
	}

	return file.Delete()
}

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
