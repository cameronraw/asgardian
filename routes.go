package asgardian

import (
	"net/http"
)

type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

type RouteConfig struct {
	url        string
	methods    []HTTPMethod
	handler    http.HandlerFunc
	middleware []Middleware
}

func DefineRoute(url string, methods []HTTPMethod, handler http.HandlerFunc) RouteConfig {
	return RouteConfig{
		url:        url,
		methods:    methods,
		handler:    handler,
		middleware: []Middleware{},
	}
}

func (routeCfg RouteConfig) createMiddlewarePolicy(middleware []Middleware) RouteConfig {
	routeCfg.middleware = append(routeCfg.middleware, middleware...)
	return routeCfg
}

func SetRoute(routeCfg RouteConfig) http.HandlerFunc {
	return routeCfg.applyMiddleware()
}

func (app *Application) RoutesConfig() []RouteConfig {
  securityMiddleware := CreateSecurityMiddleware(CreateApiKeySecurityStrategy(app.config.Key))
	return collect(
		[][]RouteConfig{
			{
				DefineRoute("/api/v1/health", []HTTPMethod{GET}, app.healthCheck).
          createMiddlewarePolicy([]Middleware{
            &securityMiddleware,
          }),
			},
			DefineRouteGroup("/just/a/test", []HTTPMethod{GET, POST}, func() []RouteConfig {
				return []RouteConfig{
					DefineRoute("/key", []HTTPMethod{GET}, app.healthCheck).
            createMiddlewarePolicy([]Middleware{
              &securityMiddleware,
            }),
				}
			}),
			DefineRouteGroup("/api/v1", []HTTPMethod{GET, POST}, func() []RouteConfig {
				return []RouteConfig{
					DefineRoute("/health", []HTTPMethod{GET}, app.healthCheck),
				}
			}),
		},
	)
}

func collect(routes [][]RouteConfig) []RouteConfig {
	combinedRoutes := []RouteConfig{}
	for _, routeGroup := range routes {
		combinedRoutes = append(combinedRoutes, routeGroup...)
	}
	return combinedRoutes
}

type RouteGroup func() []RouteConfig

func DefineRouteGroup(url string, methods []HTTPMethod, routeGroup RouteGroup) []RouteConfig {
	amendedRouteGroup := []RouteConfig{}
	for _, routeCfg := range routeGroup() {
		routeCfg.url = url + routeCfg.url
		routeCfg.methods = methods
		amendedRouteGroup = append(amendedRouteGroup, routeCfg)
	}
	return amendedRouteGroup
}

func (routeCfg *RouteConfig) applyMiddleware() http.HandlerFunc {
	if len(routeCfg.middleware) == 0 {
		return routeCfg.handler
	}

	whatToWrap := routeCfg.middleware[len(routeCfg.middleware)-1].Wrap(routeCfg.handler)

	if len(routeCfg.middleware) == 1 {
		return whatToWrap
	}

	for i := len(routeCfg.middleware) - 2; i >= 0; i-- {
		whatToWrap = routeCfg.middleware[i].Wrap(whatToWrap)
	}

	return whatToWrap
}
