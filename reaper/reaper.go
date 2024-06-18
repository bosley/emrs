package reaper

import (
	"github.com/bosley/nerv-go"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Config struct {
	Name   string
	Engine *nerv.Engine
	Grace  int
	Wg     *sync.WaitGroup
}

type Trigger func()

func Spawn(cfg *Config) (Trigger, error) {

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	var watch atomic.Bool

	watch.Store(false)

	producer, err := cfg.Engine.AddRoute(cfg.Name, func(c *nerv.Context) {

		if watch.Load() {
			return
		}

		watch.Store(true)

		slog.Warn("kill requested over nerv", "from", c.Event.Producer)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	})

	if err != nil {
		return nil, err
	}

	performShutdownBroadcast := func() {

		watch.Store(true)

		producer(nil)

		t := cfg.Grace
		for t != 0 {
			slog.Warn("shutdown alert", "t", t)
			t -= 1
			time.Sleep(1 * time.Second)
		}
	}

	cfg.Wg.Add(1)
	go func() {
		defer cfg.Wg.Done()
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			performShutdownBroadcast()
			return
		case syscall.SIGTERM:
			slog.Warn("REAPER: SIGTERM")
			os.Exit(-1)
		default:
			return
		}
	}()

	return func() {
		producer(nil)
	}, nil
}
