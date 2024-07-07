package main

import (
	"emrs/core"
	"fmt"
)

const (
	OpAdd = "opAdd"
	OpDel = "opDel"
)

const (
	SubjectSector  = "sector"
	SubjectAsset   = "asset"
	SubjectSignal  = "signal"
	SubjectAction  = "action"
	SubjectMapping = "mapping"
	SubjectTopo    = "topo"
)

type ApiMsg struct {
	Op      string `json:op`
	Subject string `json:subject`
	Data    string `json:data`
}

func (m ApiMsg) String() string {
	return fmt.Sprintf(
		"op:%s, subject:%s, data:%s", m.Op, m.Subject, m.Data)
}

type ApiAddAsset struct {
	Sector string     `json:sector`
	Asset  core.Asset `json:asset`
}
