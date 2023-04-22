package secret

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	"soul/apis/dto"
	"soul/apis/service"
	"soul/utils/httputil"
)

// GetSecretByName
//
//	@description	获取Secret信息
//	@tags			K8s,Secret
//	@summary		获取Secret信息
//	@produce		json
//	@param			secretName		path	string					true	"Secret名称"
//	@param			namespace		path	string					true	"Namespace"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回Secret信息"
//	@router			/api/v1/k8s/secret/{namespace}/{secretName} [get]
func GetSecretByName(c *gin.Context) {
	name := c.Param("secretName")
	namespace := c.Param("namespace")
	if name == "" || namespace == "" {
		httputil.Error(c, "secret名和名称空间不能为空")
		return
	}

	secret, err := service.K8sSecret.GetSecretByName(name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, secret, "获取成功")
}

// GetSecretList
//
//	@description	获取Secret列表
//	@tags			K8s,Secret
//	@summary		获取Secret列表
//	@produce		json
//	@param			namespace		path	string						false	"Namespace 不填为全部"
//	@Param			Authorization	header	string						true	"Authorization token"
//	@Param			filter			query	string						false	"根据Secret名字模糊查询"
//	@Param			limit			query	string						false	"一页获取多少条数据,默认十条"
//	@Param			page			query	string						false	"获取第几页的数据,默认第一页"
//	@success		200				object	httputil.PageResponseBody	"成功返回Secret列表"
//	@router			/api/v1/k8s/secret/{namespace} [get]
func GetSecretList(c *gin.Context) {
	namespace := c.Param("namespace")
	params := new(struct {
		FilterName string `form:"filter"`
		Limit      int    `form:"limit,default=10"`
		Page       int    `form:"page,default=1"`
	})

	if err := c.ShouldBind(params); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, params).Error())
		return
	}

	secrets, err := service.K8sSecret.GetSecretList(params.FilterName, namespace, params.Limit, params.Page)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.Page(c, secrets, "获取成功")
}

// DeleteSecretByName
//
//	@description	删除Secret
//	@tags			K8s,Secret
//	@summary		删除Secret
//	@produce		json
//	@param			secretName		path	string	true	"Secret名称"
//	@param			namespace		path	string	true	"Namespace"
//	@Param			Authorization	header	string	true	"Authorization token"
//	@success		200				object	nil		"成功返回"
//	@router			/api/v1/k8s/secret/{namespace}/{secretName} [delete]
func DeleteSecretByName(c *gin.Context) {
	name := c.Param("secretName")
	namespace := c.Param("namespace")
	if name == "" || namespace == "" {
		httputil.Error(c, "secret名和名称空间不能为空")
		return
	}

	_, err := service.K8sSecret.GetSecretByName(name, namespace)
	if err != nil {
		switch {
		case errors.IsNotFound(err):
			httputil.Error(c, fmt.Sprintf(`Secret "%s" 在 "%s" 中未找到`, name, namespace))
		default:
			httputil.Error(c, err.Error())
		}
		return
	}

	err = service.K8sSecret.DeleteSecretByName(name, namespace)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "删除成功")
}

// CreateSecret
//
//	@description	创建 Opaque 类型的 Secret
//	@tags			K8s,Secret
//	@summary		创建 Opaque 类型的 Secret
//	@produce		json
//	@produce		json
//	@Param			Authorization	header	string					true	"Authorization token"
//	@param			data			body	dto.K8sSecretCreate		true	"K8sSecretCreate 对象"
//	@success		200				object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/k8s/secret/ [post]
func CreateSecret(c *gin.Context) {
	// 初始化默认值
	secret := dto.K8sSecretCreate{}

	if err := c.ShouldBindJSON(&secret); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &secret).Error())
		return
	}

	err := service.K8sSecret.CreateSecret(&secret)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "创建成功")
}

// UpdateSecret
//
//	@description	更新 Opaque 类型的 Secret
//	@tags			K8s,Secret
//	@summary		更新 Opaque 类型的 Secret
//	@produce		json
//	@produce		json
//	@Param			Authorization	header	string					true	"Authorization token"
//	@param			data			body	dto.K8sSecretCreate		true	"K8sSecretCreate 对象"
//	@success		200				object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/k8s/secret/ [put]
func UpdateSecret(c *gin.Context) {
	// 初始化默认值
	secret := dto.K8sSecretCreate{}

	if err := c.ShouldBindJSON(&secret); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &secret).Error())
		return
	}

	err := service.K8sSecret.UpdateSecret(&secret)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "更新成功")
}

// CreateSecretForDockerRegistry
//
//	@description	创建 docker-registry 类型的 Secret
//	@tags			K8s,Secret
//	@summary		创建 docker-registry 类型的 Secret
//	@produce		json
//	@produce		json
//	@Param			Authorization	header	string									true	"Authorization token"
//	@param			data			body	dto.K8sSecretForDockerRegistryCreate	true	"K8sSecretForDockerRegistryCreate 对象"
//	@success		200				object	httputil.ResponseBody					"成功返回"
//	@router			/api/v1/k8s/secret/_docker-registry [post]
func CreateSecretForDockerRegistry(c *gin.Context) {
	// 初始化默认值
	secret := dto.K8sSecretForDockerRegistryCreate{}

	if err := c.ShouldBindJSON(&secret); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &secret).Error())
		return
	}

	err := service.K8sSecret.CreateSecretForDockerRegistry(&secret)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "创建成功")
}

// UpdateSecretForDockerRegistry
//
//	@description	更新 docker-registry 类型的 Secret
//	@tags			K8s,Secret
//	@summary		更新 docker-registry 类型的 Secret
//	@produce		json
//	@produce		json
//	@Param			Authorization	header	string									true	"Authorization token"
//	@param			data			body	dto.K8sSecretForDockerRegistryCreate	true	"K8sSecretForDockerRegistryCreate 对象"
//	@success		200				object	httputil.ResponseBody					"成功返回"
//	@router			/api/v1/k8s/secret/_docker-registry [put]
func UpdateSecretForDockerRegistry(c *gin.Context) {
	// 初始化默认值
	secret := dto.K8sSecretForDockerRegistryCreate{}

	if err := c.ShouldBindJSON(&secret); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &secret).Error())
		return
	}

	err := service.K8sSecret.UpdateSecretForDockerRegistry(&secret)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "更新成功")
}

// CreateSecretForTls
//
//	@description	创建 tls 类型的 Secret
//	@tags			K8s,Secret
//	@summary		创建 tls 类型的 Secret
//	@produce		json
//	@produce		json
//	@Param			Authorization	header	string						true	"Authorization token"
//	@param			data			body	dto.K8sSecretForTlsCreate	true	"K8sSecretForTlsCreate 对象"
//	@success		200				object	httputil.ResponseBody		"成功返回"
//	@router			/api/v1/k8s/secret/_tls [post]
func CreateSecretForTls(c *gin.Context) {
	// 初始化默认值
	secret := dto.K8sSecretForTlsCreate{}

	if err := c.ShouldBindJSON(&secret); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &secret).Error())
		return
	}

	err := service.K8sSecret.CreateSecretForTls(&secret)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "创建成功")
}

// UpdateSecretForTls
//
//	@description	更新 tls 类型的 Secret
//	@tags			K8s,Secret
//	@summary		更新 tls 类型的 Secret
//	@produce		json
//	@produce		json
//	@Param			Authorization	header	string						true	"Authorization token"
//	@param			data			body	dto.K8sSecretForTlsCreate	true	"K8sSecretForTlsCreate 对象"
//	@success		200				object	httputil.ResponseBody		"成功返回"
//	@router			/api/v1/k8s/secret/_tls [put]
func UpdateSecretForTls(c *gin.Context) {
	// 初始化默认值
	secret := dto.K8sSecretForTlsCreate{}

	if err := c.ShouldBindJSON(&secret); err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &secret).Error())
		return
	}

	err := service.K8sSecret.UpdateSecretForTls(&secret)

	if err != nil {
		httputil.Error(c, err.Error())
		return
	}

	httputil.OK(c, nil, "更新成功")
}
