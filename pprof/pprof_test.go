/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-01 13:50:19
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-01 15:32:35
 * @FilePath: \go-middleware\pprof\pprof_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package pprof

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetPrefix(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"default value", nil, "/debug/pprof"},
		{"test user input value", []string{"test/pprof"}, "test/pprof"},
		{"test user input value", []string{"test/pprof", "pprof"}, "test/pprof"},
	}
	for _, tt := range tests {
		if got := getPrefix(tt.args...); got != tt.want {
			t.Errorf("%q. getPrefix() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestRegisterAndRouteRegister(t *testing.T) {
	bearerToken := "Bearer token"
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	Register(r)
	adminGroup := r.Group("/admin", func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != bearerToken {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	})
	PprofRouteRegister(adminGroup, "pprof")

	req, _ := http.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if expected, got := http.StatusOK, rw.Code; expected != got {
		t.Errorf("expected: %d, got: %d", expected, got)
	}

	req, _ = http.NewRequest(http.MethodGet, "/admin/pprof/", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if expected, got := http.StatusForbidden, rw.Code; expected != got {
		t.Errorf("expected: %d, got: %d", expected, got)
	}

	req, _ = http.NewRequest(http.MethodGet, "/admin/pprof/", nil)
	req.Header.Set("Authorization", bearerToken)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if expected, got := http.StatusOK, rw.Code; expected != got {
		t.Errorf("expected: %d, got: %d", expected, got)
	}
}
