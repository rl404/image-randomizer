package api

import (
	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/monitoring/newrelic/middleware"
	"github.com/rl404/image-randomizer/internal/service"
	"github.com/rl404/image-randomizer/internal/utils"
)

// API contains all functions for api endpoints.
type API struct {
	service       service.Service
	accessSecret  string
	refreshSecret string
}

// New to create new api endpoints.
func New(service service.Service, accessSecret string, refreshSecret string) *API {
	return &API{
		service:       service,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

// Register to register api routes.
func (api *API) Register(r chi.Router, nrApp *newrelic.Application) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.NewHTTP(nrApp))
		r.Use(log.MiddlewareWithLog(utils.GetLogger(0), log.MiddlewareConfig{Error: true}))
		r.Use(log.MiddlewareWithLog(utils.GetLogger(1), log.MiddlewareConfig{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
			RawPath:        true,
			Error:          true,
		}))
		r.Use(utils.Recoverer)

		r.Post("/register", api.handleRegister)
		r.Post("/login", api.handleLogin)

		r.Post("/token/check", api.jwtAuth(api.handleTokenCheck))
		r.Post("/token/refresh", api.jwtAuth(api.handleTokenRefresh, tokenRefresh))

		r.Get("/images", api.jwtAuth(api.handleGetImages))
		r.Post("/images", api.jwtAuth(api.handleCreateImage))
		r.Patch("/images/{image_id}", api.jwtAuth(api.handleUpdateImage))
		r.Delete("/images/{image_id}", api.jwtAuth(api.handleDeleteImage))

		r.Get("/user/{username}/image.jpg", api.handleRandomImage)
	})
}
