package prometheus

import (
	"context"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"soul/apis/service/k8s"
	"soul/global"
	"soul/utils/httputil"
)

var (
	serviceMonitorGVR = monitoringv1.SchemeGroupVersion.WithResource("servicemonitors")
)

type ServiceMonitor struct{}

func (s *ServiceMonitor) toStruct(unStructObj map[string]any, obj any) error {
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unStructObj, obj)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceMonitor) toCells(serviceMonitors []*monitoringv1.ServiceMonitor) []k8s.DataCell {
	cells := make([]k8s.DataCell, len(serviceMonitors))
	for i, item := range serviceMonitors {
		cells[i] = k8s.DataCell(serviceMonitorCell(*item))
	}
	return cells
}

func (s *ServiceMonitor) fromCells(cells []k8s.DataCell) []monitoringv1.ServiceMonitor {
	serviceMonitors := make([]monitoringv1.ServiceMonitor, len(cells))
	for i, item := range cells {
		serviceMonitors[i] = monitoringv1.ServiceMonitor(item.(serviceMonitorCell))
	}
	return serviceMonitors
}

func (s *ServiceMonitor) GetServiceMonitorByName(clusterName, name, namespace string) (map[string]any, error) {
	unStructObj, err := global.K8s.Use(clusterName).DynamicClient.
		Resource(serviceMonitorGVR).
		Namespace(namespace).
		Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return unStructObj.UnstructuredContent(), nil
}

func (s *ServiceMonitor) GetServiceMonitorList(clusterName, filterName, namespace string, limit, page int) (*httputil.PageResp, error) {
	serviceMonitors, err := global.K8s.Use(clusterName).DynamicClient.
		Resource(serviceMonitorGVR).
		Namespace(namespace).
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	serviceMonitorList := &monitoringv1.ServiceMonitorList{}
	err = s.toStruct(serviceMonitors.UnstructuredContent(), serviceMonitorList)
	if err != nil {
		return nil, err
	}

	selectableData := k8s.DataSelect{
		GenericDataList: s.toCells(serviceMonitorList.Items),
		DataSelect: &k8s.DataSelectQuery{
			Filter: &k8s.FilterQuery{
				Name: filterName,
			},
			Paginate: &k8s.PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	total := len(selectableData.Filter().GenericDataList)
	data := selectableData.Sort().Paginate()

	return &httputil.PageResp{
		Limit: limit,
		Page:  page,
		Total: total,
		Items: data.GenericDataList,
	}, nil
}

func (s *ServiceMonitor) DeleteServiceMonitorByName(clusterName, name, namespace string) (err error) {
	err = global.K8s.Use(clusterName).DynamicClient.
		Resource(serviceMonitorGVR).
		Namespace(namespace).
		Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceMonitor) CreateSimpleServiceMonitor(clusterName string) (err error) {
	//TODO 更轻松的创建serviceMonitor

	//_, err = global.K8s.Use(clusterName).DynamicClient.
	//	Resource(serviceMonitorGVR).
	//	Namespace(namespace).
	//	Create(context.TODO(), serviceMonitor, metav1.CreateOptions{})
	//if err != nil {
	//	return err
	//}
	return nil
}
