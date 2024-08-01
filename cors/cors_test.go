/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 17:40:22
 * @FilePath: \go-middleware\cors\cors_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(CorsMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test endpoint")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	// You can add more assertions based on your specific requirements

	// Example of additional assertions:
	assert.Equal(t, "POST, GET, OPTIONS, PUT, DELETE,UPDATE", w.Header().Get("Access-Control-Allow-Methods"))

	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}
