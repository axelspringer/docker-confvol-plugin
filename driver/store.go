package driver

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

// Store helper struct
type Store struct {
	Client store.Store
	logger *logrus.Logger
}

// NewStore creates a new store. suprise ..
func NewStore(addrList []string, logger *logrus.Logger) (*Store, error) {
	s := &Store{}

	// Initialize a new store with etcd
	kv, err := libkv.NewStore(
		store.ETCD,
		addrList,
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
