package app

/*

   Submission of all events/ status/ changes/ requests/ etc.


   For now `/event` is all that will be developed on as it is
   how all external assets will trigger actions

*/

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (a *App) setupSubmit(gins *gin.Engine) {

	grp := gins.Group("/submit")
	grp.Use(a.SubmitAuthentication())
	grp.POST("/event", a.submitEvent)
}

func (a *App) SubmitAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.GetHeader("token")
		origin := c.GetHeader("origin")

		if err := a.validateRequest(origin, token); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "bad request",
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		slog.Debug("origin validated", "origin", origin)
	}
}

func (a *App) submitEvent(c *gin.Context) {

	slog.Debug("EVENT SUBMITTED", "body", c.Request.Body)

	origin := c.GetHeader("origin")
	route := c.GetHeader("route")

	data := new(bytes.Buffer)
	data.ReadFrom(c.Request.Body)

	slog.Info("event submission request", "route", route, "body", data)

  // Submit the job
  //
  if err := a.runner.SubmitJob(&Job{
    Ctx: nil,
    Origin: origin,
    Destination: route,
    Data: data.Bytes(),
  }); err != nil {
	  c.JSON(500, gin.H{
	  	"status": "failed to submit job for execution",
      "error": err.Error(),
	  })
    return
  }

  // Complete
  //
	c.JSON(200, gin.H{
		"status": "under construction",
	})
}
