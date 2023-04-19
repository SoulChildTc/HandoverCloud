package global

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type k8s struct {
	ClientSet      *kubernetes.Clientset
	Config         *rest.Config
	CacheDiscovery discovery.DiscoveryInterface
	DynamicClient  *dynamic.DynamicClient
}

var (
	K8s = &k8s{}
)
