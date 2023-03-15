package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"soul/utils"
	"soul/utils/httputil"
)

func JwtAuth(c *gin.Context) {
	token := c.Request.Header.Get("X-Token")
	if token == "" {
		httputil.ErrorWithCode(c, http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}

	tokenObj, err := utils.ParseJwtToken(token)
	if err != nil {
		httputil.Error(c, err.Error())
		c.Abort()
		return
	}

	c.Set("userId", tokenObj.UserID)
	//c.Set("token", tokenObj)
}
