package k8s

import (
	"github.com/gin-gonic/gin"
	k8scluster "soul/apis/controller/k8s/cluster"
	k8sdeployment "soul/apis/controller/k8s/deployment"
	k8singress "soul/apis/controller/k8s/ingress"
	k8snamespace "soul/apis/controller/k8s/namespace"
	k8spod "soul/apis/controller/k8s/pod"
	k8ssecret "soul/apis/controller/k8s/secret"
	k8ssvc "soul/apis/controller/k8s/svc"
	"soul/middleware"
)

// k8s模块路由

func RegisterRoute(r *gin.RouterGroup) {

	clusterResource := r.Group("/cluster")
	{
		clusterResource.GET("/", k8scluster.GetClusterList)
		clusterResource.GET("/:clusterName", k8scluster.GetClusterByName)
		clusterResource.POST("/:clusterName", k8scluster.AddCluster)
		clusterResource.PUT("/:clusterName", k8scluster.UpdateCluster)
		clusterResource.DELETE("/:clusterName", k8scluster.DeleteCluster)
	}

	cluster := r.Group("/:clusterName")
	cluster.Use(middleware.ClusterExists)
	pod := cluster.Group("/pod")
	{
		pod.GET("/", k8spod.GetPodList)
		pod.GET("/:namespace", k8spod.GetPodList)
		pod.GET("/:namespace/:podName", k8spod.GetPodByName)
		pod.DELETE("/:namespace/:podName", k8spod.DeletePodByName)
		pod.GET("/:namespace/:podName/log", k8spod.GetPodLog)
		pod.GET("/:namespace/:podName/containers", k8spod.GetPodContainers)
		pod.GET("/:namespace/:podName/shell", k8spod.ExecContainer)
	}

	deployment := cluster.Group("/deployment")
	{
		deployment.GET("/", k8sdeployment.GetDeploymentList)
		deployment.GET("/:namespace", k8sdeployment.GetDeploymentList)
		deployment.GET("/:namespace/:deploymentName", k8sdeployment.GetDeploymentByName)
		deployment.DELETE("/:namespace/:deploymentName", k8sdeployment.DeleteDeploymentByName)
		deployment.PUT("/", k8sdeployment.UpdateK8sDeployment)
		deployment.PUT("/:namespace/:deploymentName/image", k8sdeployment.SetDeploymentImage)
		deployment.PUT("/:namespace/:deploymentName/scale", k8sdeployment.ScaleDeployment)
		deployment.PUT("/:namespace/:deploymentName/restart", k8sdeployment.RestartDeployment)
		deployment.GET("/:namespace/:deploymentName/pods", k8sdeployment.GetDeploymentPods)
		deployment.POST("/", k8sdeployment.CreateDeployment)
	}

	ingress := cluster.Group("/ingress")
	{
		ingress.GET("/", k8singress.GetIngressList)
		ingress.GET("/:namespace", k8singress.GetIngressList)
		ingress.GET("/:namespace/:ingressName", k8singress.GetIngressByName)
		ingress.DELETE("/:namespace/:ingressName", k8singress.DeleteIngressByName)
		ingress.POST("/", k8singress.CreateSimpleIngress)
		ingress.PUT("/", k8singress.UpdateSimpleIngress)
	}

	namespace := cluster.Group("/namespace")
	{
		namespace.GET("/", k8snamespace.GetNamespaceList)
		namespace.GET("/:namespaceName", k8snamespace.GetNamespaceByName)
		namespace.POST("/", k8snamespace.CreateNamespace)
		namespace.DELETE("/:namespaceName", k8snamespace.DeleteNamespaceByName)
	}

	svc := cluster.Group("/svc")
	{
		svc.GET("/", k8ssvc.GetSvcList)
		svc.GET("/:namespace", k8ssvc.GetSvcList)
		svc.GET("/:namespace/:svcName", k8ssvc.GetSvcByName)
		svc.DELETE("/:namespace/:svcName", k8ssvc.DeleteSvcByName)
		svc.POST("/", k8ssvc.CreateSimpleSvc)
		svc.PUT("/", k8ssvc.UpdateSimpleSvc)
	}

	secret := cluster.Group("/secret")
	{
		secret.GET("/", k8ssecret.GetSecretList)
		secret.GET("/:namespace", k8ssecret.GetSecretList)
		secret.GET("/:namespace/:secretName", k8ssecret.GetSecretByName)
		secret.DELETE("/:namespace/:secretName", k8ssecret.DeleteSecretByName)
		secret.POST("/", k8ssecret.CreateSecret)
		secret.PUT("/", k8ssecret.UpdateSecret)

		secret.POST("/_docker-registry", k8ssecret.CreateSecretForDockerRegistry)
		secret.PUT("/_docker-registry", k8ssecret.UpdateSecretForDockerRegistry)
		secret.POST("/_tls", k8ssecret.CreateSecretForTls)
		secret.PUT("/_tls", k8ssecret.UpdateSecretForTls)
	}
}
