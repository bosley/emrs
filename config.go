package main

import (
	"crypto/tls"
	"emrs/badger"
	"emrs/core"
	"encoding/json"
	"log/slog"
	"os"
	"time"
)

type RuntimeInfo struct {
	Mode string `json:mode`
}

type HostingInfo struct {
	ApiAddress string   `json:api_address`
	ApiKeys    []string `json:api_keys`
	Key        string   `json:https_key`
	Cert       string   `json:https_cert`
}

type Config struct {
	Runtime  RuntimeInfo `json:runtime`
	Hosting  HostingInfo `json:hosting`
	EmrsCore core.Config `json:core` // Consider having this be a byte array, and b64 encoding the identity before saving/ decoding before handing to core
}

func (c *Config) Validate() error {

	slog.Warn("NEED TO VALIDATE CONFIG => ENSURE THAT ALL TOKENS IN HOSTING INFO BELONG TO, AND CAN BE AUTHED BY, THE EXISTING IDENTITY")
	return nil

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
	apiKeys, err := generateApiKeys(5, badge)
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		Runtime: RuntimeInfo{
			Mode: "debug",
		},
		Hosting: HostingInfo{
			ApiAddress: "localhost:20000",
			ApiKeys:    apiKeys,
			Key:        "./dev/keys/server.key",
			Cert:       "./dev/keys/server.crt",
		},
		EmrsCore: core.Config{
			Identity: badge.EncodeIdentityString(),
			Network:  core.BlankTopo(),
		},
	}
}

func generateApiKeys(n int, badge badger.Badge) ([]string, error) {
	if n <= 0 || n > 100 {
		panic("n was set to a bad number when generating api keys")
	}

	var err error
	result := make([]string, n)

	for i := range n {
		result[i], err = badger.NewVoucher(
			badge,
			time.Hour*24*365)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}
