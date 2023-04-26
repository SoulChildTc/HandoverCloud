package system

import "soul/model/common"

type Lock struct {
	common.ID
}

func (l *Lock) TableName() string {
	return "t_system_lock"
}
