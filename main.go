package main

import (
	"os"
	"strconv"

	"github.com/axelspringer/docker-conf-volume/driver"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	debug := os.Getenv("DEBUG")
	if ok, _ := strconv.ParseBool(debug); ok {
		logger.SetLevel(logrus.DebugLevel)
	}

	/*
		configFilePath := "/etc/docker/docker-confvol-plugin"
		config, cerr := driver.LoadConfigurationFromFile(configFilePath)
		if cerr != nil {
			logger.Fatal(cerr)
		}
	*/

	volumeStore, serr := driver.NewStore([]string{"172.17.0.2:4001"}, logger)
	if serr != nil {
		logger.Fatal(serr)
	}

	volumeDriver, verr := driver.NewConfigVolume(logger, volumeStore)
	if verr != nil {
		logger.Fatal(verr)
	}

	volumeHandler := volume.NewHandler(volumeDriver)
	if err := volumeHandler.ServeUnix("confvol", 0); err != nil {
		logger.Fatal(err)
	}
}
