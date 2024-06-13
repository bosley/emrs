package main

import (
	"internal/reaper"
	"internal/vault"
	"log/slog"
	"os"
	"sync"
)

const (
	appDebugExtra                  = true
	defaultAppGracefulShutdownSecs = 5
)

func main() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	app := &App{
		config: &AppConfig{
			Vault: &vault.Config{
				DbPath: "/tmp/emrs.db",
			},
			Reaper: reaper.Config{
				WaitGroup:    new(sync.WaitGroup),
				ShutdownSecs: defaultAppGracefulShutdownSecs,
			},
		},
	}

	app.Exec()
}
