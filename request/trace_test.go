/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 20:17:04
 * @FilePath: \go-middleware\request\trace_test.go
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

func TestTraceMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	req, _ := http.NewRequest("GET", "/test", nil)
	// req.Header.Set("X-Request-Id", "123456")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	middleware := TraceMiddleware()

	middleware(c)

	ctx := c.Request.Context()
	traceID, ok := GetTraceID(ctx)

	if !ok {
		t.Errorf("Expected traceID in context, got nothing")
	}

	if traceID != "123456" {
		t.Errorf("Expected traceID to be '123456', got %s", traceID)
	}

	expectedHeader := "X-Trace-Id"
	if header := w.Header().Get(expectedHeader); header != "123456" {
		t.Errorf("Expected %s header to be '123456', got %s", expectedHeader, header)
	}
}
