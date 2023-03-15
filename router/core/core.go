package core

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"soul/apis/controller/core"
	_ "soul/docs"
	"soul/middleware"
)

func RegisterRoute(r *gin.RouterGroup) {
	r.GET("/ping", core.Ping)
	swag := r.Group("/swagger").Use(middleware.BasicAuth)
	swag.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
