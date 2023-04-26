package dto

import (
	"soul/apis/dto/k8s"
	"soul/apis/dto/system"
)

type (
	SystemRegister                   = system.Register
	SystemAdd                        = system.Add
	SystemLogin                      = system.Login
	SystemUserInfo                   = system.UserInfo
	SystemRoleInfo                   = system.RoleInfo
	K8sDeploymentCreate              = k8s.DeploymentCreate
	K8sSetImage                      = k8s.SetImage
	K8sIngressSimpleCreate           = k8s.IngressSimpleCreate
	K8sSvcSimpleCreate               = k8s.SvcSimpleCreate
	K8sSecretCreate                  = k8s.SecretCreate
	K8sSecretForDockerRegistryCreate = k8s.SecretForDockerRegistryCreate
	K8sSecretForTlsCreate            = k8s.SecretForTlsCreate
	K8sClusterInfo                   = k8s.ClusterInfo
	//SystemUserInfo system.UserInfo
)
