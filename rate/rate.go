/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-10-16 11:12:26
 * @FilePath: \go-middleware\rate\rate.go
 * @Description: 限流中间件
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package rate

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kamalyes/go-config/global"
	"github.com/kamalyes/go-core/response"
	"github.com/kamalyes/go-middleware/constants"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

type limiterOptions struct {
	KeyPrefix string
}

func newLimiter(formatted string, opts limiterOptions) *limiter.Limiter {
	rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		panic(err)
	}

	var store limiter.Store
	if global.REDIS != nil {
		store, err = redis.NewStoreWithOptions(global.REDIS, limiter.StoreOptions{
			Prefix: opts.KeyPrefix,
		})
		if err != nil {
			panic(err.Error())
		}
	}

	if store == nil {
		store = memory.NewStore()
	}

	return limiter.New(store, rate)
}

func rateHandler(c *gin.Context, l *limiter.Limiter, key string) {
	context, err := l.Get(c, key)
	if err != nil {
		respOption := response.ResponseOption{Message: constants.ErrGainClientKey, HttpCode: response.StatusInternalServerError, Code: response.Fail}
		respOption.Sub(c)
		c.Abort()
		return
	}
	setRateHeaders(context, c)
	if context.Reached {
		respOption := response.ResponseOption{Message: constants.ErrTooManyRequests, HttpCode: response.StatusTooManyRequests, Code: response.RateLimit}
		respOption.Sub(c)
		c.Abort()
		return
	}

	c.Next()
}

func setRateHeaders(context limiter.Context, c *gin.Context) {
	c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
	c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))
}

// Rate 限流中间件（对每个client限流）
func Rate(formatted string) gin.HandlerFunc {
	l := newLimiter(formatted, limiterOptions{KeyPrefix: "Rate"})
	return func(c *gin.Context) {
		rateHandler(c, l, c.ClientIP()+":"+c.Request.RequestURI)
	}
}

// Rate0 限流中间件（总访问量限流）
func Rate0(formatted string) gin.HandlerFunc {
	l := newLimiter(formatted, limiterOptions{KeyPrefix: "Rate0"})
	return func(c *gin.Context) {
		rateHandler(c, l, "ALL:"+c.Request.RequestURI)
	}
}
