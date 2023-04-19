package namespace

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type namespaceCell corev1.Namespace

func (p namespaceCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p namespaceCell) GetName() string {
	return p.Name
}
