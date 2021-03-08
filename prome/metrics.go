package prome

import "github.com/prometheus/client_golang/prometheus"

var (
	ChanName  = make(chan string)
	ChanLevel = make(chan string)

	PromeUserInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "user_info",
		Help: "some fixed user information",
		ConstLabels: map[string]string{
			"name":  <-ChanName,
			"level": <-ChanLevel,
		},
	})

	PromeRemainTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "remain_time",
		Help: "the remaining time on the next settlement day",
	})

	PromeRemainFlow = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "remain_flow",
		Help: "remaining flow of billing period",
	})

	PromeMaxBandwidth = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "max_bandwidth",
	})

	PromeTodayUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "today_used",
		Help: "traffic used today",
	})

	PromeOnlineDeviceCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "online_device_count",
		Help: "number of current online devices",
	})

	PromeBalance = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "current_balance",
	})
)

var (
	ChanBalance           = make(chan float32, 1)
	ChanMaxBandwidth      = make(chan int, 1)
	ChanOnlineDeviceCount = make(chan int, 1)
	ChanRemainFlow        = make(chan int, 1)
	ChanRemainTime        = make(chan int, 1)
	ChanTodayUsed         = make(chan int, 1)
)
