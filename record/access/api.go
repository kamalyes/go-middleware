/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2023-07-28 09:05:05
 * @FilePath: \go-middleware\record\access\api.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package access

import (
	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/global"
	"github.com/kamalyes/go-core/page"
	"github.com/kamalyes/go-core/result"
	"go.uber.org/zap"
)

type AccessRecordApi struct{}

// GetAccessRecordPage 分页查询操作记录
func (s *AccessRecordApi) GetAccessRecordPage(c *gin.Context) {
	pageInfo := page.PageParam(c)
	if pageInfo == nil {
		result.FailMsg("获取失败,解析请求参数异常", c)
		return
	}
	err, pageBean := AccessRecordServiceApp.GetAccessRecordPage(pageInfo)
	if err != nil {
		global.LOG.Error("获取接口访问记录失败:", zap.Any("err", err))
		result.FailMsg("获取接口访问记录失败", c)
	} else {
		result.OkDataMsg(pageBean, "获取成功", c)
	}
}
