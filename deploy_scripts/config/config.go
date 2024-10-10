package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var VMConfigs *Config

type Config struct {
	Worker     int     `yaml:"worker"`
	Pool       int     `yaml:"pool"`
	StatsDir   string  `yaml:"stats-dir"`
	AvgDelay   int64   `yaml:"avg-delay"`
	Failure    float64 `yaml:"failure"`
	MaxWorkers int     `yaml:"max_workers"`
}

func LoadConfig() {
	var err error
	VMConfigs, err = readConfig("./config.yaml")
	if err != nil {
		log.Fatal("Error in reading configs", err)
	}
	VMConfigs.Failure = VMConfigs.Failure / 100

	// To avoid bankruptcy (jk)
	if VMConfigs.Worker > 10 {
		VMConfigs.Worker = 10
	}
	if VMConfigs.MaxWorkers > 10 {
		VMConfigs.MaxWorkers = 10
	}
}

func readConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode YAML: %v", err)
	}

	return &config, nil
}