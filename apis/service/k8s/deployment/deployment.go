package deployment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"soul/apis/dto"
	"soul/apis/service/k8s"
	"soul/global"
	"soul/utils/httputil"
	"time"
)

type Deployment struct{}

func (d *Deployment) toCells(deployments []appsv1.Deployment) []k8s.DataCell {
	cells := make([]k8s.DataCell, len(deployments))
	for i, item := range deployments {
		cells[i] = k8s.DataCell(deploymentCell(item))
	}
	return cells
}

func (d *Deployment) fromCells(cells []k8s.DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i, item := range cells {
		deployments[i] = appsv1.Deployment(item.(deploymentCell))
	}
	return deployments
}

func (d *Deployment) GetDeploymentByName(name, namespace string) (*appsv1.Deployment, error) {
	pod, err := global.K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (d *Deployment) GetDeploymentList(filterName, namespace string, limit, page int) (*httputil.PageResp, error) {
	deployments, err := global.K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	selectableData := k8s.DataSelect{
		GenericDataList: d.toCells(deployments.Items),
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

func (d *Deployment) GetDeploymentPods(name, namespace string) (pods *corev1.PodList, err error) {
	deployment, err := d.GetDeploymentByName(name, namespace)
	if err != nil {
		return nil, err
	}

	// 将map转换为类似这种形式 app=client-go-deploy,name
	selector := metav1.FormatLabelSelector(deployment.Spec.Selector)

	pods, err = global.K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector,
	})

	if err != nil {
		return nil, err
	}

	return pods, nil
}

func (d *Deployment) CreateDeployment(deploymentCreate *dto.K8sDeploymentCreate) (err error) {
	deploymentCreate.RevisionHistoryLimit = 10
	deploymentCreate.Strategy.MaxUnavailable = "20%"
	deploymentCreate.Strategy.MaxSurge = "20%"
	maxUnavailable := intstr.Parse(deploymentCreate.Strategy.MaxUnavailable)
	maxSurge := intstr.Parse(deploymentCreate.Strategy.MaxSurge)
	deploymentCreate.Label["handovercloud.soulchild.cn/app"] = deploymentCreate.Name

	// 组装containerPort
	var containerPort []corev1.ContainerPort
	for _, item := range deploymentCreate.ContainerPort {
		containerPort = append(containerPort, corev1.ContainerPort{
			Name:          item.Name,
			ContainerPort: item.Port,
			Protocol:      corev1.Protocol(item.Protocol),
		})
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deploymentCreate.Name,
			Namespace:   deploymentCreate.Namespace,
			Labels:      deploymentCreate.Label,
			Annotations: map[string]string{"created-by": global.K8sManager},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &deploymentCreate.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels:      deploymentCreate.Label,
				MatchExpressions: nil,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deploymentCreate.Label,
					//Annotations:                nil,
				},
				Spec: corev1.PodSpec{
					Volumes:        nil,
					InitContainers: nil,
					Containers: []corev1.Container{
						{
							Name:  deploymentCreate.Name,
							Image: deploymentCreate.Image,
							Ports: containerPort,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(deploymentCreate.Cpu),
									corev1.ResourceMemory: resource.MustParse(deploymentCreate.Memory),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(deploymentCreate.Cpu),
									corev1.ResourceMemory: resource.MustParse(deploymentCreate.Memory),
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					ImagePullSecrets: []corev1.LocalObjectReference{
						{Name: deploymentCreate.ImagePullSecret},
					},
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    &deploymentCreate.RevisionHistoryLimit,
			ProgressDeadlineSeconds: nil, // 如果n秒后还没更新好则停止滚动更新，deployment不会主动再发出更新操作
		},
	}

	if deploymentCreate.ContainerPort != nil && deploymentCreate.HttpHealthCheck.HttpHealthPath != "" {
		if deploymentCreate.HttpHealthCheck.HttpHealthPort == "" {
			deploymentCreate.HttpHealthCheck.HttpHealthPort = deploymentCreate.ContainerPort[0].Name
		}
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   deploymentCreate.HttpHealthCheck.HttpHealthPath,
					Port:   intstr.Parse(deploymentCreate.HttpHealthCheck.HttpHealthPort),
					Scheme: "HTTP",
				},
			},
			InitialDelaySeconds:           10,
			TimeoutSeconds:                2,
			PeriodSeconds:                 10,
			SuccessThreshold:              1,
			FailureThreshold:              3,
			TerminationGracePeriodSeconds: nil, // 检查到异常时给程序优雅停止的时间，需要启用ProbeTerminationGracePeriod特性
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   deploymentCreate.HttpHealthCheck.HttpHealthPath,
					Port:   intstr.Parse(deploymentCreate.HttpHealthCheck.HttpHealthPort),
					Scheme: "HTTP",
				},
			},
			InitialDelaySeconds:           10,
			TimeoutSeconds:                2,
			PeriodSeconds:                 10,
			SuccessThreshold:              1,
			FailureThreshold:              3,
			TerminationGracePeriodSeconds: nil, // 检查到异常时给程序优雅停止的时间，需要启用ProbeTerminationGracePeriod特性
		}
		deployment.Spec.Template.Spec.Containers[0].StartupProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   deploymentCreate.HttpHealthCheck.HttpHealthPath,
					Port:   intstr.Parse(deploymentCreate.HttpHealthCheck.HttpHealthPort),
					Scheme: "HTTP",
				},
			},
			// 给予105s的启动时间, or (10+3) * 10 + 5 = 135 ?
			InitialDelaySeconds:           5,
			TimeoutSeconds:                3,
			PeriodSeconds:                 10,
			SuccessThreshold:              1,
			FailureThreshold:              10,
			TerminationGracePeriodSeconds: nil, // 检查到异常时给程序优雅停止的时间，需要启用ProbeTerminationGracePeriod特性
		}
	}

	_, err = global.K8s.ClientSet.AppsV1().Deployments(deploymentCreate.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{
		FieldManager: global.K8sManager,
	})
	if err != nil {
		return err
	}
	return
}

func (d *Deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int32) (err error) {

	autoScale, err := global.K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// 设置副本数
	autoScale.Spec.Replicas = scaleNum

	_, err = global.K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, autoScale, metav1.UpdateOptions{
		FieldManager: global.K8sManager,
	})
	if err != nil {
		return err
	}

	return
}

func (d *Deployment) DeleteDeploymentByName(deploymentName, namespace string, force bool) (err error) {
	opt := metav1.DeleteOptions{}
	if force {
		opt.GracePeriodSeconds = pointer.Int64(0)
	}

	err = global.K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, opt)
	if err != nil {
		return err
	}

	return
}

func (d *Deployment) SetDeploymentImage(deploymentName, namespace string, image dto.K8sSetImage) (err error) {
	opt := metav1.PatchOptions{
		FieldManager: global.K8sManager,
	}

	var pt types.PatchType
	var data []byte
	switch {
	case len(image) == 1 && image[0].Name == "":
		// 只提供了一个image,并且没有提供name,修改第0个容器
		pt = types.JSONPatchType
		data = []byte(fmt.Sprintf(`[{"op": "replace", "path": "/spec/template/spec/containers/0/image", "value":"%s"}]`, image[0].Image))

	case len(image) > 1 && !image.NameNotEmpty():
		return errors.New("更新多个容器, 容器名必填")

	default:
		pt = types.MergePatchType
		data = []byte(fmt.Sprintf(`{"spec": {"template": {"spec": {"containers": %s}}}}`, image.String()))
	}

	_, err = global.K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, pt, data, opt)

	if err != nil {
		return err
	}

	return
}

func (d *Deployment) RestartDeployment(deploymentName string, namespace string) (err error) {
	//deployment, err := d.GetDetail(deployName, namespace)
	//if err != nil {
	//	return err
	//}
	//Update对象
	//if deployment.Spec.Template.ObjectMeta.Annotations == nil {
	//	deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	//}
	//deployment.Spec.Template.ObjectMeta.Annotations["HandoverCloud.soulchild.cn/restartedAt"] = time.Now().Format("2006-01-02 15:04:05")

	// 使用Patch
	patchData := fmt.Sprintf(
		`{"spec":{"template":{"metadata":{"annotations":{"HandoverCloud.soulchild.cn/restartedAt":"%v"}}}}}`,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	patchByte := []byte(patchData)

	_, err = global.K8s.ClientSet.AppsV1().Deployments(namespace).Patch(
		context.TODO(),
		deploymentName,
		types.StrategicMergePatchType,
		patchByte,
		metav1.PatchOptions{},
	)

	if err != nil {
		return err
	}
	return nil
}

func (d *Deployment) UpdateK8sDeployment(content string) (err error) {
	deploy := &appsv1.Deployment{}
	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		return errors.New("反序列化失败,请检查yaml。" + err.Error())
	}

	//deploy.ObjectMeta.ManagedFields = []metav1.ManagedFieldsEntry{
	//	{
	//		Manager:    global.K8sManager,
	//		Operation:  metav1.ManagedFieldsOperationUpdate,
	//		APIVersion: deploy.APIVersion,
	//	},
	//}
	_, err = global.K8s.ClientSet.AppsV1().Deployments(deploy.Namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		return errors.New("更新Deployment失败," + err.Error())
	}
	return nil
}
