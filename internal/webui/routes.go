package webui

import (
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	sessionKeyUserId = "user-id"
)

func (ui *WebUi) routeIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Home"),
	})
}

func (ui *WebUi) routeLogin(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"PageHeader": buildPageHeader("Login"),
		"NavData":    buildNavData(c),
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

	// Show the login page
	ui.routeLogin(c)
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
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
