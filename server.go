package asgardian

import (
	"log"
	"net/http"
	"os"
)

type AsgardianServer struct {
	server *http.Server
}

type AsgardianServerBuilder struct{
  app Application
  mux *http.ServeMux
}

type ApplicationSettings struct{
  config *Config
  logger *log.Logger
}

func (asb *AsgardianServerBuilder) CreateApplication() *AsgardianServerBuilder {
  app := CreateApplication(
    CreateConfig(), 
    log.New(os.Stdout, "", log.Ldate|log.Ltime),
    )

  asb.app = *app
  return asb
}

func (asb *AsgardianServerBuilder) CreateApplicationWith(settingsFunc func() ApplicationSettings) *AsgardianServerBuilder {
  settings := settingsFunc()
  app := CreateApplication(
    *settings.config, 
    settings.logger,
    )

  asb.app = *app
  return asb
}

func (asb *AsgardianServerBuilder) CreateRoutes(routeConfig func() [][]RouteConfig ) *AsgardianServerBuilder {

  mux := http.NewServeMux()

  combinedRoutes := []RouteConfig{}
  for _, routes := range routeConfig() {
    combinedRoutes = append(combinedRoutes, routes...)
  }

  for _, route := range combinedRoutes {
    mux.HandleFunc(route.url, MapRoute(route.methods, route.handler))
  }

  asb.mux = mux

  return asb
}
