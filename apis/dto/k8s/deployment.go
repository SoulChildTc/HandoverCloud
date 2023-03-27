package k8s

type containerPort struct {
	Name     string `json:"name"`
	Port     int32  `json:"port"`
	Protocol string `json:"protocol"`
}

type strategy struct {
	MaxUnavailable string
	MaxSurge       string
}

type httpHealthCheck struct {
	HttpHealthPath string `json:"httpHealthPath" binding:"required" msg:"健康检查api路径不能为空"`
	HttpHealthPort string `json:"httpHealthPort"` // 默认使用容器的0号端口
}

type DeploymentCreate struct {
	Name                 string            `json:"name" binding:"required" msg:"Deployment名称不能为空"`
	Namespace            string            `json:"namespace" binding:"required" msg:"Namespace不能为空"`
	Replicas             int32             `json:"replicas,default=1"`
	Image                string            `json:"image" binding:"required" msg:"容器镜像不能为空"`
	ImagePullSecret      string            `json:"imagePullSecret"`
	Label                map[string]string `json:"label"`
	Cpu                  string            `json:"cpu,default='500m'"`
	Memory               string            `json:"memory,default='512Mi'"`
	ContainerPort        []containerPort   `json:"containerPort"`
	HttpHealthCheck      httpHealthCheck   `json:"httpHealthCheck"`
	RevisionHistoryLimit int32             `json:"-"` // 预留
	Strategy             strategy          `json:"-"` // 预留
}
