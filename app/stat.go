package app

/*

   /     JSON dump of server status
             - uptime,


             Later we will add more than uptime, but
             this is added now because its how the api
             will check if the server is up

*/

import (
	"github.com/gin-gonic/gin"
)

func (a *App) setupStat(gins *gin.Engine) {

	grp := gins.Group("/stat")
	grp.GET("/", a.statRoot)
}

func (a *App) statRoot(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "under construction",
	})
}
