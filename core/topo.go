package core

import (
	"fmt"
	"path/filepath"
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
	Header  HeaderData `json:header`
	Assets  []Asset    `json:assets`
	Sectors []Sector   `json:sectors`
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
	Sectors []Sector          `json:sectors`
	Signals []Signal          `json:signals`
	Actions []Action          `json:actions`
	SigMap  map[string]string `json:signal_map`
}

type NetworkMap struct {
	fpAssets  Set // Set of all sectors given their fully qual name (full path)
	fpSectors Set // Set of all assets given their fully qual name
	actions   map[string]*Action
	signals   map[string]Signal
	sigmap    map[string]*Action
}

func (nm *NetworkMap) validateSectors(root string, s Sector) error {

	currentSectorFullPath := filepath.Join(root, s.Header.Name)

	if nm.fpSectors.Contains(currentSectorFullPath) {
		return NErr("Duplicate sector").
			Push(
				fmt.Sprintf("Sector (%s) is not unique",
					currentSectorFullPath))
	}
	nm.fpSectors.Insert(currentSectorFullPath)

	if err := forEach(s.Assets, func(i int, a Asset) error {
		assetFull := filepath.Join(currentSectorFullPath, a.Header.Name)
		if nm.fpAssets.Contains(assetFull) {
			return NErr("Duplicate asset").
				Push(fmt.Sprintf("Sector (%s) is not unique", assetFull))
		}
		nm.fpAssets.Insert(assetFull)
		return nil
	}); err != nil {
		return err
	}
	return forEach(s.Sectors, func(i int, x Sector) error {
		return nm.validateSectors(currentSectorFullPath, x)
	})
}

func NetworkMapFromTopo(topo Topo) (*NetworkMap, error) {

	nm := NetworkMap{
		fpAssets:  NewSet(),
		fpSectors: NewSet(),
		actions:   make(map[string]*Action),
		signals:   make(map[string]Signal),
	}

	if err := forEach(topo.Sectors, func(i int, x Sector) error {
		return nm.validateSectors("/", x)
	}); err != nil {
		return nil, err
	}

	if err := forEach(topo.Actions, func(i int, x Action) error {
		_, exists := nm.actions[x.Header.Name]
		if exists {
			return NErr("Duplicate action name").
				Push(fmt.Sprintf("Action (%s) does not have a unique name", x.Header.Name))
		}
		nm.actions[x.Header.Name] = &x
		return nil
	}); err != nil {
		return nil, err
	}

	if err := forEach(topo.Signals, func(i int, x Signal) error {
		_, exists := nm.signals[x.Header.Name]
		if exists {
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
