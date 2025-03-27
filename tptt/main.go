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

type Config struct {
	Projects map[string]Project `toml:"project"`
}

type Projects struct {
	Project Project
}

type Project struct {
	Name     string `toml:"-"`
	Binary   string `toml:"binary"`
	PlanFile string `toml:"planFile"`
	MdFile   string `toml:"mdFile"`
	Verbose  bool   `toml:"verbose"`
}

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
	// MakeTemp()
	repo, err := repository.Current()
	if err != nil {
		log.Fatal(err)
	}

	projectName := "esacteksab/gh-tp"
	binary = "terraform"
	planFile = "plan.out"
	mdFile = "plan.md"
	verbose = false

	p2Name := repo.Owner + "/" + repo.Name
	p2binary := "tofu"
	p2planFile := "tfplan.out"
	p2mdFile := "tfplan.md"
	p2Verbose := true

	fmt.Println(projectName)

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

	b, _ := toml.Marshal(data)
	fmt.Println("v2:\n" + string(b))

	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(true)
	enc.Encode(data)
	fmt.Println("v2 Encoder:\n" + buf.String())

	err = os.WriteFile("config.toml", buf.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Println("TOML file created successfully")
	// fmt.Println(data)

	cfg := ReadConfig("config.toml")

	// fmt.Printf("%s binary is: %s", p2Name, cfg.String(p2Name+".binary"))
	foo := cfg.String(p2Name)
	fmt.Println(foo.binary)
	//	for k, v := range foo {
	//		fmt.Printf("Key: %v, Value: %v\n", k, v)
	//
	// }
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

// func createConfigFile(config Config) error {
//
// 	filename := "config.toml"
// 	f, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
//
// 	data, err := toml.Mar(f)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(encoder)
// 	return encoder.Encode(config)
// }

func ReadConfig(path string) *koanf.Koanf {
	var k = koanf.New(".")

	if err := k.Load(file.Provider("./config.toml"), ktoml.Parser()); err != nil {
		log.Fatalf("error loading config file: %v", err)
	}
	// if err := k.Load(file.Provider("./config.toml"), toml.Parser()); err != nil {
	// 	log.Fatalf("error reading from config: %v", err)
	// }

	//fmt.Println("config.toml is ", k)
	// fmt.Println(k.String("1.binary"))
	// fmt.Println(k.String("2.binary"))
	// for key, value := range k.All() {
	// fmt.Printf("Key: %s, Value: %v\n", key, value)
	// }
	return k
}
