package database

import (
	"database/sql"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"soul/global"
	log "soul/internal/logger"
	"time"
)

func InitDB() (*gorm.DB, *sql.DB) {
	dbDriver := global.Config.Database.Driver

	switch dbDriver {
	case "mysql":
		fmt.Printf("[Init] 使用%s数据库驱动\n", dbDriver)
		return initMySqlGorm()
	case "sqlite":
		fmt.Printf("[Init] 使用%s数据库驱动\n", dbDriver)
		return initSqliteGorm()
	default:
		fmt.Println("[Init] 使用默认数据库驱动 - Sqlite")
		return initSqliteGorm()

	}
}

func initMySqlGorm() (*gorm.DB, *sql.DB) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		global.Config.Database.Username,
		global.Config.Database.Password,
		global.Config.Database.Host,
		global.Config.Database.Port,
		global.Config.Database.Database,
		global.Config.Database.Charset,
	)

	mysqlConfig := mysql.Config{DSN: dsn}

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: log.NewGormLogger(),
	}

	gormDB, err := gorm.Open(mysql.New(mysqlConfig), gormConfig)
	if err != nil {
		os.Exit(1)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic("数据库连接池初始化失败")
	}
	sqlDB.SetMaxIdleConns(global.Config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(global.Config.Database.MaxOpenConns)

	connMaxIdleTime := time.Duration(global.Config.Database.ConnMaxIdleTime) * time.Minute
	connMaxLifetime := time.Duration(global.Config.Database.ConnMaxLifetime) * time.Minute
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	fmt.Println("[Init] 数据库连接初始化成功")
	return gormDB, sqlDB
}

func initSqliteGorm() (*gorm.DB, *sql.DB) {
	/*
		github.com/mattn/go-sqlite3底层驱动是C实现的，所以需要开启CGO，兼容性可能会不太好。
		所以选择了性能稍弱的 github.com/glebarez/sqlite
	*/
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: log.NewGormLogger(),
	}

	gormDB, err := gorm.Open(sqlite.Open(global.Config.Database.Path), gormConfig)
	if err != nil {
		os.Exit(1)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		panic("数据库连接池初始化失败")
	}
	sqlDB.SetMaxIdleConns(global.Config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(global.Config.Database.MaxOpenConns)

	connMaxIdleTime := time.Duration(global.Config.Database.MaxOpenConns) * time.Minute
	connMaxLifetime := time.Duration(global.Config.Database.MaxOpenConns) * time.Minute
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	fmt.Println("[Init] 数据库连接初始化成功")
	return gormDB, sqlDB
}
