package core

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"strings"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSel[V any](source []V) V {
	return source[rand.Intn(len(source))]
}

func randSeq[V any](source []V, n int, joiner func([]V) V) V {
	v := make([]V, n)
	for i := range v {
		v[i] = randSel[V](source)
	}
	return joiner(v)
}

func generateList[V any](size int, generator func(idx int) V) []V {
	result := make([]V, size)
	for i := 0; i < size; i++ {
		result[i] = generator(i)
	}
	return result
}

func randString(n int) string {
	return randSeq[string](func() []string {
		s := make([]string, len(letters))
		forEach(letters, func(i int, r rune) error {
			s[i] = string(r)
			return nil
		})
		return s
	}(), n, func(x []string) string {
		return strings.Join(x, "")
	})
}

func randHeader() HeaderData {
	return HeaderData{
		Name:        randString(32),
		Description: randString(64),
		Tags: generateList[string](16, func(idx int) string {
			return randString(10)
		}),
	}
}

func randAsset() *Asset {
	return &Asset{
		Header: randHeader(),
	}
}

func randAction() *Action {
	return &Action{
		Header:        randHeader(),
		ExecutionData: make([]byte, 0),
	}
}

func randSignal() *Signal {
	return &Signal{
		Header: randHeader(),
		Trigger: randSel[string]([]string{
			TriggerOnEvent,
			TriggerOnTimeout,
			TriggerOnBumpTimeout,
			TriggerOnShutdownNotify,
			TriggerOnSchedule,
			TriggerOnEmit,
		}),
	}
}

func randSector() *Sector {
	return &Sector{
		Header: randHeader(),
		Assets: generateList[*Asset](4, func(i int) *Asset { return randAsset() }),
	}
}

func randTopoMap(nSectors int, nSignals int, nActions int, nMapped int) Topo {

	sectors := generateList[*Sector](nSectors, func(idx int) *Sector {
		return randSector()
	})

	signals := generateList[*Signal](nSignals, func(idx int) *Signal {
		return randSignal()
	})

	actions := generateList[*Action](nActions, func(idx int) *Action {
		return randAction()
	})

	sigmap := make(map[string]string)

	// TODO: MAP SOME SIGNALS TO SOME ACTIONS AT RANDOM BY nMAPPED

	return Topo{
		Sectors: sectors,
		Signals: signals,
		Actions: actions,
		SigMap:  sigmap,
	}
}

func randTopoSmall() Topo {
	return randTopoMap(1, 2, 1, 0)
}

func randTopoLarge() Topo {
	return randTopoMap(10, 100, 244, 50)
}

func randTopoDuplicateSectorNames() Topo {
	s := randSector()
	r := randTopoLarge()
	iterate(r.Sectors, func(it Iter[*Sector]) error {
		r.Sectors[it.Idx] = s
		return nil
	})
	return r
}

func TestCoreTopo(t *testing.T) {

	rand.Seed(time.Now().UnixNano())

	iterate(
		[]Topo{
			randTopoSmall(),
			randTopoLarge(),
		},
		func(it Iter[Topo]) error {
			_, err := NetworkMapFromTopo(it.Value)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			return nil
		})

	fmt.Println("hokay!")
}

func TestCoreBadTopoSector(t *testing.T) {

	tmap := randTopoDuplicateSectorNames()
	_, err := NetworkMapFromTopo(tmap)
	if err == nil {
		t.Fatal("Expected error for duplicate sector names")
	}
}

func TestCoreNetworkMap(t *testing.T) {

	assert := func(msg string, cond bool) {
		if !cond {
			t.Fatalf("test failure: %s", msg)
		}
	}

	nm := BlankNetworkMap()

	assert("nil map", nm != nil)

  for z := 0; z < 10; z++ {
	  s := randSector()

	  assert("failed to add sector", nm.AddSector(s) == nil)

	  for i := 0; i < 10; i++ {
	  	a := randAsset()
	  	assert("failed to add asset", nm.AddAsset(s.Header.Name, a) == nil)
	  	assetfull := makeAssetFullName(s.Header.Name, a.Header.Name)
	  	sig := makeAssetOnEventSignal(assetfull)
	  	_, ok := nm.signals[sig.Header.Name]
	  	assert("(onEvent) not generated for new asset", ok)
	  }

	  xasset := randAsset()
	  xassetfull := makeAssetFullName(s.Header.Name, xasset.Header.Name)
	  xsig := makeAssetOnEventSignal(xassetfull)

	  assert("failed to add asset", nm.AddAsset(s.Header.Name, xasset) == nil)
    assert("contains failure", nm.ContainsAsset(s.Header.Name, xasset.Header.Name))
	  assert("added dup asset", nm.AddAsset(s.Header.Name, xasset) != nil)
    assert("no signal made for new asset", nm.ContainsSignal(xsig.Header.Name))
    nm.DeleteAsset(s.Header.Name, xasset.Header.Name)
    assert("contains failure", nm.ContainsAsset(s.Header.Name, xasset.Header.Name) == false)
    assert("asset's onEvent signal not deleted", nm.ContainsSignal(xsig.Header.Name) == false)
	  assert("failed to add asset", nm.AddAsset(s.Header.Name, xasset) == nil)
    assert("contains failure", nm.ContainsAsset(s.Header.Name, xasset.Header.Name))

	  _, ok := nm.signals[xsig.Header.Name]
	  assert("(onEvent) not generated for new asset", ok)

    signals := generateList[*Signal](20, func(i int) *Signal{ return randSignal(); })
    iterate(signals, func (it Iter[*Signal]) error {
      assert("generated signal unable to be added", nm.AddSignal(it.Value) == nil)
      assert("generated signal not found", nm.ContainsSignal(it.Value.Header.Name))
      assert("dup signal added", nm.AddSignal(it.Value) != nil)
      nm.DeleteSignal(it.Value.Header.Name)
      assert("deleted signal found", nm.ContainsSignal(it.Value.Header.Name) == false)
      assert("generated signal unable to be re-added", nm.AddSignal(it.Value) == nil)
      return nil
    })

    actions := generateList[*Action](20, func(i int) *Action{ return randAction(); })
    iterate(actions, func (it Iter[*Action]) error {
      assert("generated action unable to be added", nm.AddAction(it.Value) == nil)
      assert("generated action not found", nm.ContainsAction(it.Value.Header.Name))
      assert("dup action added", nm.AddAction(it.Value) != nil)
      nm.DeleteAction(it.Value.Header.Name)
      assert("deleted action found", nm.ContainsAction(it.Value.Header.Name) == false)
      assert("generated action unable to be re-added", nm.AddAction(it.Value) == nil)
      return nil
    })
  }
}

// TODO: A test where we generate a numnch of random maps, dump them to a topo, then write to disk.
//        Once written, all topos are re-loaded and convered back into maps where the original 
//        random maps are checked against the loaded maps to check for changes
