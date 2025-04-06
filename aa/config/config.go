package config

import (
	"fmt"
	"strings"

	kt "github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/esacteksab/annoyed-aardvark/internal/logger"
)

// SourceConfig holds the configuration for a single source
type SourceConfig struct {
	Path  string   `koanf:"path"`
	Files []string `koanf:"files"`
}

// Config is a map of source names to their configurations
// type Config map[string]SourceConfig

// Config holds the source configuration
type Config struct {
	Verbose bool                    `koanf:"verbose"`
	Sources map[string]SourceConfig `koanf:"sources"`
}

func LoadFromFile(path string) (Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), kt.Parser()); err != nil {
		return Config{}, err
	}
	err := k.Load(env.Provider("AA_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "AA_")), "_", ".", -1)
	}), nil)
	if err != nil {
		fmt.Errorf("Error", err)
	}

	verbose := k.Bool("verbose")
	logger.Debugf("LoadFromFile says verbose is: %v\n", verbose)

	config := Config{
		Verbose: verbose,
		Sources: make(map[string]SourceConfig),
	}

	for _, key := range k.MapKeys("") {
		if key == "verbose" {
			continue
		}

		src := SourceConfig{
			Path:  k.String(key + ".path"),
			Files: k.Strings(key + ".files"),
		}

		config.Sources[key] = src
	}

	return config, nil
}
