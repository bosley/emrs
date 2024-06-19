package main

import (
	"flag"
	"fmt"
	"github.com/bosley/nerv-go"
	"internal/webui"
	"log/slog"
	"os"
	"reaper"
	"sync"
)

const (
	defaultAppUser     = "admin"
	defaultAppPassword = "admin"
)

type App struct {
	wg     *sync.WaitGroup
	engine *nerv.Engine
	kill   reaper.Trigger
	ui     *webui.WebUi
}

func main() {

	tempLoggedInUserId := "UUID-DEV" // TODO: Remove this. This is for dev auth sys, before DB setup

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	releaseMode := flag.Bool("release", false, "Turn on debug mode")

	// TODO: NOTE:
	// Until we get vaults running and databases working we will use simple auth setup so we can
	// get development underway but still have auth framed-in
	username := flag.String("user", defaultAppUser, "Username to log in with")
	password := flag.String("pass", defaultAppPassword, "Password to require for login")

	configPath := flag.String("config", "emrs.yaml", "Server config YAML")

	flag.Parse()

	uiMode := webui.ModeDebug
	if *releaseMode {
		uiMode = webui.ModeRelease
		slog.SetDefault(
			slog.New(
				slog.NewTextHandler(os.Stdout,
					&slog.HandlerOptions{
						Level: slog.LevelWarn,
					})))
	}

	sc, err := ReadServerConfig(*configPath)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(-1)
	}

	webAddress := fmt.Sprintf("%s:%s", sc.WebUi.Host, sc.WebUi.Port)

	slog.Info("using", "host", sc.WebUi.Host, "port", sc.WebUi.Port, "full", webAddress)

	appEngine := nerv.NewEngine()

	wg := new(sync.WaitGroup)

	trigger, err := reaper.Spawn(&reaper.Config{
		Name:   sc.Reaper.Name,
		Engine: appEngine,
		Grace:  sc.Reaper.Grace,
		Wg:     wg,
	})

	if err != nil {
		slog.Error(err.Error())
		os.Exit(-1)
	}

	ui := webui.New(webui.Config{
		Engine:      appEngine,
		Address:     webAddress,
		Mode:        uiMode,
		KillChannel: sc.Reaper.Name,
		AuthenticateUser: func(user string, pass string) *string {

			// TODO: Actually check a vault for this pass, and
			//       return the user's UUID if good
			if user == *username && pass == *password {
				return &tempLoggedInUserId
			}
			return nil
		},
	})

	app := &App{
		engine: appEngine,
		wg:     wg,
		kill:   trigger,
		ui:     ui,
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

	must(app.engine.Start())
	must(app.ui.Start())

	app.wg.Wait()

	must(app.ui.Stop())
	must(app.engine.Stop())
}
