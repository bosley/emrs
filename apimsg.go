package main

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
