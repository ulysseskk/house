package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RunningTaskGauge *prometheus.GaugeVec     // 用于记录运行中的数量
	VisitCounter     *prometheus.CounterVec   // 用于记录抓取频率
	TimeUse          *prometheus.HistogramVec // 用于记录抓取时间
)

func init() {
	RunningTaskGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "running_task",
		Help: "运行中的抓取数量",
	}, []string{
		"platform",
	})
	VisitCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "visit",
		Help: "访问网页的数量，用于做访问统计",
	}, []string{"platform"})
	TimeUse = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "visit_time_use",
		Help:    "",
		Buckets: []float64{10, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 1200, 1400, 1600, 1800, 2000, 2500, 3000, 3500, 4000, 4500, 5000, 6000, 7000, 8000, 9000, 10000, 12000, 14000, 16000, 18000, 20000, 25000, 30000, 40000, 50000, 60000, 70000, 80000, 90000, 100000},
	}, []string{"platform"})
	prometheus.MustRegister(RunningTaskGauge, VisitCounter, TimeUse)
}
