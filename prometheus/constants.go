/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-12 23:03:53
 * @FilePath: \go-middleware\prometheus\constants.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package prometheus

import (
	"fmt"
	"sync"
)

var (
	Namespace    = "go_middleware" // 坑 不能是xx-xx
	mu           sync.Mutex
	IgnoredPaths = map[string]bool{
		"/metrics":     true,
		"/favicon.ico": true,
	}
)

// SetNamespace 设置Namespace，带锁
func SetNamespace(ns string) {
	mu.Lock()
	defer mu.Unlock()
	Namespace = ns
}

// GetNamespace 获取Namespace
func GetNamespace() string {
	return Namespace
}

// AddIgnoredPath 添加自定义的默认忽略路径
func AddIgnoredPath(path string) {
	mu.Lock()
	defer mu.Unlock()
	// 检查路径是否合法
	if path == "" {
		fmt.Println("Error: Path cannot be empty.")
		return
	}

	// 检查该路径是否已存在
	if _, exists := IgnoredPaths[path]; !exists {
		// 不存在则添加路径
		IgnoredPaths[path] = true
		fmt.Printf("Path '%s' has been added to the ignored paths.\n", path)
	}
}

// GetIgnoredPaths 获取忽略路径
func GetIgnoredPaths() map[string]bool {
	mu.Lock()
	defer mu.Unlock()
	return IgnoredPaths
}
