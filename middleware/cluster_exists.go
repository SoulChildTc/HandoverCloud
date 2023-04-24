package middleware

import (
	"github.com/gin-gonic/gin"
	"soul/global"
	"soul/utils/httputil"
)

func ClusterExists(c *gin.Context) {
	clusterName := c.Param("clusterName")
	_, ok := global.K8s.Clusters[clusterName]
	if !ok {
		httputil.Error(c, "集群名称不存在")
		c.Abort()
		return
	}
	c.Next()
}
