/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 11:24:23
 * @FilePath: \go-middleware\recovery\recovery.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package recovery

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/global"
	"go.uber.org/zap"
)

// GinRecoveryMiddleware 用于 recover 可能出现的 panic，并使用 zap 记录相关日志
func GinRecoveryMiddleware(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() { handlePanic(c) }()
		c.Next()
	}
}

func handlePanic(c *gin.Context) {
	if err := recover(); err != nil {
		httpRequest, _ := httputil.DumpRequest(c.Request, true)
		if isBrokenPipe(err) {
			logBrokenPipe(c, err, httpRequest)
			c.Error(err.(error))
			c.Abort()
			return
		}
		logRecovery(c, err, httpRequest)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "服务器内部错误，请稍后再试",
		})
	}
}

func isBrokenPipe(err interface{}) bool {
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			errStr := strings.ToLower(se.Error())
			return strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer")
		}
	}
	return false
}

func logBrokenPipe(c *gin.Context, err interface{}, httpRequest []byte) {
	global.LOG.Error(c.Request.URL.Path,
		zap.String("request_id", c.GetString("request_id")),
		zap.Time("time", time.Now()),
		zap.Any("error", err),
		zap.String("request", string(httpRequest)),
	)
}

func logRecovery(c *gin.Context, err interface{}, httpRequest []byte) {
	global.LOG.Error("recovery from panic",
		zap.String("request_id", c.GetString("request_id")),
		zap.Time("time", time.Now()),
		zap.Any("error", err),
		zap.String("request", string(httpRequest)),
		zap.Stack("stacktrace"),
	)
}
