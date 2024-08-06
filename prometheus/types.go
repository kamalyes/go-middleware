/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-06 16:04:12
 * @FilePath: \go-middleware\prometheus\types.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package prometheus

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// HandlerPath 定义采样路由struct
type HandlerPath struct {
	sync.Map
}

// PrometheusMonitor 定义 Prometheus 监控类
type PrometheusMonitor struct {
	Engine  *gin.Engine
	Ignored map[string]bool
	PathMap *HandlerPath
	Updated bool
}

type Option func(*PrometheusMonitor)
