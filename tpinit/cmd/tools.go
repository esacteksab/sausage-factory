package cmd

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

func getDirectories() (homeDir, configDir, cwd string, err error) {
	homeDir = xdg.Home

	configDir = xdg.ConfigHome

	cwd, cwderr := os.Getwd()
	if cwderr != nil {
		log.Errorf("Error: %s", err)
	}
	return homeDir, configDir, cwd, err
}

func genConfig(conf ConfigParams) (data []byte, err error) {
	data, err = toml.Marshal(conf)
	if err != nil {
		log.Fatalf("Failed marshalling TOML: %s", err)
	}
	return data, err
}
