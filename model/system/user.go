package system

import (
	"soul/model/common"
)

type User struct {
	common.ID
	Username string  `json:"username" gorm:"size:32;not null;comment:用户名"`
	Mobile   string  `json:"mobile" gorm:"size:24;not null;uniqueIndex;comment:用户手机号"`
	Password string  `json:"-" gorm:"not null;comment:用户密码"`
	Email    *string `json:"email"  gorm:"comment:用户邮箱"`
	common.Timestamps
}
