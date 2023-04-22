package k8s

import (
	"errors"
	"fmt"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"soul/global"
	log "soul/internal/logger"
	"time"
)

type Client struct {
	ClientSet      *kubernetes.Clientset
	Config         *rest.Config
	CacheDiscovery discovery.DiscoveryInterface
	DynamicClient  *dynamic.DynamicClient
}

// ClusterMap 用于存储多个集群的client
type ClusterMap struct {
	Clusters map[string]*Client
}

var clusters = &ClusterMap{make(map[string]*Client)}

// Use 获取某个集群的client
func (c *ClusterMap) Use(clusterName string) *Client {
	return c.Get(clusterName)
}

// Get 获取某个集群的client
func (c *ClusterMap) Get(clusterName string) *Client {
	return c.Clusters[clusterName]
}

// Add 添加集群
func (c *ClusterMap) Add(clusterName string, client *Client) error {
	if c.Clusters[clusterName] != nil {
		return errors.New("cluster exists")
	}
	c.Clusters[clusterName] = client
	return nil
}

// Update 更新某个集群的client
func (c *ClusterMap) Update(clusterName string, client *Client) error {
	c.Clusters[clusterName] = client
	return nil
}

// InitClient 初始化所有集群Client
func InitClient() *ClusterMap {
	// 初始化静态集群 - 配置文件指定的
	NewClientWithKubeConfig()

	// 初始化mysql中的集群
	//for index, _ := range "aa" {
	//}

	return clusters
}

func NewClientWithKubeConfig() {
	var (
		err error
	)

	if global.Config.KubeConfig == "" {
		fmt.Println("[Init] In Kubernetes Cluster Running...")
		config, err := rest.InClusterConfig()
		if err != nil {
			panic("[Init] Kubernetes config create failed." + err.Error())
		}

		err = clusters.Add("inCluster", newClientWithRestConfig(config))
		if err != nil {
			panic("[Init] Add Cluster failed." + err.Error())
		}
	}

	// 加载KUBECONFIG
	config, err := clientcmd.LoadFromFile(global.Config.KubeConfig)
	if err != nil {
		panic("[Init] Load kubeconfig failed." + err.Error())
	}

	// 遍历上下文
	for contextName := range config.Contexts {
		clusterName := config.Contexts[contextName].Cluster
		authInfo := config.Contexts[contextName].AuthInfo
		namespace := config.Contexts[contextName].Namespace

		// 为每个上下文创建rest config
		contextConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{
			Context: clientcmdapi.Context{
				Cluster:   clusterName,
				AuthInfo:  authInfo,
				Namespace: namespace,
			},
		}).ClientConfig()
		if err != nil {
			panic("[Init] Kubernetes config create failed." + err.Error())
		}

		// 集群名称命名为 集群名称@用户
		err = clusters.Add(
			fmt.Sprintf("%s@%s", clusterName, authInfo),
			newClientWithRestConfig(contextConfig),
		)

		if err != nil {
			panic("[Init] Add Cluster failed." + err.Error())
		}

	}
	log.Debug("%v", clusters)
}

func newClientWithRestConfig(restConf *rest.Config) *Client {
	client := &Client{
		Config: restConf,
	}
	var err error

	// ClientSet
	client.ClientSet, err = kubernetes.NewForConfig(client.Config)
	if err != nil {
		panic("[Init] Kubernetes clientSet initialization failed." + err.Error())
	} else {
		fmt.Println("[Init] Kubernetes clientSet initialization successful.")
	}

	// DynamicClient
	client.DynamicClient, err = dynamic.NewForConfig(client.Config)
	if err != nil {
		panic("[Init] Kubernetes dynamic client initialization failed." + err.Error())
	} else {
		fmt.Println("[Init] Kubernetes dynamic client initialization successful.")
	}

	// DiscoveryClient
	client.CacheDiscovery = newDiscoveryClient(client.Config)
	return client
}

func newDiscoveryClient(restConf *rest.Config) (discoveryClient discovery.DiscoveryInterface) {
	discoveryClient, err := newDiskCacheDiscoveryClient(restConf)
	if err != nil {
		fmt.Println("[Init] Kubernetes DiskCacheDiscoveryClient initialization failed. Try MemCacheDiscoveryClient." + err.Error())
	} else {
		fmt.Println("[Init] Kubernetes DiskCacheDiscoveryClient initialization successful.")
		return
	}

	discoveryClient, err = newMemCacheDiscoveryClient(restConf)
	if err != nil {
		panic("[Init] Kubernetes MemCacheDiscoveryClient initialization failed." + err.Error())
	}

	fmt.Println("[Init] Kubernetes MemCacheDiscoveryClient initialization successful.")
	return

}

func newDiskCacheDiscoveryClient(restConf *rest.Config) (discoveryClient discovery.DiscoveryInterface, err error) {
	// DiskCacheDiscoveryClient
	discoveryClient, err = disk.NewCachedDiscoveryClientForConfig(
		restConf,
		"./cache/discovery",
		"./cache/http",
		3*time.Hour,
	)
	return
}

func newMemCacheDiscoveryClient(restConf *rest.Config) (discoveryClient discovery.DiscoveryInterface, err error) {
	// MemCacheDiscoveryClient
	discoveryClient, err = discovery.NewDiscoveryClientForConfig(restConf)
	if err != nil {
		return
	}
	discoveryClient = memory.NewMemCacheClient(discoveryClient)
	return
}
