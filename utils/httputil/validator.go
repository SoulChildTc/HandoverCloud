package httputil

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"soul/utils"
	"strings"
)

func ParseValidateError(err error, obj any) error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.New("数据格式有误")
	}

	var errResult []string
	for _, e := range errs {
		errStr, err := utils.GetTagValueByNamespace(obj, e.Namespace(), e.Tag()+"_err")
		if err != nil {
			// 没有获取到就获取默认的msg
			errStr, _ = utils.GetTagValueByNamespace(obj, e.Namespace(), "msg")
		}
		errResult = append(errResult, errStr)

	}
	return errors.New(strings.Join(errResult, ","))
}

// RegisterAllValidator 注册自定义校验器
func RegisterAllValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", mobileValidate) // 指定tag名称和处理函数
		v.RegisterValidation("pem", pemValidate)       // 指定tag名称和处理函数
	}
}

func mobileValidate(fl validator.FieldLevel) bool {
	mobile, ok := fl.Field().Interface().(string)
	if ok {
		if len(mobile) != 11 {
			return false
		}
	}

	reg := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return reg.MatchString(mobile)
}

func pemValidate(fl validator.FieldLevel) bool {
	data, _ := fl.Field().Interface().(string)
	block, _ := pem.Decode([]byte(data))
	if block == nil {
		return false
	}

	pemType := fl.Param()
	if pemType == "cert" {
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return false
		}
	}

	if pemType == "key" {
		_, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			_, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return false
			}
		}
	}

	return true
}

func CheckParams(c *gin.Context, params ...string) error {
	for _, p := range params {
		if c.Param(p) == "" {
			return errors.New(fmt.Sprintf("%s不能为空", p))
		}
	}
	return nil
}
