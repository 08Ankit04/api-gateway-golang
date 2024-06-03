package logger

import (
	"log"
	"net/http"
	"time"
)

// Middleware logs the details of each incoming request
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		// Capture the response status
		rec := statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(&rec, r)

		duration := time.Since(startTime)
		log.Printf("Completed %s %s %d %v", r.Method, r.URL.Path, rec.status, duration)
	})
}

// statusRecorder is a custom ResponseWriter to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}
