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
	existingUserKey  = "user-exists"
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

func (wc *controller) routeLogin(c *gin.Context) {
	_, attempted := c.Get(loginAttemptKey)
	c.HTML(200, "login.html", gin.H{
		"PageHeader":  buildPageHeader("Login"),
		"NavData":     buildNavData(c),
		"PrevAttempt": attempted,
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

	db := wc.appCore.GetUserStore()

	if !db.Validate(username, password) {
		c.Set(loginAttemptKey, true)
		wc.routeLogin(c)
		return
	}

	session.Set(sessionKeyUserId, username)

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

func (wc *controller) routeNewUser(c *gin.Context) {

	if !wc.appCore.RequiresSetup() {
		c.HTML(http.StatusUnauthorized, "message.html", gin.H{
			"NavData":    buildNavData(c),
			"PageHeader": buildPageHeader("Not Available"),
			"Message":    "Creating users is an unauthorized action after initial setup",
			"ShowLogin":  false,
		})
		return
	}

	_, exists := c.Get(existingUserKey)
	c.HTML(200, "new_user.html", gin.H{
		"NavData":    buildNavData(c),
		"PageHeader": buildPageHeader("Create User"),
		"UserExists": exists,
	})
}

func (wc *controller) routeCreateUser(c *gin.Context) {

	if !wc.appCore.RequiresSetup() {
		c.HTML(http.StatusUnauthorized, "message.html", gin.H{
			"NavData":    buildNavData(c),
			"PageHeader": buildPageHeader("Not Available"),
			"Message":    "Creating users is an unauthorized action after initial setup",
			"ShowLogin":  false,
		})
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	slog.Warn("need to create user", "username", username, "password", password)

	db := wc.appCore.GetUserStore()

	if err := db.AddUser(username, password); err != nil {
		slog.Error("Error creating user", "err", err.Error())
		c.Set(existingUserKey, true)
		wc.routeNewUser(c)
		return
	}

	wc.appCore.IndicateSetupComplete()

	c.Redirect(http.StatusMovedPermanently, "/login")
}
