package vault

import (
	"context"
	"github.com/bosley/nerv-go"
	"log/slog"
	"sync"
)

type Vault struct {
	input   chan *nerv.Event
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	running bool
}

type Config struct {
	DbPath string
}

func (v *Vault) Store(event *nerv.Event) {
	slog.Debug("vault:Store", "spawned", event.Spawned, "topic", event.Topic, "producer", event.Producer)

	v.input <- event
}

func New() *Vault {
	vault := &Vault{
		input:   make(chan *nerv.Event),
		wg:      new(sync.WaitGroup),
		running: false,
	}
	vault.ctx, vault.cancel = context.WithCancel(context.Background())

	return vault
}

func (v *Vault) Stop() {
	if !v.running {
		return
	}
	v.cancel()
	close(v.input)
	v.wg.Wait()
}

func (v *Vault) Start() {
	if v.running {
		panic("vault already started")
	}
	v.wg.Add(1)
	go func() {
		defer func() {
			v.wg.Done()
			v.running = false
		}()

		for {
			select {
			case <-v.ctx.Done():
				return
			case event := <-v.input:
				v.storeEvent(event)
				return
			}
		}
	}()
	v.running = true
}

func (v *Vault) storeEvent(event *nerv.Event) {

	slog.Debug("vault:storeEvent - TODO")
}
