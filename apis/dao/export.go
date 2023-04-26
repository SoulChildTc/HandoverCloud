package dao

import (
	"soul/apis/dao/k8s"
	"soul/apis/dao/system"
)

var (
	SystemUser     system.User
	SystemRole     system.Role
	SystemInitData system.InitData
	K8sCluster     k8s.Cluster
)
