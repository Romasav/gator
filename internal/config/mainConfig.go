package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Username string `json:"current_user_name"`
	DbUrl    string `json:"db_url"`
}

func (con *Config) SetUpUser(username string) error {
	con.Username = username
	err := write(*con)
	if err != nil {
		return fmt.Errorf("could not write to config: %w", err)
	}
	return nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("could not find config file: %w", err)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("could not unmarshal config file: %w", err)
	}

	return config, nil
}

func write(config Config) error {
	jsonData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal config to JSON: %w", err)
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("could not find config file: %w", err)
	}

	err = os.WriteFile(configFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("could not write to config file: %w", err)
	}

	return nil
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not find current working directory: %w", err)
	}

	configFilePath := filepath.Join(currentDir, configFileName)
	return configFilePath, nil
}
