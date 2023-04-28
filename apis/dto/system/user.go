package system

import (
	"soul/model"
	"soul/model/common"
)

type Register struct {
	Username string  `json:"username" binding:"required" msg:"用户名必填"`
	Nickname string  `json:"nickname" binding:"required" msg:"昵称必填"`
	Mobile   *string `json:"mobile" binding:"omitempty,mobile" mobile_err:"手机号格式有误"`
	Avatar   *string `json:"avatar"`
	Password string  `json:"password" binding:"required" msg:"密码必填"`
	Email    *string `json:"email" binding:"omitempty,email" email_err:"邮箱格式错误"`
}

type Login struct {
	Username string `json:"username" binding:"required" required_err:"账号不能为空"`
	Password string `json:"password" binding:"required" required_err:"密码不能为空"`
}

type Add struct {
	Register
	Roles []uint `json:"roles"`
}

type UserInfo struct {
	UserId   uint     `json:"userId"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Mobile   *string  `json:"mobile"`
	Avatar   *string  `json:"avatar,omitempty"`
	Email    *string  `json:"email,omitempty"`
	Roles    []string `json:"roles"`
	common.Timestamps
}

func (u *UserInfo) FromModel(user *model.SystemUser) *UserInfo {
	u.UserId = user.ID.ID
	u.Username = user.Username
	u.Nickname = user.Nickname
	u.Mobile = user.Mobile
	u.Avatar = user.Avatar
	u.Email = user.Email
	u.Roles = user.RolesToList()
	u.Timestamps = user.Timestamps

	return u
}
