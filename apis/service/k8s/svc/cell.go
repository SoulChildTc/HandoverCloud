package svc

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type svcCell corev1.Service

func (p svcCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p svcCell) GetName() string {
	return p.Name
}
