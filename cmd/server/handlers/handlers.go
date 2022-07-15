package handlers

import (
	"errors"
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
		write.Header().Set("content-type", "text/plain; charset=utf-8")

		metricName := chi.URLParam(request, "metricName")
		metricType := chi.URLParam(request, "metricType")

		resVal, err := metricRepository.GetValue(metricName)
		if err != nil {
			http.Error(write, err.Error(), http.StatusNotFound)
			return
		}

		valueMetricType, errGetType := helper.GetType(resVal)
		if errGetType != nil {
			http.Error(write, errGetType.Error(), http.StatusNotFound)
			return
		}

		if valueMetricType != metricType {
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
		if request.Method != http.MethodPost {
			http.Error(write, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		write.Header().Set("content-type", "text/plain; charset=utf-8")

		log.Println(request.Method, request.Host, request.URL.Path)

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
			http.Error(write, writeError, http.StatusBadRequest)
			log.Println(writeError)

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
			}
		default:
			writeError := "Type not found"
			http.Error(write, writeError, http.StatusNotImplemented)
			log.Println(writeError)

			return
		}

		if metricFromRequestType == metricTypeCounterConst {
			metricValueNew, ok := metricValue.(metrics.Counter)
			if !ok {
				err := errors.New("type change error")
				http.Error(write, err.Error(), http.StatusBadRequest)
				log.Println(err)
			}

			_, err := metricRepository.IncreaseValue("PollCount", metricValueNew)
			if err != nil {
				http.Error(write, err.Error(), http.StatusBadRequest)
				log.Println(err)
			}
		} else {
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
