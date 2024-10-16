/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 23:45:43
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
	"github.com/stretchr/testify/assert"
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
	expectedTraceID := "trace-id-"

	if !ok {
		t.Errorf("Expected traceID in context, got nothing")
	}
	assert.Contains(t, traceID, expectedTraceID, "Expected TraceID should be contained in the actual TraceID")

	expectedHeader := "X-Trace-Id"
	header := w.Header().Get(expectedHeader)
	assert.Contains(t, header, expectedTraceID, "Expected header TraceID should be contained in the actual TraceID")
}
