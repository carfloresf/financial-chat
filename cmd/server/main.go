package main

import (
	"flag"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
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
		log.Errorf("error executing server")
		os.Exit(1)
	}
}
