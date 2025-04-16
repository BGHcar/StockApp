// File: backend/main.go
// Description: This is the main entry point for the backend application.

package main

import (
	"backend/api"
	"backend/config"
	"backend/db"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Cargar configuración
	cfg := config.GetConfig()

	// Configurar logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Iniciando aplicación Stock Analyzer")

	// Initialize the router
	router := chi.NewRouter()

	// Middleware básico
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Middleware CORS
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	// Initialize the database connection
	db.InitDB()
	// Coloca el defer aquí, justo después de inicializar la BD
	defer db.DB.Close()
	log.Println("Database connection will be closed when the application exits")

	// Set up routes for the API
	api.SetupRoutes(router)

	// Crear servidor con timeouts
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Iniciar servidor en goroutine separada para manejar graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Escuchar señales para graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Recibida señal de apagado, cerrando servidor...")

		// Crear contexto con timeout para shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel() // Ensure the context is canceled to prevent leaks

		// Cerrar servidor
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Error durante el apagado del servidor: %v", err)
		}

		serverStopCtx()
	}()

	// Start the HTTP server
	log.Printf("Starting server on port %s...", cfg.Server.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

	// Esperar a que el servidor se cierre completamente
	<-serverCtx.Done()
	log.Println("Servidor apagado correctamente")
}
