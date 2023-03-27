package pod

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}
