package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	configFileKey     = "configFile"
	defaultConfigFile = "config/config.yml"
	configFileUsage   = "config file path"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
	flag.Parse()

	if err := Execute(configFile); err != nil {
		log.Errorf("error executing bot")
		os.Exit(1)
	}
}
