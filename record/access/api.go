/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 11:15:37
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
	"github.com/kamalyes/go-core/response"
	"github.com/kamalyes/go-middleware/constants"
	"go.uber.org/zap"
)

type AccessRecordApi struct{}

// GetAccessRecordPage 分页查询操作记录
func (s *AccessRecordApi) GetAccessRecordPage(ctx *gin.Context) {
	pageInfo := database.PageParam(ctx)
	if pageInfo == nil {
		response.GenGinResponse(ctx, &response.ResponseOption{Code: response.ServerError, HttpCode: response.StatusInternalServerError, Message: constants.ErrParseRequestData})
		return
	}
	err, pageBean := AccessRecordServiceApp.GetAccessRecordPage(pageInfo)
	if err != nil {
		global.LOG.Error(constants.ErrGainRecordResponse, zap.Any("err", err))
		response.GenGinResponse(ctx, &response.ResponseOption{Code: response.ServerError, HttpCode: response.StatusInternalServerError, Message: constants.ErrGainRecordResponse})
	} else {
		response.GenGinResponse(ctx, &response.ResponseOption{Data: pageBean})

	}
}
