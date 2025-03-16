package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

var (
	cfgFile    string
	configDir  string
	configName string
	createFile bool
	cwd        string
	homeDir    string
	noConfig   bool
)

type ConfigParams struct {
	Binary   string `toml:"binary" comment:"binary: (type: string) The name of the binary, either 'tofu' or 'terraform'. Must exist on your $PATH."`
	PlanFile string `toml:"planFile" comment:"planFile: (type: string) The name of the plan file created by 'gh tp'."`
	MdFile   string `toml:"mdFile" comment:"mdFile: (type: string) The name of the Markdown file created by 'gh tp'."`
	Verbose  bool   `toml:"verbose" comment:"verbose: (type: bool) Enable Verbose Logging. Default is false."`
}

func main() {

	homeDir, configDir, cwd, err := getDirectories()
	if err != nil {
		log.Fatal("Bad things happened here!")
	}

	configName = ".tp.toml"

	huh.NewSelect[string]().
		Title("Where would you like to save your .tp.toml config file?").
		Options(
			huh.NewOption("Home Config Directory: "+configDir+"/"+configName,
				configDir).Selected(true),
			huh.NewOption("Home Directory: "+homeDir+"/"+configName, homeDir),
			huh.NewOption("Project Root:"+"/"+configName, cwd),
		).
		Value(&cfgFile).
		Run()

	fmt.Printf("Inside main(), and config is %s/%s\n", cfgFile, configName)

	// fmt.Println(cfgFile)
	noConfig, createFile := mkFile(cfgFile)

	binary := "terraform"
	planFile := "plan.out"
	mdFile := "plan.md"

	conf := ConfigParams{
		Binary:   binary,
		PlanFile: planFile,
		MdFile:   mdFile,
		Verbose:  false,
	}

	config, err := genConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	if createFile {
		fmt.Printf("Inside main() if createFile and createFile is %t\n", createFile)
		if !noConfig {
			existingConfigFile := cfgFile + "/" + configName
			err := backupFile(existingConfigFile, existingConfigFile)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = os.WriteFile(cfgFile+"/.tp.toml", config, 0o600) //nolint:mnd    // https://go.dev/ref/spec#Integer_literals
		if err != nil {
			log.Fatalf("Error writing Config file: %s", err)
		}
		log.Infof("Config file %s/.tp.toml created", cfgFile)
	} else if !createFile {
		fmt.Printf("Inside main if !createFile and createFile is %t\n", createFile)
		fmt.Printf("Not writing file, writing to stdout: \n%s\n", string(config))
	}

}

func getDirectories() (homeDir, configDir, cwd string, err error) {
	homeDir = xdg.Home

	configDir = xdg.ConfigHome

	cwd, cwderr := os.Getwd()
	if cwderr != nil {
		log.Errorf("Error: %s", err)
	}
	return homeDir, configDir, cwd, err
}

// returns true if file doesn't exist
func doesNotExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}
	return false
}

// takes cfgFile, appends ".tp.toml" to it
// checks to see if file exists
// based on existence, asks to create (doesn't exist)
// or overwrite (exists)
func mkFile(cfgFile string) (exists, createFile bool) {
	configName = ".tp.toml"
	noConfig := doesNotExist(cfgFile + "/" + configName)
	fmt.Println(cfgFile + configName)
	if noConfig {
		// log.Info("File Exists!")
		fmt.Printf("%s/%s doesn't exist\n", cfgFile, configName)
		huh.NewConfirm().
			Title("Create new file?").
			Affirmative("Yes").
			Negative("No").
			Value(&createFile).
			Run()
		fmt.Printf("Inside mkFile() and config is %s/%s\n", cfgFile, configName)
		return noConfig, createFile
	} else if !noConfig {
		huh.NewConfirm().
			Title("Overwrite existing config file?").
			Affirmative("Yes").
			Negative("No").
			Value(&createFile).Run()
		fmt.Printf("Inside mkFile if exists, config is %s/%s\n", cfgFile, configName)
		return noConfig, createFile

	}
	fmt.Printf("inside mkFile() noConfig is %t\n", noConfig)
	fmt.Printf("Inside mkFile() config is %s/%s\n", cfgFile, configName)
	return noConfig, createFile
}

func genConfig(conf ConfigParams) (data []byte, err error) {
	data, err = toml.Marshal(conf)
	if err != nil {
		log.Fatalf("Failed marshalling TOML: %s", err)
	}
	return data, err
}

// backupFile copies the file at source to dest
func backupFile(source, dest string) error {
	// epoch as an int64
	e := time.Now().Unix()

	// epoch as string
	epoch := strconv.FormatInt(e, 10)

	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest + "-backup-" + epoch)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	err = destFile.Sync()
	return err
}
