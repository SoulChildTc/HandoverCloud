package k8s

import (
	"encoding/base64"
	"fmt"
)

type SecretCreate struct {
	Name      string            `json:"name" binding:"required" msg:"Secret名称不能为空"`
	Namespace string            `json:"namespace" binding:"required" msg:"Namespace不能为空"`
	Data      map[string]string `json:"data" binding:"required" msg:"Data不能为空"`
}

type SecretForDockerRegistryCreate struct {
	Name      string          `json:"name" binding:"required" msg:"Secret名称不能为空"`
	Namespace string          `json:"namespace" binding:"required" msg:"Namespace不能为空"`
	Auths     map[string]Auth `json:"auths" binding:"required" msg:"仓库地址和认证信息不能为空"`
}

func (s *SecretForDockerRegistryCreate) ToDockerconfig() map[string]map[string]string {
	data := map[string]map[string]string{
		"auths": {},
	}

	for k, v := range s.Auths {
		b64auth := []byte(fmt.Sprintf("%s:%s", v.User, v.Password))
		data["auths"][k] = base64.StdEncoding.EncodeToString(b64auth)
	}

	return data
}

type Auth struct {
	User     string `json:"user" binding:"required" msg:"用户名不能为空"`
	Password string `json:"password" binding:"required" msg:"密码不能为空"`
}

type SecretForTlsCreate struct {
	Name        string `json:"name" binding:"required" msg:"Secret名称不能为空"`
	Namespace   string `json:"namespace" binding:"required" msg:"Namespace不能为空"`
	Certificate string `json:"certificate" binding:"required" msg:"证书不能为空"`
	PrivateKey  string `json:"privateKey" binding:"required" msg:"私钥不能为空"`
}
