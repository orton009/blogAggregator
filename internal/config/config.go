package config

import (
	"encoding/json"
	"os"
)

const FileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(c)
}

func Read() (Config, error) {
	config := Config{}

	path, err := getConfigFilePath()
	if err != nil {
		return config, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	return config, nil
}
