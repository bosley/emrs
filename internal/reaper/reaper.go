/*
  This nerv module is meant to be used as a signal watchdog that is used to
  dictate the life of the program.

  When a SIGINT is raised, the reaper begins a timed shutdown sequence, where
  it will message a countdown over the event engine on its configured topic.

*/

package reaper

import (
	"github.com/bosley/nerv-go"
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
	WaitGroup    *sync.WaitGroup
	ShutdownSecs int
}

type Reaper struct {
	waitGroup    *sync.WaitGroup
	shutdownSecs int
	submitter    *nerv.ModuleSubmitter
}

// Create a reaper pointer, which is a valid nerv.Module
func New(cfg Config) *Reaper {
	return &Reaper{
		waitGroup:    cfg.WaitGroup,
		shutdownSecs: cfg.ShutdownSecs,
	}
}

// Trigger reaper shutdown from anywhere if running
func Interrupt() {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func (r *Reaper) Start() error {
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

func (r *Reaper) SetSubmitter(submitter *nerv.ModuleSubmitter) {
	r.submitter = submitter
}

func (r *Reaper) sigShutdown() {
	t := r.shutdownSecs
	for t != 0 {
		r.submitter.SubmitData(&ReaperMsg{
			SecondsRemaining: t,
		})
		t -= 1
		time.Sleep(1 * time.Second)
	}
}

func (r *Reaper) sigKill() {
	os.Exit(ReturnCodeForceKill)
}
