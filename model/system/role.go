package system

import "soul/model/common"

type Role struct {
	common.ID
	RoleName string `json:"roleName" gorm:"uniqueIndex;comment:角色名称"`
	Users    []User `json:"users,omitempty" gorm:"many2many:system_user_roles;"`
}

func (r Role) TableName() string {
	return "t_system_role"
}

func (r Role) UsersToList() []string {
	var roleNames []string
	for _, u := range r.Users {
		roleNames = append(roleNames, u.Username)
	}
	return roleNames
}
