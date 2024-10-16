/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 10:53:43
 * @FilePath: \go-middleware\record\access\service.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package access

import (
	"github.com/golang-module/carbon/v2"
	"github.com/kamalyes/go-core/database"
	"github.com/kamalyes/go-core/global"
	"github.com/kamalyes/go-middleware/constants"
	"go.uber.org/zap"
)

type AccessRecordService struct{}

var AccessRecordServiceApp = new(AccessRecordService)

// CreateAccessRecord 创建记录
func (opt *AccessRecordService) CreateAccessRecord(record AccessRecordModel, retainDays int) (err error) {
	err = global.DB.Create(&record).Error
	go func() {
		// 默认保留7天
		if retainDays < 7 {
			retainDays = 7
		}
		time := carbon.Now().SubDays(retainDays).ToDateTimeString()
		err = global.DB.Where("create_time < ?", time).Delete(&AccessRecordModel{}).Error
		if err != nil {
			global.LOG.Error(constants.ErrDeleteRecordData, zap.Any("err", err))
		}
	}()
	return err
}

// GetAccessRecordPage 分页获取操作记录列表
func (opt *AccessRecordService) GetAccessRecordPage(pageInfo *database.PageInfo) (err error, pageBean *database.PageBean) {
	pageBean = &database.PageBean{Page: pageInfo.Current, PageSize: pageInfo.RowCount}
	rows := make([]*AccessRecordModel, 0)
	err, pageBean = database.FindPage(&AccessRecordModel{}, &rows, pageInfo)
	return
}
