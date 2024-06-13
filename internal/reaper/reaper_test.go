package reaper

import (
  "os"
  "fmt"
  "sync"
  "time"
  "syscall"
  "testing"
  "log/slog"
	"github.com/bosley/nerv-go"
)

func TestReaper(t *testing.T) {

	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				})))

  shutdownSec   := 5
  shutdownAlertRecv := 0
  shutdownTopic := "emrs.internal.shutdown"

  wg := new(sync.WaitGroup)

  config := Config {
    WaitGroup:    wg,
    ShutdownSecs: shutdownSec,
  }

  reaper := New(config)

  engine := nerv.NewEngine()

	topic := nerv.NewTopic(shutdownTopic).
		UsingDirect().
		UsingRoundRobinSelection()

	consumers := []nerv.Consumer{
		nerv.Consumer{
			Id: "emrs.internal.watchdog.reaper",
      Fn: func (event *nerv.Event) {
        slog.Debug("shutdown event", "sec-remaining", event.Data.(*ReaperMsg).SecondsRemaining)
        shutdownAlertRecv += 1
      },
    },
	}

	if err := engine.UseModule(
		reaper,
		topic,
		consumers); err != nil {
		t.Fatalf("err:%v", err)
	}

	fmt.Println("starting engine")
	if err := engine.Start(); err != nil {
		t.Fatalf("err: %v", err)
	}

  time.Sleep(1 * time.Second)

  fmt.Println("sending SIGINT")

  syscall.Kill(syscall.Getpid(), syscall.SIGINT)

  // No _need_ to wait here, Stop() fn of reaper waits on the
  // group when the module is shutdown

	fmt.Println("stopping engine")
	if err := engine.Stop(); err != nil {
		t.Fatalf("err: %v", err)
	}

	fmt.Println("[ENGINE STOPPED]")

  if shutdownAlertRecv != shutdownSec {
    t.Fatal("failed to receive countdown timer")
  }
}
