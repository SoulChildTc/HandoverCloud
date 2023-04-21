package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"soul/utils"
	"soul/utils/httputil"
	"strings"
)

const bearerPrefix = "Bearer "

func extractToken(c *gin.Context) (string, error) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("authorization token is missing")
	}
	if !strings.HasPrefix(token, bearerPrefix) {
		return "", fmt.Errorf("invalid token")
	}
	return strings.TrimSpace(token[len(bearerPrefix):]), nil
}

func JwtAuth(c *gin.Context) {
	token, err := extractToken(c)
	if err != nil {
		httputil.ErrorWithCode(c, http.StatusUnauthorized, err.Error())
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
