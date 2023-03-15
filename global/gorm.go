package global

import (
	"database/sql"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	SqlDB *sql.DB
)
