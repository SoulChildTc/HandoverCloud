package system

import (
	"soul/model"
)

type RoleInfo struct {
	RoleId   uint     `json:"roleId"`
	RoleName string   `json:"roleName"`
	User     []string `json:"users"`
}

func (r *RoleInfo) FromModel(role *model.SystemRole) *RoleInfo {
	r.RoleId = role.ID.ID
	r.RoleName = role.RoleName
	r.User = role.UsersToList()

	return r
}
