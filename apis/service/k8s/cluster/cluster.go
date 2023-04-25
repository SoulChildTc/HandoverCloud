package cluster

import (
	"database/sql"
	"errors"
	"k8s.io/client-go/rest"
	"soul/apis/dao"
	"soul/apis/dto"
	"soul/apis/dto/k8s"
	"soul/global"
	log "soul/internal/logger"
	"soul/model"
	"soul/utils"
)

type Cluster struct{}

func (c *Cluster) GetClusterByName(clusterName string) *dto.K8sClusterInfo {
	config := global.K8s.Get(clusterName).Config
	info := &dto.K8sClusterInfo{
		ClusterName: clusterName,
		Host:        config.Host,
		BearerToken: config.BearerToken,
		TLSClientConfig: k8s.TlsClientConfig{
			Insecure: config.TLSClientConfig.Insecure,
			CertData: string(config.TLSClientConfig.CertData),
			KeyData:  string(config.TLSClientConfig.KeyData),
			CAData:   string(config.TLSClientConfig.CAData),
		},
	}

	return info
}

func (c *Cluster) GetClusterList() (clusterInfos []dto.K8sClusterInfo) {
	clusters := global.K8s.ListName()

	for _, clusterName := range clusters {
		info := c.GetClusterByName(clusterName)
		clusterInfos = append(clusterInfos, *info)
	}

	return clusterInfos
}

func (c *Cluster) AddCluster(info dto.K8sClusterInfo) error {
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

func (c *Cluster) UpdateCluster(info dto.K8sClusterInfo) error {
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
	err := dao.K8sCluster.DeleteClusterByName(clusterName)
	if err != nil {
		log.Error(err.Error())
		return errors.New("集群删除失败")
	}
	global.K8s.Remove(clusterName)
	return nil
}
