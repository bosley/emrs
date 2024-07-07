/*
  This file describes the structures that we use to represent items on the network.
  We have "Topo" Which is a flat representation meant to be dumped to JSON, and then
  there is NetworkMap, the structure that utlizes the representation for runtime operations.

  A NetworkMap can be generated from any valid Topo representation

  Validity Requirements:
    All sector names unique
      All assets inside sector are unique to the sector
    All signals and actions must have names unique to their repspective categories,
    All signal->action-list pairs mapped in sigmap in a Topo must have values
      that correspond to defined signals and actions. (The signals and actions that
        are pointed to must have been defined in the signals and actions lists)
    All action_list items in sigmap are unique within their list (no duplicates)



*/

package core

import (
	"fmt"
	"log/slog"
)

const (
	TriggerOnEvent          = "onEvent"
	TriggerOnTimeout        = "onTimeout"
	TriggerOnBumpTimeout    = "onBumpTimeout"
	TriggerOnShutdownNotify = "onShutdown"
	TriggerOnSchedule       = "onSchedule"
	TriggerOnEmit           = "onEmit"
)

const (
	ExecutionTypeFile     = "file"
	ExecutionTypeEmbedded = "embedded"
)

type HeaderData struct {
	Name        string   `json:name`
	Description string   `json:description`
	Tags        []string `json:tags`
}

type Sector struct {
	Header HeaderData `json:header`
	Assets []*Asset   `json:assets`
}

type Asset struct {
	Header HeaderData `json:header`
}

type Signal struct {
	Header  HeaderData `json:header`
	Trigger string     `json:trigger`
}

type Action struct {
	Header HeaderData `json:header`
	Type   string     `json:type`
	Info   string     `json:info`
}

type Topo struct {
	Sectors []*Sector           `json:sectors`
	Signals []*Signal           `json:signals`
	Actions []*Action           `json:actions`
	SigMap  map[string][]string `json:signal_map`
}

type NetworkMap struct {
	sectors map[string]*Sector
	assets  map[string]*Asset
	actions map[string]*Action
	signals map[string]*Signal
	sigmap  map[string][]*Action
}

func (nm *NetworkMap) Matches(o *NetworkMap) bool {
	good := true
	itermap[string, *Sector](nm.sectors, func(name string, value *Sector) {
		good = o.ContainsSector(name)
	})
	itermap[string, *Asset](nm.assets, func(name string, value *Asset) {
		good = o.ContainsAssetByFullName(name)
	})
	itermap[string, *Action](nm.actions, func(name string, value *Action) {
		good = o.ContainsAction(name)
	})
	itermap[string, *Signal](nm.signals, func(name string, value *Signal) {
		good = o.ContainsSignal(name)
	})
	if !good {
		return false
	}
	itermap[string, []*Action](nm.sigmap, func(name string, actlist []*Action) {
		if !mapContains(o.sigmap, name) {
			good = false
			return
		}
		oactlist := o.sigmap[name]
		if len(oactlist) != len(actlist) {
			good = false
			return
		}
		for i, _ := range actlist {
			if oactlist[i].Header.Name != actlist[i].Header.Name {
				good = false
				return
			}
		}
	})
	return good
}

func BlankTopo() Topo {
	return Topo{
		Sectors: make([]*Sector, 0),
		Signals: make([]*Signal, 0),
		Actions: make([]*Action, 0),
		SigMap:  make(map[string][]string),
	}
}

func BlankNetworkMap() *NetworkMap {
	return &NetworkMap{
		sectors: make(map[string]*Sector),
		assets:  make(map[string]*Asset),
		actions: make(map[string]*Action),
		signals: make(map[string]*Signal),
		sigmap:  make(map[string][]*Action),
	}
}

func NetworkMapFromTopo(topo Topo) (*NetworkMap, error) {

	nm := BlankNetworkMap()

	validTriggers := SetFrom([]string{
		TriggerOnEvent,
		TriggerOnTimeout,
		TriggerOnBumpTimeout,
		TriggerOnShutdownNotify,
		TriggerOnSchedule,
		TriggerOnEmit,
	})

	validActionTypes := SetFrom([]string{
		ExecutionTypeFile,
		ExecutionTypeEmbedded,
	})

	if err := forEach(topo.Sectors, func(i int, x *Sector) error {
		return nm.validateSectors(x)
	}); err != nil {
		return nil, err
	}

	if err := forEach(topo.Actions, func(i int, x *Action) error {

		slog.Debug("action", "name", x.Header.Name, "type", x.Type, "data", x.Info)
		if mapContains(nm.actions, x.Header.Name) {
			return NErr("Duplicate action name").
				Push(fmt.Sprintf("Action (%s) does not have a unique name", x.Header.Name))
		}
		if !validActionTypes.Contains(x.Type) {
			return NErr("Invalid trigger specified").
				Push(fmt.Sprintf("Action (%s) has invalid execution type: %s", x.Header.Name, x.Type))
		}

		nm.actions[x.Header.Name] = x

		return nil
	}); err != nil {
		return nil, err
	}

	if err := forEach(topo.Signals, func(i int, x *Signal) error {

		if !validTriggers.Contains(x.Trigger) {
			return NErr("Invalid trigger specified").
				Push(fmt.Sprintf("Signal (%s) has invalid trigger: %s", x.Header.Name, x.Trigger))
		}

		if mapContains(nm.signals, x.Header.Name) {
			return NErr("Duplicate signal name").
				Push(fmt.Sprintf("Signal (%s) does not have a unique name", x.Header.Name))
		}
		nm.signals[x.Header.Name] = x
		return nil
	}); err != nil {
		return nil, err
	}

	for signame, actlist := range topo.SigMap {
		if !mapContains(nm.signals, signame) {
			return nil, NErr("Unknown signal present in signal map").
				Push(fmt.Sprintf("Signal (%s) does not have a definition", signame))
		}
		if err := iterate(actlist, func(it Iter[string]) error {
			if !mapContains(nm.actions, it.Value) {
				return NErr("Unknown action present in signal map's action list").
					Push(fmt.Sprintf("Signal (%s) mapped to unknown action (%s)", signame, it.Value))
			}
			return nil
		}); err != nil {
			return nil, err
		}

		nm.sigmap[signame] = make([]*Action, 0)
		iterate(actlist, func(it Iter[string]) error {
			nm.sigmap[signame] = append(nm.sigmap[signame], nm.actions[it.Value])
			return nil
		})
	}

	// NOTE: In this we don't check to see if every asset contained also has an onEvent signal
	//       as its possible that users want to disable this event and it should be considered
	//       valid to do so
	return nm, nil
}

func (nm *NetworkMap) validateSectors(s *Sector) error {
	if mapContains(nm.sectors, s.Header.Name) {
		return NErr("Duplicate sector").
			Push(
				fmt.Sprintf("Sector (%s) is not unique",
					s.Header.Name))
	}
	nm.sectors[s.Header.Name] = s

	if err := forEach(s.Assets, func(i int, a *Asset) error {
		fullName := makeAssetFullName(s.Header.Name, a.Header.Name)
		if mapContains(nm.assets, fullName) {
			return NErr("Duplicate asset").
				Push(fmt.Sprintf("Asset (%s) is not unique in sector (%s) [%s]", a.Header.Name, s.Header.Name, fullName))
		}
		nm.assets[fullName] = a
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func makeAssetFullName(s string, a string) string {
	return fmt.Sprintf("%s.%s", s, a)
}

func formatAssetOnEventSignalName(asset string) string {
	return fmt.Sprintf("%s:onEvent", asset)
}

func makeAssetOnEventSignal(asset string) *Signal {
	return &Signal{
		Header: HeaderData{
			Name:        formatAssetOnEventSignalName(asset),
			Description: "Signal that fires when associated asset emits event to network",
			Tags:        []string{},
		},
		Trigger: TriggerOnEvent,
	}
}

func (nm *NetworkMap) AddSector(s *Sector) error {
	if mapContains(nm.sectors, s.Header.Name) {
		return NErr("Duplicate sector name").
			Push(fmt.Sprintf("Sector (%s) is not unique", s.Header.Name))
	}
	nm.sectors[s.Header.Name] = s
	return nil
}

func (nm *NetworkMap) AddAsset(sectorName string, asset *Asset) error {
	targetSector, hasSec := nm.sectors[sectorName]
	if !hasSec {
		return NErr("Unknown sector given").
			Push(fmt.Sprintf("Sector (%s) is not a known sector", sectorName))
	}

	full := makeAssetFullName(sectorName, asset.Header.Name)
	if mapContains(nm.assets, full) {
		return NErr("Duplicate asset").
			Push(fmt.Sprintf("Asset name (%s) is not unique in sector (%s) [%s]",
				asset.Header.Name, sectorName, full))
	}

	// Add asset to sector, and set asset map value to new sector entry
	targetSector.Assets = append(targetSector.Assets, asset)
	nm.sectors[sectorName] = targetSector
	nm.assets[full] = asset

	// Setup onEvent for asset
	signal := makeAssetOnEventSignal(full)
	nm.signals[signal.Header.Name] = signal
	return nil
}

func (nm *NetworkMap) AddAction(action *Action) error {
	if mapContains(nm.actions, action.Header.Name) {
		return NErr("Duplicate action name").
			Push(fmt.Sprintf("Action name (%s) is not unique", action.Header.Name))
	}
	nm.actions[action.Header.Name] = action
	return nil
}

func (nm *NetworkMap) AddSignal(signal *Signal) error {
	if mapContains(nm.signals, signal.Header.Name) {
		return NErr("Duplicate signal name").
			Push(fmt.Sprintf("Signal name (%s) is not unique", signal.Header.Name))
	}
	nm.signals[signal.Header.Name] = signal
	return nil
}

func (nm *NetworkMap) ContainsSector(sector string) bool {
	return mapContains(nm.sectors, sector)
}

func (nm *NetworkMap) ContainsAsset(sector string, asset string) bool {
	return mapContains(nm.assets, makeAssetFullName(sector, asset))
}

func (nm *NetworkMap) ContainsAssetByFullName(assetFullPath string) bool {
	return mapContains(nm.assets, assetFullPath)
}

func (nm *NetworkMap) ContainsAction(action string) bool {
	return mapContains(nm.actions, action)
}

func (nm *NetworkMap) ContainsSignal(signal string) bool {
	return mapContains(nm.signals, signal)
}

func (nm *NetworkMap) DeleteSector(sector string) {
	if !nm.ContainsSector(sector) {
		slog.Debug("nm doesn't contain sector", "name", sector)
		return
	}

	s := nm.sectors[sector]
	forEach(s.Assets, func(i int, a *Asset) error {
		fullName := makeAssetFullName(s.Header.Name, a.Header.Name)
		delete(nm.assets, fullName)
		sig := makeAssetOnEventSignal(fullName)
		delete(nm.signals, sig.Header.Name)
		return nil
	})

	delete(nm.sectors, sector)
}

func (nm *NetworkMap) DeleteAsset(sector string, asset string) {
	if !nm.ContainsSector(sector) {
		return
	}

	nm.sectors[sector].Assets = deleteIf[*Asset](nm.sectors[sector].Assets, func(a *Asset) bool {
		return a.Header.Name == asset
	})

	if !nm.ContainsAsset(sector, asset) {
		return
	}

	name := makeAssetFullName(sector, asset)
	delete(nm.assets, name)

	sig := makeAssetOnEventSignal(name)
	delete(nm.signals, sig.Header.Name)
}

func (nm *NetworkMap) DeleteAction(action string) {
	if !nm.ContainsAction(action) {
		return
	}
	delete(nm.actions, action)
}

func (nm *NetworkMap) DeleteSignal(signal string) {
	if !nm.ContainsSignal(signal) {
		return
	}
	delete(nm.signals, signal)
}

func (nm *NetworkMap) MapAction(action string, signal string) error {
	if !nm.ContainsAction(action) {
		return NErr("Unknown action")
	}
	if !nm.ContainsSignal(signal) {
		return NErr("Unknown signal")
	}
	if contains[*Action](nm.sigmap[signal], nm.actions[action],
		func(i int, l *Action, r *Action) bool {
			return l.Header.Name == r.Header.Name
		}) {
		return NErr("Action already mapped to signal")
	}
	nm.sigmap[signal] = append(nm.sigmap[signal], nm.actions[action])
	return nil
}

func (nm *NetworkMap) UnMapSignal(action string, signal string) {
	if !nm.ContainsAction(action) {
		return
	}
	if !nm.ContainsSignal(signal) {
		return
	}
	nm.sigmap[signal] = deleteIf(nm.sigmap[signal], func(v *Action) bool {
		return v.Header.Name == action
	})
}

func (nm *NetworkMap) ToTopo() Topo {
	topo := BlankTopo()
	for _, s := range nm.sectors {
		topo.Sectors = append(topo.Sectors, s)
	}
	itermap[string, *Action](nm.actions, func(name string, value *Action) {
		topo.Actions = append(topo.Actions, value)
	})
	itermap[string, *Signal](nm.signals, func(name string, value *Signal) {
		topo.Signals = append(topo.Signals, value)
	})
	itermap[string, []*Action](nm.sigmap, func(name string, mapped []*Action) {
		topo.SigMap[name] = make([]string, 0)
		iterate[*Action](mapped, func(it Iter[*Action]) error {
			topo.SigMap[name] = append(topo.SigMap[name], it.Value.Header.Name)
			return nil
		})
	})
	return topo
}
