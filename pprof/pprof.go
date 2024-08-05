/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-01 13:50:19
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-05 18:27:05
 * @FilePath: \go-middleware\pprof\pprof.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package pprof

import (
	"fmt"
	"net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// pprof的默认url前缀
	DefaultPrefix = "/debug/pprof"
)

var startTime = time.Now()

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

func Handler(c *gin.Context) {
	m := NewSystemInfo(startTime)
	info := fmt.Sprintf("%s:%s\n", "服务器", m.ServerName)
	info += fmt.Sprintf("%s:%s\n", "运行时间", m.Runtime)
	info += fmt.Sprintf("%s:%s\n", "goroutine数量", m.GoroutineNum)
	info += fmt.Sprintf("%s:%s\n", "CPU核数", m.CPUNum)
	info += fmt.Sprintf("%s:%s\n", "当前内存使用量", m.UsedMem)
	info += fmt.Sprintf("%s:%s\n", "当前堆内存使用量", m.HeapInuse)
	info += fmt.Sprintf("%s:%s\n", "总分配的内存", m.TotalMem)
	info += fmt.Sprintf("%s:%s\n", "系统内存占用量", m.SysMem)
	info += fmt.Sprintf("%s:%s\n", "指针查找次数", m.Lookups)
	info += fmt.Sprintf("%s:%s\n", "内存分配次数", m.Mallocs)
	info += fmt.Sprintf("%s:%s\n", "内存释放次数", m.Frees)
	info += fmt.Sprintf("%s:%s\n", "距离上次GC时间", m.LastGCTime)
	info += fmt.Sprintf("%s:%s\n", "下次GC内存回收量", m.NextGC)
	info += fmt.Sprintf("%s:%s\n", "GC暂停时间总量", m.PauseTotalNs)
	info += fmt.Sprintf("%s:%s\n", "上次GC暂停时间", m.PauseNs)
	_, _ = fmt.Fprint(c.Writer, info)
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
		prefixRouter.GET("/sysinfo", Handler)
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
