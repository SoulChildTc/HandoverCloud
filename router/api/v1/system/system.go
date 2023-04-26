package system

import (
	"github.com/gin-gonic/gin"
	"soul/apis/controller/system/dbInitializer"
	"soul/apis/controller/system/user"
	"soul/middleware"
)

// system模块路由

func RegisterRoute(r *gin.RouterGroup) {
	// 用户
	userGroup := r.Group("/user")
	{
		userGroup.POST("/login", user.Login)
		//userGroup.POST("/register", user.Register) // 关闭注册
	}

	userAuthGroup := r.Group("/user").Use(middleware.JwtAuth)
	{
		userAuthGroup.POST("/", user.AddUser)
		userAuthGroup.GET("/info", user.Info)
		userAuthGroup.POST("/:userId/roles", user.AssignRole)
	}

	// 数据初始化
	dbinit := r.Group("/dbInitializer")
	{
		dbinit.GET("/", dbInitializer.IsInit)
		dbinit.POST("/", dbInitializer.InitData)
	}

}
