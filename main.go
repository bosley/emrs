package main

import (
	"github.com/bosley/nerv-go"
	"internal/reaper"
	"internal/vault"
	"log/slog"
	"os"
	"sync"
)

const (
	defaultAppGracefulShutdownSecs = 2
	defaultAppVaultPath            = ".emrs.vault.db"
)

type AppConfig struct {
	Vault  *vault.Config
	Reaper reaper.Config
}

type App struct {
	wg     *sync.WaitGroup
	engine *nerv.Engine
	config *AppConfig
	vault  *vault.Vault
	webUi  *WebUi
}

func main() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	app := &App{
		webUi: CreateWebUi(
			defaultWebUiMode,
			defaultWebUiAddr,
		),
		config: &AppConfig{
			Vault: &vault.Config{
				DbPath: defaultAppVaultPath,
			},
			Reaper: reaper.Config{
				WaitGroup:    new(sync.WaitGroup),
				ShutdownSecs: defaultAppGracefulShutdownSecs,
			},
		},
	}

	slog.Warn("TODO: NEED TO CONVERT WEB UI TO CONFORM TO NERV MODULE")

	app.Exec()
}

func must(e error) {
	if e != nil {
		slog.Error(e.Error())
		os.Exit(-1)
	}
}

func (app *App) Exec() {
	slog.Debug("app:run")

	app.createEngine()

	PopulateModules(app.engine, app.config)

	must(app.engine.Start())

	must(app.webUi.Start())

	app.config.Reaper.WaitGroup.Wait()

	must(app.engine.Stop())

	if app.vault != nil {
		app.vault.Stop()
	}
}

func (app *App) createEngine() {

	app.engine = nerv.NewEngine()

	if app.config.Vault != nil {
		app.setupVault()
	}
}

func (app *App) setupVault() {
	slog.Debug("vault config detected, setting up..")

	app.vault = vault.New()

	app.engine = app.engine.WithCallbacks(
		nerv.EngineCallbacks{
			RegisterCb: app.vault.Store,
			NewTopicCb: app.vault.Store,
			ConsumeCb:  app.vault.Store,
			SubmitCb:   app.vault.Store,
		})
}
