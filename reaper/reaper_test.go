package reaper

import (
	"fmt"
	"github.com/bosley/nerv-go"
	"sync"
	"testing"
)

func TestReaper(t *testing.T) {

	check := func(err error) {
		if err != nil {
			t.Fatalf("err:%v", err)
		}
	}

	engine := nerv.NewEngine()

	testRouteName := "test.route"

	wg := new(sync.WaitGroup)

	cfg := Config{
		Name:   testRouteName,
		Engine: engine,
		Grace:  1,
		Wg:     wg,
	}

	trigger, err := Spawn(&cfg)

	check(err)

	type listener struct {
		consumer nerv.Consumer
		hit      bool
	}

	nListeners := 65535
	listeners := make([]*listener, nListeners)

	buildConsumer := func(n int) nerv.Consumer {
		return nerv.Consumer{
			Id: fmt.Sprintf("consumer-%d", n),
			Fn: func(event *nerv.Event) {
				listeners[n].hit = true
			},
		}
	}

	for n, _ := range listeners {
		consumer := buildConsumer(n)
		listeners[n] = &listener{
			consumer: consumer,
			hit:      false,
		}

		engine.Register(consumer)
		check(engine.SubscribeTo(testRouteName, consumer.Id))
	}

	check(engine.Start())

	trigger()

	fmt.Println("\ntriggered")

	wg.Wait()

	check(engine.Stop())

	for n, listener := range listeners {
		if !listener.hit {
			t.Fatalf("listner %d did not receive the shutdown warning", n)
		}
	}
}
