package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Sahil2k07/Blog-App-Go/src/config"
	"github.com/Sahil2k07/Blog-App-Go/src/routes"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	origins := strings.Split(allowedOrigins, ",")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins(origins),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	config.InitBloomFilter(1000000, 0.01)

	// Database
	config.DBConnect()
	defer config.DBDisconnect()

	// Routes
	router := routes.AppRoutes()

	log.Printf("Server is running on %s", PORT)

	err := http.ListenAndServe(PORT, corsHandler(router))
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}

}
