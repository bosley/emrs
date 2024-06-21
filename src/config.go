package main

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Home string `yaml:home`
	Hostname      string `yaml:hostname`
	Port          int    `yaml:port`
	Key  string `yaml:key`
	Cert string `yaml:cert`
	Datastore string `yaml:datastore`
}

func LoadConfig(path string) (*Config, error) {

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sc Config
	if err := yaml.Unmarshal(f, &sc); err != nil {
		return nil, err
	}
	return &sc, nil
}

func (c *Config) GetDabasePath() string {
	return filepath.Join(
		c.Home,
		c.Datastore)
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

func (c *Config) LoadTlsCert() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(c.Cert, c.Key)
}
