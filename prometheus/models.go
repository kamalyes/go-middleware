/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-08-06 16:12:23
 * @FilePath: \go-middleware\prometheus\models.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package prometheus

import (
	"github.com/kamalyes/go-toolbox/system"
	"github.com/prometheus/client_golang/prometheus"
)

// HTTPHistogram prometheus 模型
var (
	LabelNames    = []string{"service", "method", "code", "endpoint", "client_ip", "size"}
	subSystem     = system.SafeGetHostName()
	HTTPHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   Namespace,
		Subsystem:   subSystem,
		Name:        "response_info",
		Help:        "Histogram of response of http handlers.",
		ConstLabels: nil,
		Buckets:     []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1}, // 代表duration的分布区间
	}, LabelNames)

	SlowApplies = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: Namespace,
		Subsystem: subSystem,
		Name:      "slow_apply_total",
		Help:      "The total number of slow apply requests (likely overloaded from slow disk).",
	})

	ApplySec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: Namespace,
		Subsystem: subSystem,
		Name:      "apply_duration_seconds",
		Help:      "The latency distributions of v2 apply called by backend.",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 20),
	},
		[]string{"version", "op", "success"})
)
