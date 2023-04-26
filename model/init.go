package model

import (
	"gorm.io/gorm"
)

func Init(db *gorm.DB) error {
	var MigrateModels = []any{
		//your model. eg: SystemUser{},
		&SystemUser{},
		&SystemLock{},
		&K8sCluster{},
	}
	err := db.AutoMigrate(MigrateModels...)

	return err
}
