package driver

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// Configuration
type Configuration struct {
	Driver    DriverSettings    `json:"driver"`
	Backend   BackendSettings   `json:"backend"`
	Generator GeneratorSettings `json:"generator,omitempty"`
}

// DriverSettings
type DriverSettings struct {
	RootPath string `json:"rootpath"`
}

// BackendSettings holds the settings for the libkv backend
type BackendSettings struct {
	Type      string `json:"type"`
	Endpoints string `json:"endpoints"`
}

// GeneratorSettings
type GeneratorSettings struct {
	Disabled bool `json:"disabled,omitempty"`
}

// LoadFromString loads a configuration from json string
func (c *Configuration) LoadFromString(d string) error {
	if d == "" {
		return errors.New("Loading empty json data")
	}

	if err := json.Unmarshal([]byte(d), &c); err != nil {
		return err
	}

	return nil
}

// LoadFromFile loads a configuration from json file
func (c *Configuration) LoadFromFile(p string) error {
	dataBuffer, err := ioutil.ReadFile(p)

	if err != nil {
		return err
	}

	return c.LoadFromString(string(dataBuffer))
}

// CheckIntegrity tests the configuration integrity
func (c *Configuration) CheckIntegrity() (bool, []error) {
	errorList := []error{}

	// check root path
	if stat, err := os.Stat(c.Driver.RootPath); err != nil || stat.IsDir() == false {
		errorList = append(errorList, errors.New("driver.rootpath directory did not exist"))
	}

	// check backend type
	if c.Backend.Type != "etcd" {
		errorList = append(errorList, errors.New("backend.type only supports 'etcd' at the moment"))
	}

	// check backend endpoints
	if c.Backend.Endpoints == "" {
		errorList = append(errorList, errors.New("backend.endpoints is a neccessary field"))
	}

	res := len(errorList) == 0
	return res, errorList
}

// GetBackendEndpointList returns the endpoints as list
func (c *Configuration) GetBackendEndpointList() []string {
	return strings.Split(strings.Replace(c.Backend.Endpoints, " ", "", -1), ",")
}

// NewConfiguration creates a new Configuration
func NewConfiguration() *Configuration {
	c := &Configuration{}
	c.Backend.Type = "etcd"
	return c
}
