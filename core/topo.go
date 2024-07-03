package core

import (
	"fmt"
)

const (
	TriggerOnEvent          = "onEvent"
	TriggerOnTimeout        = "onTimeout"
	TriggerOnBumpTimeout    = "onBumpTimeout"
	TriggerOnShutdownNotify = "onShutdown"
	TriggerOnSchedule       = "onSchedule"
	TriggerOnEmit           = "onEmit"
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
	Header        HeaderData `json:header`
	ExecutionData []byte     `json:data`
}

type Topo struct {
	Sectors []*Sector         `json:sectors`
	Signals []*Signal         `json:signals`
	Actions []*Action         `json:actions`
	SigMap  map[string]string `json:signal_map`
}

type NetworkMap struct {
	sectors map[string]*Sector
	assets  map[string]*Asset
	actions map[string]*Action
	signals map[string]*Signal
	sigmap  map[string]*Action
}

func BlankNetworkMap() *NetworkMap {
	return &NetworkMap{
		sectors: make(map[string]*Sector),
		assets:  make(map[string]*Asset),
		actions: make(map[string]*Action),
		signals: make(map[string]*Signal),
	}
}

func NetworkMapFromTopo(topo Topo) (*NetworkMap, error) {

	nm := BlankNetworkMap()

	if err := forEach(topo.Sectors, func(i int, x *Sector) error {
		return nm.validateSectors(x)
	}); err != nil {
		return nil, err
	}

	if err := forEach(topo.Actions, func(i int, x *Action) error {
		if mapContains(nm.actions, x.Header.Name) {
			return NErr("Duplicate action name").
				Push(fmt.Sprintf("Action (%s) does not have a unique name", x.Header.Name))
		}
		nm.actions[x.Header.Name] = x
		return nil
	}); err != nil {
		return nil, err
	}

	if err := forEach(topo.Signals, func(i int, x *Signal) error {
		if mapContains(nm.signals, x.Header.Name) {
			return NErr("Duplicate signal name").
				Push(fmt.Sprintf("Signal (%s) does not have a unique name", x.Header.Name))
		}
		nm.signals[x.Header.Name] = x
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
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

func makeAssetOnEventSignal(a string) *Signal {
	return &Signal{
		Header: HeaderData{
			Name:        fmt.Sprintf("%s:onEvent", a),
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
	if mapContains(nm.signals, action.Header.Name) {
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

func (nm *NetworkMap) ContainsAction(action string) bool {
	return mapContains(nm.actions, action)
}

func (nm *NetworkMap) ContainsSignal(signal string) bool {
	return mapContains(nm.signals, signal)
}

func (nm *NetworkMap) DeleteSector(sector string) {
	if !nm.ContainsSector(sector) {
		return
	}
	s := nm.sectors[sector]
	forEach(s.Assets, func(i int, a *Asset) error {
		fullName := makeAssetFullName(s.Header.Name, a.Header.Name)
		delete(nm.assets, fullName)
		sig := makeAssetOnEventSignal(fullName)
		delete(nm.signals, sig.Header.Name)
		delete(nm.sectors, sector)
		return nil
	})
}

func (nm *NetworkMap) DeleteAsset(sector string, asset string) {
	if !nm.ContainsSector(sector) {
		return
	}

	nm.sectors[sector].Assets = deleteIf[*Asset](nm.sectors[sector].Assets, func(a *Asset) bool {
		return a.Header.Name == asset
	})

	name := makeAssetFullName(sector, asset)
	if !nm.ContainsAsset(sector, name) {
		return
	}
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
