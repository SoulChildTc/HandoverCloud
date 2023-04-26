package dbInitializer

import (
	"github.com/gin-gonic/gin"
	"soul/apis/service"
	"soul/utils/httputil"
)

// IsInit
//
//	@description	是否已经初始化数据库
//	@tags			DbInitializer
//	@summary		是否已经初始化数据库
//	@produce		json
//	@success		200	object	string	"成功返回"
//	@router			/api/v1/system/dbInitializer/ [get]
func IsInit(c *gin.Context) {
	err, res := service.SystemInitData.IsInit()
	if err != nil {
		httputil.ErrorWithCode(c, 500, "内部错误 "+err.Error())
		return
	}
	httputil.OK(c, gin.H{"isInit": res}, "查询成功")
}

// InitData
//
//	@description	初始化数据库数据
//	@tags			DbInitializer
//	@summary		初始化数据库数据
//	@produce		json
//	@success		200	object	string	"成功返回"
//	@router			/api/v1/system/dbInitializer/ [post]
func InitData(c *gin.Context) {
	err := service.SystemInitData.InitData()
	if err != nil {
		httputil.Error(c, err.Error())
		return
	}
	httputil.OK(c, nil, "操作成功")
}
