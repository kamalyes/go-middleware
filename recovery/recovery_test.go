/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 17:28:42
 * @FilePath: \go-middleware\recovery\recovery_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package recovery

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	kgoConfig "github.com/kamalyes/go-config"
	kgoGlobal "github.com/kamalyes/go-core/global"
	"github.com/kamalyes/go-core/zap"
	"github.com/stretchr/testify/assert"
)

func TestGinRecoveryMiddleware(t *testing.T) {
	// 获取全局配置
	kgoGlobal.CONFIG = kgoConfig.GlobalConfig()
	// 初始化zap日志
	kgoGlobal.LOG = zap.Zap()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(GinRecoveryMiddleware(true))

	router.GET("/ok", func(c *gin.Context) {
		// 返回 200 状态码
		c.AbortWithStatus(200)
	})

	okRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ok", nil)
	router.ServeHTTP(okRecorder, req)
	assert.Equal(t, http.StatusOK, okRecorder.Code, "Status code is 200")

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	panicRecorder := httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/panic", nil)
	router.ServeHTTP(panicRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, panicRecorder.Code, "Status code is not 500")
	expectedResponse := "{\"message\":\"服务器内部错误，请稍后再试\"}"
	assert.Contains(t, panicRecorder.Body.String(), expectedResponse, "Response body does not contain expected message")
}
