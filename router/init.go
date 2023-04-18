package router

import (
	"github.com/gin-gonic/gin"
	"soul/middleware"
	"soul/router/api/v1/k8s"
	"soul/router/api/v1/system"
	"soul/router/core"
)

func registerRoute(r gin.IRouter, groupName string, register func(r *gin.RouterGroup), middlewares ...gin.HandlerFunc) {
	group := r.Group(groupName)
	group.Use(middlewares...)
	register(group)
}

func InitRouter(r *gin.Engine) {
	/*
		功能模块路由注册
	*/

	// core api
	{
		registerRoute(r, "", core.RegisterRoute)
	}

	// /api/v1
	apiV1 := r.Group("/api/v1")
	{
		registerRoute(apiV1, "/system", system.RegisterRoute)
	}

	// /api/v1 - Auth
	apiV1Auth := r.Group("/api/v1")
	apiV1Auth.Use(middleware.JwtAuth)
	{
		registerRoute(apiV1Auth, "/k8s", k8s.RegisterRoute)
	}

	// /api/v2
	{

	}
}
