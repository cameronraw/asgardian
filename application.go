package asgardian

import (
	"fmt"
	"log"
	"net/http"
)

type Application struct {
	config     Config
	logger     *log.Logger
	Middleware []Middleware
}

func CreateApplication(config Config, logger *log.Logger) *Application {
	app := &Application{
		config:     config,
		logger:     logger,
		Middleware: []Middleware{},
	}

	ConfigureMiddleware(app)

	return app
}

func (app *Application) CreateAddr() string {
	return fmt.Sprintf(":%d", app.config.Port)
}

func (app *Application) CreatePreconfiguredMux() *http.ServeMux {
	mux := http.NewServeMux()
  for _, route := range app.RoutesConfig() {
    mux.HandleFunc(route.url, MapRoute(route.methods, app.applyMiddleware(route.handler)))
  }
	return mux
}

func MapRoute(methods []HTTPMethod, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(methods) > 0 && !Contains(methods, HTTPMethod(r.Method)) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func Contains[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}
