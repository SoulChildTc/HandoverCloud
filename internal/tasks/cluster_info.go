package tasks

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"soul/internal/k8s"
	log "soul/internal/logger"
	"time"
)

func ClusterInfoTask() {
	for {
		clusters := k8s.GetClusterMap()
		for _, cluster := range clusters.ListName() {
			version, err := clusters.Get(cluster).CacheDiscovery.ServerVersion()
			if err != nil {
				clusters.Get(cluster).Status = err.Error()
			} else {
				clusters.Get(cluster).Version = version.String()

				nodes, err := clusters.Get(cluster).ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
				clusters.Get(cluster).NodeNum = uint(len(nodes.Items))
				if err != nil && clusters.Get(cluster).Status != err.Error() {
					clusters.Get(cluster).Status += ". " + err.Error()
				} else {
					clusters.Get(cluster).Status = "运行中"
				}
			}
		}
		log.Debug("Update cluster information completed.")
		time.Sleep(time.Hour * 1)
	}
}
