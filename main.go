package main

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Mode      string `yaml:mode`
	Hostname  string `yaml:home`
	Port      int    `yaml:port`
	Key       string `yaml:key`
	Cert      string `yaml:cert`
	Datastore string `yaml:datastore`
}

func MustLoadConfig(path string) Config {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	var c Config
	if err := yaml.Unmarshal(f, &c); err != nil {
		panic(err.Error())
	}
	return c
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

func (c *Config) LoadTlsCert() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(c.Cert, c.Key)
}

// --

func main() {

	fmt.Println("yo")
}
