package driver

import (
	"log"

	"github.com/docker/go-plugins-helpers/volume"
)

// ConfigVolume driver
type ConfigVolume struct {
}

func (v *ConfigVolume) Create(r *volume.CreateRequest) error {
	log.Printf("Create %s\n", r.Name)

	return nil
}

func (v *ConfigVolume) List() (*volume.ListResponse, error) {
	log.Printf("List\n")

	return &volume.ListResponse{}, nil
}

func (v *ConfigVolume) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	log.Printf("Get %s\n", r.Name)

	return &volume.GetResponse{}, nil
}

func (v *ConfigVolume) Remove(r *volume.RemoveRequest) error {
	log.Printf("Remove %s\n", r.Name)

	return nil
}

func (v *ConfigVolume) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	log.Printf("Path %s\n", r.Name)

	return &volume.PathResponse{}, nil
}

func (v *ConfigVolume) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	log.Printf("Mounting volume %s\n", r.Name)

	return &volume.MountResponse{}, nil
}

func (v *ConfigVolume) Unmount(r *volume.UnmountRequest) error {
	log.Printf("Unmounting volume %s\n", r.Name)

	return nil
}

func (v *ConfigVolume) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{}
}

// NewConfigVolume creates a new ConfigVolume
func NewConfigVolume() (*ConfigVolume, error) {
	return &ConfigVolume{}, nil
}
