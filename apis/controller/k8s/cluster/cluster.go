package cluster

import (
	"github.com/gin-gonic/gin"
	"soul/apis/dto"
	"soul/apis/service"
	log "soul/internal/logger"
	"soul/utils/httputil"
)

// GetClusterByName
//
//	@description	获取集群信息
//	@tags			K8s,Cluster
//	@summary		获取集群信息
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回集群信息"
//	@router			/api/v1/k8s/cluster/{clusterName}/ [get]
func GetClusterByName(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	cluster := service.K8sCluster.GetClusterByName(clusterName)

	httputil.OK(c, cluster, "获取成功")
}

// GetClusterList
//
//	@description	获取集群列表
//	@tags			K8s,Cluster
//	@summary		获取集群列表
//	@produce		json
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回集群列表"
//	@router			/api/v1/k8s/cluster/ [get]
func GetClusterList(c *gin.Context) {
	cluster := service.K8sCluster.GetClusterList()

	httputil.OK(c, cluster, "获取成功")
}

// AddCluster
//
//	@description	添加集群
//	@tags			K8s,Cluster
//	@summary		添加集群
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			clusterInfo		body	dto.K8sClusterInfo		true	"集群信息"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回集群信息"
//	@router			/api/v1/k8s/cluster/{clusterName}/ [post]
func AddCluster(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	cluster := dto.K8sClusterInfo{}

	if err := c.ShouldBindJSON(&cluster); err != nil {
		log.Debug(err.Error())
		httputil.Error(c, httputil.ParseValidateError(err, &cluster).Error())
		return
	}

	cluster.ClusterName = clusterName

	err := service.K8sCluster.AddCluster(cluster)
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, cluster, "添加成功")
}

// UpdateCluster
//
//	@description	更新集群
//	@tags			K8s,Cluster
//	@summary		更新集群
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@param			clusterInfo		body	dto.K8sClusterInfo		true	"集群信息"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回集群信息"
//	@router			/api/v1/k8s/cluster/{clusterName}/ [put]
func UpdateCluster(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	cluster := dto.K8sClusterInfo{}

	if err := c.ShouldBindJSON(&cluster); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &cluster).Error())
		return
	}

	cluster.ClusterName = clusterName

	err := service.K8sCluster.UpdateCluster(cluster)
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, cluster, "更新成功")
}

// DeleteCluster
//
//	@description	删除集群
//	@tags			K8s,Cluster
//	@summary		删除集群
//	@produce		json
//	@param			clusterName		path	string					true	"Cluster Name"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回集群信息"
//	@router			/api/v1/k8s/cluster/{clusterName}/ [delete]
func DeleteCluster(c *gin.Context) {
	if err := httputil.CheckParams(c, "clusterName"); err != nil {
		httputil.Error(c, err.Error())
		return
	}

	clusterName := c.Param("clusterName")

	err := service.K8sCluster.DeleteCluster(clusterName)
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "删除成功")
}
