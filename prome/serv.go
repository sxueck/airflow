package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func StartPromeServ() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
