/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 15:17:36
 * @FilePath: \go-middleware\cors\cors.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package cors

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-config/cors"
	"github.com/kamalyes/go-toolbox/convert"
)

// CorsMiddleware 跨域中间件
func CorsMiddleware(config cors.Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if !config.AllowedAllOrigins && !isOriginAllowed(origin, config.AllowedOrigins) {
			c.AbortWithStatus(config.OptionsResponseCode)
			return
		}

		setCorsHeaders(c, origin, config)
		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(config.OptionsResponseCode)
			return
		}

		defer recoverFromPanic()
		c.Next()
	}
}

// setCorsHeaders 设置Cors头部信息
func setCorsHeaders(c *gin.Context, origin string, config cors.Cors) {
	// 接收客户端发送的origin （重要）
	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	// 服务器支持的所有跨域请求的方法
	c.Header("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ","))
	// 允许跨域设置可以返回其他子段，可以自定义字段
	c.Header("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ","))
	// 允许浏览器（客户端）可以解析的头部 （重要）
	c.Header("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ","))
	// 设置缓存时间
	c.Header("Access-Control-Max-Age", config.MaxAge)
	// 允许客户端传递校验信息比如 cookie (重要)
	c.Header("Access-Control-Allow-Credentials", convert.MustString(config.AllowCredentials))
}

// isOriginAllowed 检查请求的 Origin 是否在允许的列表中
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if len(allowedOrigins) == 0 {
		return false
	}
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || strings.TrimSpace(allowedOrigin) == origin {
			return true
		}
	}
	return false
}

func recoverFromPanic() {
	if err := recover(); err != nil {
		log.Printf("Panic info: %v", err)
	}
}
