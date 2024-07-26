package main

import (
	"encoding/json"
	"fmt"
	"github.com/bosley/emrs/badger"
	"github.com/bosley/emrs/datastore"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultServerName        = "EMRS Server"
	defaultEnvHome           = "EMRS_HOME"
	defaultBinding           = "localhost:8080"
	defaultStoragePath       = "storage"
	defaultConfigName        = "server.cfg"
	defaultUiKeyDuration     = "8760h" // 1 year
	defaultUserGivenDuration = "4320h" // ~6 months

	defaultActionsDir = "actions"
)

const (
	ttlEphemeralVoucher = "30s"
)

type Config struct {
	Binding  string            `yaml:binding`
	Key      string            `yaml:key`
	Cert     string            `yaml:cert`
	Identity string            `yaml:identity`
	Actions  map[string]string `yaml:actions`
}

func mustFindHome(potential string) string {
	if potential == "" {
		fromEnv := os.Getenv(defaultEnvHome)
		if fromEnv == "" {
			slog.Error("unable to determine emrs home directory from environment")
			os.Exit(1)
		}
		return fromEnv
	}
	return potential
}

func mustLoadCfgAndBadge(home string) (Config, badger.Badge) {
	cfg := getConfig(home)
	badge, err := badger.DecodeIdentityString(cfg.Identity)
	if err != nil {
		slog.Error("badger failed to decode server identity", "error", err.Error())
		os.Exit(1)
	}
	return cfg, badge
}

func mustLoadDataStore(dir string) datastore.DataStore {

	dataStrj, err := datastore.Load(dir)
	if err != nil {
		slog.Error("failed to load datastore", "error", err.Error())
		os.Exit(1)
	}
	return dataStrj
}

func mustLoadDefaultDataStore(home string) datastore.DataStore {
	return mustLoadDataStore(
		filepath.Join(home, defaultStoragePath))
}

func getConfig(home string) Config {
	var config Config
	target, err := os.ReadFile(filepath.Join(home, defaultConfigName))
	if err != nil {
		slog.Error("failed to load config", "error", err.Error())
		os.Exit(1)
	}
	err = yaml.Unmarshal(target, &config)
	if err != nil {
		slog.Error("failed to load config", "error", err.Error())
		os.Exit(1)
	}
	slog.Debug("loaded config", "binding", config.Binding, "key", config.Key, "cert", config.Cert)
	return config
}

func generateVouchers(badge badger.Badge, n int, durr time.Duration) {
	vouchers := make([]string, n)
	for i := range n {
		voucher, err := badger.NewVoucher(badge, durr)
		if err != nil {
			slog.Error("failed to generate vouchers")
			os.Exit(1)
		}
		vouchers[i] = voucher
	}

	b, _ := json.Marshal(vouchers)
	fmt.Println(string(b))
}
