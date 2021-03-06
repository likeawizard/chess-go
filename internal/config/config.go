package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const CONFIG_PATH = "config.yml"

type Config struct {
	Init    Init    `yaml:"init"`
	Engine  Engine  `yaml:"engine"`
	Render  Render  `yaml:"render"`
	Lichess Lichess `yaml:"lichess"`
}

type Init struct {
	StartingFen string `yaml:"startingFen"`
}

type Engine struct {
	Algorithm     string `yaml:"algorithm"`
	MaxDepth      int    `yaml:"maxDepth"`
	MaxGoRoutines int    `yaml:"maxGoRoutines"`
	Debug         bool   `yaml:"debug"`
}

type Render struct {
	Mode string `yaml:"mode"`
}

type Lichess struct {
	APIToken string `yaml:"apiToken"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	cfgFile, err := os.Open(CONFIG_PATH)
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	d := yaml.NewDecoder(cfgFile)
	err = d.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
