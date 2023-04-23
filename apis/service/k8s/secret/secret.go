package secret

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"soul/apis/dto"
	"soul/apis/service/k8s"
	"soul/global"
	"soul/utils/httputil"
)

type Secret struct{}

func (s *Secret) toCells(secrets []corev1.Secret) []k8s.DataCell {
	cells := make([]k8s.DataCell, len(secrets))
	for i, item := range secrets {
		cells[i] = k8s.DataCell(secretCell(item))
	}
	return cells
}

func (s *Secret) fromCells(cells []k8s.DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i, item := range cells {
		secrets[i] = corev1.Secret(item.(secretCell))
	}
	return secrets
}

func (s *Secret) GetSecretByName(clusterName, name, namespace string) (*corev1.Secret, error) {
	secret, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (s *Secret) GetSecretList(clusterName, filterName, namespace string, limit, page int) (*httputil.PageResp, error) {
	secrets, err := global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	selectableData := k8s.DataSelect{
		GenericDataList: s.toCells(secrets.Items),
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

func (s *Secret) DeleteSecretByName(clusterName, secretName, namespace string) (err error) {
	err = global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Secret) CreateSecretForDockerRegistry(clusterName string, secretForDockerRegistryCreate *dto.K8sSecretForDockerRegistryCreate) (err error) {
	// 格式转换
	data := secretForDockerRegistryCreate.ToDockerconfig()
	secretStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretForDockerRegistryCreate.Name,
			Namespace:   secretForDockerRegistryCreate.Namespace,
			Labels:      nil,
			Annotations: map[string]string{"created-by": global.K8sManager},
		},
		StringData: map[string]string{".dockerconfigjson": string(secretStr)},
		Type:       corev1.SecretTypeDockerConfigJson,
	}

	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(secretForDockerRegistryCreate.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
func (s *Secret) UpdateSecretForDockerRegistry(clusterName string, secretForDockerRegistryCreate *dto.K8sSecretForDockerRegistryCreate) (err error) {
	// 格式转换
	data := secretForDockerRegistryCreate.ToDockerconfig()
	secretStr, err := json.Marshal(data)

	if err != nil {
		return err
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretForDockerRegistryCreate.Name,
			Namespace:   secretForDockerRegistryCreate.Namespace,
			Labels:      nil,
			Annotations: map[string]string{"created-by": global.K8sManager},
		},
		StringData: map[string]string{".dockerconfigjson": string(secretStr)},
		Type:       corev1.SecretTypeDockerConfigJson,
	}

	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(secretForDockerRegistryCreate.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Secret) CreateSecretForTls(clusterName string, secretForTlsCreate *dto.K8sSecretForTlsCreate) (err error) {
	if err != nil {
		return err
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretForTlsCreate.Name,
			Namespace:   secretForTlsCreate.Namespace,
			Labels:      nil,
			Annotations: map[string]string{"created-by": global.K8sManager},
		},
		StringData: map[string]string{
			"tls.crt": secretForTlsCreate.Certificate,
			"tls.key": secretForTlsCreate.PrivateKey,
		},
		Type: corev1.SecretTypeTLS,
	}

	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(secretForTlsCreate.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
func (s *Secret) UpdateSecretForTls(clusterName string, secretForTlsCreate *dto.K8sSecretForTlsCreate) (err error) {
	if err != nil {
		return err
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretForTlsCreate.Name,
			Namespace:   secretForTlsCreate.Namespace,
			Labels:      nil,
			Annotations: map[string]string{"created-by": global.K8sManager},
		},
		StringData: map[string]string{
			"tls.crt": secretForTlsCreate.Certificate,
			"tls.key": secretForTlsCreate.PrivateKey,
		},
		Type: corev1.SecretTypeTLS,
	}

	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(secretForTlsCreate.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Secret) CreateSecret(clusterName string, secretCreate *dto.K8sSecretCreate) (err error) {
	if err != nil {
		return err
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretCreate.Name,
			Namespace:   secretCreate.Namespace,
			Labels:      nil,
			Annotations: map[string]string{"created-by": global.K8sManager},
		},
		StringData: secretCreate.Data,
		Type:       corev1.SecretTypeOpaque,
	}

	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(secretCreate.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
func (s *Secret) UpdateSecret(clusterName string, secretCreate *dto.K8sSecretCreate) (err error) {
	if err != nil {
		return err
	}
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        secretCreate.Name,
			Namespace:   secretCreate.Namespace,
			Labels:      nil,
			Annotations: map[string]string{"created-by": global.K8sManager},
		},
		StringData: secretCreate.Data,
		Type:       corev1.SecretTypeOpaque,
	}

	_, err = global.K8s.Use(clusterName).ClientSet.CoreV1().Secrets(secretCreate.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}
