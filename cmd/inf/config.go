package main

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
)

// TODO: manage multiple installations
type Config struct {
	Token            string `yaml:"token,omitempty"`
	DefaultNamespace string `yaml:"defaultNamespace,omitempty"`
}

func (c *Config) Write() error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".inf")
	_ = os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config")

	file, err := os.OpenFile(configPath, os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(&c)
}

func ReadConfig() (c *Config, err error) {
	file, err := os.OpenFile(getDefaultConfigPath(), os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getDefaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".inf", "config")
}
