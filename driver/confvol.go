package driver

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/docker/libkv/store"
)

// VolumeMount
type VolumeMount struct {
	Root             string
	Relative         string
	ReferenceCounter int
}

// ConfigVolume driver
type ConfigVolume struct {
	logger     *logrus.Logger
	volumes    map[string]*VolumeMount
	m          *sync.Mutex
	mountPoint string
	store      *Store
}

func (v *ConfigVolume) syncFolder(kvEntries []*store.KVPair, basePath string, relativePath string) {
	kv := v.store.Client

	for _, pair := range kvEntries {
		v.logger.Infof("Sync pair %#v", pair.Key)

		p, e := kv.List(pair.Key)
		v.logger.Infof("List key %s %#v %#v", pair.Key, p, e)
		g, e2 := kv.Get(pair.Key)
		v.logger.Infof("Get key %s %#v %#v", pair.Key, g, e2)

		filePath := pair.Key[len(relativePath):]

		v.logger.Infof("Reduce to %s", filePath)

	}
}

func (v *ConfigVolume) syncMountPoint(vm *VolumeMount) error {
	kv := v.store.Client
	entries, err := kv.List(vm.Relative)

	v.logger.Infof("Sync mount point %s", vm.Root)

	if err != nil {
		v.logger.Error(err)
		return err
	}

	v.logger.Infof("Success in pulling data")
	v.syncFolder(entries, vm.Root, vm.Relative)

	return nil
}

// Create This function is called each time a client wants to create a volume
func (v *ConfigVolume) Create(r *volume.CreateRequest) error {
	v.logger.Infof("Create volume %s", r.Name)
	v.logger.Infof("Dump %#v", r)
	v.m.Lock()
	defer v.m.Unlock()

	// already loaded
	if _, ok := v.volumes[r.Name]; ok {
		return nil
	}

	// create base dir
	volumePath := filepath.Join(v.mountPoint, r.Name)
	os.MkdirAll(volumePath, os.ModePerm)

	_, err := os.Lstat(volumePath)
	if err != nil {
		logrus.Errorf("Error %s %v", volumePath, err.Error())
		return err
	}

	vm := &VolumeMount{
		Root:             volumePath,
		Relative:         r.Name,
		ReferenceCounter: 0,
	}

	v.volumes[r.Name] = vm
	return nil
}

// List returns a list of all mounted volumes
func (v *ConfigVolume) List() (*volume.ListResponse, error) {
	v.logger.Printf("List\n")

	volumes := []*volume.Volume{}

	for name, vol := range v.volumes {
		volumes = append(volumes, &volume.Volume{
			Name:       name,
			Mountpoint: vol.Root,
		})
	}

	return &volume.ListResponse{Volumes: volumes}, nil
}

// Returns a volume by name
func (v *ConfigVolume) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	v.logger.Printf("Get %#v\n", r)

	if vol, ok := v.volumes[r.Name]; ok {
		v.logger.Printf("Found entry %#v\n", vol)

		return &volume.GetResponse{
			Volume: &volume.Volume{
				Name:       r.Name,
				Mountpoint: vol.Root,
			},
		}, nil
	}

	return nil, errors.New("Element not found")
}

// Remove is called when the client wants to remove the vol
func (v *ConfigVolume) Remove(r *volume.RemoveRequest) error {
	v.logger.Printf("Remove %#v\n", r)

	v.m.Lock()
	defer v.m.Unlock()

	if _, ok := v.volumes[r.Name]; ok {
		delete(v.volumes, r.Name)
	}

	return nil
}

// Path
func (v *ConfigVolume) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	v.logger.Printf("Path %s\n", r.Name)

	if vm, ok := v.volumes[r.Name]; ok {
		return &volume.PathResponse{
			Mountpoint: vm.Root,
		}, nil
	}

	return &volume.PathResponse{}, nil
}

// Mount can be used for ressource allocation
func (v *ConfigVolume) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	v.logger.Printf("Mounting volume %s\n%#v\n", r.Name, r)

	if vm, ok := v.volumes[r.Name]; ok {
		v.syncMountPoint(vm)
		return &volume.MountResponse{
			Mountpoint: vm.Root,
		}, nil
	}

	return &volume.MountResponse{}, nil
}

// Unmount
func (v *ConfigVolume) Unmount(r *volume.UnmountRequest) error {
	v.logger.Printf("Unmounting volume %#v\n", r)

	return nil
}

func (v *ConfigVolume) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}

// NewConfigVolume creates a new ConfigVolume
func NewConfigVolume(l *logrus.Logger, s *Store) (*ConfigVolume, error) {
	return &ConfigVolume{
		logger:     l,
		volumes:    make(map[string]*VolumeMount),
		m:          &sync.Mutex{},
		mountPoint: "/tmp/confvol/",
		store:      s,
	}, nil
}
