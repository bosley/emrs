package webui

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strings"
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

	slog.Debug("request for dashboard data")
	type TableEntry struct {
		Col1 string
		Col2 string
	}

	type Response struct {
		Assets  []TableEntry
		Actions []TableEntry
		Signals []TableEntry
	}

	assetDb := wc.appCore.GetAssetStore()

	response := Response{
		Assets:  make([]TableEntry, 0),
		Actions: make([]TableEntry, 0),
		Signals: make([]TableEntry, 0),
	}

	if assetDb == nil {
		panic("WJY")
	}

	stored_assets := assetDb.GetAssets()

	for _, asset := range stored_assets {
		slog.Debug("Adding", "name", asset.Name)
		response.Assets = append(
			response.Assets,
			TableEntry{
				Col1: asset.Name,
				Col2: "[under construction]",
			})
	}

	// TODO: Return actions and signals once their database stuff
	//        is setup

	c.JSON(200, gin.H{
		"asset":  response.Assets, // Note: The key matches the UI's dashboard view names in dashboard.js
		"action": response.Actions,
		"signal": response.Signals,
	})
}

func (wc *controller) routeSettings(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "under construction",
	})
}

func (wc *controller) routeCreateItem(c *gin.Context) {

	type pb struct {
		Classification string `json:classification`
		Name           string `json:name`
	}

	var post pb
	arr, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Failed to create record",
			"error":  err.Error(),
		})
		slog.Error("error adding asset", "err", err.Error())
	}

	json.Unmarshal(arr, &post)

	classification := post.Classification
	name := post.Name

	slog.Debug("create item", "name", name)

	if strings.Trim(classification, " ") == "" || strings.Trim(name, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		slog.Error("posted parameters were empty")
		return
	}

	if classification == "asset" {
		slog.Debug("Add asset")
		db := wc.appCore.GetAssetStore()
		if err := db.AddAsset(name, "FIELD CURRENTLY UNUSED"); err != nil {
			c.JSON(500, gin.H{
				"status": "Failed to create record",
				"error":  err.Error(),
			})
			slog.Error("error adding asset", "err", err.Error())
			return
		}
	} else if classification == "action" {
		c.JSON(503, gin.H{
			"status": "under construction",
		})
	} else if classification == "signal" {
		c.JSON(503, gin.H{
			"status": "under construction",
		})
	} else {
		c.JSON(400, gin.H{
			"status": "unknown classification",
		})
	}
	slog.Debug("record created")
	c.JSON(200, gin.H{
		"status": "record created",
	})
}

func (wc *controller) routeDeleteItem(c *gin.Context) {

	type pb struct {
		Classification string `json:classification`
		Name           string `json:name`
	}

	var post pb
	arr, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Failed to create record",
			"error":  err.Error(),
		})
		slog.Error("error adding asset", "err", err.Error())
	}

	json.Unmarshal(arr, &post)

	classification := post.Classification
	name := post.Name

	slog.Debug("create item", "name", name)

	if strings.Trim(classification, " ") == "" || strings.Trim(name, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		slog.Error("posted parameters were empty")
		return
	}

	if classification == "asset" {
		slog.Debug("Add asset")
		db := wc.appCore.GetAssetStore()
		if err := db.DeleteAsset(name); err != nil {
			c.JSON(500, gin.H{
				"status": "Failed to delete record",
				"error":  err.Error(),
			})
			slog.Error("error removing asset", "err", err.Error())
			return
		}
	} else if classification == "action" {
		c.JSON(503, gin.H{
			"status": "under construction",
		})
	} else if classification == "signal" {
		c.JSON(503, gin.H{
			"status": "under construction",
		})
	} else {
		c.JSON(400, gin.H{
			"status": "unknown classification",
		})
	}
	slog.Debug("record deleted")
	c.JSON(200, gin.H{
		"status": "record deleted",
	})
}

func (wc *controller) routeUpdateItem(c *gin.Context) {

	type pb struct {
		Classification string `json:classification`
		Name           string `json:name`
	}

	var post pb
	arr, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "Failed to create record",
			"error":  err.Error(),
		})
		slog.Error("error adding asset", "err", err.Error())
	}

	json.Unmarshal(arr, &post)

	classification := post.Classification
	name := post.Name

	slog.Debug("create item", "name", name)

	if strings.Trim(classification, " ") == "" || strings.Trim(name, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		slog.Error("posted parameters were empty")
		return
	}

	if classification == "asset" {
		slog.Debug("Add asset")
		db := wc.appCore.GetAssetStore()
		if err := db.UpdateAsset(name, "[FIELD NOT CURRENTLY USED]"); err != nil {
			c.JSON(500, gin.H{
				"status": "Failed to update record",
				"error":  err.Error(),
			})
			slog.Error("error removing asset", "err", err.Error())
			return
		}
	} else if classification == "action" {
		c.JSON(503, gin.H{
			"status": "under construction",
		})
	} else if classification == "signal" {
		c.JSON(503, gin.H{
			"status": "under construction",
		})
	} else {
		c.JSON(400, gin.H{
			"status": "unknown classification",
		})
	}
	slog.Debug("record updated")
	c.JSON(200, gin.H{
		"status": "record updated",
	})
}
