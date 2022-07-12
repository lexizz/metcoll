package handlers

import (
	"errors"
	"fmt"
	"github.com/lexizz/metcoll/internal/metrics"
	"github.com/lexizz/metcoll/internal/repository/metricMemoryRepository"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func UpdateMetric() http.HandlerFunc {
	metricRepository := metricMemoryRepository.New()

	return func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(response, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		response.Header().Set("content-type", "text/plain")

		log.Println(request.Method, request.Host, request.URL.Path)

		url := request.URL.Path
		url = strings.Trim(url, "/")
		urlData := strings.Split(url, "/")

		log.Println("---urlData:", urlData)

		if len(urlData) < 4 {
			writeError := "Params not found"
			http.Error(response, writeError, http.StatusNotFound)
			log.Println(writeError)

			return
		}

		metricFromRequestName := urlData[2]
		metricFromRequestType := urlData[1]
		metricFromRequestValue := urlData[3]
		log.Println("---metricFromRequestName:", metricFromRequestName)
		log.Println("---metricFromRequestType:", metricFromRequestType)
		log.Println("---metricFromRequestValue:", metricFromRequestValue)

		reg, err := regexp.Compile(`^[0-9\.]+$`)

		if !reg.MatchString(metricFromRequestValue) || err != nil {
			writeError := "Value of metric not correct"
			http.Error(response, writeError, http.StatusBadRequest)
			log.Println(writeError)

			return
		}

		var metricValue interface{}

		if metricFromRequestType == "gauge" {
			metricValue = metrics.Gauge(0)
			if s, err := strconv.ParseFloat(metricFromRequestValue, 64); err == nil {
				metricValue = metrics.Gauge(s)
			}
		} else if metricFromRequestType == "counter" {
			metricValue = metrics.Counter(0)
			if s, err := strconv.Atoi(metricFromRequestValue); err == nil {
				metricValue = metrics.Counter(s)
			}
		} else {
			writeError := "Type not found"
			http.Error(response, writeError, http.StatusNotImplemented)
			log.Println(writeError)

			return
		}

		if metricFromRequestType == "counter" {
			metricValueNew, ok := metricValue.(metrics.Counter)
			if !ok {
				err := errors.New("type change error")
				http.Error(response, err.Error(), http.StatusBadRequest)
				log.Println(err)
			}

			_, err := metricRepository.IncreaseValue("PollCount", metricValueNew)
			if err != nil {
				http.Error(response, err.Error(), http.StatusBadRequest)
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

		resStr := "OK: Value `" + fmt.Sprint(recorderValue) + "` in field `" + metricFromRequestName + "` was recorded\n | PollCount: " + fmt.Sprint(pollCount, err)

		response.WriteHeader(http.StatusOK)

		_, writeError := response.Write([]byte(resStr))
		if writeError != nil {
			log.Println(writeError)
		}
	}
}
