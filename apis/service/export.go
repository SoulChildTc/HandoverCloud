package service

import (
	"soul/apis/service/k8s/deployment"
	"soul/apis/service/k8s/pod"
	"soul/apis/service/system/user"
)

var (
	SystemUser    user.User
	K8sPod        pod.Pod
	K8sDeployment deployment.Deployment
)
