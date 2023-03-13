package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Group(func(mux chi.Router) {
		// Public routes
		mux.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
		mux.Use(middleware.Heartbeat("/ping"))

		mux.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("doc.json"), // The url pointing to API definition
		))

		mux.Get("/v1/projects", app.BaseHandler.ProjectRead)
		mux.Post("/v1/projects", app.BaseHandler.ProjectCreate)
		mux.Patch("/v1/projects", app.BaseHandler.ProjectUpdate)
		mux.Delete("/v1/projects", app.BaseHandler.ProjectDelete)
	})

	return mux
}
