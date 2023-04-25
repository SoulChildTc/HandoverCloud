package middleware

import (
	"github.com/gin-gonic/gin"
	"soul/utils/httputil"
)

func CheckPermission(c *gin.Context) {
	ok := hasPermission(c)

	if !ok {
		httputil.Error(c, "您没有访问此资源的权限")
		c.Abort()
		return
	}

	c.Next()
}

func hasPermission(c *gin.Context) bool {
	// TODO: 检查用户是否具有权限
	return true
}
