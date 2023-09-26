package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App
	Coingecko
	Etherscan
	Bscscan
}

type App struct {
	APIAddress string `yaml:"address"`
	LogFile    string `yaml:"log_file"`
	LogLevel   string `yaml:"log_level"`
}

type Coingecko struct {
	APIAddress string `yaml:"api"`
}

type Etherscan struct {
	APIAddress      string `yaml:"api"`
	APIKey          string `yaml:"apikey"`
	ContractAddress string `yaml:"contractaddress"`
}

type Bscscan struct {
	APIAddress      string `yaml:"api"`
	APIKey          string `yaml:"apikey"`
	ContractAddress string `yaml:"contractaddress"`
}

func New(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	var config Config

	if err = yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
