/*
  This nerv module is meant to be used as a signal watchdog that is used to
  dictate the life of the program.

  When a SIGINT is raised, the reaper begins a timed shutdown sequence, where
  it will message a countdown over the event engine on its configured topic.

*/

package reaper

import (
	"github.com/bosley/nerv-go"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	ReturnCodeForceKill = 42
)

// Message type sent on event-to-topic from module
type ReaperMsg struct {
	SecondsRemaining int
}

// Configuration for the reaper module;
// The waitgroup can be used to wait and ensure that the
type Config struct {
	WaitGroup       *sync.WaitGroup
	ShutdownSecs    int
	DesignatedTopic string
}

// NOTE: We could add a `func(int) bool` to this
//
//	as a simple way to request more time
type Listener func(remainingSecs int)

type Reaper struct {
	waitGroup    *sync.WaitGroup
	shutdownSecs int
	pane         *nerv.ModulePane
	listeners    []Listener
	mu           sync.Mutex
	topic        string
}

// Create a reaper pointer, which is a valid nerv.Module
func New(cfg Config) *Reaper {
	return &Reaper{
		waitGroup:    cfg.WaitGroup,
		shutdownSecs: cfg.ShutdownSecs,
		listeners:    make([]Listener, 0),
		topic:        cfg.DesignatedTopic,
	}
}

// Add a listener function that cares about pending death
func (r *Reaper) AddListener(target Listener) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listeners = append(r.listeners, target)
}

// Trigger reaper shutdown from anywhere if running
func Interrupt() {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func (r *Reaper) GetName() string {
	return "mod.reaper"
}

func (r *Reaper) Start() error {
	slog.Info("reaper:Start")

	if err := r.pane.SubscribeTo(r.topic, []nerv.Consumer{
		nerv.Consumer{
			Id: r.GetName(),
			Fn: r.killCmd,
		}},
		true); err != nil {
		return err
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	r.waitGroup.Add(1)
	go func() {
		defer r.waitGroup.Done()

		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			r.sigShutdown()
		case syscall.SIGTERM:
			r.sigKill()
		default:
			return
		}
	}()
	return nil
}

func (r *Reaper) Shutdown() {
	Interrupt()
	r.waitGroup.Wait()
}

func (r *Reaper) RecvModulePane(pane *nerv.ModulePane) {
	r.pane = pane
}

func (r *Reaper) sigShutdown() {
	t := r.shutdownSecs
	for t != 0 {
		r.pane.SubmitTo(
			r.topic,
			&ReaperMsg{
				SecondsRemaining: t,
			},
		)
		t -= 1
		time.Sleep(1 * time.Second)
	}
}

func (r *Reaper) sigKill() {
	os.Exit(ReturnCodeForceKill)
}

// The reaper is set to receive its own kill command
// as we don't need to inundate the engine with broadcasts
// to every module, though, we can. Instead, we receive the
// same message we sent and THEN directly call listeners.
//
// This also allows the reaper to indicate shutdown to non-nerv
// related services that may or may not be running somewhere
func (r *Reaper) killCmd(event *nerv.Event) {
	remaining := event.Data.(*ReaperMsg).SecondsRemaining

	slog.Debug("shutdown imminent", "seconds", remaining)

	for _, fn := range r.listeners {
		go fn(remaining)
	}
}
