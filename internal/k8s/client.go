package k8s

import (
	"fmt"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"soul/global"
	"time"
)

func NewK8sClient() {
	var (
		err error
	)

	// Rest Config
	if global.Config.KubeConfig == "" {
		fmt.Println("[Init] In Kubernetes Cluster Running...")
		global.K8s.Config, err = rest.InClusterConfig()
	} else {
		fmt.Printf("[Init] Using kubeconfig file: %s .\n", global.Config.KubeConfig)
		global.K8s.Config, err = clientcmd.BuildConfigFromFlags("", global.Config.KubeConfig)
	}

	if err != nil {
		panic("[Init] Kubernetes config parse failed." + err.Error())
	}

	// ClientSet
	global.K8s.ClientSet, err = kubernetes.NewForConfig(global.K8s.Config)
	if err != nil {
		panic("[Init] Kubernetes clientSet initialization failed." + err.Error())
	} else {
		fmt.Println("[Init] Kubernetes clientSet initialization successful.")
	}

	// DynamicClient
	global.K8s.DynamicClient, err = dynamic.NewForConfig(global.K8s.Config)
	if err != nil {
		panic("[Init] Kubernetes dynamic client initialization failed." + err.Error())
	} else {
		fmt.Println("[Init] Kubernetes dynamic client initialization successful.")
	}

	// DiscoveryClient
	global.K8s.CacheDiscovery = newDiscoveryClient()

}

func newDiscoveryClient() (discoveryClient discovery.DiscoveryInterface) {
	discoveryClient, err := newDiskCacheDiscoveryClient()
	if err != nil {
		fmt.Println("[Init] Kubernetes DiskCacheDiscoveryClient initialization failed. Try MemCacheDiscoveryClient." + err.Error())
	} else {
		fmt.Println("[Init] Kubernetes DiskCacheDiscoveryClient initialization successful.")
		return
	}

	discoveryClient, err = newMemCacheDiscoveryClient()
	if err != nil {
		panic("[Init] Kubernetes MemCacheDiscoveryClient initialization failed." + err.Error())
	}

	fmt.Println("[Init] Kubernetes MemCacheDiscoveryClient initialization successful.")
	return

}

func newDiskCacheDiscoveryClient() (discoveryClient discovery.DiscoveryInterface, err error) {
	// DiskCacheDiscoveryClient
	discoveryClient, err = disk.NewCachedDiscoveryClientForConfig(
		global.K8s.Config,
		"./cache/discovery",
		"./cache/http",
		3*time.Hour,
	)
	return
}

func newMemCacheDiscoveryClient() (discoveryClient discovery.DiscoveryInterface, err error) {
	// MemCacheDiscoveryClient
	discoveryClient, err = discovery.NewDiscoveryClientForConfig(global.K8s.Config)
	if err != nil {
		return
	}
	discoveryClient = memory.NewMemCacheClient(discoveryClient)
	return
}
