package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func NewConfig() (*Config, error) {
	cfg, err := read()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read config: %v\nPlease Check the Readme!", err)
	}
	return &cfg, nil
}
func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	return write(*c)
}

func read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("Open File: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("Decoding File: %v", err)
	}

	return cfg, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Creating File: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return fmt.Errorf("Encoding: %v", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getConfigFilePath: Can't get home dir %v", err)
	}
	fullFilePath := filepath.Join(homedir, configFileName)
	return fullFilePath, nil

}
