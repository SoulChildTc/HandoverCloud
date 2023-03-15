package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"soul/utils/logutil"
	"time"
)

import log "soul/internal/logger"

func Logger(c *gin.Context) {
	// 开始时间
	start := time.Now()

	// Process request
	c.Next()

	// 请求耗时
	duration := logutil.GetDurationInMillisecond(start)

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}

	entry := log.GetEntry().WithFields(logrus.Fields{
		"duration":    duration,
		"client_ip":   c.ClientIP(),
		"method":      c.Request.Method,
		"status_code": c.Writer.Status(),
		"body_size":   c.Writer.Size(),
		"referer":     c.Request.Referer(),
		"path":        path,
		"request_id":  c.Writer.Header().Get("X-Request-ID"),
	})

	if c.Writer.Status() >= 500 {
		entry.Error(c.Errors.String())
	} else {
		entry.Info()
	}
}
