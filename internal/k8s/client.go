package k8s

import (
	"fmt"
	"k8s.io/client-go/discovery/cached/disk"
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

	global.K8s.ClientSet, err = kubernetes.NewForConfig(global.K8s.Config)
	if err != nil {
		panic("[Init] Kubernetes client initialization failed." + err.Error())
	} else {
		fmt.Println("[Init] Kubernetes clientSet initialization successful.")
	}

	global.K8s.CacheDiscovery, err = disk.NewCachedDiscoveryClientForConfig(
		global.K8s.Config,
		"./cache/discovery",
		"./cache/http",
		1*time.Hour,
	)
	if err != nil {
		fmt.Println("[Init] Kubernetes CacheDiscovery Client initialization failed." + err.Error())
	}
}
