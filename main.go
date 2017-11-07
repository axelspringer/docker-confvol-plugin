package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/axelspringer/docker-conf-volume/driver"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
)

var (
	configuration  *driver.Configuration
	configFilePath string
)

// process flags
func init() {
	// args
	flag.StringVar(&configFilePath, "config", "", "Path to the configuration file")
	// parse
	flag.Parse()
}

func main() {
	logger := logrus.New()
	debug := os.Getenv("DEBUG")
	if ok, _ := strconv.ParseBool(debug); ok {
		logger.SetLevel(logrus.DebugLevel)
	}

	// configuration
	configuration = driver.NewConfiguration()
	if len(configFilePath) > 0 {
		configuration.LoadFromFile(configFilePath)
	}

	// check configuration integrity
	if integer, errList := configuration.CheckIntegrity(); integer == false {
		for _, err := range errList {
			logger.Error(err)
		}

		os.Exit(1)
	}

	// create kv store
	volumeStore, serr := driver.NewStore(configuration, logger)
	if serr != nil {
		logger.Fatal(serr)
	}

	// create volume driver
	volumeDriver, verr := driver.NewConfigVolume(configuration, logger, volumeStore)
	if verr != nil {
		logger.Fatal(verr)
	}

	// create docker plugin socket
	volumeHandler := volume.NewHandler(volumeDriver)
	if err := volumeHandler.ServeUnix("confvol", 0); err != nil {
		logger.Fatal(err)
	}
}
