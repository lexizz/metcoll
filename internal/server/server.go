package server

import (
	"log"
	"net/http"
	"time"

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
	init *http.Server
}

func New() *server {
	metricRepository := metricmemoryrepository.New()

	return &server{
		init: &http.Server{
			Addr:              HOST + ":" + PORT,
			Handler:           GetRoutes(metricRepository),
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (httpServer *server) Listen() error {
	return httpServer.init.ListenAndServe()
}

func GetRoutes(metricRepository metricrepository.Interface) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", handlers.ShowPossibleValue(metricRepository))
	router.Post("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	router.Route("/update", func(routerUrlUpdate chi.Router) {
		routerUrlUpdate.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			log.Println("=== Part url was detected `/update` === ")

			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		})

		routerUrlUpdate.Route("/{metricType:[a-zA-Z]+}", func(routerUrlUpdateType chi.Router) {
			routerUrlUpdateType.Post("/", func(writer http.ResponseWriter, request *http.Request) {
				log.Printf("=== Part url was detected `/%v` === ", chi.URLParam(request, "metricType"))

				http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			})

			routerUrlUpdateType.Route("/{metricName:[a-zA-Z0-9]+}", func(routerUrlUpdateTypeName chi.Router) {
				routerUrlUpdateTypeName.Post("/", func(writer http.ResponseWriter, request *http.Request) {
					log.Printf("=== Part url was detected `/%v` === ", chi.URLParam(request, "metricName"))

					http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				})

				routerUrlUpdateTypeName.Route("/{metricValue:[a-zA-Z0-9.]+}", func(routerUrlUpdateTypeNameValue chi.Router) {
					routerUrlUpdateTypeNameValue.Post("/", handlers.UpdateMetric(metricRepository))

					routerUrlUpdateTypeNameValue.Get("/", func(writer http.ResponseWriter, request *http.Request) {
						http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					})
				})
			})
		})
	})

	router.Route("/value", func(routerUrlValue chi.Router) {
		routerUrlValue.Get("/", func(writer http.ResponseWriter, request *http.Request) {
			log.Println("=== Part url was detected `/value` === ")

			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		})

		routerUrlValue.Route("/{metricType:[a-zA-Z]+}", func(routerUrlValueType chi.Router) {
			routerUrlValueType.Get("/", func(writer http.ResponseWriter, request *http.Request) {
				log.Printf("=== Part url was detected `/%v` === ", chi.URLParam(request, "metricType"))

				http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			})

			routerUrlValueType.Route("/{metricName:[a-zA-Z0-9]+}", func(routerUrlValueTypeName chi.Router) {
				routerUrlValueTypeName.Get("/", handlers.ShowValueMetric(metricRepository))
			})
		})
	})

	return router
}
