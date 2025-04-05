package main

import (
	"fmt"
	"io"
	"log"
	"os"

	kt "github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func main() {
	GetConfig()
	// This is a simple Go program that prints "Hello, World!" to the console.
	// It serves as a basic example of how to structure a Go application.
}

func GetConfig() {
	// Create a new Koanf instance
	k := koanf.New(".")

	// Load the configuration file
	if err := k.Load(file.Provider(".aardvark.toml"), kt.Parser()); err != nil {
		panic(err)
	}

	tables := append([]string{}, k.MapKeys("")...)
	// fmt.Printf("Tables in config file: %s", tables)
	for _, table := range tables {
		// The Table
		fmt.Println(table)
		// Get Path for each table
		path := k.String(table + ".path")
		fmt.Printf(" Path: %s\n", path)

		// Get Files for each table
		files := k.Strings(table + ".files")
		fmt.Printf(" Files: %s\n", files)

		CopyFile(path, files)

	}
}

// CopyFile will copy files from the source path to the destination path.
func CopyFile(path string, files []string) {
	// This function will copy files from the source path to the destination path.
	for _, file := range files {
		fmt.Printf("Copying file from %s: %s\n", path, file)
	}
	sourceFile, err := os.Open(path)
	if err != nil {
		log.Fatalf(err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(path)
	if err != nil {
		log.Fatalf(err)
	}
	defer destinationFile.Close()
	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		log.Fatalf(err)
	}
}
