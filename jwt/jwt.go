/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 19:34:23
 * @FilePath: \go-middleware\jwt\jwt.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/jwt"
	"github.com/kamalyes/go-middleware/internal"
)

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if strings.Contains(path, "swagger") {
			ctx.Next()
			return
		}
		if strings.Contains(path, "login") || strings.Contains(path, "health") || strings.Contains(path, "captcha") {
			ctx.Next()
			return
		}
		token := ctx.Request.Header.Get("ACCESS_TOKEN")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": internal.ErrUnauthorized,
			})
			ctx.Abort()
			return
		}
		j := jwt.NewJWT()
		// 解析token包含的信息
		claims, err := j.ResolveToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
