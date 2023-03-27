package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	"time"
)

type deploymentCell appsv1.Deployment

func (p deploymentCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p deploymentCell) GetName() string {
	return p.Name
}
