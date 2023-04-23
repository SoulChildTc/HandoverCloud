package discovery

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"soul/global"
)

type Discovery struct{}

// GetResourceByGV
// params: gv example "traefik.containo.us/v1alpha1"
func (d *Discovery) GetResourceByGV(clusterName, gv string) (resources []metav1.APIResource) {

	_, apiResourceList, _ := global.K8s.Use(clusterName).CacheDiscovery.ServerGroupsAndResources()

	for _, apiResourceList := range apiResourceList {
		// apiResourceList包含一个gv + n个apiResource
		if apiResourceList.GroupVersion == gv {
			for _, apiResource := range apiResourceList.APIResources {
				resources = append(resources, apiResource)
			}
		}
	}
	return
}
