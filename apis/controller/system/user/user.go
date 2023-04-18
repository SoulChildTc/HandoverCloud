package user

import (
	"github.com/gin-gonic/gin"
	"soul/apis/dto"
	"soul/apis/service"
	"soul/utils/httputil"
)

// Login
//
//	@description	用户登录
//	@tags			User
//	@summary		用户登录
//	@accept			json
//	@produce		json
//	@param			data	body	dto.SystemLogin			true	"手机号,密码"
//	@success		200		object	httputil.ResponseBody	"成功返回token"
//	@router			/api/v1/system/user/login [post]
func Login(c *gin.Context) {
	var u dto.SystemLogin
	err := c.ShouldBindJSON(&u)
	if err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &u).Error())
		return
	}
	res, ok := service.SystemUser.Login(u)
	if ok {
		httputil.OK(c, gin.H{"token": res}, "登录成功")
	} else {
		httputil.Error(c, res)
	}

}

// Register
//
//	@description	用户注册
//	@tags			User
//	@summary		用户注册
//	@accept			json
//	@produce		json
//	@param			data	body	dto.SystemRegister		true	"用户信息"
//	@success		200		object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/system/user/register [post]
func Register(c *gin.Context) {
	var u dto.SystemRegister
	err := c.ShouldBindJSON(&u)
	if err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &u).Error())
		return
	}
	msg, ok := service.SystemUser.Register(u)
	if ok {
		httputil.OK(c, nil, msg)
		return
	}
	httputil.Error(c, msg)
}

// Info
//
//	@description	用户信息
//	@tags			User
//	@summary		用户信息
//	@Param			X-Token	header	string					true	"Authorization token"
//	@success		200		object	httputil.ResponseBody	"成功返回用户信息"
//	@router			/api/v1/system/user/info [get]
func Info(c *gin.Context) {
	userId := c.GetUint("userId")
	res, ok := service.SystemUser.Info(userId)
	if ok {
		httputil.OK(c, res, "操作成功")
	} else {
		httputil.Error(c, "获取用户信息失败")
	}

}
