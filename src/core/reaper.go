package core

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type trigger func()

func initReaperIntercept(wg *sync.WaitGroup, countDown time.Duration, alertCb func()) trigger {

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	var watch atomic.Bool

	watch.Store(false)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			alertCb()
			time.Sleep(countDown)
			return
		case syscall.SIGTERM:
			os.Exit(-1)
		default:
			return
		}
	}()

	return func() {
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
}
