package app

/*


   TODO:


   The CNC is the command and control for the server. All requests MUST be using the ui key within
   the database assigned to the ring-0 user.

   For now each request will validate the voucher key received, but in the future we can cache it.
   Its not critical enough at first to care



   Every single operation that changes the state of the server itself must go through CNC.

   The entire GUI will act as a fancy front to work with the CNC API, so we can
   start off with makign the cli tool to do the CRUD operations on assets and signals.


   Signals should not be directly tied to assets, rather, signals should be a sort of API
   that permits the transit of data along with a context for that data.

   Assets should be defined as thus:

     an entity that submits signals to the network

     All signals evented to the network from an asset must be accompanied by a voucher
     generated by the server


  Signals are defined as a call-to-action, thus, they are directly tied to a
  specific action.

  Actions are not yet fully defined, but they will take the shape of
  "some system that executes a route based on the data input".

  Kind of a convoluted RPC



*/

import (
	"github.com/gin-gonic/gin"
)

func (a *App) setupCNC(gins *gin.Engine) {

	priv := gins.Group("/cnc")
	priv.Use(a.CNCAuthentication())
	{
		priv.POST("/add", a.cncAdd)
	}
}

func (a *App) CNCAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO:
		//
		//      Retrieve cnc key from the request, and match
		//      it against the one that we have stored in the DB. (different auth than the other endpoints)
		//
		//      Reminder: on startup we should identify if the voucher is still good (in addition to each request)
		//      as it might have expired
		//
	}
}

func (a *App) cncAdd(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "under construction",
	})
}