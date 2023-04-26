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

func (s *User) GetUserByAccount(account any) *model.SystemUser {
	var user model.SystemUser
	if err := global.DB.
		Preload("Roles").
		Where("username = ? or mobile = ? or email = ?", account, account, account).
		First(&user).
		Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Error(err.Error())
		return nil
	}

	return &user
}

func (s *User) GetUserByUserName(userName string) *model.SystemUser {
	var user model.SystemUser
	if err := global.DB.Where("user_name = ?", userName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Error(err.Error())
		return nil
	}

	return &user
}

func (s *User) GetUserByMobile(mobile any) *model.SystemUser {
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
	if err := global.DB.Debug().Model(&user).Preload("Roles").First(&user).Error; err != nil {
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

func (s *User) AssignRoleById(roleId, userId uint) error {
	role := model.SystemRole{}
	result := global.DB.First(&role, roleId)
	if result.Error != nil {
		return result.Error
	}

	user := model.SystemUser{}
	result = global.DB.First(&user, userId)
	if result.Error != nil {
		return result.Error
	}

	err := global.DB.Debug().Model(&user).Association("Roles").Append(&role)
	if err != nil {
		return err
	}
	return nil
}
