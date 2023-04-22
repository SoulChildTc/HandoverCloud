package secret

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type secretCell corev1.Secret

func (s secretCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s secretCell) GetName() string {
	return s.Name
}
