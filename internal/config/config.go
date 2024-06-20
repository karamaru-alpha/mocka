package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Packages []map[string]*Package `yaml:"packages"`
}

type Package struct {
	All        bool     `yaml:"all,omitempty"`
	Recursive  bool     `yaml:"recursive,omitempty"`
	Interfaces []string `yaml:"interfaces,omitempty"`
}

func Load() (*Config, error) {
	var configPath string
	flag.StringVar(&configPath, "c", ".mocka.yaml", "Config yaml file path")
	flag.Parse()

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("place the config file on %s. err: %+v", configPath, err)
	}
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("invalid config file on %s. err: %+v", configPath, err)
	}
	return &config, nil
}

func (p *Package) IsTarget(t string) bool {
	if p == nil || p.All {
		return true
	}
	for _, i := range p.Interfaces {
		if i == t {
			return true
		}
	}
	return false
}
