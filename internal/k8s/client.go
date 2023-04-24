package k8s

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"soul/model"
	"time"
)

type Client struct {
	ClientSet      *kubernetes.Clientset
	Config         *rest.Config
	CacheDiscovery discovery.DiscoveryInterface
	DynamicClient  *dynamic.DynamicClient
	Static         bool
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
func (c *ClusterMap) Update(clusterName string, client *Client) {
	c.Clusters[clusterName] = client
}

// List 列出所有集群名称
func (c *ClusterMap) List() (clusterList []*Client) {
	for _, client := range c.Clusters {
		clusterList = append(clusterList, client)
	}
	return
}

// ListName 列出所有集群名称
func (c *ClusterMap) ListName() (clusterList []string) {
	for name := range c.Clusters {
		clusterList = append(clusterList, name)
	}
	return
}

func (c *ClusterMap) Remove(clusterName string) {
	delete(c.Clusters, clusterName)
}

// IsStatic 集群类型
func (c *ClusterMap) IsStatic(clusterName string) bool {
	cluster, ok := c.Clusters[clusterName]
	if ok {
		return cluster.Static
	}
	return false
}

// AddClientWithKubeConfigOrInCluster 使用kubeconfig和incluster生成client, 并添加到clusters中
func (c *ClusterMap) AddClientWithKubeConfigOrInCluster(configPath string, inCluster bool) error {
	if inCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			return errors.New("[Init] Kubernetes config create failed." + err.Error())
		}

		client, err := c.NewClientWithRestConfig(config)
		err = c.Add("in-cluster", client)
		// 静态集群禁止修改
		c.Clusters["in-cluster"].Static = true
		if err != nil {
			return errors.New("[Init] Add Cluster failed." + err.Error())
		}
	}

	if configPath == "" {
		return nil
	}
	// 加载KUBECONFIG
	config, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		return errors.New("[Init] Load kubeconfig failed." + err.Error())
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
			return errors.New("[Init] Kubernetes config create failed." + err.Error())
		}

		// 集群名称命名
		client, err := c.NewClientWithRestConfig(contextConfig)
		err = c.Add(contextName, client)
		// 静态集群禁止修改
		c.Clusters[contextName].Static = true
		if err != nil {
			return errors.New("[Init] Add Cluster failed." + err.Error())
		}

	}
	return nil
}

func (c *ClusterMap) NewClientWithRestConfig(restConf *rest.Config) (*Client, error) {
	client := &Client{
		Config: restConf,
	}
	var err error

	// ClientSet
	client.ClientSet, err = kubernetes.NewForConfig(client.Config)
	if err != nil {
		return nil, errors.New("Kubernetes clientSet created failed. " + err.Error())
	}

	// DynamicClient
	client.DynamicClient, err = dynamic.NewForConfig(client.Config)
	if err != nil {
		return nil, errors.New("Kubernetes dynamic client created failed. " + err.Error())
	}

	// DiscoveryClient
	client.CacheDiscovery = c.newDiscoveryClient(client.Config)
	return client, nil
}

func (c *ClusterMap) newDiscoveryClient(restConf *rest.Config) (discoveryClient discovery.DiscoveryInterface) {
	discoveryClient, err := c.newDiskCacheDiscoveryClient(restConf)
	if err != nil {
		fmt.Println("Kubernetes DiskCacheDiscoveryClient created failed. Try MemCacheDiscoveryClient. " + err.Error())
	} else {
		//fmt.Println("Kubernetes DiskCacheDiscoveryClient created successful.")
		return
	}

	discoveryClient, err = c.newMemCacheDiscoveryClient(restConf)
	if err != nil {
		panic("Kubernetes MemCacheDiscoveryClient created failed. " + err.Error())
	}

	//fmt.Println("Kubernetes MemCacheDiscoveryClient created successful.")
	return

}

func (c *ClusterMap) newDiskCacheDiscoveryClient(restConf *rest.Config) (discoveryClient discovery.DiscoveryInterface, err error) {
	// DiskCacheDiscoveryClient
	discoveryClient, err = disk.NewCachedDiscoveryClientForConfig(
		restConf,
		"./cache/discovery",
		"./cache/http",
		3*time.Hour,
	)
	return
}

func (c *ClusterMap) newMemCacheDiscoveryClient(restConf *rest.Config) (discoveryClient discovery.DiscoveryInterface, err error) {
	// MemCacheDiscoveryClient
	discoveryClient, err = discovery.NewDiscoveryClientForConfig(restConf)
	if err != nil {
		return
	}
	discoveryClient = memory.NewMemCacheClient(discoveryClient)
	return
}

// InitClient 初始化所有集群Client
func InitClient(db *gorm.DB, configPath string, inCluster bool) *ClusterMap {
	// 初始化静态集群 - kubeconfig + in cluster
	err := clusters.AddClientWithKubeConfigOrInCluster(configPath, inCluster)
	if err != nil {
		panic(err.Error())
	}

	// 初始化DB中的集群
	var clusterList []model.K8sCluster
	res := db.Find(&clusterList)
	if res.Error != nil {
		panic(res.Error)
	}
	for _, cluster := range clusterList {
		restConf := &rest.Config{
			Host:        cluster.Host,
			BearerToken: cluster.BearerToken.String,
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: cluster.Insecure,
				CertData: []byte(cluster.CertData.String),
				KeyData:  []byte(cluster.KeyData.String),
				CAData:   []byte(cluster.CAData.String),
			},
		}

		client, err := clusters.NewClientWithRestConfig(restConf)
		if err != nil {
			panic(fmt.Sprintf("Cluster: %s. %s", cluster.ClusterName, err.Error()))
		}

		err = clusters.Add(cluster.ClusterName, client)
		if err != nil {
			panic(err.Error())
		}
	}

	return clusters
}
