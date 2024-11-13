package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/ddosify/alaz/aggregator"
	"github.com/ddosify/alaz/cri"
	"github.com/ddosify/alaz/datastore"
	"github.com/ddosify/alaz/ebpf"
	"github.com/ddosify/alaz/k8s"
	"github.com/ddosify/alaz/log"
	"github.com/ddosify/alaz/logstreamer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	debug.SetGCPercent(80)
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()

	var k8sCollector *k8s.K8sCollector
	kubeEvents := make(chan interface{}, 1000)
	var k8sVersion string

	var k8sCollectorEnabled bool = true
	k8sEnabled, err := strconv.ParseBool(os.Getenv("K8S_COLLECTOR_ENABLED"))
	if err == nil && !k8sEnabled {
		k8sCollectorEnabled = false
	}

	if k8sCollectorEnabled {
		// k8s collector
		var err error
		k8sCollector, err = k8s.NewK8sCollector(ctx)
		if err != nil {
			panic(err)
		}
		k8sVersion = k8sCollector.GetK8sVersion()
		log.Logger.Info().Msgf("Kubertenes version: %s", k8sVersion)
		go k8sCollector.Init(kubeEvents)
	}

	tracingEnabled, err := strconv.ParseBool(os.Getenv("TRACING_ENABLED"))
	// logsEnabled, _ := strconv.ParseBool(os.Getenv("LOGS_ENABLED"))

	// Temporarily closed until otlp export's cpu performance problem is resolved
	// https://github.com/getanteon/alaz/tree/feat/logs-in-otlp
	// https://github.com/open-telemetry/opentelemetry-go/issues/5196
	logsEnabled := false

	// datastore backend
	reg := prometheus.NewRegistry()
	dsBackend := datastore.NewOodleDS(reg)

	var ct *cri.CRITool
	ct, err = cri.NewCRITool(ctx)
	if err != nil {
		log.Logger.Error().Err(err).Msg("failed to create cri tool")
	}

	// deploy ebpf programs
	var ec *ebpf.EbpfCollector
	if tracingEnabled {
		ec = ebpf.NewEbpfCollector(ctx, ct)

		a := aggregator.NewAggregator(ctx, ct, kubeEvents, ec.EbpfEvents(), ec.EbpfProcEvents(), ec.EbpfTcpEvents(), ec.TlsAttachQueue(), dsBackend)
		a.Run()

		ec.Init()
		go ec.ListenEvents()
	}

	var ls *logstreamer.LogStreamer
	if logsEnabled {
		if ct != nil {
			go func() {
				backoff := 5 * time.Second
				for {
					// retry creating LogStreamer with backoff
					// it will throw an error if connection to backend is not established
					log.Logger.Info().Msg("creating logstreamer")
					ls, err = logstreamer.NewLogStreamer(ctx, ct)
					if err != nil {
						log.Logger.Error().Err(err).Msg("failed to create logstreamer")
						select {
						case <-time.After(backoff):
						case <-ctx.Done():
							return
						}
						backoff *= 2
					} else {
						log.Logger.Info().Msg("logstreamer successfully created")
						break
					}
				}

				err := ls.StreamLogs()
				if err != nil {
					log.Logger.Error().Err(err).Msg("failed to stream logs")
				}
			}()

		} else {
			log.Logger.Error().Msg("logs enabled but cri tool not available")
		}
	}

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go http.ListenAndServe(":8181", nil)

	if k8sCollectorEnabled {
		<-k8sCollector.Done()
		log.Logger.Info().Msg("k8sCollector done")
	}

	if tracingEnabled {
		<-ec.Done()
		log.Logger.Info().Msg("ebpfCollector done")
	}

	if logsEnabled && ls != nil {
		<-ls.Done()
		log.Logger.Info().Msg("cri done")
	}

	<-ctx.Done()
}
