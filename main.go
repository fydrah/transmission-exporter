package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Krast76/transmission"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			url := os.Getenv("TRANSMISSION_URL")
			user := os.Getenv("TRANSMISSION_USER")
			password := os.Getenv("TRANSMISSION_PASSWORD")
			client := transmission.New(url, user, password)
			torrents, err := client.GetTorrents()
			if err != nil {
				log.Panic(err)
			}

			totalTorrent := len(torrents)
			torrentCount.With(prometheus.Labels{"instance": url}).Set(float64(totalTorrent))

			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	torrentCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "transmission_count_torrents",
		Help: "The number of total torrents",
	}, []string{"instance"})
)

func main() {
	recordMetrics()
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
