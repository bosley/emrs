package webui

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getLoggedInUser(c *gin.Context) interface{} {
	session := sessions.Default(c)
	return session.Get(sessionKeyUserId)
}
