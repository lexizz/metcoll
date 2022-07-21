package handlers

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/lexizz/metcoll/internal/helper"
	"github.com/lexizz/metcoll/internal/metrics"
	"github.com/lexizz/metcoll/internal/repository/interfaces/metricrepository"
)

const (
	metricTypeCounterConst string = "counter"
	patternConst           string = `[0-9\.]+$`
)

func ShowPossibleValue(metricRepository metricrepository.Interface) http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		log.Println(request.Method, request.Host, request.URL.Path)
		log.Println("=== Work handler - ShowPossibleValue ===")

		if request.Method != http.MethodGet {
			log.Println("---ERR http-method")
			http.Error(write, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		write.Header().Set("content-type", "text/html; charset=utf-8")

		listMetrics, err := metricRepository.GetAll()
		if err != nil {
			http.Error(write, err.Error(), http.StatusNotFound)
			return
		}

		htmlMetrics := ""
		for metricName, metricValue := range listMetrics {
			htmlMetrics += metricName + ": " + fmt.Sprint(metricValue) + "; | <br>\n"
		}

		write.WriteHeader(http.StatusOK)

		_, writeError := write.Write([]byte(htmlMetrics))
		if writeError != nil {
			log.Println(writeError)
		}
	}
}

func ShowValueMetric(metricRepository metricrepository.Interface) http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		log.Println(request.Method, request.Host, request.URL.Path)
		log.Println("=== Work handler - ShowValueMetric ===")

		if request.Method != http.MethodGet {
			http.Error(write, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		write.Header().Set("content-type", "text/plain; charset=utf-8")

		metricName := chi.URLParam(request, "metricName")
		metricType := chi.URLParam(request, "metricType")

		log.Println("---metricType:", metricName)
		log.Println("---metricType:", metricType)

		resVal, err := metricRepository.GetValue(metricName)
		if err != nil {
			log.Println("---ERR GetValue:", err)
			http.Error(write, err.Error(), http.StatusNotFound)
			return
		}

		valueMetricType, errGetType := helper.GetType(resVal)
		if errGetType != nil {
			log.Println("---ERR GetType:", errGetType)
			http.Error(write, errGetType.Error(), http.StatusNotFound)
			return
		}

		if valueMetricType != metricType {
			log.Println("---ERR Wrong metric type")
			http.Error(write, "Wrong metric type", http.StatusMethodNotAllowed)
			return
		}

		result := fmt.Sprint(resVal)

		write.WriteHeader(http.StatusOK)

		_, writeError := write.Write([]byte(result))
		if writeError != nil {
			log.Println(writeError)
		}
	}
}

func UpdateMetric(metricRepository metricrepository.Interface) http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		log.Println(request.Method, request.Host, request.URL.Path)
		log.Println("=== Work handler - UpdateMetric ===")

		if request.Method != http.MethodPost {
			log.Println("---ERR http-method")
			http.Error(write, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		write.Header().Set("content-type", "text/plain; charset=utf-8")

		url := request.URL.Path
		url = strings.Trim(url, "/")
		urlData := strings.Split(url, "/")

		log.Println("---urlData:", urlData)

		if len(urlData) < 4 {
			writeError := "Params not found"
			http.Error(write, writeError, http.StatusNotFound)
			log.Println(writeError)

			return
		}

		metricFromRequestName := urlData[2]
		metricFromRequestType := urlData[1]
		metricFromRequestValue := urlData[3]
		log.Println("---metricFromRequestName:", metricFromRequestName)
		log.Println("---metricFromRequestType:", metricFromRequestType)
		log.Println("---metricFromRequestValue:", metricFromRequestValue)

		regex := regexp.MustCompile(patternConst)

		if !regex.MatchString(metricFromRequestValue) {
			writeError := "Value of metric not correct"
			log.Println(writeError)
			http.Error(write, writeError, http.StatusBadRequest)

			return
		}

		var metricValue interface{}

		switch metricFromRequestType {
		case "gauge":
			metricValue = metrics.Gauge(0)
			if s, err := strconv.ParseFloat(metricFromRequestValue, 64); err == nil {
				metricValue = metrics.Gauge(s)
			}
		case metricTypeCounterConst:
			metricValue = metrics.Counter(0)
			if s, err := strconv.Atoi(metricFromRequestValue); err == nil {
				metricValue = metrics.Counter(s)
				_, errIncrease := metricRepository.IncreaseValue(metricFromRequestName, metrics.Counter(s))
				if errIncrease != nil {
					log.Println(errIncrease)
					http.Error(write, errIncrease.Error(), http.StatusBadRequest)

					return
				}
			}
		default:
			writeError := "Type not found"
			log.Println(writeError)
			http.Error(write, writeError, http.StatusNotImplemented)

			return
		}

		if metricFromRequestType != metricTypeCounterConst {
			metricRepository.InsertValue(metricFromRequestName, metricValue)
		}

		recorderValue, err := metricRepository.GetValue(metricFromRequestName)

		if err != nil {
			log.Println("---metricRepositoryValue-Error:", err)
		} else {
			log.Println("---metricRepositoryValue:", recorderValue)
		}

		pollCount, err := metricRepository.GetValue("PollCount")

		resStr := "OK: Value `" + fmt.Sprint(recorderValue) + "` in field `"
		resStr += metricFromRequestName + "` was recorded\n | PollCount: " + fmt.Sprint(pollCount, err)

		write.WriteHeader(http.StatusOK)

		_, writeError := write.Write([]byte(resStr))
		if writeError != nil {
			log.Println(writeError)
		}
	}
}
