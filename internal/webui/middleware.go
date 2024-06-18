package webui

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

func (ui *WebUi) EmrsAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if getLoggedInUser(c) == nil {
			c.HTML(http.StatusUnauthorized, "denied.html", gin.H{
				"NavData":    buildNavData(c),
				"PageHeader": buildPageHeader("Access denied"),
			})
			c.Abort()
			return
		}
	}
}

func (ui *WebUi) ReaperMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//  If the application reaper has indicated that we are
		//  about to be killed, then we might as well stop serving
		//  requests. This will stop requsts from users and submitter
		//  programs. We may not need it, but this will give us extra
		//  time to finish processing before the long goodnight
		if ui.killOtw.Load() {
			c.AbortWithStatusJSON(500, gin.H{"message": "server is shutting down"})
		}
	}
}

func (ui *WebUi) RequestProfiler() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Keep track of how many requests we have served (because why not?)
		ui.metrics.requests.Add(1)

		t := time.Now()

		c.Next()

		latency := time.Since(t)

		//  TODO:
		//  In the future I would like to push the
		//  request number and the latency through
		//  the bus to a metrics recording consumer
		//  that will keep a running average of latency
		//  and track metadata about the slowest requests
		//  This isn't the time to do it, but this here
		//  is the place to do it
		//ui.record(fmt.Sprintf("latency: %s", latency.String()))
		slog.Info("req complete", "latency", latency)
	}
}
