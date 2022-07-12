package main

import (
	"context"
	"fmt"
	"github.com/lexizz/metcoll/internal/helper"
	"github.com/lexizz/metcoll/internal/metrics"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	urlDestinationConst string        = "http://127.0.0.1:8080"
	pollInterval        time.Duration = 2
	reportInterval      time.Duration = 10
)

type listUrls []string

type exporter struct {
	httpClient *http.Client
	metrics    metrics.Metrics
}

func main() {
	exp := exporter{
		httpClient: &http.Client{},
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context, cancelContext context.CancelFunc) {
		signalChanel := make(chan os.Signal, 1)
		signal.Notify(signalChanel,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		for {
			select {
			case <-signalChanel:
				cancelContext()
			case <-ctx.Done():
				os.Exit(0)
			}
		}
	}(ctx, cancel)

	tickerGetMetrics := time.NewTicker(pollInterval * time.Second)
	tickerSendData := time.NewTicker(reportInterval * time.Second)

	for {
		select {
		case <-tickerGetMetrics.C:
			exp.metrics = metrics.CollectData()
			fmt.Println("Get metrics...")
		case <-tickerSendData.C:
			exp.sendDataToServer()
		case <-ctx.Done():
			tickerGetMetrics.Stop()
			tickerSendData.Stop()
			os.Exit(0)
		}
	}
}

func (exporter *exporter) sendDataToServer() {
	for _, url := range exporter.getListUrls() {
		log.Printf("Request to url: %v\n", url)

		exporter.sendRequest(url)
	}
}

func (exporter *exporter) sendRequest(url string) {
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "text/plain; charset=UTF-8")

	response, err := exporter.httpClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	errorResponseClose := response.Body.Close()
	if errorResponseClose != nil {
		log.Fatal(errorResponseClose)
	}

	fmt.Printf("Response: %v | Body: %v\n-----------\n", response.Status, string(body))
}

func (exporter *exporter) getListUrls() listUrls {
	urls := make(listUrls, 0, 50)

	for metricName, metricValue := range exporter.metrics {
		metricType, err := helper.GetType(metricValue)

		if err != nil {
			log.Printf("--- metricType: %v, metricType; ERR: %v", metricType, err)
		}

		log.Printf("--- metricName: %v | metricType: %v | metricValue: %v\n", metricName, metricType, metricValue)

		url := urlDestinationConst + "/update/" + metricType + "/" + metricName + "/" + exporter.convertValueToString(metricValue, metricType)

		urls = append(urls, url)
	}

	return urls
}

func (exporter *exporter) convertValueToString(metricValue interface{}, metricType string) string {
	switch metricType {
	case "gauge":
		return fmt.Sprintf("%3.2f", metricValue)
	case "counter":
		return fmt.Sprintf("%d", metricValue)
	default:
		return fmt.Sprint(metricValue)
	}
}
