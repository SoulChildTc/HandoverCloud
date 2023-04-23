package deployment

import (
	"github.com/gin-gonic/gin"
	"io"
	"soul/apis/dto"
	"soul/apis/service"
	"soul/utils/httputil"
	"strconv"
)

// GetDeploymentByName
//
//	@description	获取Deployment信息
//	@tags			K8s,Deployment
//	@summary		获取Deployment信息
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			deploymentName	path	string					true	"deployment名称"
//	@param			namespace		path	string					true	"Namespace"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回Deployment信息"
//	@router			/api/v1/k8s/{clusterName}/deployment/{namespace}/{deploymentName} [get]
func GetDeploymentByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "deploymentName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("deploymentName")
	namespace := c.Param("namespace")

	deployment, err := service.K8sDeployment.GetDeploymentByName(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, deployment, "获取成功")
}

// GetDeploymentList
//
//	@description	获取Deployment列表
//	@tags			K8s,Deployment
//	@summary		获取Deployment列表
//	@produce		json
//	@param			clusterName		path	string						true	"Cluster Name"
//	@param			namespace		path	string						false	"Namespace 不填为全部"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@Param			filter			query	string						false	"根据Deployment名字模糊查询"
//	@Param			limit			query	string						false	"一页获取多少条数据,默认十条"
//	@Param			page			query	string						false	"获取第几页的数据,默认第一页"
//	@success		200				object	httputil.PageResponseBody	"成功返回Deployment列表"
//	@router			/api/v1/k8s/{clusterName}/deployment/{namespace} [get]
func GetDeploymentList(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace"); err != nil {
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

	deployments, err := service.K8sDeployment.GetDeploymentList(clusterName, params.FilterName, namespace, params.Limit, params.Page)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.Page(c, deployments, "获取成功")
}

// GetDeploymentPods
//
//	@description	获取 Deployment 管理的 Pod 信息
//	@tags			K8s,Deployment
//	@summary		获取 Deployment 管理的 Pod 信息
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			deploymentName	path	string	true	"Deployment名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@router			/api/v1/k8s/{clusterName}/deployment/{namespace}/{deploymentName}/pods [get]
func GetDeploymentPods(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "deploymentName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("deploymentName")
	namespace := c.Param("namespace")

	pods, err := service.K8sDeployment.GetDeploymentPods(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	data := map[string]interface{}{
		"total": len(pods.Items),
		"items": pods.Items,
	}

	httputil.OK(c, data, "获取成功")
}

// CreateDeployment
//
//	@description	创建 Deployment
//	@tags			K8s,Deployment
//	@summary		创建 Deployment
//	@Accept			json
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@param			data			body	dto.K8sDeploymentCreate	true	"Deployment对象"
//	@success		200				object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/k8s/{clusterName}/deployment/ [post]
func CreateDeployment(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	// 初始化默认值
	deploymentCreate := dto.K8sDeploymentCreate{
		Replicas: 1,
		Cpu:      "300m",
		Memory:   "512Mi",
	}

	if err := c.ShouldBindJSON(&deploymentCreate); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &deploymentCreate).Error())
		return
	}

	err := service.K8sDeployment.CreateDeployment(clusterName, &deploymentCreate)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "创建成功")
}

// ScaleDeployment
//
//	@description	修改 Deployment 副本数
//	@tags			K8s,Deployment
//	@summary		修改 Deployment 副本数
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			deploymentName	path	string	true	"Deployment名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			replicas		body	int		true	"副本数"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@router			/api/v1/k8s/{clusterName}/deployment/{namespace}/{deploymentName}/scale [put]
func ScaleDeployment(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "deploymentName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("deploymentName")
	namespace := c.Param("namespace")

	params := new(struct {
		Replicas int `json:"replicas" binding:"required" msg:"副本数不能为空"`
	})

	if err := c.ShouldBind(params); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, params).Error())
		return
	}

	err := service.K8sDeployment.ScaleDeployment(clusterName, name, namespace, int32(params.Replicas))

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "修改成功")
}

// DeleteDeploymentByName
//
//	@description	删除 Deployment
//	@tags			K8s,Deployment
//	@summary		删除 Deployment
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			deploymentName	path	string					true	"deployment名称"
//	@param			namespace		path	string					true	"Namespace"
//	@param			force			query	bool					false	"是否强制删除"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/k8s/{clusterName}/deployment/{namespace}/{deploymentName} [delete]
func DeleteDeploymentByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "deploymentName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("deploymentName")
	namespace := c.Param("namespace")

	forceStr, _ := c.GetQuery("force")
	force, err := strconv.ParseBool(forceStr)
	if err != nil {
		force = false
	} else {
		force = true
	}

	if name == "" || namespace == "" {
		httputil.Error(c, "deployment和namespace不能为空")
		return
	}

	err = service.K8sDeployment.DeleteDeploymentByName(clusterName, name, namespace, force)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "删除成功")
}

// SetDeploymentImage
//
//	@description	修改 Deployment 容器镜像
//	@tags			K8s,Deployment
//	@summary		修改 Deployment 容器镜像
//	@produce		json
//	@param			clusterName		path	string			true	"Cluster Name"
//	@param			deploymentName	path	string			true	"Deployment名称"
//	@param			namespace		path	string			true	"Namespace"
//	@Param			container		body	dto.K8sSetImage	true	"新的容器镜像,只更新第一个容器时 name参数可忽略"
//	@Param			Authorization	header	string			true	"Authorization token"
//	@router			/api/v1/k8s/{clusterName}/deployment/{namespace}/{deploymentName}/image [put]
func SetDeploymentImage(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "deploymentName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("deploymentName")
	namespace := c.Param("namespace")

	params := dto.K8sSetImage{}
	if err := c.ShouldBind(&params); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &params).Error())
		return
	}

	err := service.K8sDeployment.SetDeploymentImage(clusterName, name, namespace, params)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "修改成功")
}

// RestartDeployment
//
//	@description	重启 Deployment 管理的 Pod
//	@tags			K8s,Deployment
//	@summary		重启 Deployment 管理的 Pod
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			deploymentName	path	string	true	"Deployment名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@router			/api/v1/k8s/{clusterName}/deployment/{namespace}/{deploymentName}/restart [put]
func RestartDeployment(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName", "namespace", "deploymentName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")
	name := c.Param("deploymentName")
	namespace := c.Param("namespace")

	err := service.K8sDeployment.RestartDeployment(clusterName, name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "操作成功")
}

// UpdateK8sDeployment
//
//	@description	使用原生 deployment api 对象更新
//	@tags			K8s,Deployment
//	@summary		使用原生 deployment api 对象更新
//	@produce		json
//	@param			clusterName		path	string	true	"Cluster Name"
//	@param			deploymentName	body	object	true	"Deployment Api Object"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@router			/api/v1/k8s/{clusterName}/deployment [put]
func UpdateK8sDeployment(c *gin.Context) {
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
	err = service.K8sDeployment.UpdateK8sDeployment(clusterName, string(content))
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "更新成功")
}

//TODO CreateK8sDeployment 使用 K8s 1:1 Api 创建Deployment

//TODO UpdateDeployment 使用 dto.K8sDeploymentCreate 对象更新Deployment
