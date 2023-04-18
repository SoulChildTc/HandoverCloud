package global

import (
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type k8s struct {
	ClientSet      *kubernetes.Clientset
	Config         *rest.Config
	CacheDiscovery *disk.CachedDiscoveryClient
}

var (
	K8s = &k8s{}
)
