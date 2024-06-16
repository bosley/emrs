package webui

import (
	"log/slog"
)

func (ui *WebUi) procCmdInfo(cmd *MsgInfo) {
	slog.Info("webui:info", "msg", cmd.Info)
}
