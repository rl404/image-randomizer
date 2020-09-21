package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/internal"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// startHTTP to start serving HTTP.
func startHTTP(cfg internal.Config) error {
	r := chi.NewRouter()

	// Set default recommended go-chi router middlewares.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.AllowAll().Handler)

	// Register base routes.
	internal.RegisterBaseRoutes(r)

	// Prepare user routes.
	userRoute, err := internal.GetUserRoutes(cfg)
	if err != nil {
		return err
	}

	// Register user routes.
	r.Mount("/user", userRoute)

	fmt.Println("server listen at " + cfg.Port)
	return http.ListenAndServe(cfg.Port, r)
}

func main() {
	err := startHTTP(internal.GetConfig())
	if err != nil {
		log.Fatal("error starting HTTP", " - ", err.Error())
		return
	}
}
