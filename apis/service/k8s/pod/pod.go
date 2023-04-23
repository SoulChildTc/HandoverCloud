package pod

import (
	"bytes"
	"context"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/remotecommand"
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

func (p *Pod) GetPodByName(clusterName, name, namespace string) (*corev1.Pod, error) {
	pod, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (p *Pod) GetPodList(clusterName, filterName, namespace string, limit, page int) (*httputil.PageResp, error) {
	pods, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
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

func (p *Pod) DeletePodByName(clusterName, podName, namespace string) (err error) {
	err = global.K8s.Use(clusterName).ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (p *Pod) GetPodLog(clusterName, podName, containerName, namespace string, line int64) (log string, err error) {
	if containerName == "" {
		pod, err := p.GetPodByName(clusterName, podName, namespace)
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

	req := global.K8s.Use(clusterName).ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)
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

func (p *Pod) GetPodContainers(clusterName, podName, namespace string) (containers []corev1.Container, err error) {
	pod, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod.Spec.Containers, nil
}

func (p *Pod) StartTerminal(clusterName, namespace, podName, containerName, shell string) (string, error) {
	sessionID, err := genTerminalSessionId()
	if err != nil {
		return "", err
	}

	terminalSessions.Set(sessionID, TerminalSession{
		id:       sessionID,
		bound:    make(chan error),
		sizeChan: make(chan remotecommand.TerminalSize),
	})

	// {"Op":"bind","SessionID":"db1888b4dd29e3c61540c56a5f7cfc22"}
	// {"Op":"stdin","Data":"ls\r","Cols":164,"Rows":41}
	go WaitForTerminal(clusterName, namespace, podName, containerName, shell, sessionID)
	return sessionID, err
}
