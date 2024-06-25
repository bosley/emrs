package webui

import (
	"github.com/gin-gonic/gin"
)

func (wc *controller) routeIndex(c *gin.Context) {
	if wc.appCore.RequiresSetup() {
		wc.routeNewUser(c)
	} else {
		c.HTML(200, "index.html", gin.H{
			"NavData":    buildNavData(c),
			"PageHeader": buildPageHeader("Home"),
		})
	}
}

func (wc *controller) routeStatus(c *gin.Context) {
	c.HTML(200, "status.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Status"),
	})
}

func (wc *controller) routeDashboard(c *gin.Context) {
	c.HTML(200, "dashboard.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Dashboard"),
	})
}

func (wc *controller) routeSettings(c *gin.Context) {
	c.HTML(200, "settings.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Settings"),
	})
}
