package ingress

import (
	"context"
	ingressv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"soul/apis/dto"
	"soul/apis/service/k8s"
	"soul/global"
	"soul/utils/httputil"
)

type Ingress struct{}

func (i *Ingress) toCells(ingresses []ingressv1.Ingress) []k8s.DataCell {
	cells := make([]k8s.DataCell, len(ingresses))
	for index, item := range ingresses {
		cells[index] = k8s.DataCell(ingressCell(item))
	}
	return cells
}

func (i *Ingress) fromCells(cells []k8s.DataCell) []ingressv1.Ingress {
	ingress := make([]ingressv1.Ingress, len(cells))
	for i, item := range cells {
		ingress[i] = ingressv1.Ingress(item.(ingressCell))
	}
	return ingress
}

func (i *Ingress) GetIngressByName(name, namespace string) (*ingressv1.Ingress, error) {
	ingress, err := global.K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return ingress, nil
}

func (i *Ingress) GetIngressList(filterName, namespace string, limit, page int) (*httputil.PageResp, error) {
	ingresses, err := global.K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	selectableData := k8s.DataSelect{
		GenericDataList: i.toCells(ingresses.Items),
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

func (i *Ingress) CreateSimpleIngress(ingressSimpleCreate *dto.K8sIngressSimpleCreate) (err error) {
	var rules []ingressv1.IngressRule
	var tlses []ingressv1.IngressTLS

	// 默认值
	pathType := ingressv1.PathTypePrefix
	if ingressSimpleCreate.Annotations != nil {
		ingressSimpleCreate.Annotations["created-by"] = global.K8sManager
	} else {
		ingressSimpleCreate.Annotations = map[string]string{"created-by": global.K8sManager}
	}

	rule := ingressv1.IngressRule{
		IngressRuleValue: ingressv1.IngressRuleValue{
			HTTP: &ingressv1.HTTPIngressRuleValue{Paths: []ingressv1.HTTPIngressPath{
				{
					Path:     ingressSimpleCreate.Rule.Path,
					PathType: &pathType,
					Backend: ingressv1.IngressBackend{
						Service: &ingressv1.IngressServiceBackend{
							Name: ingressSimpleCreate.Rule.Service,
							Port: ingressv1.ServiceBackendPort{
								Number: ingressSimpleCreate.Rule.ServicePort,
							},
						},
					},
				},
			}},
		},
	}

	// 将hosts生成多个IngressRule和TLS
	if len(ingressSimpleCreate.Rule.Hosts) != 0 {
		for _, host := range ingressSimpleCreate.Rule.Hosts {
			rule.Host = host
			rules = append(rules, rule)
		}
	} else {
		rules = append(rules, rule)
	}

	// 生成tls配置
	if ingressSimpleCreate.Tls != nil {
		for _, item := range ingressSimpleCreate.Tls {
			tls := ingressv1.IngressTLS{
				Hosts:      item.Hosts,
				SecretName: item.SecretName,
			}
			tlses = append(tlses, tls)
		}
	}

	ing := &ingressv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ingressSimpleCreate.Name,
			Namespace:   ingressSimpleCreate.Namespace,
			Labels:      map[string]string{},
			Annotations: ingressSimpleCreate.Annotations,
		},
		Spec: ingressv1.IngressSpec{
			IngressClassName: pointer.String(ingressSimpleCreate.IngressClassName),
			Rules:            rules,
			TLS:              tlses,
		},
	}

	_, err = global.K8s.ClientSet.NetworkingV1().Ingresses(ing.Namespace).Create(context.TODO(), ing, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return
}

func (i *Ingress) UpdateSimpleIngress(ingressSimpleCreate *dto.K8sIngressSimpleCreate) (err error) {
	var rules []ingressv1.IngressRule
	var tlses []ingressv1.IngressTLS

	// 默认值
	pathType := ingressv1.PathTypePrefix
	if ingressSimpleCreate.Annotations != nil {
		ingressSimpleCreate.Annotations["created-by"] = global.K8sManager
	} else {
		ingressSimpleCreate.Annotations = map[string]string{"created-by": global.K8sManager}
	}

	rule := ingressv1.IngressRule{
		IngressRuleValue: ingressv1.IngressRuleValue{
			HTTP: &ingressv1.HTTPIngressRuleValue{Paths: []ingressv1.HTTPIngressPath{
				{
					Path:     ingressSimpleCreate.Rule.Path,
					PathType: &pathType,
					Backend: ingressv1.IngressBackend{
						Service: &ingressv1.IngressServiceBackend{
							Name: ingressSimpleCreate.Rule.Service,
							Port: ingressv1.ServiceBackendPort{
								Number: ingressSimpleCreate.Rule.ServicePort,
							},
						},
					},
				},
			}},
		},
	}

	// 将hosts生成多个IngressRule和TLS
	if len(ingressSimpleCreate.Rule.Hosts) != 0 {
		for _, host := range ingressSimpleCreate.Rule.Hosts {
			rule.Host = host
			rules = append(rules, rule)
		}
	} else {
		rules = append(rules, rule)
	}

	// 生成tls配置
	if ingressSimpleCreate.Tls != nil {
		for _, item := range ingressSimpleCreate.Tls {
			tls := ingressv1.IngressTLS{
				Hosts:      item.Hosts,
				SecretName: item.SecretName,
			}
			tlses = append(tlses, tls)
		}
	}

	ing := &ingressv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ingressSimpleCreate.Name,
			Namespace:   ingressSimpleCreate.Namespace,
			Labels:      map[string]string{},
			Annotations: ingressSimpleCreate.Annotations,
		},
		Spec: ingressv1.IngressSpec{
			IngressClassName: pointer.String(ingressSimpleCreate.IngressClassName),
			Rules:            rules,
			TLS:              tlses,
		},
	}

	_, err = global.K8s.ClientSet.NetworkingV1().Ingresses(ing.Namespace).Update(context.TODO(), ing, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return
}
