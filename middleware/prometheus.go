package middleware

import (
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/phantom-atom/file-explorer/internal/stats"
	"github.com/prometheus/client_golang/prometheus"
)

//Prometheus *
type Prometheus struct {
	name      string
	instance  string
	gatherer  *prometheus.Registry
	paramsFn  func() (string, int)
	onceStats *sync.Once
}

//DefaultMetricsParamsFn 默认参数
func DefaultMetricsParamsFn() (string, int) {
	return "http://localhost:9091", 15
}

//NewPrometheus *
func NewPrometheus(name, instance string, gatherer *prometheus.Registry,
	paramsFn func() (string, int)) *Prometheus {
	if paramsFn == nil {
		paramsFn = DefaultMetricsParamsFn
	}
	return &Prometheus{
		name:      name,
		instance:  instance,
		paramsFn:  paramsFn,
		gatherer:  gatherer,
		onceStats: new(sync.Once),
	}
}

func (p *Prometheus) startPrometheus() {
	go func() {
		stats.StartPushingMetric(p.name, p.instance, p.gatherer, p.paramsFn)
	}()
}

//HandlerFunc *
func (p *Prometheus) HandlerFunc() gin.HandlerFunc {
	p.onceStats.Do(func() {
		p.startPrometheus()
	})
	return func(c *gin.Context) {
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		url := c.Request.URL.String()
		stats.HTTPRequestCounter.WithLabelValues(status, c.Request.Method, url, c.Request.Host, c.HandlerName()).Inc()
	}
}
