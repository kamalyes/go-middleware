/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 23:05:41
 * @FilePath: \go-middleware\prometheus\handler.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package prometheus

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// 提供选项功能来忽略特定路径
func Ignore(path ...string) func(*PrometheusMonitor) {
	return func(pm *PrometheusMonitor) {
		for _, p := range path {
			pm.Ignored[p] = true
		}
	}
}

// 创建 Prometheus 监控实例并初始化
func NewPrometheusMonitorWithMetrics(e *gin.Engine, options ...func(*PrometheusMonitor)) *PrometheusMonitor {
	if e == nil {
		return nil
	}

	pm := &PrometheusMonitor{
		Engine:  e,
		Ignored: GetIgnoredPaths(),
		PathMap: &HandlerPath{},
	}

	// 初始化 Prometheus 模型
	prometheus.MustRegister(HTTPHistogram)
	prometheus.MustRegister(ApplySec)
	prometheus.MustRegister(SlowApplies)

	for _, o := range options {
		o(pm)
	}

	return pm
}

// GetHandlerPath 获取path
func (hp *HandlerPath) GetHandlerPath(handler string) string {
	v, ok := hp.Load(handler)
	if !ok {
		return ""
	}
	return v.(string)

}

// SetHandlerPath 保存path到sync.Map
func (hp *HandlerPath) SetHandlerPath(ri gin.RouteInfo) {
	hp.Store(ri.Handler, ri.Path)
}

// 更新路由信息
func (pm *PrometheusMonitor) updatePath() {
	pm.Updated = true
	for _, ri := range pm.Engine.Routes() {
		pm.PathMap.SetHandlerPath(ri)
	}
}

// 中间件函数，用于收集 Prometheus 监控信息
func (pm *PrometheusMonitor) PrometheusMiddleware() gin.HandlerFunc {
	start := time.Now()
	// 初始化 Ignored map，避免空指针异常
	if pm.Ignored == nil {
		pm.Ignored = make(map[string]bool)
	}
	return func(c *gin.Context) {
		if !pm.Updated {
			pm.updatePath()
		}

		// 检查c.Request.URL是否为nil
		if c.Request == nil {
			c.Next()
			return
		}

		// 过滤忽略的路径
		if pm.Ignored[c.Request.URL.String()] {
			c.Next()
			return
		}

		c.Next()
		HTTPHistogram.WithLabelValues(
			"prometheus",
			c.Request.Method,
			strconv.Itoa(c.Writer.Status()),
			pm.PathMap.GetHandlerPath(c.HandlerName()),
			c.ClientIP(),
			strconv.Itoa(c.Writer.Size()),
		).Observe(time.Since(start).Seconds())
	}
}
