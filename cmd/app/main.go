package cmd

import (
	"soul/global"
	"soul/internal/config"
	"soul/internal/database"
	"soul/internal/logger"
	"soul/internal/server"
)

func init() {
	//初始化配置文件
	global.V = config.LoadConfig()

	// 初始化logrus
	logger.InitLogger()

	// 初始化gorm
	global.DB, global.SqlDB = database.InitDB()

	// 数据迁移
	if global.V.GetBool("migrate") {
		database.InitDBMigrate()
	}

}

func Execute() {
	server.StartServer()
}
