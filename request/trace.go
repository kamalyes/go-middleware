/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 20:18:01
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
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	incrNum uint64
	pid     = os.Getpid()
)

type traceIDKey struct{}

func generateTraceID() string {
	return fmt.Sprintf("trace-id-%d-%s-%d", pid, time.Now().Format("20060102150405999"), atomic.AddUint64(&incrNum, 1))
}

func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func GetTraceID(ctx context.Context) (string, bool) {
	if v := ctx.Value(traceIDKey{}); v != nil {
		if id, ok := v.(string); ok && id != "" {
			return id, true
		}
	}
	return "", false
}

func TraceMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = generateTraceID()
		}

		ctx := NewTraceIDContext(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Trace-Id", traceID)

		c.Next()
	}
}
