package k8s

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"soul/global"
)

func NewK8sClient() (*rest.Config, *kubernetes.Clientset) {
	var (
		err       error
		config    *rest.Config
		clientSet *kubernetes.Clientset
	)

	if global.Config.KubeConfig == "" {
		fmt.Println("[Init] In Kubernetes Cluster Running...")
		config, err = rest.InClusterConfig()
	} else {
		fmt.Printf("[Init] Using kubeconfig file: %s .\n", global.Config.KubeConfig)
		config, err = clientcmd.BuildConfigFromFlags("", global.Config.KubeConfig)
	}

	if err != nil {
		panic("[Init] Kubernetes config parse failed." + err.Error())
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic("[Init] Kubernetes client initialization failed." + err.Error())
	} else {
		fmt.Println("[Init] Kubernetes clientSet initialization successful.")
	}

	return config, clientSet
}
