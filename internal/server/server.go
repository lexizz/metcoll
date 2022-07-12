package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lexizz/metcoll/cmd/server/handlers"
	"github.com/lexizz/metcoll/internal/repository/interfaces/metricrepository"
	"github.com/lexizz/metcoll/internal/repository/metricmemoryrepository"
)

const (
	HOST string = "127.0.0.1"
	PORT string = "8080"
)

type server struct {
	Init *http.Server
}

func New() *server {
	metricRepository := metricmemoryrepository.New()

	return &server{
		Init: &http.Server{
			Addr:    HOST + ":" + PORT,
			Handler: GetRoutes(metricRepository),
		},
	}
}

func GetRoutes(metricRepository metricrepository.Interface) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/*", func(response http.ResponseWriter, request *http.Request) {
		http.Error(response, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})

	router.Post("/*", handlers.UpdateMetric(metricRepository))

	router.Get("/", handlers.ShowPossibleValue(metricRepository))
	router.Get("/value/{metricType:[gauge|counter]+}/{metricName:[a-zA-Z]+}", handlers.ShowValueMetric(metricRepository))

	return router
}
