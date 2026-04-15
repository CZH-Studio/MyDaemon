package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Value  ValueConfig   `json:"value"`
	Emails []EmailConfig `json:"emails"`
}

type EmailConfig struct {
	From string `json:"from"`
	Pwd  string `json:"pwd"`
	To   string `json:"to"`
}

type ValueConfig struct {
	BufferSize   int    `json:"buffer_size"`
	LogName      string `json:"log_name"`
	FromTitle    string `json:"from_title"`
	SubjectTitle string `json:"subject_title"`
}

func DefaultConfig() Config {
	return Config{
		Value: ValueConfig{
			BufferSize:   20,
			LogName:      "mydaemon.log",
			FromTitle:    "MyDaemon",
			SubjectTitle: "Process Result",
		},
		Emails: []EmailConfig{},
	}
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return home + string(os.PathSeparator) + ".mydaemon.json"
}

func LoadConfig() (*Config, error) {
	path := GetConfigPath()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig()
			if err := SaveConfig(&cfg); err != nil {
				return nil, err
			}
			return &cfg, nil
		}
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&cfg)
	if err != nil {
		cfg = DefaultConfig()
		if err := SaveConfig(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}
	return &cfg, nil
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

func HandleConfig(args *Args) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println(ErrorReadingFile)
		os.Exit(1)
	}
	if args.ConfigBuffer != 0 {
		config.Value.BufferSize = args.ConfigBuffer
	}
	if args.ConfigLog != "" {
		config.Value.LogName = args.ConfigLog
	}
	if args.ConfigFrom != "" {
		config.Value.FromTitle = args.ConfigFrom
	}
	if args.ConfigSubject != "" {
		config.Value.SubjectTitle = args.ConfigSubject
	}
	err = SaveConfig(config)
	if err != nil {
		fmt.Println(ErrorWritingFile)
		os.Exit(1)
	}
}
