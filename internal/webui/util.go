package webui

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func getLoggedInUser(c *gin.Context) interface{} {
	session := sessions.Default(c)
	return session.Get(sessionKeyUserId)
}

type PageHeader struct {
	Title string
}

func buildPageHeader(name string) PageHeader {
	return PageHeader{
		Title: fmt.Sprintf("E.M.R.S %s", name),
	}
}

type NavData struct {
	LoggedIn bool
	UserId   string
}

func buildNavData(c *gin.Context) NavData {

	user := getLoggedInUser(c)

	if user == nil {
		return NavData{
			LoggedIn: false,
			UserId:   "",
		}
	}
	return NavData{
		LoggedIn: true,
		UserId:   user.(string),
	}
}
