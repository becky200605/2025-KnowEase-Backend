package utils

import (
	"github.com/prometheus/client_golang/prometheus"
)

// 定义指标-记录请求数据
var (
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "count request totals",
			Help: "counts of requests",
		},
		[]string{"Method", "Path", "StatusCode"},
	)
)

// 注册指标
func Init() {
	prometheus.MustRegister(RequestCounter)
}
