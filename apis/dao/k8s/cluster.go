package k8s

import (
	"errors"
	"gorm.io/gorm"
	"soul/global"
	log "soul/internal/logger"
	"soul/model"
	"soul/model/common"
	"soul/utils"
)

type Cluster struct{}

func (c *Cluster) CreateCluster(cluster *model.K8sCluster) error {
	result := global.DB.Create(cluster)
	return result.Error
}

func (c *Cluster) UpdateCluster(cluster *model.K8sCluster) error {
	result := global.DB.
		Model(&model.K8sCluster{}).
		Select("*").
		Where("cluster_name = ?", cluster.ClusterName).
		Updates(cluster)
	return result.Error
}

func (c *Cluster) ListCluster() (clusters []model.K8sCluster) {
	global.DB.Find(&clusters)
	return
}

func (c *Cluster) GetClusterByName(clusterName string) *model.K8sCluster {
	cluster := &model.K8sCluster{ClusterName: clusterName}
	if err := global.DB.First(cluster).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(err.Error())
		}
		return nil
	}
	return cluster
}

func (c *Cluster) GetClusterById(id uint) error {
	cluster := &model.K8sCluster{ID: common.ID{ID: id}}
	if err := global.DB.First(cluster).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error(err.Error())
		}
		return nil
	}
	return nil
}

func (c *Cluster) DeleteClusterByName(clusterName string) error {
	bakName := clusterName + "_" + utils.RandStringBytesMaskImprSrc(6)
	result := global.DB.
		Model(&model.K8sCluster{}).
		Where("cluster_name = ?", clusterName).
		Select("cluster_name").
		UpdateColumns(&model.K8sCluster{ClusterName: bakName})

	if result.Error != nil {
		return result.Error
	}

	result = global.DB.
		Model(&model.K8sCluster{}).
		Where("cluster_name = ?", bakName).
		Delete(&model.K8sCluster{})
	return result.Error
}
