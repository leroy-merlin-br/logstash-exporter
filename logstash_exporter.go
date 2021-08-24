package main

import (
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/leroy-merlin-br/logstash_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	scrapeDurations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: collector.Namespace,
			Subsystem: "exporter",
			Name:      "scrape_duration_seconds",
			Help:      "logstash_exporter: Duration of a scrape job.",
		},
		[]string{"collector", "result"},
	)
)

// LogstashCollector collector type
type LogstashCollector struct {
	collectors map[string]collector.Collector
}

// NewLogstashCollector register a logstash collector
func NewLogstashCollector(logstashEndpoint string) (*LogstashCollector, error) {
	nodeStatsCollector, err := collector.NewNodeStatsCollector(logstashEndpoint)
	if err != nil {
		log.Error().AnErr("Cannot register a new collector", err)
	}

	nodeInfoCollector, err := collector.NewNodeInfoCollector(logstashEndpoint)
	if err != nil {
		log.Error().AnErr("Cannot register a new collector", err)
	}

	return &LogstashCollector{
		collectors: map[string]collector.Collector{
			"node": nodeStatsCollector,
			"info": nodeInfoCollector,
		},
	}, nil
}

func listen(exporterBindAddress string) {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})

	log.Info().Str("port", exporterBindAddress).Msg("Exporter started.")
	if err := http.ListenAndServe(exporterBindAddress, nil); err != nil {
		log.Error().Err(err).Msg("Exiting...")
	}
}

// Describe logstash metrics
func (coll LogstashCollector) Describe(ch chan<- *prometheus.Desc) {
	scrapeDurations.Describe(ch)
}

// Collect logstash metrics
func (coll LogstashCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(coll.collectors))
	for name, c := range coll.collectors {
		go func(name string, c collector.Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	scrapeDurations.Collect(ch)
}

func execute(name string, c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Collect(ch)
	duration := time.Since(begin)
	var result string

	if err != nil {
		log.Info().Msg("Failed to fetch metrics.")
		log.Debug().Str("collector", name).Float64("duration", duration.Seconds()).Err(err).Msg("Failed to fetch metrics.")
		result = "error"
	} else {
		log.Debug().Str("collector", name).Float64("duration", duration.Seconds()).Msg("Collected metrics with success.")
		result = "success"
	}
	scrapeDurations.WithLabelValues(name, result).Observe(duration.Seconds())
}

func init() {
	prometheus.MustRegister(version.NewCollector("logstash_exporter"))
}

func main() {
	var (
		logstashEndpoint    = kingpin.Flag("logstash.endpoint", "The protocol, host and port on which logstash metrics API listens.").Default("http://localhost:9600").String()
		exporterBindAddress = kingpin.Flag("listen.address", "Address on which to expose metrics and web interface.").Default(":9198").String()
		logLevel            = kingpin.Flag("log.level", "The logging level to be defined.").Default("info").String()
	)

	kingpin.Version(version.Print("logstash_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	switch *logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		log.Fatal().Msg("Log level needs to be set to \"info\" or \"debug\".")
	}

	logstashCollector, err := NewLogstashCollector(*logstashEndpoint)
	if err != nil {
		log.Error().Err(err).Msg("Cannot register a new logstash collector")
	}

	prometheus.MustRegister(logstashCollector)

	log.Info().Msg("Starting logstash exporter...")
	log.Debug().Msg(version.Info())
	log.Debug().Msg(version.BuildContext())

	listen(*exporterBindAddress)
}
