package logutil

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"runtime"
	"strings"
	"time"
)

func CallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 0
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func GetDurationInMillisecond(start time.Time) (millisecond float64) {
	end := time.Now()
	duration := end.Sub(start)
	millisecond = float64(duration) / float64(time.Millisecond)

	// 四舍五入保留两位小数
	millisecond = math.Round(millisecond*100) / 100

	return
}

func GetClientIP(c *gin.Context) string {
	/*
		这里优先获取X-Forwarded-For(暂未使用此方法)
		防止X-Forwarded-For伪造:
			1. 在最外层代理添加类似如下配置,使用$remote_addr覆盖客户端携带的X-Forwarded-For
				proxy_set_header X-Forwarded-For $remote_addr;
			2. 如果有多层代理，内层的代理配置如下, 将在X-Forwarded-For中追加一个$remote_addr
				proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	*/

	// 获取X-Forwarded-For
	clientIP := c.Request.Header.Get("X-Forwarded-For")

	// 不存在X-Forwarded-For, 获取X-Real-IP
	if clientIP == "" {
		clientIP = c.Request.Header.Get("X-Real-IP")
	}

	// 不存在X-Real-IP, 获取socket连接地址
	if clientIP == "" {
		clientIP = c.Request.RemoteAddr
	}

	// 如果IP包含 "," 代表是个IP列表, 取第一个作为clientIP
	if strings.Contains(clientIP, ",") {
		clientIP = strings.Split(clientIP, ",")[0]
	}

	return clientIP
}
