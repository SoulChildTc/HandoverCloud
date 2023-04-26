package role

import (
	"soul/apis/dao"
	"soul/apis/dto"
)

type Role struct {
}

func (r *Role) Info(roleId uint) (*dto.SystemRoleInfo, bool) {
	role := dao.SystemRole.GetRoleById(roleId)

	if role == nil {
		return nil, false
	}

	roleInfo := &dto.SystemRoleInfo{}
	return roleInfo.FromModel(role), true
}
