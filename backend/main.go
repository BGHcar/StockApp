// File: backend/main.go
// Description: This is the main entry point for the backend application.

package main

import (
	"backend/config"
	"backend/db"
	"backend/routes"
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

	// Configurar cierre de la conexión
	defer func() {
		if err := db.CloseDB(); err != nil {
			log.Printf("Error al cerrar la conexión DB: %v", err)
		}
		log.Println("Database connection closed")
	}()

	// Set up routes for the API
	routes.SetupRoutes(router)

	// Contexto para gestionar el apagado graceful
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Escuchar señales para graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	// Crear el servidor HTTPS (puerto 443)
	httpServer := &http.Server{
		Addr:         "0.0.0.0:8081", // Puerto HTTPS estándar
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Crear el servidor HTTP (puerto 80)
/* 	httpServer := &http.Server{
		Addr: "0.0.0.0:8082", // Puerto HTTP estándar
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Redirigir todo el tráfico HTTP a HTTPS
			target := "https://" + r.Host + r.URL.Path
			if len(r.URL.RawQuery) > 0 {
				target += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	} */

	// Goroutine para manejar el servidor HTTP (puerto 80)
	go func() {
		log.Printf("Iniciando servidor HTTP en puerto 8081 (redirección a HTTPS)...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor HTTP: %v", err)
		}
	}()

/* 	// Goroutine para manejar el servidor HTTPS (puerto 443)
	go func() {
		log.Printf("Iniciando servidor HTTPS en puerto 443...")

		// Para entorno de producción con certificados:
		certFile := os.Getenv("CERT_FILE")
		keyFile := os.Getenv("KEY_FILE")

		var err error
		if certFile != "" && keyFile != "" {
			// Si hay certificados configurados, usar HTTPS
			err = httpsServer.ListenAndServeTLS(certFile, keyFile)
		} else {
			// Si no hay certificados, usar HTTP en puerto 443
			err = httpsServer.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor HTTPS: %v", err)
		}
	}() */

	// Goroutine para manejar señales de apagado
	go func() {
		<-signalChan
		log.Println("Recibida señal de apagado, cerrando servidores...")

		// Crear contexto con timeout para shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownCancel()

		// Cerrar servidor HTTP
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error durante el apagado del servidor HTTP: %v", err)
		}

		// Cerrar servidor HTTPS
/* 		if err := httpsServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Error durante el apagado del servidor HTTPS: %v", err)
		} */

		serverStopCtx()
	}()

	// Esperar a que los servidores se cierren completamente
	<-serverCtx.Done()
	log.Println("Servidores apagados correctamente")
}
