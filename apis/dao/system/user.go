package system

import (
	"errors"
	"gorm.io/gorm"
	"soul/global"
	log "soul/internal/logger"
	"soul/model"
	"soul/model/common"
)

type User struct{}

func (s *User) GetUserByMobile(mobile string) *model.SystemUser {
	var user model.SystemUser
	if err := global.DB.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Error(err.Error())
		return nil
	}

	return &user
}

func (s *User) GetUserById(id uint) *model.SystemUser {
	user := model.SystemUser{
		ID: common.ID{ID: id},
	}
	if err := global.DB.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Error(err.Error())
		return nil
	}

	return &user
}

func (s *User) CreateUser(user *model.SystemUser) error {
	result := global.DB.Create(user)
	return result.Error
}
