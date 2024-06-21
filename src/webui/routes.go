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

func (wc *controller) routeIndex(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Home"),
	})
}

func (wc *controller) routeLogin(c *gin.Context) {

	_, ok := c.Get(loginAttemptKey)

	c.HTML(200, "login.html", gin.H{
		"PageHeader":  buildPageHeader("Login"),
		"NavData":     buildNavData(c),
		"PrevAttempt": ok,
	})
}

func (wc *controller) routeLogout(c *gin.Context) {

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

func (wc *controller) routeAuth(c *gin.Context) {

	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	uuid := wc.appCore.ValidateUserAndGetId(username, password)
	if uuid == nil {
		c.Set(loginAttemptKey, true)
		wc.routeLogin(c)
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

func (wc *controller) routeStatus(c *gin.Context) {
	c.HTML(200, "status.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Status"),
	})
}

func (wc *controller) routeDashboard(c *gin.Context) {
	c.HTML(200, "dashboard.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Home"),
	})
}

func (wc *controller) routeSettings(c *gin.Context) {
	c.HTML(200, "settings.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Settings"),
	})
}
