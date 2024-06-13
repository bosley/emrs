package main

import (
	"github.com/bosley/nerv-go"
	"internal/reaper"
	"log/slog"
	"os"
	"sync"
)

const (
	appDebugExtra = true

	defaultAppGracefulShutdownSecs = 5
)

type AppConfig struct {
	Reaper reaper.Config
}

func must(e error) {
	if e != nil {
		slog.Error(e.Error())
		os.Exit(-1)
	}
}

func main() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

	appWaitGroup := new(sync.WaitGroup)

	appConfig := &AppConfig{
		Reaper: reaper.Config{
			WaitGroup:    appWaitGroup,
			ShutdownSecs: defaultAppGracefulShutdownSecs,
		},
	}

	eventEngine := CreateEngine()

	PopulateModules(eventEngine, appConfig)

	must(eventEngine.Start())

	appWaitGroup.Wait()

	must(eventEngine.Stop())
}

func CreateEngine() *nerv.Engine {

	engine := nerv.NewEngine()

	if !appDebugExtra {
		return engine
	}

	createCallback := func(id string) nerv.EventRecvr {
		return func(event *nerv.Event) {
			slog.Debug(
				"nerv engine cb",
				"cb id", id,
				"topic", event.Topic,
				"prod", event.Producer,
				"spawned", event.Spawned)
		}
	}

	return engine.WithCallbacks(
		nerv.EngineCallbacks{
			RegisterCb: createCallback("registration"),
			NewTopicCb: createCallback("new_topic"),
			ConsumeCb:  createCallback("consumed"),
			SubmitCb:   createCallback("submission"),
		})
}
