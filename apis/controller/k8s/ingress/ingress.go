package ingress

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	"soul/apis/dto"
	"soul/apis/service"
	"soul/utils/httputil"
)

// GetIngressByName
//
//	@description	获取 Ingress 信息
//	@tags			K8s,Ingress
//	@summary		获取 Ingress 信息
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			ingressName		path	string					true	"ingress名称"
//	@param			namespace		path	string					true	"Namespace"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回 Ingress 信息"
//	@router			/api/v1/k8s/{clusterName}/ingress/{namespace}/{ingressName} [get]
func GetIngressByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "ingressName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("ingressName")
	namespace := c.Param("namespace")

	deployment, err := service.K8sIngress.GetIngressByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, deployment, "获取成功")
}

// GetIngressList
//
//	@description	获取 Ingress 列表
//	@tags			K8s,Ingress
//	@summary		获取 Ingress 列表
//	@produce		json
//	@param			clusterName		path	string						true	"Cluster Name"
//	@param			namespace		path	string						false	"Namespace 不填为全部"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@Param			filter			query	string						false	"根据 Ingress 名字模糊查询"
//	@Param			limit			query	string						false	"一页获取多少条数据,默认十条"
//	@Param			page			query	string						false	"获取第几页的数据,默认第一页"
//	@success		200				object	httputil.PageResponseBody	"成功返回 Ingress 列表"
//	@router			/api/v1/k8s/{clusterName}/ingress/ [get]
//	@router			/api/v1/k8s/{clusterName}/ingress/{namespace} [get]
func GetIngressList(c *gin.Context) {
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

	deployments, err := service.K8sIngress.GetIngressList(clusterName, params.FilterName, namespace, params.Limit, params.Page)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.Page(c, deployments, "获取成功")
}

// CreateSimpleIngress
//
//	@description	创建简单 Ingress ,仅支持一个路由前缀,支持多个Host, 对应一个Service
//	@tags			K8s,Ingress
//	@summary		创建简单 Ingress
//	@produce		json
//	@param			clusterName	path	string	true	"Cluster Name"
//	@produce		json
//	@param			clusterName		path	string						true	"Cluster Name"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@param			data			body	dto.K8sIngressSimpleCreate	true	"K8sIngressSimpleCreate 对象"
//	@success		200				object	httputil.ResponseBody		"成功返回"
//	@router			/api/v1/k8s/{clusterName}/ingress/ [post]
func CreateSimpleIngress(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	// 初始化默认值
	ingress := dto.K8sIngressSimpleCreate{}
	ingress.Rule.Path = "/"

	if err := c.ShouldBindJSON(&ingress); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &ingress).Error())
		return
	}

	err := service.K8sIngress.CreateSimpleIngress(clusterName, &ingress)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "创建成功")
}

// UpdateSimpleIngress
//
//	@description	更新简单 Ingress
//	@tags			K8s,Ingress
//	@summary		更新简单 Ingress
//	@produce		json
//	@param			clusterName	path	string	true	"Cluster Name"
//	@produce		json
//	@param			clusterName		path	string						true	"Cluster Name"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@param			data			body	dto.K8sIngressSimpleCreate	true	"K8sIngressSimpleCreate 对象"
//	@success		200				object	httputil.ResponseBody		"成功返回"
//	@router			/api/v1/k8s/{clusterName}/ingress/ [put]
func UpdateSimpleIngress(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	// 初始化默认值
	ingress := dto.K8sIngressSimpleCreate{}
	ingress.Rule.Path = "/"

	if err := c.ShouldBindJSON(&ingress); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &ingress).Error())
		return
	}

	err := service.K8sIngress.UpdateSimpleIngress(clusterName, &ingress)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "更新成功")
}

// DeleteIngressByName
//
//	@description	删除 Ingress
//	@tags			K8s,Ingress
//	@summary		删除 Ingress
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			ingressName		path	string	true	"Ingress名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@success		200				object	nil		"成功返回"
//	@router			/api/v1/k8s/{clusterName}/ingress/{namespace}/{ingressName} [delete]
func DeleteIngressByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "ingressName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("ingressName")
	namespace := c.Param("namespace")

	_, err := service.K8sIngress.GetIngressByName(clusterName, name, namespace)
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			httputil.Error(c, fmt.Sprintf(`Service "%s" 在 "%s" 中未找到`, name, namespace))
		default:
			httputil.Error(c, err.Error())
		}
		return
	}

	err = service.K8sIngress.DeleteIngressByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "删除成功")
}
