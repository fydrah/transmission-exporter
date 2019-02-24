package main

import (
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
	url := os.Getenv("TRANSMISSION_URL")
	if len(url) > 0 {
		log.Printf("Transmission URL set to %v", url)
	}

	user := os.Getenv("TRANSMISSION_USER")
	if len(user) > 0 {
		log.Printf("Transmission User is set to %v", user)
	}

	password := os.Getenv("TRANSMISSION_PASSWORD")
	if len(password) > 0 {
		log.Printf("Transmission password is set")
	}

	client := transmission.New(url, user, password)
	torrents, err := client.GetTorrents()
	if err != nil {
		log.Panic(err)
	}

	go func() {
		totalTorrent := len(torrents)
		torrentCount.With(prometheus.Labels{"instance": url}).Set(float64(totalTorrent))
	}()

	go func() {
		for _, torrent := range torrents {
			seedersCount.With(prometheus.Labels{"instance": url, "torrent": torrent.Name}).Set(float64(torrent.Seeders))
		}
	}()

	go func() {
		for _, torrent := range torrents {
			rateUpload.With(prometheus.Labels{"instance": url, "torrent": torrent.Name}).Set(float64(torrent.RateUpload))
		}
	}()

	go func() {
		for {
			for _, torrent := range torrents {
				rateDownload.With(prometheus.Labels{"instance": url, "torrent": torrent.Name}).Set(float64(torrent.RateDownload))
			}
		}
	}()
}

var (
	torrentCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "transmission_count_torrents",
		Help: "The number of total torrents",
	}, []string{"instance"})
	seedersCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "transmission_count_seeders",
		Help: "The number of seeders per torrent",
	}, []string{"instance", "torrent"})
	rateUpload = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "transmission_byte_rateupload",
		Help: "rate upload in bytes/second",
	}, []string{"instance", "torrent"})
	rateDownload = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "transmission_byte_ratedownload",
		Help: "rate download in bytes/second",
	}, []string{"instance", "torrent"})
)

func main() {
	for {
		recordMetrics()
		time.Sleep(20 * time.Second)
	}
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}
