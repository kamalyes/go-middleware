/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 11:08:43
 * @FilePath: \go-middleware\record\access\model.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package access

import (
	"github.com/kamalyes/go-core/global"
	"go.uber.org/zap"
)

// 自动建表
func AutoCreateTables() {
	if global.DB != nil {
		err := global.DB.AutoMigrate(
			AccessRecordModel{},
		)
		if err != nil {
			global.LOG.Error("创建访问记录表失败", zap.Any("err", err))
		}
	}
}

// TableName 自定义表名
func (AccessRecordModel) TableName() string {
	return "access_record"
}

type AccessRecordModel struct {
	global.Model
	Ip             string `json:"ip"                     form:"ip"                   gorm:"column:ip;comment:ip;type:varchar(40);"`
	Method         string `json:"method"                 form:"method"               gorm:"column:method;comment:请求方法;type:varchar(10);"`
	Path           string `json:"path"                   form:"path"                 gorm:"column:path;comment:请求路径;type:text;"`
	Status         int    `json:"status"                 form:"status"               gorm:"column:status;comment:状态码;type:varchar(4);"`
	Latency        int64  `json:"latency"                form:"latency"              gorm:"column:latency;comment:延迟;type:varchar(5);"`
	UserAgent      string `json:"userAgent"              form:"userAgent"            gorm:"column:user_agent;comment:User-Agent;type:text;"`
	Error          string `json:"error"                  form:"error"                gorm:"column:error;comment:错误;type:varchar(255);"`
	Body           string `json:"body"                   form:"body"                 gorm:"column:body;comment:Body;type:text;"`
	Query          string `json:"query"                  form:"query"                gorm:"column:query;comment:Query;type:text;"`
	Header         string `json:"header"                 form:"header"               gorm:"column:header;comment:Header;type:text;"`
	Response       string `json:"response"               form:"response"             gorm:"column:response;comment:响应;type:text;"`
	Location       string `json:"location"               form:"location"             gorm:"column:location;comment:位置;type:varchar(30);"`
	Country        string `json:"country"                form:"country"              gorm:"column:country;comment:国家;type:varchar(255)"`
	Province       string `json:"province"               form:"province"             gorm:"column:province;comment:省;type:varchar(20);"`
	City           string `json:"city"                   form:"city"                 gorm:"column:city;comment:市;type:varchar(20);"`
	District       string `json:"district"               form:"district"             gorm:"column:district;comment:区;type:varchar(20);"`
	Isp            string `json:"isp"                    form:"isp"                  gorm:"column:isp;comment:运营商;type:varchar(20);"`
	Platform       string `json:"platform"               form:"platform"             gorm:"column:platform;comment:平台;type:varchar(50);"`
	OS             string `json:"os"                     form:"os"                   gorm:"column:os;comment:系统;type:varchar(50);"`
	Engine         string `json:"engine"                 form:"engine"               gorm:"column:engine;comment:浏览器引擎;type:varchar(100);"`
	BrowserName    string `json:"browserName"            form:"browserName"          gorm:"column:browser_name;comment:浏览器;type:varchar(200);"`
	BrowserVersion string `json:"browserVersion"         form:"browserVersion"       gorm:"column:browser_version;comment:浏览器版本;type:varchar(50);"`
	Brand          string `json:"brand"                  form:"brand"                gorm:"column:brand;comment:品牌;type:varchar(255);"`
	ProductModel   string `json:"productModel"           form:"productModel"         gorm:"column:product_model;comment:型号;type:varchar(255);"`
	ContentType    string `json:"contentType"            form:"contentType"          gorm:"column:content_type;comment:内容类型;type:varchar(255);"`
	MerchantNo     string `json:"merchantNo"             form:"merchantNo"           gorm:"column:merchant_no;comment:商户号;type:varchar(32);"`
	PlatformType   int32  `json:"platformType"           form:"platformType"         gorm:"column:platform_type;comment:平台类型;type:int(3);"`
	AppProductId   int32  `json:"appProductId"           form:"appProductId"         gorm:"column:app_product_id;comment:应用Id;type:int(3);"`
	UserId         string `json:"userId"                 form:"userId"               gorm:"column:user_id;comment:用户id;type:varchar(64);"`
	UserName       string `json:"userName"               form:"userName"             gorm:"column:username;comment:用户名称;type:varchar(32);"`
}
