package main

import (
	"log"

	"github.com/esacteksab/annoyed-aardvark/config"
	"github.com/esacteksab/annoyed-aardvark/internal"
	"github.com/esacteksab/annoyed-aardvark/internal/logger"
)

func main() {
	// Read Config
	cfg, err := config.LoadFromFile(".aardvark.toml")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	logger.SetVerbose(cfg.Verbose)
	logger.Debug("Gonna sync some shit!")
	if logger.Verbose {
		logger.Debug("Debug logging enabled")
		logger.Debugf("Main Logger verbose mode: %v", logger.Verbose)
		logger.Debugf("Config verbose is: %v", cfg.Verbose)
	}

	// Get a list of files from various sources and copy them
	for _, source := range cfg.Sources {
		internal.CopyFile(source.Path, source.Files)
	}
}
