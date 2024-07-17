package app

/*

   Submission of all events/ status/ changes/ requests/ etc.


   For now `/event` is all that will be developed on as it is
   how all external assets will trigger actions

*/

import (
	"github.com/gin-gonic/gin"
)

func (a *App) setupSubmit(gins *gin.Engine) {

	grp := gins.Group("/submit")
	grp.Use(a.SubmitAuthentication())
	grp.POST("/event", a.submitEvent)
}

func (a *App) SubmitAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO:
		//
		//      Retrieve voucher from the request and validate it against the
		//      server identity
		//
		//
		//        For now, loading and checking every time will be fine, but eventually
		//        checking the voucher once, storing its hash in a map, and then making
		//        a timed callback to remove the item the moment it would expire,
		//        may give a performance gain.
		//            For now though, its a debt we can take because it doesn't matter yet
		//
		//
	}
}

func (a *App) submitEvent(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "under construction",
	})
}
