package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type WebServerCfg struct {
	Host       string `yaml:host`
	Port       string `yaml:port`
	CertPath   string `yaml:certFile`
	KeyPath    string `yaml:keyFile`
	CACertPath string `yaml:caCertPath`
  LocalUser  string `yaml:localUser`
  LocalPass  string `yaml:localPass`
}

type ReaperServerCfg struct {
	Name  string `yaml:name`
	Grace int    `yaml:grace`
}

type ServerCfg struct {
	WebUi  WebServerCfg    `yaml:webui`
	Reaper ReaperServerCfg `yaml:reaper`
}

func ReadServerConfig(path string) (*ServerCfg, error) {

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sc ServerCfg
	if err := yaml.Unmarshal(f, &sc); err != nil {
		return nil, err
	}

	return &sc, nil
}
