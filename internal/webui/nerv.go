package webui

const (
	MsgTypeInfo = iota
)

// Wrapper for message in event
type MsgCommand struct {
	Type int
	Msg  interface{}
}

type MsgInfo struct {
	Info string
}
