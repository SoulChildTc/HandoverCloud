package ingress

import (
	ingressv1 "k8s.io/api/networking/v1"
	"time"
)

type ingressCell ingressv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressCell) GetName() string {
	return i.Name
}
