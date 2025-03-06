package cmd

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/pelletier/go-toml/v2"
)

// func getRepo() (repo string, err error) {}

func getDirectories() (homeDir, configDir string) {
	homeDir = xdg.Home
	configDir = xdg.ConfigHome
	return homeDir, configDir
}

func getRepo() (name string) {
	repo, err := repository.Current()
	if err != nil {
		log.Fatal(err)
	}
	return repo.Owner + "/" + repo.Name
}

func getCWD() (cwd string, err error) {
	cwd, cwderr := os.Getwd()
	if cwderr != nil {
		log.Errorf("Error: %s", err)
	}
	return cwd, cwderr
}

func genConfig(conf map[string]any) (data []byte, err error) {
	data, err = toml.Marshal(conf)
	if err != nil {
		log.Fatalf("Failed marshalling TOML: %s", err)
	}
	return data, err
}
