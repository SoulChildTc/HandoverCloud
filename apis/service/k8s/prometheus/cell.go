package prometheus

import (
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"time"
)

type serviceMonitorCell monitoringv1.ServiceMonitor

func (s serviceMonitorCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s serviceMonitorCell) GetName() string {
	return s.Name
}
