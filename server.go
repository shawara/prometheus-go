package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Number of get requests.",
},
	[]string{"path", "code"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests."},
	[]string{"path", "code"},
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := NewResponseWriter(w)

		start := time.Now()
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		statusCodeStr := strconv.Itoa(rw.statusCode)

		totalRequests.WithLabelValues(path, statusCodeStr).Inc()
		httpDuration.WithLabelValues(path, statusCodeStr).Observe(duration.Seconds())

	})
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(httpDuration)
}

func main() {
	router := mux.NewRouter()
	router.Use(prometheusMiddleware)

	// Prometheus endpoint
	router.Path("/metrics").Handler(promhttp.Handler())

	// Serving / route
	router.PathPrefix("/").Handler(http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("Hello From Go server"))
		}))

	fmt.Println("Serving requests on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}
