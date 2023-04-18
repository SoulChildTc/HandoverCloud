package dto

import (
	"soul/apis/dto/k8s"
	"soul/apis/dto/system"
)

type (
	SystemRegister      = system.Register
	SystemLogin         = system.Login
	K8sDeploymentCreate = k8s.DeploymentCreate
	K8sSetImage         = k8s.SetImage
	//SystemUserInfo system.UserInfo
)