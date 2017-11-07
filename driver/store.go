package driver

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

type StoreKVPair = store.KVPair

// Store interface
type Store interface {
	Get(key string) (*StoreKVPair, error)
	List(key string) ([]*StoreKVPair, error)	
}

// LibKVStore helper struct
type LibKVStore struct {
	Client store.Store
	logger *logrus.Logger
}

// Get a kv entry by key
func (s *LibKVStore) Get(key string) (*StoreKVPair, error) {
	kv := s.Client
	return kv.Get(key)
}

// List kv entries by key
func (s *LibKVStore) List(key string) ([]*StoreKVPair, error) {
	kv := s.Client
	return kv.List(key)
}

// NewStore creates a new store. suprise ..
func NewStore(c *Configuration, logger *logrus.Logger) (Store, error) {
	s := &LibKVStore{}

	// Initialize a new store with etcd
	kv, err := libkv.NewStore(
		store.Backend(c.Backend.Type),
		c.GetBackendEndpointList(),
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)

	if err != nil {
		return nil, err
	}

	s.Client = kv
	s.logger = logger

	return s, nil
}

// register the backend(s)
func init() {
	etcd.Register()
}
