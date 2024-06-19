package webui

import (
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	sessionKeyUserId = "user-id" // We may want to change this?
	loginAttemptKey  = "login-failure"
)

func (ui *WebUi) routeIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Home"),
	})
}

func (ui *WebUi) routeLogin(c *gin.Context) {

	_, ok := c.Get(loginAttemptKey)

	slog.Debug("Looks like there was a previous login attempt")

	c.HTML(200, "login.html", gin.H{
		"PageHeader":  buildPageHeader("Login"),
		"NavData":     buildNavData(c),
		"PrevAttempt": ok,
	})
}

func (ui *WebUi) routeLogout(c *gin.Context) {

	session := sessions.Default(c)
	token := session.Get(sessionKeyUserId)
	if token == nil {
		slog.Debug("invalid session token for logout")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}

	session.Delete(sessionKeyUserId)
	if err := session.Save(); err != nil {
		slog.Debug("Unable to save session")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.HTML(200, "message.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Logged Out"),
		"Message":    "You have been logged out",
		"ShowLogin":  true,
	})
}

func (ui *WebUi) routeAuth(c *gin.Context) {

	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	uuid := ui.authUserFn(username, password)
	if uuid == nil {
		c.Set(loginAttemptKey, true)
		ui.routeLogin(c)

		return
	}

	session.Set(sessionKeyUserId, *uuid)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	slog.Info("user logged in", "user", username)

	c.Redirect(http.StatusMovedPermanently, "/emrs/dashboard")
}

func (ui *WebUi) routeStatus(c *gin.Context) {
	c.HTML(200, "status.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Status"),
	})
}

func (ui *WebUi) routeDashboard(c *gin.Context) {
	c.HTML(200, "dashboard.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Home"),
	})
}

func (ui *WebUi) routeSettings(c *gin.Context) {
	c.HTML(200, "settings.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Settings"),
	})
}
