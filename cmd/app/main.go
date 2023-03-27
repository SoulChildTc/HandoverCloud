package cmd

import (
	"soul/global"
	"soul/internal/config"
	"soul/internal/database"
	"soul/internal/k8s"
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

	// 初始化client-go
	global.K8s.Config, global.K8s.ClientSet = k8s.NewK8sClient()
}

func Execute() {
	server.StartServer()
}
