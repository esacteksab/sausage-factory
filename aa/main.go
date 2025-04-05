package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	kt "github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func main() {
	GetConfig()
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
		// fmt.Println(table)
		// Get Path for each table
		path := k.String(table + ".path")
		// fmt.Printf(" Path: %s\n", path)

		// Get Files for each table
		files := k.Strings(table + ".files")
		// fmt.Printf(" Files: %s\n", files)

		CopyFile(path, files)

	}
}

// CopyFile will copy files from the source path to the destination path.
func CopyFile(path string, files []string) {
	// This function will copy files from the source path to the destination path.
	for _, file := range files {
		// fmt.Printf("Copying file from %s: %s\n", path, file)
		filePath := path + "/" + file
		// fmt.Println(filePath)
		// Open the source file for reading
		sourceFile, err := os.Open(filePath)
		// fmt.Printf("Opening file: %s\n", filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer sourceFile.Close()

		// We need to create the destination director[y|ies] if it doesn't exist
		fpath := filepath.Dir(file)

		// If it doesn't exist, create it
		if !doesExist(fpath) {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}

		// Presumably we are copying to the current directory, so this is `file`
		destinationFile, err := os.Create(file)
		// fmt.Printf("Creating file: %s\n", file)
		if err != nil {
			log.Fatal(err)
		}
		defer destinationFile.Close()

		// Check to see if the files are the same prior to trying to copy them.
		same, err := sameFile(filePath, file)
		if err != nil {
			log.Fatal(err)
		}

		if same {
			fmt.Printf("File %s already exists, skipping copy\n", file)
		} else if !same {
			// Copy the contents of the source file to the destination file
			_, err = io.Copy(destinationFile, sourceFile)
			fmt.Println("Files are not the same, copying...")
			fmt.Printf("Copying file from %s: %s\n", path, file)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func doesExist(path string) bool {
	// Check if the file or directory exists
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

func sameFile(file1, file2 string) (bool, error) {
	// Get info about both files
	f1, err := os.Stat(file1)
	if err != nil {
		return false, err
	}

	f2, err := os.Stat(file2)
	if err != nil {
		return false, err
	}
	// Check if the file names are the same

	return os.SameFile(f1, f2), nil
}
