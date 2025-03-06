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

func genConfig(conf map[string]interface{}) {

	data, err := toml.Marshal(conf)
	if err != nil {
		log.Fatalf("Failed marshalling TOML: %s", err)
	}

	err = os.WriteFile(".tp-test.toml", data, 0644)
	if err != nil {
		log.Fatalf("Error writing Config file: %s", err)
	}

	// fmt.Printf("Binary is: %s", cp.binary)
	// fmt.Printf("Plan File is: %s", cp.planFile)
	// fmt.Printf("Markdown File is: %s", cp.mdFile)
	// fmt.Printf("Project Name is: %s", cp.projectName)
}
