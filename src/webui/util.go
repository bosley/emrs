package webui

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func getLoggedInUser(c *gin.Context) interface{} {
	session := sessions.Default(c)
	return session.Get(sessionKeyUserId)
}

func (wc *controller) getPostData(c *gin.Context, params []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, param := range params {
		target := c.PostForm(param)
		if strings.Trim(target, " ") == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Parameter '%s' can not be empty", param),
			})
			return result, errors.New("empty parameter")
		}
		result[param] = target
	}
	return result, nil
}
