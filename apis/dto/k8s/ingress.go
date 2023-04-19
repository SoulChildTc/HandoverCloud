package k8s

type ruleConfig struct {
	Hosts       []string `json:"hosts"`
	Path        string   `json:"path,default='/'"`
	Service     string   `json:"service" binding:"required" msg:"Service不能为空"`
	ServicePort int32    `json:"servicePort"`
}

type tlsConfig struct {
	Hosts      []string `json:"host"`
	SecretName string   `json:"secretName"`
}

type IngressSimpleCreate struct {
	Name             string            `json:"name" binding:"required" msg:"Deployment名称不能为空"`
	Namespace        string            `json:"namespace" binding:"required" msg:"Namespace不能为空"`
	Annotations      map[string]string `json:"annotations"`
	IngressClassName string            `json:"ingressClassName"`
	Rule             ruleConfig        `json:"rule" binding:"required" msg:"规则不能为空"`
	Tls              []tlsConfig       `json:"tls"`
}
