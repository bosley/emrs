package webui

const (
	MsgTypeSector = iota
	MsgTypeAsset
	MsgTypeLayer
	MsgTypeSignal
	MsgTypeAttribute
)

const (
	OpAdd = iota
	OpDel
	OpInfo
)

// Wrapper for message in event
type MsgCommand struct {
	Type int
	Msg  interface{}
}

type MsgSector struct {
	Op int
}

type MsgAsset struct {
	Op int
}

type MsgLayer struct {
	Op int
}

type MsgSignal struct {
	Op int
}

type MsgAttribute struct {
	Op int
}
