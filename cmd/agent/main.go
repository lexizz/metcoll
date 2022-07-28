package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lexizz/metcoll/internal/helper"
	"github.com/lexizz/metcoll/internal/metrics"
)

const (
	urlDestinationConst string        = "http://127.0.0.1:8080"
	pollInterval        time.Duration = 2
	reportInterval      time.Duration = 10
	requestTimeout      time.Duration = 300

	methodSendingGet  string = "get"
	methodSendingPost string = "post"
)

type listUrls []string

type exporter struct {
	httpClient        *http.Client
	metrics           *metrics.Metrics
	metricsData       metrics.Collection
	methodSendingData string
}

func main() {
	exp := exporter{
		httpClient:        &http.Client{},
		metrics:           metrics.New(),
		methodSendingData: methodSendingPost,
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
			exp.metricsData = exp.metrics.CollectData()
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
	switch exporter.methodSendingData {
	case methodSendingGet:
		for _, url := range exporter.getListUrls() {
			log.Printf("Request to url: %v\n", url)

			exporter.sendRequest(url, http.MethodGet, nil)
		}
	case methodSendingPost:
		metricsData := exporter.metricsData

		for metricName, value := range metricsData {
			metricType, err := helper.GetType(value)
			if err != nil {
				log.Printf("--- metricType: %v, metricType; ERR get type: %v", metricType, err)
			}

			met := metrics.Metrics{
				ID:    metricName,
				MType: metricType,
			}

			switch val := value.(type) {
			case metrics.Gauge:
				preparedValueGauge := float64(val)
				met.Value = &preparedValueGauge
			case metrics.Counter:
				preparedValueCounter := int64(val)
				met.Delta = &preparedValueCounter
			}

			requestBody, err := json.Marshal(met)
			if err != nil {
				log.Printf("--- Error json.Marshal: %v", err)

				return
			}

			exporter.sendRequest(urlDestinationConst+"/update", http.MethodPost, requestBody)
		}
	}
}

func (exporter *exporter) sendRequest(url string, method string, requestBody []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*time.Second)
	defer cancel()

	var buf io.Reader

	if len(requestBody) > 0 {
		buf = bytes.NewBuffer(requestBody)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		log.Println(err)

		return
	}

	request.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	if method == http.MethodPost {
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	}

	response, err := exporter.httpClient.Do(request)
	if err != nil {
		log.Println(err)

		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)

		return
	}

	errorResponseClose := response.Body.Close()
	if errorResponseClose != nil {
		log.Println(errorResponseClose)

		return
	}

	fmt.Printf("Response: %v | Body: %v\n-----------\n", response.Status, string(body))
}

func (exporter *exporter) getListUrls() listUrls {
	urls := make(listUrls, 0, 50)

	for metricName, metricValue := range exporter.metricsData {
		metricType, err := helper.GetType(metricValue)
		if err != nil {
			log.Printf("--- metricType: %v, metricType; ERR get type: %v", metricType, err)
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
