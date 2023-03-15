package k8s

import (
	"github.com/gin-gonic/gin"
	"soul/apis/controller/core"
)

// k8s模块路由

func RegisterRoute(r *gin.RouterGroup) {
	r.GET("/deployment", core.Ping)
}
