package pod

import (
	"bytes"
	"context"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"soul/apis/service/k8s"
	"soul/global"
	"soul/utils/httputil"
)

type Pod struct{}

// 类型转换 corev1.Pod -> podCell
func (p *Pod) toCells(pods []corev1.Pod) []k8s.DataCell {
	cells := make([]k8s.DataCell, len(pods))
	for i, item := range pods {
		cells[i] = k8s.DataCell(podCell(item))
	}
	return cells
}

// DataCell -> corev1.Pod
func (p *Pod) fromCells(cells []k8s.DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i, item := range cells {
		pods[i] = corev1.Pod(item.(podCell))
	}
	return pods
}

func (p *Pod) GetPodByName(name, namespace string) (*corev1.Pod, error) {
	pod, err := global.K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (p *Pod) GetPodList(filterName, namespace string, limit, page int) (*httputil.PageResp, error) {
	pods, err := global.K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	selectableData := k8s.DataSelect{
		GenericDataList: p.toCells(pods.Items),
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

func (p *Pod) DeletePodByName(podName, namespace string) (err error) {
	err = global.K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (p *Pod) GetPodLog(podName, containerName, namespace string, line int64) (log string, err error) {
	if containerName == "" {
		pod, err := p.GetPodByName(podName, namespace)
		if err != nil {
			return "", err
		}
		containerName = pod.Spec.Containers[0].Name
	}

	option := &corev1.PodLogOptions{
		Container:                    containerName,
		Follow:                       false,
		Previous:                     false,
		SinceSeconds:                 nil,
		SinceTime:                    nil,
		Timestamps:                   false,
		TailLines:                    &line,
		LimitBytes:                   nil,
		InsecureSkipTLSVerifyBackend: false,
	}

	req := global.K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
	logReader, err := req.Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer logReader.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logReader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (p *Pod) GetPodContainers(podName, namespace string) (containers []corev1.Container, err error) {
	pod, err := global.K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod.Spec.Containers, nil
}
