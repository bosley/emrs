package main

import (
	"flag"
	"github.com/bosley/nerv-go"
	"internal/reaper"
	"internal/webui"
	"log/slog"
	"os"
	"sync"
  "path/filepath"
)

const (
	defaultAppGracefulShutdownSecs = 5
	defaultAppUser                 = "admin"
	defaultAppPassword             = "admin"
  defaultEmrsInstanceDirName     = ".emrs"
)

type AppConfig struct {
	WebUi  webui.Config
	Reaper reaper.Config
}

type App struct {
	wg     *sync.WaitGroup
	engine *nerv.Engine
	config *AppConfig
}

func main() {

	tempLoggedInUserId := "UUID-DEV" // TODO: Remove this. This is for dev auth sys, before DB setup

  emrsInstance := filepath.Join(mustGet(os.UserHomeDir()).(string), defaultEmrsInstanceDirName)


	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	setupTarget := flag.String("emrs", emrsInstance, "Create an EMRS instance. takes <DIR>")
  createTarget := flag.Bool("new", false, "Create target defined by -emrs")

	webUiAddr := flag.String("addr", webui.DefaultWebUiAddr, "Address to bind Web UI to [address:port]")
	releaseMode := flag.Bool("release", false, "Turn on debug mode")
	gracefulSecs := flag.Int("grace", defaultAppGracefulShutdownSecs, "Graceful shutdown time (seconds)")

	// TODO: NOTE:
	// Until we get vaults running and databases working we will use simple auth setup so we can
	// get development underway but still have auth framed-in
	username := flag.String("user", defaultAppUser, "Username to log in with")
	password := flag.String("pass", defaultAppPassword, "Password to require for login")
	flag.Parse()

  if *createTarget {
    RunSetupUI(*setupTarget)
    os.Exit(0)
  }

	appEngine := nerv.NewEngine()

	uiCfg := webui.Config{
		Engine:  appEngine,
		Address: *webUiAddr,
		Mode:    webui.DefaultWebUiMode,
		AuthenticateUser: func(user string, pass string) *string {

			// TODO: Actually check a vault for this pass, and
			//       return the user's UUID if good
			if user == *username && pass == *password {
				return &tempLoggedInUserId
			}
			return nil
		},
	}

	reaperCfg := reaper.Config{
		WaitGroup:    new(sync.WaitGroup),
		ShutdownSecs: *gracefulSecs,
	}

	appCfg := AppConfig{
		WebUi:  uiCfg,
		Reaper: reaperCfg,
	}

	if *releaseMode {
		configureReleaseMode(&appCfg)
	}

	app := &App{
		engine: appEngine,
		config: &appCfg,
	}

	app.Exec()
}

func must(e error) {
	if e != nil {
		slog.Error(e.Error())
		os.Exit(-1)
	}
}

func mustGet(t any, e error) any {
  if e != nil {
    slog.Error(e.Error())
    os.Exit(-1)
  }
  return t
}

func configureReleaseMode(cfg *AppConfig) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelWarn,
				})))

	cfg.WebUi.Mode = webui.ModeRelease
}

func (app *App) Exec() {

	PopulateModules(app.engine, app.config)

	must(app.engine.Start())

	app.config.Reaper.WaitGroup.Wait()

	must(app.engine.Stop())
}
