package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BasicAuth(c *gin.Context) {
	// 获取Basic认证的用户名和密码
	username, password, ok := c.Request.BasicAuth()

	// TODO 读取配置文件
	if !ok || username != "admin" || password != "soul" {
		// 如果没有提供正确的用户名和密码，则返回401 Unauthorized状态码
		c.Header("WWW-Authenticate", `Basic realm="Please enter your username and password."`)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 继续处理下一个请求
	c.Next()
}
