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
		Col3 string
	}

	type Response struct {
		Assets  []TableEntry
		Actions []TableEntry
		Signals []TableEntry
	}

	assetDb := wc.appCore.GetAssetStore()
	actionDb := wc.appCore.GetActionStore()
	signalDb := wc.appCore.GetSignalStore()

	response := Response{
		Assets:  make([]TableEntry, 0),
		Actions: make([]TableEntry, 0),
		Signals: make([]TableEntry, 0),
	}

	stored_assets := assetDb.GetAssets()

	for _, asset := range stored_assets {
		slog.Debug("Adding", "name", asset.Name)
		response.Assets = append(
			response.Assets,
			TableEntry{
				Col1: asset.Name,
				Col2: asset.Description,
				Col3: "",
			})
	}

	stored_actions := actionDb.GetActions()

	for _, action := range stored_actions {
		slog.Debug("Adding", "name", action.Name)
		response.Actions = append(
			response.Actions,
			TableEntry{
				Col1: action.Name,
				Col2: action.Description,
				Col3: action.ExecutionInfo,
			})
	}

	stored_signals := signalDb.GetSignals()

	for _, signal := range stored_signals {
		slog.Debug("Adding", "name", signal.Name)
		response.Signals = append(
			response.Signals,
			TableEntry{
				Col1: signal.Name,
				Col2: signal.Description,
				Col3: signal.Triggers,
			})
	}

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

	slog.Debug("delete item", "name", name)

	if strings.Trim(classification, " ") == "" || strings.Trim(name, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		slog.Error("posted parameters were empty")
		return
	}
	if classification == "asset" {
		slog.Debug("Delete asset")
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
		db := wc.appCore.GetActionStore()
		if err := db.DeleteAction(name); err != nil {
			c.JSON(500, gin.H{
				"status": "Failed to delete record",
				"error":  err.Error(),
			})
			slog.Error("error removing action", "err", err.Error())
			return
		}
	} else if classification == "signal" {
		db := wc.appCore.GetSignalStore()
		if err := db.DeleteSignal(name); err != nil {
			c.JSON(500, gin.H{
				"status": "Failed to delete record",
				"error":  err.Error(),
			})
			slog.Error("error removing signal", "err", err.Error())
			return
		}
	} else {
		c.JSON(400, gin.H{
			"status": "unknown classification",
		})
		return
	}
	slog.Debug("record deleted")
	c.JSON(200, gin.H{
		"status": "record deleted",
	})
}

func (wc *controller) routeAddAsset(c *gin.Context) {

	slog.Debug("request to add asset")

	params, err := wc.getPostData(c, []string{
		"name",
		"description",
	})

	if err != nil {
		return // Error response set in getPostData
	}

	db := wc.appCore.GetAssetStore()

	if err := db.AddAsset(params["name"], params["description"]); err != nil {
		slog.Error("Error creating asset", "err", err.Error())
		c.JSON(500, gin.H{
			"status": "Failed to update record",
			"error":  err.Error(),
		})
		c.Abort()
		return
	}
	slog.Debug("record created")
	c.Redirect(http.StatusFound, emrsUrlAppRoot)
}

func (wc *controller) routeEditAsset(c *gin.Context) {

	params, err := wc.getPostData(c, []string{
		"name",
		"original_name",
		"description",
	})

	if err != nil {
		return // Error response set in getPostData
	}

	db := wc.appCore.GetAssetStore()

	if err := db.UpdateAsset(params["original_name"], params["name"], params["description"]); err != nil {
		slog.Error("Error editing asset", "err", err.Error())
		c.JSON(500, gin.H{
			"status": "Failed to update record",
			"error":  err.Error(),
		})
		return
	}
	slog.Debug("record edited")
	c.Redirect(http.StatusFound, emrsUrlAppRoot)
}

func (wc *controller) routeAddAction(c *gin.Context) {

	slog.Debug("request to add action")

	params, err := wc.getPostData(c, []string{
		"name",
		"description",
		"execution_info",
	})

	if err != nil {
		return // Error response set in getPostData
	}

	db := wc.appCore.GetActionStore()

	if err := db.AddAction(params["name"], params["description"], params["execution_info"]); err != nil {
		slog.Error("Error creating action", "err", err.Error())
		c.JSON(500, gin.H{
			"status": "Failed to update record",
			"error":  err.Error(),
		})
		c.Abort()
		return
	}
	slog.Debug("record created")
	c.Redirect(http.StatusFound, emrsUrlAppRoot)
}

func (wc *controller) routeEditAction(c *gin.Context) {

	params, err := wc.getPostData(c, []string{
		"name",
		"original_name",
		"description",
		"execution_info",
	})

	if err != nil {
		return // Error response set in getPostData
	}

	db := wc.appCore.GetActionStore()

	if err := db.UpdateAction(params["original_name"], params["name"], params["description"], params["execution_info"]); err != nil {
		slog.Error("Error editing action", "err", err.Error())
		c.JSON(500, gin.H{
			"status": "Failed to update record",
			"error":  err.Error(),
		})
		return
	}
	slog.Debug("record edited")
	c.Redirect(http.StatusFound, emrsUrlAppRoot)
}

func (wc *controller) routeAddSignal(c *gin.Context) {

	slog.Debug("request to add action")

	params, err := wc.getPostData(c, []string{
		"name",
		"description",
		"triggers",
	})

	if err != nil {
		return // Error response set in getPostData
	}

	db := wc.appCore.GetSignalStore()

	if err := db.AddSignal(params["name"], params["description"], params["triggers"]); err != nil {
		slog.Error("Error creating action", "err", err.Error())
		c.JSON(500, gin.H{
			"status": "Failed to update record",
			"error":  err.Error(),
		})
		c.Abort()
		return
	}
	slog.Debug("record created")
	c.Redirect(http.StatusFound, emrsUrlAppRoot)
}

func (wc *controller) routeEditSignal(c *gin.Context) {

	params, err := wc.getPostData(c, []string{
		"name",
		"original_name",
		"description",
		"triggers",
	})

	if err != nil {
		return // Error response set in getPostData
	}

	db := wc.appCore.GetSignalStore()

	if err := db.UpdateSignal(params["original_name"], params["name"], params["description"], params["triggers"]); err != nil {
		slog.Error("Error editing action", "err", err.Error())
		c.JSON(500, gin.H{
			"status": "Failed to update record",
			"error":  err.Error(),
		})
		return
	}
	slog.Debug("record edited")
	c.Redirect(http.StatusFound, emrsUrlAppRoot)
}
