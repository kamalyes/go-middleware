/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 09:07:41
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
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/global"
	"go.uber.org/zap"
)

// GinRecoveryMiddleware 用于recover可能出现的panic，并使用zap记录相关日志
func GinRecoveryMiddleware(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查错误类型
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errorMsg := strings.ToLower(se.Error())
						if strings.Contains(errorMsg, "broken pipe") || strings.Contains(errorMsg, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 获取请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 错误处理
				if brokenPipe {
					global.LOG.Error("Broken pipe or connection reset by peer",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					_ = c.Error(err.(error))
					c.Abort()
					return
				}

				// 日志记录
				logFields := []zap.Field{
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
				}
				if stack {
					logFields = append(logFields, zap.String("stack", string(debug.Stack())))
				}

				global.LOG.Error("Recovery from panic",
					logFields...,
				)

				// 返回500错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
