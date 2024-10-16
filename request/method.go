/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 16:29:09
 * @FilePath: \go-middleware\request\method.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package request

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-core/response"
)

// NoMethodHandler 处理请求方法不被允许的情况
func NoMethodHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response.GenGinResponse(ctx, &response.ResponseOption{Code: response.SceneCode(response.StatusMethodNotAllowed), Message: "方法不允许", HttpCode: response.StatusMethodNotAllowed})
	}
}

// NoRouteHandler 处理请求路由不存在的情况
func NoRouteHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response.GenGinResponse(ctx, &response.ResponseOption{Code: response.SceneCode(response.StatusNotFound), Message: "路由不存在", HttpCode: response.StatusNotFound})
	}
}

// SkipperFunc 定义跳过中间件的函数类型
type SkipperFunc func(*gin.Context) bool

// AllowPathPrefixSkipper 生成一个函数，用于检查请求路径是否以指定前缀开头
func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// AllowMethodAndPathPrefixSkipper 生成一个函数，用于检查请求方法和路径是否符合指定要求
func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(ctx *gin.Context) bool {
		path := JoinRouter(ctx.Request.Method, ctx.Request.URL.Path)
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// JoinRouter 拼接请求方法和路径并返回字符串
func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

// SkipHandler 执行一系列判断函数，用于跳过特定的中间件
func SkipHandler(ctx *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(ctx) {
			return true
		}
	}
	return false
}

// EmptyMiddleware 仅执行下一个中间件的处理函数
func EmptyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
