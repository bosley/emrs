package webui

import (
  "log/slog"
	"net/http"
	"github.com/gin-gonic/gin"
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
    c.Redirect(http.StatusFound, "/emrs")
    return
  }
  
  slog.Debug("login page requested")

	_, attempted := c.Get(loginAttemptKey)
	c.HTML(200, "window.html", gin.H{
    "Topic": "Login",
    "PostTo": "/auth",
    "Prompt": "EMRS Login",
    "Prompting": true,
		"PrevAttempt": attempted,
	})
}

func (wc *controller) routeStatus(c *gin.Context) {
	c.HTML(200, "status.html", gin.H{
	})
}

func (wc *controller) routeDashboard(c *gin.Context) {
	c.HTML(200, "dashboard.html", gin.H{
	})
}

func (wc *controller) routeSettings(c *gin.Context) {
	c.HTML(200, "settings.html", gin.H{
	})
}
func (wc *controller) routeDev(c *gin.Context) {
	c.HTML(200, "signin.html", gin.H{
	})
}
