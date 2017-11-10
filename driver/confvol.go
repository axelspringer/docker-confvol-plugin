package driver

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
)

// VolumeMount
type VolumeMount struct {
	Root              string
	Relative          string
	ReferenceCounter  int
	Mode              int
	TemplateGenerator bool
}

// ConfigVolume driver
type ConfigVolume struct {
	logger     *logrus.Logger
	volumes    map[string]*VolumeMount
	m          *sync.Mutex
	mountPoint string
	store      Store
}

// synchronize a list of kv entries to the fs
func (v *ConfigVolume) syncFolder(kvEntries []*StoreKVPair, basePath string, relativePath string) {
	s := v.store

	for _, pair := range kvEntries {
		v.logger.Debugf("Sync source %s", pair.Key)

		fileName := pair.Key[len(relativePath):]
		dstPath := path.Join(basePath, fileName)

		entryList, _ := s.List(pair.Key)
		entryData, _ := s.Get(pair.Key)

		//TODO find a better way to distinguish files from folders (EC empty folder or empty file)
		isFolder := len(entryData.Value) == 0

		if isFolder == true {
			os.MkdirAll(dstPath, os.ModePerm)
			v.syncFolder(entryList, dstPath, pair.Key)
		} else {
			if err := ioutil.WriteFile(dstPath, entryData.Value, 0644); err != nil {
				v.logger.Error(err)
			}
		}
	}
}

// sync mount point. Folder mounts MUST end with a slash
func (v *ConfigVolume) syncMountPoint(vm *VolumeMount) error {
	s := v.store

	syncFolder := strings.HasSuffix(vm.Relative, "/")

	if syncFolder == true {
		entries, err := s.List(vm.Relative)

		if err != nil {
			v.logger.Error(err)
			return err
		}

		os.MkdirAll(vm.Root, os.ModePerm)
		v.syncFolder(entries, vm.Root, vm.Relative)
	} else {
		entry, err := s.Get(vm.Relative)
		if err != nil {
			v.logger.Error(err)
			return err
		}

		os.MkdirAll(path.Dir(vm.Root), os.ModePerm)

		mode := 0644
		if vm.Mode > 0 {
			mode = vm.Mode
		}

		data := entry.Value
		if vm.TemplateGenerator {
			tmpl := NewTemplate(v.store)
			tmplOutput, err := tmpl.Parse(string(data), nil)

			if err != nil {
				v.logger.Error(err)
			}
			data = []byte(tmplOutput)
		}

		if err := ioutil.WriteFile(vm.Root, data, os.FileMode(mode)); err != nil {
			v.logger.Error(err)
			return err
		}
	}

	return nil
}

// volumeExist
func (v *ConfigVolume) volumeExist(name string) bool {
	_, ok := v.volumes[name]
	return ok
}

// Create is called when a volume didn't exist yet
// In this case a former Get call returned
func (v *ConfigVolume) Create(r *volume.CreateRequest) error {
	v.logger.Debugf("Create volume %s", r.Name)

	v.m.Lock()
	defer v.m.Unlock()

	// already loaded
	if v.volumeExist(r.Name) {
		return nil
	}

	// create base dir
	volumePath := filepath.Join(v.mountPoint, r.Name)

	vm := &VolumeMount{
		Root:             volumePath,
		Relative:         r.Name,
		ReferenceCounter: 0,
	}

	// template mode
	if v, ok := r.Options["tmpl"]; ok && len(v) > 0 {
		vm.TemplateGenerator = true
	}

	// mode bits
	if v, ok := r.Options["mode"]; ok && len(v) > 0 {
		if m, err := strconv.ParseInt(v, 8, 64); err == nil {
			vm.Mode = int(m)
		}
	}

	v.volumes[r.Name] = vm
	return nil
}

// List returns a list of the available volumes
func (v *ConfigVolume) List() (*volume.ListResponse, error) {
	v.logger.Debugf("List volumes")
	v.m.Lock()
	defer v.m.Unlock()

	volumes := []*volume.Volume{}

	for name, vol := range v.volumes {
		volumes = append(volumes, &volume.Volume{
			Name:       name,
			Mountpoint: vol.Root,
		})
	}

	return &volume.ListResponse{Volumes: volumes}, nil
}

/* Get is the first request from the docker engine that reach the plugin
 * If a volume, specified by name (a path in this case), is prepared and ready to mount
 * then Create is not called
 */
func (v *ConfigVolume) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	v.logger.Debugf("Get volume %s", r.Name)
	v.m.Lock()
	defer v.m.Unlock()

	if v.volumeExist(r.Name) {
		return &volume.GetResponse{
			Volume: &volume.Volume{
				Name:       r.Name,
				Mountpoint: v.volumes[r.Name].Root,
			},
		}, nil
	}

	return nil, errors.New("Element not found")
}

// Remove is called when the client explicitly remove the volume (docker volume rm)
func (v *ConfigVolume) Remove(r *volume.RemoveRequest) error {
	v.logger.Debugf("Remove volume %s", r.Name)
	v.m.Lock()
	defer v.m.Unlock()

	if v.volumeExist(r.Name) {
		os.RemoveAll(v.volumes[r.Name].Root)
		delete(v.volumes, r.Name)
	}

	return nil
}

// Path
func (v *ConfigVolume) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	v.logger.Debugf("Path %s", r.Name)
	v.m.Lock()
	defer v.m.Unlock()

	res := &volume.PathResponse{}

	if v.volumeExist(r.Name) {
		res = &volume.PathResponse{
			Mountpoint: v.volumes[r.Name].Root,
		}
	}

	return res, nil
}

// Mount can be used for ressource allocation
func (v *ConfigVolume) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	v.logger.Debugf("Remove volume %s", r.Name)
	v.m.Lock()
	defer v.m.Unlock()

	res := &volume.MountResponse{}

	if vm, ok := v.volumes[r.Name]; ok {
		v.syncMountPoint(vm)
		vm.ReferenceCounter += 1
		res = &volume.MountResponse{
			Mountpoint: vm.Root,
		}
	}

	return res, nil
}

// Unmount
func (v *ConfigVolume) Unmount(r *volume.UnmountRequest) error {
	v.logger.Debugf("Unmounting volume %s", r.Name)
	v.m.Lock()
	defer v.m.Unlock()

	if vm, ok := v.volumes[r.Name]; ok {
		vm.ReferenceCounter -= 1
	}

	return nil
}

// Capabilities define the scope of this plugin
func (v *ConfigVolume) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}

// NewConfigVolume creates a new ConfigVolume
func NewConfigVolume(c *Configuration, l *logrus.Logger, s Store) (*ConfigVolume, error) {
	return &ConfigVolume{
		logger:     l,
		volumes:    make(map[string]*VolumeMount),
		m:          &sync.Mutex{},
		mountPoint: c.Driver.RootPath,
		store:      s,
	}, nil
}
