package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestId(c *gin.Context) {
	requestId := c.Request.Header.Get("X-Request-Id")

	if requestId == "" {
		requestId = uuid.New().String()
	}
	c.Set("X-Request-Id", requestId)
	c.Writer.Header().Set("X-Request-Id", requestId)
	c.Next()
}
