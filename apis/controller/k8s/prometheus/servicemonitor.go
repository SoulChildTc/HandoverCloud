package prometheus

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	"soul/apis/service"
	"soul/utils/httputil"
)

// GetServiceMonitorByName
//
//	@description	获取ServiceMonitor信息
//	@tags			K8s,ServiceMonitor
//	@summary		获取ServiceMonitor信息
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			name			path	string					true	"ServiceMonitor名称"
//	@param			namespace		path	string					true	"Namespace"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回ServiceMonitor信息"
//	@router			/api/v1/k8s/{clusterName}/prometheus/servicemonitor/{namespace}/{name} [get]
func GetServiceMonitorByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "name"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("name")
	namespace := c.Param("namespace")

	servicemonitor, err := service.K8sPrometheusServiceMonitor.GetServiceMonitorByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, servicemonitor, "获取成功")
}

// GetServiceMonitorList
//
//	@description	获取ServiceMonitor列表
//	@tags			K8s,ServiceMonitor
//	@summary		获取ServiceMonitor列表
//	@produce		json
//	@param			clusterName		path	string						true	"Cluster Name"
//	@param			namespace		path	string						false	"Namespace 不填为全部"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@Param			filter			query	string						false	"根据ServiceMonitor名字模糊查询"
//	@Param			limit			query	string						false	"一页获取多少条数据,默认十条"
//	@Param			page			query	string						false	"获取第几页的数据,默认第一页"
//	@success		200				object	httputil.PageResponseBody	"成功返回ServiceMonitor列表"
//	@router			/api/v1/k8s/{clusterName}/prometheus/servicemonitor/ [get]
//	@router			/api/v1/k8s/{clusterName}/prometheus/servicemonitor/{namespace} [get]
func GetServiceMonitorList(c *gin.Context) {
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

	servicemonitors, err := service.K8sPrometheusServiceMonitor.GetServiceMonitorList(clusterName, params.FilterName, namespace, params.Limit, params.Page)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.Page(c, servicemonitors, "获取成功")
}

// DeleteServiceMonitorByName
//
//	@description	删除ServiceMonitor
//	@tags			K8s,ServiceMonitor
//	@summary		删除ServiceMonitor
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			name			path	string	true	"ServiceMonitor名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@success		200				object	nil		"成功返回"
//	@router			/api/v1/k8s/{clusterName}/prometheus/servicemonitor/{namespace}/{name} [delete]
func DeleteServiceMonitorByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "servicemonitorName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("name")
	namespace := c.Param("namespace")

	_, err := service.K8sPrometheusServiceMonitor.GetServiceMonitorByName(clusterName, name, namespace)
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			httputil.Error(c, fmt.Sprintf(`ServiceMonitor "%s" 在 "%s" 中未找到`, name, namespace))
		default:
			httputil.Error(c, err.Error())
		}
		return
	}

	err = service.K8sPrometheusServiceMonitor.DeleteServiceMonitorByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "删除成功")
}
