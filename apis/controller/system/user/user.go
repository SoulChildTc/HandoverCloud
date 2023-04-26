package user

import (
	"github.com/gin-gonic/gin"
	"soul/apis/dto"
	"soul/apis/service"
	"soul/utils/httputil"
	"strconv"
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
	res, err := service.SystemUser.Login(u)
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}
	httputil.OK(c, res, "登录成功")
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
//func Register(c *gin.Context) {
//	var u dto.SystemRegister
//	err := c.ShouldBindJSON(&u)
//	if err != nil {
//		httputil.Error(c, httputil.ParseValidateError(err, &u).Error())
//		return
//	}
//	msg, ok := service.SystemUser.Register(u)
//	if ok {
//		httputil.OK(c, nil, msg)
//		return
//	}
//	httputil.Error(c, msg)
//}

// Info
//
//	@description	用户信息
//	@tags			User
//	@summary		用户信息
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回用户信息"
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

// AddUser
//
//	@description	添加用户
//	@tags			User
//	@summary		添加用户
//	@accept			json
//	@produce		json
//	@param			data			body	dto.SystemAdd		true	"用户信息"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				object	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/system/user/ [post]
func AddUser(c *gin.Context) {
	var u dto.SystemAdd
	err := c.ShouldBindJSON(&u)
	if err != nil {
		httputil.Error(c, httputil.ParseValidateError(err, &u).Error())
		return
	}
	msg, ok := service.SystemUser.Add(u)
	if ok {
		httputil.OK(c, nil, msg)
		return
	}
	httputil.Error(c, msg)
}

// AssignRole
//
//	@description	为用户分配一个角色
//	@tags			User
//	@summary		为用户分配一个角色
//	@accept			json
//	@produce		json
//	@param			data			body	object		true		"用户和角色Id: {"userId": 1, roleId: 1}"
//	@Param			Authorization	header	string					true	"Authorization token"
//	@success		200				{object}	httputil.ResponseBody	"成功返回"
//	@router			/api/v1/system/user/{userId}/roles [post]
func AssignRole(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		httputil.Error(c, "用户ID错误")
		return
	}
	params := new(struct {
		RoleId uint `json:"roleId"`
	})

	err = c.ShouldBindJSON(&params)
	if params.RoleId <= 0 {
		httputil.Error(c, "请输入正确的角色ID")
		return
	}
	err = service.SystemUser.AssignRole(params.RoleId, uint(userId))
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}
	httputil.OK(c, nil, "操作成功")
}
