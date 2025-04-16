// File: backend/main.go
// Description: This is the main entry point for the backend application.

package main

import (
	"backend/api"
	"backend/db"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Carga las variables de entorno del archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env. Usando variables del entorno.")
	}

	// Initialize the router
	router := chi.NewRouter()

	// Middleware básico
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Initialize the database connection
	db.InitDB()
	// Coloca el defer aquí, justo después de inicializar la BD
	defer db.DB.Close()
	log.Println("Database connection will be closed when the application exits")

	// Set up routes for the API
	api.SetupRoutes(router)

	// Start the HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
