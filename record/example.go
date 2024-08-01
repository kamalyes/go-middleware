/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-07-28 09:05:05
 * @FilePath: \go-middleware\record\example.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package record

import (
	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-middleware/record/access"
)

type ApiGroup struct {
	access.AccessRecordApi
}

var ApiGroupApp = new(ApiGroup)

// RouterRegister 注册访问记录中间件路由
func RouterRegister(rGroup *gin.RouterGroup) {
	// 初始化表
	access.AutoCreateTables()
	// 创建路由
	{
		accessRecordApi := ApiGroupApp.AccessRecordApi
		// 私有接口
		recordGroup := rGroup.Group("record")
		{
			recordGroup.GET("getAccessRecordPage", accessRecordApi.GetAccessRecordPage)
		}
	}
}
