package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/axelspringer/docker-conf-volume/driver"
	"github.com/docker/go-plugins-helpers/volume"
)

func main() {
	debug := os.Getenv("DEBUG")
	if ok, _ := strconv.ParseBool(debug); ok {
		logrus.SetLevel(logrus.DebugLevel)
	}

	volumeDriver, err := driver.NewConfigVolume()
	if err != nil {
		log.Fatal(err)
	}

	volumeHandler := volume.NewHandler(volumeDriver)
	if err := volumeHandler.ServeUnix("confvol", 0); err != nil {
		log.Fatal(err)
	}
}
