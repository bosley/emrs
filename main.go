package main

import (
	"github.com/bosley/nerv-go"
	"internal/reaper"
	"internal/vault"
	"internal/webui"
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
	WebUi  webui.Config
	Reaper reaper.Config
}

type App struct {
	wg     *sync.WaitGroup
	engine *nerv.Engine
	config *AppConfig
	vault  *vault.Vault
}

func main() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	// These will be CLI args/ config files
	webUiAddr := webui.DefaultWebUiAddr
	webUiMode := webui.DefaultWebUiMode
	vaultDbPath := defaultAppVaultPath
	reaperGraceShutdown := defaultAppGracefulShutdownSecs

	appEngine := nerv.NewEngine()

	app := &App{
		engine: appEngine,
		config: &AppConfig{
			Vault: &vault.Config{
				DbPath: vaultDbPath,
			},
			WebUi: webui.Config{
				Engine:  appEngine,
				Address: webUiAddr,
				Mode:    webUiMode,
			},
			Reaper: reaper.Config{
				WaitGroup:    new(sync.WaitGroup),
				ShutdownSecs: reaperGraceShutdown,
			},
		},
	}

	app.Exec()
}

func must(e error) {
	if e != nil {
		slog.Error(e.Error())
		os.Exit(-1)
	}
}

func (app *App) Exec() {

	app.setupVault()

	PopulateModules(app.engine, app.config)

	must(app.engine.Start())

	app.config.Reaper.WaitGroup.Wait()

	must(app.engine.Stop())

	if app.vault != nil {
		app.vault.Stop()
	}
}

func (app *App) setupVault() {

	if app.config.Vault == nil {
		slog.Debug("no vault configuration detected - skipping")
		return
	}

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
