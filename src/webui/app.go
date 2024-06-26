package webui

import (
	"github.com/gin-gonic/gin"
)

func (wc *controller) routeAppLaunch(c *gin.Context) {

  // Here we could send information regarding the 
  // general un-changing information about the app
  // such as the links in the nab bar,
  // endpoints to acquire specific data, ect
	c.HTML(200, "app.html", gin.H{})
}
