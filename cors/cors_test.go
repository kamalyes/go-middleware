/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 19:06:29
 * @FilePath: \go-middleware\cors\cors_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package cors

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-config/cors"
	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a mock CORS configuration
	corsConfig := cors.Cors{
		AllowedAllOrigins:   false,
		AllowedAllMethods:   false,
		AllowedOrigins:      []string{"https://example.com"},
		AllowedMethods:      []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:      []string{"Content-Type", "Authorization"},
		MaxAge:              "3600",
		ExposedHeaders:      []string{"X-Custom-Header"},
		AllowCredentials:    true,
		OptionsResponseCode: http.StatusNoContent,
	}

	// Create a test handler
	testHandler := func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	}

	// Create a Gin router and apply the CORS middleware
	router := gin.New()
	router.Use(CorsMiddleware(corsConfig))
	router.GET("/", testHandler)
	router.POST("/", testHandler)

	// Test cases
	tests := []struct {
		method          string
		origin          string
		expectedCode    int
		expectedHeader  string
		expectedBody    string
		expectedHeaders map[string]string
	}{
		{"GET", "https://example.com", http.StatusOK, "https://example.com", "Hello, world!", map[string]string{"Content-Type": "text/plain; charset=utf-8"}},
		{"POST", "https://example.com", http.StatusOK, "https://example.com", "Hello, world!", map[string]string{"Content-Type": "text/plain; charset=utf-8"}},
		{"OPTIONS", "https://example.com", http.StatusNoContent, "", "", map[string]string{}},
		{"GET", "https://unauthorized.com", http.StatusNoContent, "", "", map[string]string{}},
	}

	for _, ts := range tests {
		funcName := fmt.Sprintf("%v", &ts)
		t.Run(funcName, func(t *testing.T) {
			req := httptest.NewRequest(ts.method, "/", nil)
			req.Header.Set("Origin", ts.origin)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, ts.expectedCode, w.Code)

			// Check the response body
			assert.Equal(t, ts.expectedBody, w.Body.String())

			// Check the Access-Control-Allow-Origin header if applicable
			if ts.method != "OPTIONS" {
				assert.Equal(t, ts.expectedHeader, w.Header().Get("Access-Control-Allow-Origin"))
			}

			// Check expected headers
			for key, value := range ts.expectedHeaders {
				assert.Equal(t, value, w.Header().Get(key))
			}
		})
	}
}
