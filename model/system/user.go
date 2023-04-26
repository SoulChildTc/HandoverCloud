package system

import (
	"soul/model/common"
)

type User struct {
	common.ID
	Username string  `json:"username" gorm:"size:32;not null;uniqueIndex;comment:用户名"`
	Nickname string  `json:"nickname" gorm:"size:32;not null;comment:昵称"`
	Mobile   *string `json:"mobile" gorm:"size:24;uniqueIndex;comment:用户手机号"`
	Avatar   *string `json:"avatar" gorm:"size:256;common:用户头像"`
	Password string  `json:"-" gorm:"not null;comment:用户密码"`
	Email    *string `json:"email"  gorm:"uniqueIndex;comment:用户邮箱"`
	Roles    []Role  `json:"role" gorm:"many2many:system_user_roles;comment:角色列表"`
	common.Timestamps
}

func (u *User) TableName() string {
	return "t_system_user"
}

func (u *User) RolesToList() []string {
	var roleNames []string
	for _, r := range u.Roles {
		roleNames = append(roleNames, r.RoleName)
	}
	return roleNames
}
