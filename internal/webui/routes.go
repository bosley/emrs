package webui

import (
	"github.com/gin-gonic/gin"
)

func (ui *WebUi) routeHome(c *gin.Context) {

	// For testing
	ui.submitter.SubmitData(&MsgCommand{
		Type: MsgTypeInfo,
		Msg:  "HOME PAGE HIT",
	})

	c.HTML(200, "index.html", gin.H{
		"title":   "E.M.R.S",
		"message": "WORK IN PROGRESS",
	})
}

func (ui *WebUi) routeStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "success",
	})
}
