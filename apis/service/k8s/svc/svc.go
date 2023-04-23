package svc

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"soul/apis/dto"
	"soul/apis/service/k8s"
	"soul/apis/service/k8s/deployment"
	"soul/global"
	"soul/utils/httputil"
)

type Svc struct{}

func (s *Svc) toCells(services []corev1.Service) []k8s.DataCell {
	cells := make([]k8s.DataCell, len(services))
	for i, item := range services {
		cells[i] = k8s.DataCell(svcCell(item))
	}
	return cells
}

func (s *Svc) fromCells(cells []k8s.DataCell) []corev1.Service {
	services := make([]corev1.Service, len(cells))
	for i, item := range cells {
		services[i] = corev1.Service(item.(svcCell))
	}
	return services
}

func (s *Svc) GetSvcByName(clusterName, name, namespace string) (*corev1.Service, error) {
	svc, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *Svc) GetSvcList(clusterName, filterName, namespace string, limit, page int) (*httputil.PageResp, error) {
	services, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	selectableData := k8s.DataSelect{
		GenericDataList: s.toCells(services.Items),
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
	selectableData.Sort().Paginate()

	return &httputil.PageResp{
		Limit: limit,
		Page:  page,
		Total: total,
		Items: selectableData.GenericDataList,
	}, nil
}

func (s *Svc) DeleteSvcByName(clusterName, name, namespace string) (err error) {
	err = global.K8s.Use(clusterName).ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Svc) CreateSimpleSvc(clusterName string, svcSimpleCreate *dto.K8sSvcSimpleCreate) (err error) {
	svc, err := s.simpleSvcToService(clusterName, svcSimpleCreate)
	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Services(svc.Namespace).Create(context.TODO(), svc, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return
}

func (s *Svc) UpdateSimpleSvc(clusterName string, svcSimpleCreate *dto.K8sSvcSimpleCreate) (err error) {
	svc, err := s.simpleSvcToService(clusterName, svcSimpleCreate)
	if err != nil {
		return err
	}
	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Services(svc.Namespace).Update(context.TODO(), svc, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return
}

func (s *Svc) simpleSvcToService(clusterName string, svcSimpleCreate *dto.K8sSvcSimpleCreate) (*corev1.Service, error) {
	svcSimpleCreate.Type = "ClusterIP"

	if svcSimpleCreate.DeploymentName != "" {
		d := deployment.Deployment{}
		deploy, err := d.GetDeploymentByName(clusterName, svcSimpleCreate.DeploymentName, svcSimpleCreate.Namespace)
		if err != nil {
			return nil, err
		}
		svcSimpleCreate.Selector = deploy.Spec.Template.Labels
	}

	var ports []corev1.ServicePort
	for _, port := range svcSimpleCreate.Ports {
		ports = append(ports, corev1.ServicePort{
			Name:     port.Name,
			Protocol: corev1.Protocol(port.Protocol),
			Port:     port.ContainerPort,
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: port.ContainerPort,
			},
		})
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        svcSimpleCreate.Name,
			Namespace:   svcSimpleCreate.Namespace,
			Labels:      svcSimpleCreate.Labels,
			Annotations: map[string]string{"created-by": global.K8sManager},
		},
		Spec: corev1.ServiceSpec{
			Ports:    ports,
			Selector: svcSimpleCreate.Selector,
			Type:     corev1.ServiceType(svcSimpleCreate.Type),
		},
	}
	return svc, nil
}
