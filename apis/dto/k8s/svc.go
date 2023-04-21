package k8s

type ports struct {
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	ContainerPort int32  `json:"containerPort"`
}

type SvcSimpleCreate struct {
	Name           string            `json:"name" binding:"required" msg:"Service名称不能为空"`
	Namespace      string            `json:"namespace" binding:"required" msg:"Namespace不能为空"`
	Labels         map[string]string `json:"labels"`
	DeploymentName string            `json:"deploymentName"`
	Selector       map[string]string `json:"selector"`
	Type           string            `json:"-"` // 预留 默认,ClusterIP
	Ports          []ports           `json:"ports"`
}
