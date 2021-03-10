package prome

import (
	"context"
	"fmt"
	. "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func RegisterMetrics(name, level string) {
	var (
		PromeUserInfo = NewGauge(GaugeOpts{
			Name: "user_info",
			Help: "some fixed user information",
			ConstLabels: map[string]string{
				"name":  name,
				"level": level,
			},
		})
	)

	// Metrics have to be registered to be exposed
	log.Println("init metrics")
	_ = Register(PromeBalance)
	_ = Register(PromeMaxBandwidth)
	_ = Register(PromeOnlineDeviceCount)
	_ = Register(PromeRemainFlow)
	_ = Register(PromeRemainTime)
	_ = Register(PromeTodayUsed)
	_ = Register(PromeUserInfo)
}

func StartPromeServ(ctx context.Context, name, level string) {
	fmt.Println("starting server")
	RegisterMetrics(name, level)
	var errChan = make(chan error)
	go RecvMetricsValue(ctx)

	http.Handle("/metrics", promhttp.Handler())
	select {
	case errChan <- http.ListenAndServe(":8080", nil):
		log.Fatal(errChan)
	case <-ctx.Done():
	}
	return
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
		case v := <-ChanMaxBandwidth:
			PromeMaxBandwidth.Set(float64(v))
		case v := <-ChanOnlineDeviceCount:
			PromeOnlineDeviceCount.Set(float64(v))
		case v := <-ChanRemainFlow:
			PromeRemainFlow.Set(float64(v))
		case v := <-ChanBalance:
			PromeBalance.Set(float64(v))
		}
	}
}
