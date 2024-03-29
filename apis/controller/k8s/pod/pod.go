package pod

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	"soul/apis/service"
	"soul/utils/httputil"
	"strconv"
)

// GetPodByName
//
//	@description	获取Pod信息
//	@tags			K8s,Pod
//	@summary		获取Pod信息
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			podName			path	string					true	"Pod名称"
//	@param			namespace		path	string					true	"Namespace"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回Pod信息"
//	@router			/api/v1/k8s/{clusterName}/pod/{namespace}/{podName} [get]
func GetPodByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "podName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("podName")
	namespace := c.Param("namespace")

	pod, err := service.K8sPod.GetPodByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, pod, "获取成功")
}

// GetPodList
//
//	@description	获取Pod列表
//	@tags			K8s,Pod
//	@summary		获取Pod列表
//	@produce		json
//	@param			clusterName		path	string						true	"Cluster Name"
//	@param			namespace		path	string						false	"Namespace 不填为全部"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@Param			filter			query	string						false	"根据Pod名字模糊查询"
//	@Param			limit			query	string						false	"一页获取多少条数据,默认十条"
//	@Param			page			query	string						false	"获取第几页的数据,默认第一页"
//	@success		200				object	httputil.PageResponseBody	"成功返回Pod列表"
//	@router			/api/v1/k8s/{clusterName}/pod/ [get]
//	@router			/api/v1/k8s/{clusterName}/pod/{namespace} [get]
func GetPodList(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	namespace := c.Param("namespace")

	params := new(struct {
		FilterName string `form:"filter"`
		Limit      int    `form:"limit,default=10"`
		Page       int    `form:"page,default=1"`
	})

	if err := c.ShouldBind(params); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, params).Error())
		return
	}

	pods, err := service.K8sPod.GetPodList(clusterName, params.FilterName, namespace, params.Limit, params.Page)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.Page(c, pods, "获取成功")
}

// DeletePodByName
//
//	@description	删除Pod
//	@tags			K8s,Pod
//	@summary		删除Pod
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			podName			path	string	true	"Pod名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@success		200				object	nil		"成功返回"
//	@router			/api/v1/k8s/{clusterName}/pod/{namespace}/{podName} [delete]
func DeletePodByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "podName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("podName")
	namespace := c.Param("namespace")

	_, err := service.K8sPod.GetPodByName(clusterName, name, namespace)
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			httputil.Error(c, fmt.Sprintf(`Pod "%s" 在 "%s" 中未找到`, name, namespace))
		default:
			httputil.Error(c, err.Error())
		}
		return
	}

	err = service.K8sPod.DeletePodByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "删除成功")
}

// GetPodLog
//
//	@description	获取Pod日志
//	@tags			K8s,Pod
//	@summary		获取Pod日志
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			podName			path	string	true	"Pod名称"
//	@param			namespace		path	string	true	"Namespace"
//	@param			containerName	query	string	false	"容器名,默认第1个容器"
//	@Param			line			query	int		false	"查看最后多少行日志,默认200"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@success		200				object	nil		"成功返回Pod日志"
//	@router			/api/v1/k8s/{clusterName}/pod/{namespace}/{podName}/log [get]
func GetPodLog(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "podName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("podName")
	namespace := c.Param("namespace")

	containerName := c.Query("containerName")
	line, err := strconv.Atoi(c.DefaultQuery("line", "200"))
	if err != nil {
		line = 200
	}
	if name == "" || namespace == "" {
		httputil.Error(c, "pod名和名称空间不能为空")
		return
	}

	pod, err := service.K8sPod.GetPodLog(clusterName, name, containerName, namespace, int64(line))

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, pod, "获取成功")
}

// GetPodContainers
//
//	@description	获取Pod容器信息
//	@tags			K8s,Pod
//	@summary		获取Pod容器信息
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			podName			path	string	true	"Pod名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@router			/api/v1/k8s/{clusterName}/pod/{namespace}/{podName}/containers [get]
func GetPodContainers(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "podName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("podName")
	namespace := c.Param("namespace")

	containers, err := service.K8sPod.GetPodContainers(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	data := map[string]interface{}{
		"total": len(containers),
		"items": containers,
	}

	httputil.OK(c, data, "获取成功")
}

// ExecContainer
//
//	@description	获取 exec sessionId
//	@tags			K8s,Pod
//	@summary		获取 exec sessionId
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			podName			path	string	true	"Pod名称"
//	@param			namespace		path	string	true	"Namespace"
//	@param			containerName	query	string	false	"容器名,默认第1个容器"
//	@Param			shell			query	string	false	"执行shell"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@success		200				object	nil		"成功返回 sessionId 日志"
//	@router			/api/v1/k8s/{clusterName}/pod/{namespace}/{podName}/exec [get]
func ExecContainer(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "podName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("podName")
	namespace := c.Param("namespace")

	containerName := c.Query("containerName")
	if containerName == "" {
		containers, err := service.K8sPod.GetPodContainers(clusterName, name, namespace)
		if err != nil {
			httputil.Error(c, err.Error())
			return
		}
		containerName = containers[0].Name
	}

	shell := c.Query("shell")

	sessionID, err := service.K8sPod.StartTerminal(clusterName, namespace, name, containerName, shell)
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, gin.H{"id": sessionID}, "获取成功")
}
