package k8s

import (
	"github.com/gin-gonic/gin"
	k8sdeployment "soul/apis/controller/k8s/deployment"
	k8singress "soul/apis/controller/k8s/ingress"
	k8snamespace "soul/apis/controller/k8s/namespace"
	k8spod "soul/apis/controller/k8s/pod"
)

// k8s模块路由

func RegisterRoute(r *gin.RouterGroup) {
	pod := r.Group("/pod")
	{
		pod.GET("/", k8spod.GetPodList)
		pod.GET("/:namespace", k8spod.GetPodList)
		pod.GET("/:namespace/:podName", k8spod.GetPodByName)
		pod.DELETE("/:namespace/:podName", k8spod.DeletePodByName)
		pod.GET("/:namespace/:podName/log", k8spod.GetPodLog)
		pod.GET("/:namespace/:podName/containers", k8spod.GetPodContainers)
		pod.GET("/:namespace/:podName/shell", k8spod.ExecContainer)
	}

	deployment := r.Group("/deployment")
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

	ingress := r.Group("/ingress")
	{
		ingress.GET("/", k8singress.GetIngressList)
		ingress.GET("/:namespace", k8singress.GetIngressList)
		ingress.GET("/:namespace/:ingressName", k8singress.GetIngressByName)
		ingress.POST("/", k8singress.CreateSimpleIngress)
		ingress.PUT("/", k8singress.UpdateSimpleIngress)
	}

	namespace := r.Group("/namespace")
	{
		namespace.GET("/", k8snamespace.GetNamespaceList)
		namespace.GET("/:namespaceName", k8snamespace.GetNamespaceByName)
		namespace.POST("/", k8snamespace.CreateNamespace)
		namespace.DELETE("/:namespaceName", k8snamespace.DeleteNamespaceByName)
	}

}
