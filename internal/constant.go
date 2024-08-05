/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-07-31 20:25:55
 * @FilePath: \go-middleware\internal\constant.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package internal

const (
	prefix     = "go-middleware"
	ReqBodyKey = prefix + "/req-body"
	ResBodyKey = prefix + "/res-body"
)

const (
	GainSuccess           = "获取成功"
	ErrGainClientKey      = "获取Client限流Key失败"
	ErrUnauthorized       = "请求未携带token,无访问权限！"
	ErrTooManyRequests    = "接口访问超过限制"
	ErrLimiterInit        = "限流中间件(Redis)初始化出错: "
	ErrParseRequestData   = "获取失败,解析请求参数异常"
	ErrGainRecordResponse = "获取接口访问记录失败"
	ErrCreateRecordTable  = "创建访问记录表失败"
	ErrDeleteRecordData   = "删除访问记录异常"
)
