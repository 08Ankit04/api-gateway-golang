package router

import (
	"fmt"
	"io"
	"net/http"

	"github.com/api-gateway-golang/internal/auth"
	"github.com/api-gateway-golang/internal/logger"
	"github.com/api-gateway-golang/internal/model"
	"github.com/api-gateway-golang/internal/rate_limit"
	"github.com/gorilla/mux"
)

const (
	errInternalServerError = "Internal server error"
	errServiceUnavailable  = "Service unavailable error"
)

// InitializeRouter sets up the router with the given routes and applies middleware
func InitializeRouter(routes []model.Route) *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes {
		handler := logger.Middleware(auth.Middleware(rate_limit.Middleware(proxy(route.Service, route.ServicePort))))
		router.Handle(route.Path, handler).Methods("GET", "POST", "PUT", "DELETE")
	}

	return router
}

// proxy function to forward requests to the appropriate microservice
func proxy(service string, port int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := "http://" + service + ":" + fmt.Sprint(port) + r.RequestURI
		req, err := http.NewRequest(r.Method, url, r.Body)
		if err != nil {
			http.Error(w, errInternalServerError, http.StatusInternalServerError)
			return
		}

		req.Header = r.Header
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, errServiceUnavailable, http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		respBody, _ := io.ReadAll(resp.Body)

		w.WriteHeader(resp.StatusCode)
		_, err = w.Write(respBody)
		if err != nil {
			http.Error(w, errInternalServerError, http.StatusInternalServerError)
			return
		}
	})
}
