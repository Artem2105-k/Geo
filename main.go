package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/internal/auth"
	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/internal/controller"
	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/internal/service/rpcclient"
)

// @title Geo Service API
// @version 1.0
// @description API для поиска адресов и геокодирования
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	r := chi.NewRouter()

	jwtSecret := os.Getenv("JWT_SECRET")
	authHandler := auth.NewAuthHandler(jwtSecret)

	// Подключаемся к RPC-сервису
	rpcClient, err := rpcclient.NewRPCClient("rpcserver:8081")
	if err != nil {
		log.Fatalf("Failed to connect to RPC server: %v", err)
	}

	geoController := controller.NewGeoController(rpcClient)

	// Роуты
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(authHandler.TokenAuth))
		r.Use(authHandler.Authenticator)
		r.Post("/api/address/search", geoController.SearchAddress)
		r.Post("/api/address/geocode", geoController.GeoCode)
	})

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	log.Println("Starting API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
