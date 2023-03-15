package httputil

import (
	"errors"
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
		errStr, err := utils.GetTagValue(obj, e.Field(), e.Tag()+"_err")
		if err != nil {
			// 没有获取到就获取默认的msg
			errStr, _ = utils.GetTagValue(obj, e.Field(), "msg")
		}
		errResult = append(errResult, errStr)

	}
	return errors.New(strings.Join(errResult, ","))
}

func RegisterAllValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mobile", mobileValidate) // 指定tag名称和处理函数
	}
}

func mobileValidate(fl validator.FieldLevel) bool {
	mobile, ok := fl.Field().Interface().(string)
	if ok {
		if len(mobile) != 11 {
			return false
		}
		return false
	}

	reg := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return reg.MatchString(mobile)
}
