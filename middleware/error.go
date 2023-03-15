package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	log "soul/internal/logger"
	"soul/utils/logutil"
	"strings"
)

func ErrorInterceptor(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				var se *os.SyscallError
				if errors.As(ne, &se) {
					seStr := strings.ToLower(se.Error())
					if strings.Contains(seStr, "broken pipe") ||
						strings.Contains(seStr, "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			stack := logutil.CallerInfo(6)
			httpRequest, _ := httputil.DumpRequest(c.Request, false)
			headers := strings.Split(string(httpRequest), "\r\n")
			for idx, header := range headers {
				current := strings.Split(header, ":")
				if current[0] == "Authorization" {
					headers[idx] = current[0] + ": *"
				}
			}

			if brokenPipe {
				log.Error("%s, %s", err, headers)
			} else if gin.IsDebugging() {
				log.Error("panic recovered: %s, %s, %s", headers, err, stack)
			} else {
				log.Error("panic recovered: %s, %s", err, stack)
			}

			if brokenPipe {
				// If the connection is dead, we can't write a status to it.
				c.Error(err.(error)) //nolint: errcheck
				c.Abort()
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器异常"})
			}
		}
	}()
	c.Next()
}
