package core

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
)

type proc struct {
	signal  *Signal
	action  *Action
	ctx     context.Context
	stream  chan Event
	running atomic.Bool
}

func newProc(signal *Signal, action *Action, ctx context.Context, stream chan Event) *proc {
	return &proc{
		signal: signal,
		action: action,
		ctx:    ctx,
		stream: stream,
	}
}

func (p *proc) exec(wg *sync.WaitGroup) {
	if p.running.Load() {
		panic("attempt to run a running runner")
	}
	defer p.running.Store(true)
	slog.Debug("starting processor",
		"signal", p.signal.Header.Name,
		"action", p.action.Header.Name)

	wg.Add(1)
	go func() {
		for {
			select {
			case <-p.ctx.Done():
				wg.Done()
				p.running.Store(false)
				slog.Debug("processor complete",
					"signal", p.signal.Header.Name,
					"action", p.action.Header.Name)
				return
			case event := <-p.stream:
				p.handle(event)
			}
		}
	}()
}

func (p *proc) handle(event Event) {
	slog.Error("NOT YET COMPLETE")
	slog.Info("Handle Event", "assignment", p.signal.Header.Name)

	// Need to determine if we need to re/load teh source file
	// for execution based on the actionType
}
