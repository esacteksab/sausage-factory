package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/cli/go-gh/v2/pkg/repository"
	ktoml "github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/pelletier/go-toml/v2"
)

// Config represents the config file
type Config struct {
	Projects map[string]Project `toml:"project"`
}

// A config file consists of one more Projects
type Projects struct {
	Project Project
}

// A Project has a:
// Name: string, repo.Owner/repo.Name
// Binary: string, expects Tofu or Terraform, only required if both binaries are in $PATH
// PlanFile: string, required
// MdFile: string, required
type Project struct {
	Name     string `toml:"-"`
	Binary   string `toml:"binary"`
	PlanFile string `toml:"planFile"`
	MdFile   string `toml:"mdFile"`
	Verbose  bool   `toml:"verbose"`
}

// tp has two core files, a plan output file and the markdown file
type TPFile struct {
	Name    string
	Purpose string
}

var (
	binary   string
	planFile string
	mdFile   string
	verbose  bool
)

func main() {
	// We need a unique identifier, so we're using repo.Owner/repo.Name as the key
	repo, err := repository.Current()
	if err != nil {
		log.Fatal(err)
	}

	// A Project
	projectName := "esacteksab/gh-tp"
	binary = "terraform"
	planFile = "plan.out"
	mdFile = "plan.md"
	verbose = false

	// Another project
	p2Name := repo.Owner + "/" + repo.Name
	p2binary := "tofu"
	p2planFile := "tfplan.out"
	p2mdFile := "tfplan.md"
	p2Verbose := true

	// This is a map[string]interface{} representing the projects above
	data := map[string]interface{}{
		projectName: map[string]any{
			"binary":   binary,
			"planFile": planFile,
			"mdFile":   mdFile,
			"verbose":  verbose,
		},
		p2Name: map[string]any{
			"binary":   p2binary,
			"planFile": p2planFile,
			"mdFile":   p2mdFile,
			"verbose":  p2Verbose,
		},
	}

	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)

	// I prefer indents in TOML for readability TOML doesn't give a shit.
	enc.SetIndentTables(true)
	err = enc.Encode(data)
	if err != nil {
		fmt.Println(err)
	}

	// Create the file
	err = os.WriteFile("config.toml", buf.Bytes(), 0o644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	// This essentially reads the file created above
	// The file was created with pellettier/toml/v2
	// The file is read with knadh/koanf
	GetProjectParameters(p2Name)
	GetProjectParameters(projectName)
}

// func MakeTemp() (date string, err error) {
//
// 	tpFiles := []TPFile{
// 		{"planFile", "Plan"},
// 		{"mdFile", "Markdown"},
// 	}
//
// 	date = time.Now().Local().Format("20060102150405")
//
// 	for i, v := range tpFiles {
// 		fmt.Printf("Index: %d, Name: %s, Purpose %s\n", i, v.Name, v.Purpose)
// 	}
//
// 	fmt.Println(date)
//
// 	tmpDir := os.TempDir()
//
// 	exists, err := doesExist(tmpDir)
// 	if err != nil {
// 		fmt.Errorf("Error: %s", err)
// 	}
//
// 	if exists {
// 		fmt.Printf("%s exists\n", tmpDir)
// 	} else {
// 		fmt.Print("No Temp Dir Found")
// 	}
//
// 	return date, err
// }
//
// func doesExist(file string) (exists bool, err error) {
// 	f, err := os.Stat(file)
// 	if err != nil {
// 		fmt.Errorf("%s does not exist. Error: %s", f.Name(), err)
// 		return false, err
// 	}
// 	return true, nil
// }

// GetProjectParameters gets a project's parameters from the defined config file
// More than one project may be defined in the config file
func GetProjectParameters(project string) {
	k := koanf.New(".")

	if err := k.Load(file.Provider("./config.toml"), ktoml.Parser()); err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	// The config file may consist of any projects, we just want this project
	config := k.Get(project)

	// Type assert the interface{} config to map[string]any
	fubar, ok := config.(map[string]any)
	if !ok {
		fmt.Println("interface is not a map")
	}
	binary, ok := fubar["binary"].(string)
	if !ok {
		fmt.Println("Something wrong binary: ", binary)
	}
	planFile, ok := fubar["planFile"].(string)
	if !ok {
		fmt.Println("something wrong with planFile: ", planFile)
	}
	mdFile, ok := fubar["mdFile"].(string)
	if !ok {
		fmt.Println("something wrong with mdFile: ", mdFile)
	}
	verbose, ok := fubar["verbose"].(bool)
	if !ok {
		fmt.Println("something wrong with verbose: ", verbose)
	}

	fmt.Println("\n========================")
	fmt.Println(project)
	fmt.Println(binary)
	fmt.Println(planFile)
	fmt.Println(mdFile)
	fmt.Println(verbose)
}
