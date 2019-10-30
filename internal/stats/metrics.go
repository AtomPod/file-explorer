package stats

import (
	"time"

	"github.com/phantom-atom/file-explorer/internal/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	//HTTPGather *
	HTTPGather = prometheus.NewRegistry()
	//HTTPRequestCounter *
	HTTPRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "FileService",
			Subsystem: "HTTP",
			Name:      "request_total",
			Help:      "Counter of http requests.",
		},
		[]string{"code", "method", "url", "host", "handler"},
	)
)

func init() {
	HTTPGather.MustRegister(HTTPRequestCounter)
}

//StartPushingMetric 开始推送指标信息
func StartPushingMetric(name, instance string, gatherer *prometheus.Registry, paramsFn func() (addr string, intervalSeconds int)) {
	if paramsFn == nil {
		return
	}

	addr, interval := paramsFn()
	pusher := push.New(addr, name).Gatherer(gatherer).Grouping("instance", instance)
	currentAddr := addr

	for {
		if currentAddr != "" {
			if err := pusher.Push(); err != nil {
				log.Info("msg", "count not push metrics to prometheus", "addr", currentAddr, "err", err.Error())
			}
		}

		if interval <= 0 {
			interval = 15
		}

		time.Sleep(time.Duration(interval) * time.Second)
		addr, interval = paramsFn()
		if addr != currentAddr {
			pusher = push.New(addr, name).Gatherer(gatherer).Grouping("instance", instance)
			currentAddr = addr
		}
	}
}
