package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aimotrens/scrutzone/cmd"
	"github.com/aimotrens/scrutzone/conf"
	"github.com/aimotrens/scrutzone/version"
)

func init() {
	version.RawBuildInfo = func() (string, string) {
		return compileDate, scrutzoneVersion
	}
}

var (
	scrutzoneVersion = "vX.X.X"
	compileDate      = "unknown"
)

func main() {
	versionString := version.String()
	fmt.Println(versionString)

	configFile := os.Getenv("SCRUTZONE_CONFIG_FILE")
	if configFile == "" {
		configFile = filepath.Join("config", "scrutzone.yml")
	}

	// Load the configuration
	config, err := conf.Load(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	// set the defaults for the checks
	for _, check := range config.Checks {
		check.SetDefaults(config.CheckDefaults)
	}

	// validate the configuration
	for _, check := range config.Checks {
		if err := check.Validate(); err != nil {
			log.Fatalf("Error validating check %s: %v", check.Name, err)
		}
	}

	runQueue := make(chan func() func(*cmd.Notification))
	for _, check := range config.Checks {
		check.InitTicker(runQueue)
	}

	if config.StartupNotification != nil {
		go config.Notification.Notify(config.StartupNotification.Targets, "scrutzone Startup", "scrutzone started")
	}

	for checkRun := range runQueue {
		go func() {
			if notify := checkRun(); notify != nil {
				notify(config.Notification)
			}
		}()
	}
}