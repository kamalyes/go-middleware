/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-06 16:21:16
 * @FilePath: \go-middleware\prometheus\handler_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package prometheus

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMonitorMiddleware(t *testing.T) {
	// 创建一个请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)
	pm := NewPrometheusMonitorWithMetrics(engine)

	// 定义一个处理该请求的 handler
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// 模拟请求
	engine.ServeHTTP(w, req)

	// 断言
	assert.Equal(t, http.StatusOK, w.Code)

	// 运行 Prometheus 中间件
	middleware := pm.PrometheusMiddleware()
	middleware(ctx)

	// 你可以根据具体情况添加更多断言

}
