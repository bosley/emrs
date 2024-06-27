package webui

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

/*
If its the very first time running we route to
new user, otherwise we show login page
*/
func (wc *controller) routeIndex(c *gin.Context) {
	if wc.appCore.RequiresSetup() {
		slog.Debug("system requires setup - displaying user creation")
		wc.routeNewUser(c)
		return
	}

	user := getLoggedInUser(c)
	if user != nil {
		slog.Debug("user already logged in, redirecting to app")
		c.Redirect(http.StatusFound, emrsUrlAppRoot)
		return
	}

	slog.Debug("login page requested")

	_, attempted := c.Get(loginAttemptKey)
	c.HTML(200, "window.html", gin.H{
		"Topic":       "Login",
		"PostTo":      emrsUrlAuth,
		"Prompt":      "EMRS Login",
		"Prompting":   true,
		"PrevAttempt": attempted,
	})
}

func (wc *controller) routeNotificationPoll(c *gin.Context) {

	// All notifications/ alerts (like KILL OTW) should be
	// queued into an area and then dumped out to the
	// caller over JSON jere

	c.JSON(200, gin.H{
		"status": "under construction",
	})
}

func (wc *controller) routeStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "under construction",
	})
}

func (wc *controller) routeDashboard(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "under construction",
	})
}

func (wc *controller) routeSettings(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "under construction",
	})
}
