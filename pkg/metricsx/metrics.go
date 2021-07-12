package metricsx

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics 结构体
type Metrics struct {
	Name  string
	Help  string
	Label []string
}

// MetricsCounter 结构体
type MetricsCounter struct {
	Metrics
	CounterVec *prometheus.CounterVec
}

// MetricsSummary 结构体
type MetricsSummary struct {
	Metrics
	Objectives map[float64]float64
	SummaryVec *prometheus.SummaryVec
}

// MetricsHistogram 结构体
type MetricsHistogram struct {
	Metrics
	Buckets      []float64
	HistogramVec *prometheus.HistogramVec
}

// MetricsGauge 结构体
type MetricsGauge struct {
	Metrics
	GaugeVec *prometheus.GaugeVec
}

type MetricValue struct {
	Label []string
	Num   float64
}

// MetricsAdd 累加Counter类型数值
func (c *MetricsCounter) MetricsAdd(s []string, num float64) {
	c.CounterVec.WithLabelValues(s...).Add(num)
}

// MetricsAdd 增加Summary类型样例
func (s *MetricsSummary) MetricsAdd(str []string, num float64) {
	s.SummaryVec.WithLabelValues(str...).Observe(num)
}

// MetricsAdd 增加Histogram类型样例
func (s *MetricsHistogram) MetricsAdd(str []string, num float64) {
	s.HistogramVec.WithLabelValues(str...).Observe(num)
}

// MetricsAdd 增加Gauge类型样例
func (g *MetricsGauge) MetricsAdd(str []string, num float64) {
	g.GaugeVec.WithLabelValues(str...).Add(num)
}

// MetricsSet 给Gauge类型设置值
func (g *MetricsGauge) MetricsSet(str []string, num float64) {
	// g.GaugeVec.WithLabelValues(str...).Set(num)
	g.GaugeVec.WithLabelValues(str...).Set(num)
}

// NewCounterVec 第一次调用会初始化CounterVec，后面会返回*prometheus.CounterVec
func (g *MetricsGauge) NewGaugeVec() {
	g.GaugeVec = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: g.Name,
			Help: g.Help,
		},
		g.Label,
	)
	// c.chCounter = make(chan *MetricValue, 100000)
	prometheus.MustRegister(g.GaugeVec)
}

// NewCounterVec 第一次调用会初始化CounterVec，后面会返回*prometheus.CounterVec
func (c *MetricsCounter) NewCounterVec() {
	c.CounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: c.Name,
			Help: c.Help,
		},
		c.Label,
	)
	prometheus.MustRegister(c.CounterVec)
}

// NewSummaryVec 第一次调用会初始化SummaryVec，后面会返回*prometheus.SummaryVec
func (s *MetricsSummary) NewSummaryVec() {
	s.SummaryVec = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       s.Name,
			Help:       s.Help,
			Objectives: s.Objectives,
		},
		s.Label,
	)
	prometheus.MustRegister(s.SummaryVec)
}

// NewHistogramVec 第一次调用会初始化SummaryVec，后面会返回*prometheus.SummaryVec
func (s *MetricsHistogram) NewHistogramVec() {
	s.HistogramVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    s.Name,
			Help:    s.Help,
			Buckets: s.Buckets,
		},
		s.Label,
	)
	prometheus.MustRegister(s.HistogramVec)
}
