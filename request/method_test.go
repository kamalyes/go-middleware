/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 20:36:18
 * @FilePath: \go-middleware\request\method_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNoMethodHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := NoMethodHandler()
	handler(c)

	// 在这里添加断言来验证处理函数的行为
}

func TestNoRouteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := NoRouteHandler()
	handler(c)

	// 在这里添加断言来验证处理函数的行为
}

func TestJoinRouter(t *testing.T) {
	expected := "GET/example"
	result := JoinRouter("GET", "example")

	if result != expected {
		t.Errorf("JoinRouter returned unexpected result. Expected: %s, Got: %s", expected, result)
	}
}

func TestAllowPathPrefixSkipper(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	skipper := AllowPathPrefixSkipper("/example")

	result := skipper(c)
	if !result {
		t.Error("AllowPathPrefixSkipper did not skip as expected")
	}
}

func TestAllowPathPrefixNoSkipper(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	skipper := AllowPathPrefixSkipper("/test")

	result := skipper(c)
	if result {
		t.Error("AllowPathPrefixNoSkipper skipped unexpectedly")
	}
}

func TestAllowMethodAndPathPrefixSkipper(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	skipper := AllowMethodAndPathPrefixSkipper("GET/example")

	result := skipper(c)
	if !result {
		t.Error("AllowMethodAndPathPrefixSkipper did not skip as expected")
	}
}
