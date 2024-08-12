/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 17:56:40
 * @FilePath: \go-middleware\record\access\api.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package access

import (
	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-config/global"
	"github.com/kamalyes/go-core/database"
	"github.com/kamalyes/go-core/pkg/response"
	"github.com/kamalyes/go-middleware/internal"
	"go.uber.org/zap"
)

type AccessRecordApi struct{}

// GetAccessRecordPage 分页查询操作记录
func (s *AccessRecordApi) GetAccessRecordPage(ctx *gin.Context) {
	pageInfo := database.PageParam(ctx)
	if pageInfo == nil {
		response.GenResponse(ctx, &response.ResponseOption{Code: response.ServerError, HttpCode: response.FAIL, Message: internal.ErrParseRequestData})
		return
	}
	err, pageBean := AccessRecordServiceApp.GetAccessRecordPage(pageInfo)
	if err != nil {
		global.LOG.Error(internal.ErrGainRecordResponse, zap.Any("err", err))
		response.GenResponse(ctx, &response.ResponseOption{Code: response.ServerError, HttpCode: response.FAIL, Message: internal.ErrGainRecordResponse})
	} else {
		response.GenResponse(ctx, &response.ResponseOption{Data: pageBean})

	}
}
