package system

import (
	"errors"
	"gorm.io/gorm"
	"soul/global"
	"soul/model"
)

type InitData struct{}

func (d *InitData) Locked() error {
	lock := &model.SystemLock{}

	if err := global.DB.Model(&model.SystemLock{}).Create(lock); err.Error != nil {
		return err.Error
	}

	return nil
}

func (d *InitData) HasLock() (error, bool) {
	lock := &model.SystemLock{}
	result := global.DB.Model(&model.SystemLock{}).First(lock)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error, false
	}

	if result.RowsAffected == 0 {
		return nil, false
	}

	return nil, true
}
