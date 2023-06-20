package main

import (
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	cpuTemp      prometheus.Gauge
	requestTotal prometheus.Counter
	hdFailures   *prometheus.CounterVec
	summary      prometheus.Summary
	summaryVec   *prometheus.SummaryVec
	histogram    prometheus.Histogram
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}),

		requestTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "htt_request_total",
			Help: "Number of request.",
		}),

		hdFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hd_errors_total",
				Help: "Number of hard-disk errors.",
			},
			[]string{"device"},
		),

		summary: prometheus.NewSummary(prometheus.SummaryOpts{
			Name:       "pond_temperature_celsius",
			Help:       "The temperature of the frog pond.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}),

		summaryVec: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "pond_temperature",
				Help:       "The temperature of pond.",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
			},
			[]string{"species"},
		),

		histogram: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "pond_temperature_celsius",
			Help:    "The temperature of the frog pond.", // Sorry, we can't measure how badly it smells.
			Buckets: prometheus.LinearBuckets(20, 5, 5),  // 5 buckets, each 5 centigrade wide.
		}),
	}
	reg.MustRegister(m.cpuTemp)
	reg.MustRegister(m.requestTotal)
	reg.MustRegister(m.hdFailures)
	reg.MustRegister(m.summary)
	reg.MustRegister(m.summaryVec)

	return m
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func SimulateStat(reg prometheus.Registerer) {
	go func() {
		// Create new metrics and register them using the custom registry.
		m := NewMetrics(reg)
		for {
			// gauge update
			m.cpuTemp.Set(rand.Float64() * 100.0)

			// counter update
			if rand.Intn(10)%2 == 0 {
				m.requestTotal.Inc()
			} else {
				m.requestTotal.Add(7)
			}

			// counterVec update
			m.hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
			m.hdFailures.WithLabelValues("/dev/sdb").Inc()

			// summary update
			for i := 0; i < 10; i++ {
				m.summary.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
			}

			//summaryVec update
			for i := 0; i < 10; i++ {
				m.summaryVec.WithLabelValues("litoria-caerulea").Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
				m.summaryVec.WithLabelValues("lithobates-catesbeianus").Observe(32 + math.Floor(100*math.Cos(float64(i)*0.11))/10)
			}

			// Simulate some observations.
			for i := 0; i < 10; i++ {
				m.histogram.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
			}

			// Just for demonstration, let's check the state of the summary by
			// (ab)using its Write method (which is usually only used by Prometheus
			// internally).
			// metric := &dto.Metric{}
			// m.summary.Write(metric)
			// fmt.Println(metric.String())
			time.Sleep(time.Second * 2)
		}
	}()
}

func main() {

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.

	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	SimulateStat(reg)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.HandleFunc("/hello", HelloHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
