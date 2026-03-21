package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	From     string `json:"from"`
	Password string `json:"password"`
	To       string `json:"to"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return home + string(os.PathSeparator) + ".mydaemon.json"
}

func LoadConfig() (*Config, error) {
	path := GetConfigPath()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	return &cfg, err
}

func SaveConfig(cfg *Config) error {
	path := GetConfigPath()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}
