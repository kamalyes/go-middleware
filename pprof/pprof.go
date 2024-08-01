/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-01 13:50:19
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 15:31:46
 * @FilePath: \go-middleware\pprof\pprof.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package pprof

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

const (
	// pprof的默认url前缀
	DefaultPrefix = "/debug/pprof"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

// 从. net/http/pprof包注册标准HandlerFuncs
// 提供的gin.IRouter。prefixOptions是可选的。如果不是prefixOptions，
// 使用默认的路径前缀，否则第一个prefixOptions将是路径前缀。
func Register(r gin.IRouter, prefixOptions ...string) {
	PprofRouteRegister(r, prefixOptions...)
}

// 将标准HandlerFuncs从net/http/pprof包中注册到
// 提供的gin.IRouter。prefixOptions是可选的。如果不是prefixOptions，
// 使用默认的路径前缀，否则第一个prefixOptions将是路径前缀。
func PprofRouteRegister(rg gin.IRouter, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)

	prefixRouter := rg.Group(prefix)
	{
		prefixRouter.GET("/", gin.WrapF(pprof.Index))
		prefixRouter.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		prefixRouter.GET("/profile", gin.WrapF(pprof.Profile))
		prefixRouter.POST("/symbol", gin.WrapF(pprof.Symbol))
		prefixRouter.GET("/symbol", gin.WrapF(pprof.Symbol))
		prefixRouter.GET("/trace", gin.WrapF(pprof.Trace))
		prefixRouter.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		prefixRouter.GET("/block", gin.WrapH(pprof.Handler("block")))
		prefixRouter.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		prefixRouter.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		prefixRouter.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		prefixRouter.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}
