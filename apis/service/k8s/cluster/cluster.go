package cluster

import (
	"context"
	"database/sql"
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sort"
	"soul/apis/dao"
	"soul/apis/dto"
	"soul/apis/dto/k8s"
	"soul/global"
	log "soul/internal/logger"
	"soul/model"
	"soul/utils"
)

type Cluster struct{}

func (c *Cluster) GetClusterByName(clusterName string, force bool) *dto.K8sClusterInfo {
	cluster := global.K8s.Get(clusterName)
	if force {
		version, err := cluster.CacheDiscovery.ServerVersion()
		if err != nil {
			cluster.Status = err.Error()
		} else {
			cluster.Version = version.String()
			nodes, err := cluster.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			cluster.NodeNum = uint(len(nodes.Items))
			if err != nil && cluster.Status != err.Error() {
				cluster.Status += ". " + err.Error()
			} else {
				cluster.Status = "运行中"
			}
		}
	}
	info := &dto.K8sClusterInfo{
		ClusterCreate: k8s.ClusterCreate{
			ClusterName: clusterName,
			Host:        cluster.Config.Host,
			BearerToken: cluster.Config.BearerToken,
			TLSClientConfig: k8s.TlsClientConfig{
				Insecure: cluster.Config.TLSClientConfig.Insecure,
				CertData: string(cluster.Config.TLSClientConfig.CertData),
				KeyData:  string(cluster.Config.TLSClientConfig.KeyData),
				CAData:   string(cluster.Config.TLSClientConfig.CAData),
			},
		},
		Version: cluster.Version,
		Status:  cluster.Status,
		NodeNum: cluster.NodeNum,
	}

	return info
}

func (c *Cluster) GetClusterList(force bool) (clusterInfos []dto.K8sClusterInfo) {
	clusters := global.K8s.ListName()

	for _, clusterName := range clusters {
		info := c.GetClusterByName(clusterName, force)
		if info == nil {
			continue
		}
		clusterInfos = append(clusterInfos, *info)
	}

	sort.Slice(clusterInfos, func(i, j int) bool {
		return clusterInfos[i].ClusterName < clusterInfos[j].ClusterName
	})
	return clusterInfos
}

func (c *Cluster) AddCluster(info dto.K8sClusterCreate) error {
	if global.K8s.Get(info.ClusterName) != nil {
		return errors.New("集群已存在")
	}

	// 创建reset client
	restConf := &rest.Config{
		Host:        info.Host,
		BearerToken: info.BearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: info.TLSClientConfig.Insecure,
			CertData: []byte(info.TLSClientConfig.CertData),
			KeyData:  []byte(info.TLSClientConfig.KeyData),
			CAData:   []byte(info.TLSClientConfig.CAData),
		},
	}
	client, err := global.K8s.NewClientWithRestConfig(restConf)
	if err != nil {
		log.Error(err.Error())
		return errors.New("创建集群失败")
	}

	_ = global.K8s.Add(info.ClusterName, client)

	// 存入数据库
	cluster := &model.K8sCluster{
		ClusterName: info.ClusterName,
		Host:        info.Host,
		BearerToken: sql.NullString{
			String: info.BearerToken,
			Valid:  true,
		},
		Insecure: info.TLSClientConfig.Insecure,
		CertData: sql.NullString{
			String: info.TLSClientConfig.CertData,
			Valid:  true,
		},
		KeyData: sql.NullString{
			String: info.TLSClientConfig.KeyData,
			Valid:  true,
		},
		CAData: sql.NullString{
			String: info.TLSClientConfig.CAData,
			Valid:  true,
		},
	}
	err = dao.K8sCluster.CreateCluster(cluster)
	if err != nil {
		log.Error(err.Error())
		global.K8s.Remove(info.ClusterName)
		return errors.New("集群创建失败")
	}

	return nil
}

func (c *Cluster) UpdateCluster(info dto.K8sClusterCreate) error {
	if global.K8s.Get(info.ClusterName) == nil {
		return errors.New("集群不存在")
	}
	if global.K8s.IsStatic(info.ClusterName) {
		return errors.New("静态集群无法修改")
	}

	cluster := &model.K8sCluster{
		ClusterName: info.ClusterName,
		Host:        info.Host,
		BearerToken: utils.ScanNullString(info.BearerToken),
		Insecure:    info.TLSClientConfig.Insecure,
		CertData:    utils.ScanNullString(info.TLSClientConfig.CertData),
		KeyData:     utils.ScanNullString(info.TLSClientConfig.KeyData),
		CAData:      utils.ScanNullString(info.TLSClientConfig.CAData),
	}
	err := dao.K8sCluster.UpdateCluster(cluster)
	if err != nil {
		log.Error(err.Error())
		return errors.New("集群更新失败")
	}
	restConf := &rest.Config{
		Host:        info.Host,
		BearerToken: info.BearerToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: info.TLSClientConfig.Insecure,
			CertData: []byte(info.TLSClientConfig.CertData),
			KeyData:  []byte(info.TLSClientConfig.KeyData),
			CAData:   []byte(info.TLSClientConfig.CAData),
		},
	}
	client, err := global.K8s.NewClientWithRestConfig(restConf)
	if err != nil {
		log.Error(err.Error())
		return errors.New("集群更新失败")
	}
	global.K8s.Update(info.ClusterName, client)
	return nil
}

func (c *Cluster) DeleteCluster(clusterName string) error {
	if global.K8s.Get(clusterName).Static == true {
		return errors.New("静态集群不能删除")
	}
	err := dao.K8sCluster.DeleteClusterByName(clusterName)
	if err != nil {
		log.Error(err.Error())
		return errors.New("集群删除失败" + err.Error())
	}
	global.K8s.Remove(clusterName)
	return nil
}
