package app

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (a *App) setupCNC(gins *gin.Engine) {

	priv := gins.Group("/cnc")
	priv.Use(a.CNCAuthentication())
	{
		priv.POST("/shutdown", a.cncShutdown)
	}
}

func (a *App) CNCAuthentication() gin.HandlerFunc {

	return func(c *gin.Context) {
		token := c.GetHeader("token")

		if !a.badge.ValidateVoucher(token) {
			slog.Error("cnc auth failure: invalid voucher")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "invalid token",
			})
			c.Abort()
			return
		}

		o, e := a.db.GetOwner()
		if e != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "unable to retrieve access code",
			})
			c.Abort()
			return
		}

		if o.UiKey != o.UiKey {
			slog.Error("cnc auth failure: incorrect key for current server instance")
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "invalid token",
			})
			c.Abort()
			return
		}
	}
}

func (a *App) cncShutdown(c *gin.Context) {

	slog.Info("CNC SHUTDOWN REQUEST")

	println(`

    TODO:

      The CNC server has received a shutdown request.
      For now we are going to exit without grace.

      We need to instrument the server with a coordinated,
      gentle, and correct shutdown sequence


  `)

	go func() {
		time.Sleep(2 * time.Second)
		println("TIMED SHUTDOWN TRIGGERED")
		os.Exit(55)
	}()

	c.JSON(200, gin.H{
		"status": "shutdown imminent",
	})
}
