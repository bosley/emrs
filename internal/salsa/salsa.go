package salsa

/*

chips
  id
  name
  description
  type (0 = sector, 1 = asset, 2 = layer, 3 = signal, 4 = attribute)

associations
  id
  chip_container      ; sector 7
  chip                ; asset  43   (sector 7 contains asset 42)

tokens
  id
  value
  revoked



sector X contains asset Y
layer A contains attribute B
sector G contains signal H
asset P contains signal Q

asset contains layer..

signal contains sector... (when signal is raised, all items in sector are notified)

attributes are members of layers, which can be applied to groups
of items, but assets, sectors, signals, and even layers can have
non-layer-based attributes.

Signals are events that can trigger responses within a sector/ asset.
items "contained" by a signal are essentially subscribers to that signal
and can respond based on the severity attribute of the signal


*/

const (
  SignalSeverityLow = iota
  SignalSeverityMed
  SignalSeverityHigh
)

type Chip struct {
  Name  string
  Description string
}

type Sector struct {
  Tag Chip
  IsVirtual bool
}

type Asset struct {
  Tag Chip
  IsVirtual bool
}

type Layer struct {
  Tag Chip
}

type Signal struct {
  Tag Chip
  Severity    int
}

type Attribute struct {
  Tag Chip
}

