package main

import (
	"crypto/tls"
	"emrs/badger"
	"emrs/core"
	"encoding/json"
	"os"
)

type RuntimeInfo struct {
	Mode string `json:mode`
}

type HostingInfo struct {
	ApiAddress string `json:api_address`
	Key        string `json:https_key`
	Cert       string `json:https_cert`
}

type Config struct {
	Runtime  RuntimeInfo `json:runtime`
	Hosting  HostingInfo `json:hosting`
	EmrsCore core.Config `json:core` // Consider having this be a byte array, and b64 encoding the identity before saving/ decoding before handing to core
}

func (c *Config) LoadTLSCert() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(c.Hosting.Cert, c.Hosting.Key)
}

func (c *Config) WriteTo(path string) error {
	encoded, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, encoded, 0644)
}

func CreateConfigTemplate() *Config {
	badge, err := badger.New(
		badger.Config{
			Nickname: "EmrsServer",
		})
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		Runtime: RuntimeInfo{
			Mode: "debug",
		},
		Hosting: HostingInfo{
			ApiAddress: "localhost:20000",
			Key:        "./dev/keys/server.key",
			Cert:       "./dev/keys/server.crt",
		},
		EmrsCore: core.Config{
			Identity: badge.EncodeIdentityString(),
			Network:  core.BlankTopo(),
		},
	}
}
