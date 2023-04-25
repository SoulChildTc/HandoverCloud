package svc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	"soul/apis/dto"
	"soul/apis/service"
	"soul/utils/httputil"
)

// GetSvcByName
//
//	@description	获取Svc信息
//	@tags			K8s,Svc
//	@summary		获取Svc信息
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			svcName			path	string					true	"Svc名称"
//	@param			namespace		path	string					true	"Namespace"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回 service 信息"
//	@router			/api/v1/k8s/{clusterName}/svc/{namespace}/{svcName} [get]
func GetSvcByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "svcName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("svcName")
	namespace := c.Param("namespace")

	svc, err := service.K8sSvc.GetSvcByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, svc, "获取成功")
}

// GetSvcList
//
//	@description	获取Svc列表
//	@tags			K8s,Svc
//	@summary		获取Svc列表
//	@produce		json
//	@param			clusterName		path	string						true	"Cluster Name"
//	@param			namespace		path	string						false	"Namespace 不填为全部"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@Param			filter			query	string						false	"根据service名字模糊查询"
//	@Param			limit			query	string						false	"一页获取多少条数据,默认十条"
//	@Param			page			query	string						false	"获取第几页的数据,默认第一页"
//	@success		200				object	httputil.PageResponseBody	"成功返回Service列表"
//	@router			/api/v1/k8s/{clusterName}/svc/ [get]
//	@router			/api/v1/k8s/{clusterName}/svc/{namespace} [get]
func GetSvcList(c *gin.Context) {
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

	services, err := service.K8sSvc.GetSvcList(clusterName, params.FilterName, namespace, params.Limit, params.Page)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.Page(c, services, "获取成功")
}

// DeleteSvcByName
//
//	@description	删除Svc
//	@tags			K8s,Svc
//	@summary		删除Svc
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			svcName			path	string	true	"Svc名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@success		200				object	nil		"成功返回"
//	@router			/api/v1/k8s/{clusterName}/svc/{namespace}/{svcName} [delete]
func DeleteSvcByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "svcName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("svcName")
	namespace := c.Param("namespace")

	_, err := service.K8sSvc.GetSvcByName(clusterName, name, namespace)
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			httputil.Error(c, fmt.Sprintf(`Service "%s" 在 "%s" 中未找到`, name, namespace))
		default:
			httputil.Error(c, err.Error())
		}
		return
	}

	err = service.K8sSvc.DeleteSvcByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "删除成功")
}

// CreateSimpleSvc
//
//	@description	创建简单 Svc, 支持通过DeploymentName 转换为 Selector, 也可以自定义Selector, 二选一。 仅支持创建ClusterIP类型
//	@tags			K8s,Svc
//	@summary		创建简单 Svc
//	@produce		json
//	@param			clusterName	path	string	true	"Cluster Name"
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@param			data			body	dto.K8sSvcSimpleCreate	true	"K8sSvcSimpleCreate 对象"
//	@success		200				object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/k8s/{clusterName}/svc/ [post]
func CreateSimpleSvc(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	svc := dto.K8sSvcSimpleCreate{}

	if err := c.ShouldBindJSON(&svc); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &svc).Error())
		return
	}

	err := service.K8sSvc.CreateSimpleSvc(clusterName, &svc)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "创建成功")
}

// UpdateSimpleSvc
//
//	@description	更新简单 Svc
//	@tags			K8s,Svc
//	@summary		更新简单 Svc
//	@produce		json
//	@param			clusterName	path	string	true	"Cluster Name"
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@param			data			body	dto.K8sSvcSimpleCreate	true	"K8sSvcSimpleCreate 对象"
//	@success		200				object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/k8s/{clusterName}/svc/ [put]
func UpdateSimpleSvc(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	svc := dto.K8sSvcSimpleCreate{}

	if err := c.ShouldBindJSON(&svc); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &svc).Error())
		return
	}

	err := service.K8sSvc.UpdateSimpleSvc(clusterName, &svc)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "更新成功")
}
