/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-06 16:16:58
 * @FilePath: \go-middleware\prometheus\constants.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package prometheus

var (
	Namespace           = "go_middleware" // 坑 不能是xx-xx
	DefaultIgnoredPaths = map[string]bool{
		"/metrics":     true,
		"/favicon.ico": true,
	}
)

// SetNamespace 设置Namespace
func SetNamespace(ns string) {
	Namespace = ns
}
