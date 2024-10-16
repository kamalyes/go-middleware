/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 23:49:58
 * @FilePath: \go-middleware\request\trace.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package request

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	incrNum       uint64
	pid           = os.Getpid()
	mu            sync.Mutex
	traceIdPrefix string
)

// SetTraceIdPrefix 设置TraceId Prefix，带锁
func SetTraceIdPrefix(ns string) {
	mu.Lock()
	defer mu.Unlock()
	traceIdPrefix = ns
}

// GetTraceIdPrefix 获取prefix
func GetTraceIdPrefix() string {
	return traceIdPrefix
}

// traceIDKey 用于在上下文中存储追踪ID的键
type traceIDKey struct{}

// generateTraceID 生成一个唯一的追踪ID
func generateTraceID() string {
	return fmt.Sprintf("%s-%d-%s-%d", GetTraceIdPrefix(), pid, time.Now().Format("20060102150405999"), atomic.AddUint64(&incrNum, 1))
}

// NewTraceIDContext 将追踪ID存储到上下文中
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// GetTraceID 从上下文中获取追踪ID
func GetTraceID(ctx context.Context) (string, bool) {
	if v := ctx.Value(traceIDKey{}); v != nil {
		if id, ok := v.(string); ok && id != "" {
			return id, true
		}
	}
	return "", false
}

// TraceMiddleware 追踪中间件处理函数
func TraceMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过中间件
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		// 从请求头中获取追踪ID，若为空则生成一个新的追踪ID
		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = generateTraceID()
		}

		// 将追踪ID存储到上下文中，并设置响应头中的追踪ID
		ctx := NewTraceIDContext(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Trace-Id", traceID)

		c.Next()
	}
}
