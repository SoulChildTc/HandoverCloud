package namespace

import (
	"context"
	"encoding/json"
	"errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"soul/apis/service/k8s"
	"soul/global"
	"soul/utils/httputil"
)

type Namespace struct{}

func (n *Namespace) toCells(namespaces []corev1.Namespace) []k8s.DataCell {
	cells := make([]k8s.DataCell, len(namespaces))
	for i, item := range namespaces {
		cells[i] = k8s.DataCell(namespaceCell(item))
	}
	return cells
}

func (n *Namespace) fromCells(cells []k8s.DataCell) []corev1.Namespace {
	namespaces := make([]corev1.Namespace, len(cells))
	for i, item := range cells {
		namespaces[i] = corev1.Namespace(item.(namespaceCell))
	}
	return namespaces
}

func (n *Namespace) GetNamespaceByName(clusterName, name string) (*corev1.Namespace, error) {
	namespace, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return namespace, nil
}

func (n *Namespace) GetNamespaceList(clusterName, filterName string, limit, page int) (*httputil.PageResp, error) {
	namespaces, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	selectableData := k8s.DataSelect{
		GenericDataList: n.toCells(namespaces.Items),
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

func (n *Namespace) DeleteNamespaceByName(clusterName, namespace string) (err error) {
	err = global.K8s.Use(clusterName).ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (n *Namespace) CreateNamespace(clusterName, content string) error {
	ns := &corev1.Namespace{}
	err := json.Unmarshal([]byte(content), ns)
	if err != nil {
		return errors.New("反序列化失败,请检查yaml。" + err.Error())
	}

	if ns.Annotations != nil {
		ns.Annotations["created-by"] = global.K8sManager
	} else {
		ns.Annotations = map[string]string{"created-by": global.K8sManager}
	}

	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return errors.New("创建namespace失败," + err.Error())
	}

	return nil
}
