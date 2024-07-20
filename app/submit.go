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

	route := c.GetHeader("route")

	data := new(bytes.Buffer)
	data.ReadFrom(c.Request.Body)

	slog.Info("event submission request", "route", route, "body", data)

	println(`



      TODO: 

          This is not yet completed.

          The request has been authorized, but needs to be executed.


          We need to setup the event system that takes the given "route" and passes the "data"

          The first item in the route (which we have yet to decompose/ validate) must be an ingestion
          node. This node will define and populate the data type that the trailered nodes will process
          based on the data handed to us by the user (if any)




  `)

	c.JSON(200, gin.H{
		"status": "under construction",
	})
}
