package main

import (
	"context"
	"fmt"
	"github.com/lexizz/metcoll/internal/metrics"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

const (
	urlDestinationConst string        = "http://127.0.0.1:8080"
	pollInterval        time.Duration = 2
	reportInterval      time.Duration = 10
)

func main() {
	var (
		stats metrics.Metrics
	)

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
			stats = metrics.Collect()
			fmt.Println("Get metrics")
		case <-tickerSendData.C:
			sendData(stats)
		case <-ctx.Done():
			tickerGetMetrics.Stop()
			tickerSendData.Stop()
			os.Exit(0)
		}
	}
}

func sendData(metrics metrics.Metrics) {
	client := http.Client{}

	var url string

	elem := reflect.ValueOf(&metrics).Elem()
	for i := 0; i < elem.NumField(); i++ {
		valueField := elem.Type().Field(i)
		metricName := valueField.Name
		metricType := valueField.Type.Name()
		metricValue := elem.Field(i).Interface()

		var metricValuePrepare string

		switch metricType {
		case "gauge":
			metricValuePrepare = fmt.Sprintf("%.2f", metricValue)
		case "counter":
			metricValuePrepare = fmt.Sprintf("%d", metricValue)
		default:
			metricValuePrepare = fmt.Sprint(metricValue)
		}

		url = urlDestinationConst + "/update/" + metricType + "/" + metricName + "/" + metricValuePrepare

		fmt.Printf("Request to url: %v\n", url)

		request, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		request.Header.Set("Content-Type", "text/plain; charset=UTF-8")

		response, err := client.Do(request)
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

}
