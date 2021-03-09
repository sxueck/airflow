package prome

import (
	"context"
	. "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func init() {
	//var (
	//	PromeUserInfo = NewGauge(GaugeOpts{
	//		Name: "user_info",
	//		Help: "some fixed user information",
	//		ConstLabels: map[string]string{
	//			"name":  <-ChanName,
	//			"level": <-ChanLevel,
	//		},
	//	})
	//)

	// Metrics have to be registered to be exposed
	log.Println("init metrics")
	_ = Register(PromeBalance)
	_ = Register(PromeMaxBandwidth)
	_ = Register(PromeOnlineDeviceCount)
	_ = Register(PromeRemainFlow)
	_ = Register(PromeRemainTime)
	_ = Register(PromeTodayUsed)
	// _ = Register(PromeUserInfo)
}

func StartPromeServ() {
	ctx, cancel := context.WithCancel(context.Background())
	http.Handle("/metrics", promhttp.Handler())
	go RecvMetricsValue(ctx)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
	cancel()
	time.Sleep(2 * time.Second)
}

func RecvMetricsValue(ctx context.Context) {
	log.Println("start accepting metrics...")
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-ChanTodayUsed:
			PromeTodayUsed.Set(float64(v))
		case v := <-ChanRemainTime:
			PromeRemainTime.Set(float64(v))
		}
	}
}
