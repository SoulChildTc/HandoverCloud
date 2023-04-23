package namespace

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"k8s.io/apimachinery/pkg/api/errors"
	"soul/apis/service"
	"soul/utils/httputil"
)

// GetNamespaceByName
//
//	@description	获取 Namespace 信息
//	@tags			K8s,Namespace
//	@summary		获取 Namespace 信息
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			namespaceName	path	string					true	"Namespace名称"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回 Namespace 信息"
//	@router			/api/v1/k8s/{clusterName}/namespace/{namespaceName} [get]
func GetNamespaceByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespaceName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("deploymentName")

	namespace, err := service.K8sNamespace.GetNamespaceByName(clusterName, name)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, namespace, "获取成功")
}

// GetNamespaceList
//
//	@description	获取 Namespace 列表
//	@tags			K8s
//	@summary		获取 Namespace 列表
//	@produce		json
//	@param			clusterName		path	string						true	"Cluster Name"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@Param			filter			query	string						false	"根据Namespace名字模糊查询"
//	@Param			limit			query	string						false	"一页获取多少条数据,默认十条"
//	@Param			page			query	string						false	"获取第几页的数据,默认第一页"
//	@success		200				object	httputil.PageResponseBody	"成功返回Namespace列表"
//	@router			/api/v1/k8s/{clusterName}/namespace/ [get]
func GetNamespaceList(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	params := new(struct {
		FilterName string `form:"filter"`
		Limit      int    `form:"limit,default=10"`
		Page       int    `form:"page,default=1"`
	})

	if err := c.ShouldBind(params); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, params).Error())
		return
	}

	namespaces, err := service.K8sNamespace.GetNamespaceList(clusterName, params.FilterName, params.Limit, params.Page)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.Page(c, namespaces, "获取成功")
}

// DeleteNamespaceByName
//
//	@description	删除Namespace
//	@tags			K8s
//	@summary		删除Namespace
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			namespaceName	path	string	true	"Namespace名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@success		200				object	nil		"成功返回"
//	@router			/api/v1/k8s/{clusterName}/namespace/{namespaceName}/ [delete]
func DeleteNamespaceByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespaceName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("deploymentName")

	_, err := service.K8sNamespace.GetNamespaceByName(clusterName, name)
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			httputil.Error(c, fmt.Sprintf(`Namespace %s 不存在`, name))
		default:
			httputil.Error(c, err.Error())
		}
		return
	}

	err = service.K8sNamespace.DeleteNamespaceByName(clusterName, name)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "删除成功")
}

// CreateNamespace
//
//	@description	创建 Namespace
//	@tags			K8s,Namespace
//	@summary		创建 Namespace
//	@Accept			json
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@param			data			body	object					true	"Namespace对象"
//	@success		200				object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/k8s/{clusterName}/namespace/ [post]
func CreateNamespace(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	content, err := io.ReadAll(c.Request.Body)
	if err != nil || len(content) == 0 {
		httputil.Error(c, "参数异常")
		return
	}
	err = service.K8sNamespace.CreateNamespace(clusterName, string(content))
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "创建成功")
}
