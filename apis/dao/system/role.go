package system

import (
	"errors"
	"gorm.io/gorm"
	"soul/global"
	log "soul/internal/logger"
	"soul/model"
	"soul/model/common"
)

type Role struct{}

func (r *Role) CreateRole(role *model.SystemRole) error {
	result := global.DB.Create(role)
	return result.Error
}

func (r *Role) GetRoleById(id uint) *model.SystemRole {
	role := model.SystemRole{
		ID: common.ID{ID: id},
	}
	if err := global.DB.Debug().Model(&role).Preload("Users").First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Error(err.Error())
		return nil
	}

	return &role
}
