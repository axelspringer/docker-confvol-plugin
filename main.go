package main

import (
	"flag"
	"os"

	"github.com/axelspringer/docker-conf-volume/driver"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
)

var (
	configuration  *driver.Configuration
	configFilePath string
	debugFlag      bool
)

// process flags
func init() {
	// args
	flag.StringVar(&configFilePath, "config", "", "Path to the configuration file")
	flag.BoolVar(&debugFlag, "debug", false, "Set debug mode")
	// parse
	flag.Parse()
}

func main() {
	logger := logrus.New()

	// configuration
	configuration = driver.NewConfiguration()

	// env root path
	if p := os.Getenv("CONFVOL_DRIVER_ROOT"); len(p) > 0 {
		configuration.Driver.RootPath = p
	}

	// env etcd instances
	if e := os.Getenv("CONFVOL_BACKEND_ENDPOINTS"); len(e) > 0 {
		configuration.Backend.Endpoints = e
	}

	// env root path
	if d := os.Getenv("CONFVOL_DEBUG"); len(d) > 0 {
		debugFlag = true
	}

	// set debug mode
	//if debugFlag == true {
	logger.SetLevel(logrus.DebugLevel)
	//}

	// load configuration from file
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
