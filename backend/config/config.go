package config

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Config contiene la configuración general de la aplicación
type Config struct {
	Database struct {
		URL             string
		User            string
		Password        string
		Host            string
		Port            string
		Name            string
		MaxOpenConns    int
		MaxIdleConns    int
		ConnMaxLifetime time.Duration
	}
	API struct {
		URL   string
		Token string
	}
	Server struct {
		Host         string
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
}

var (
	config     *Config
	configOnce sync.Once
)

// GetConfig devuelve la configuración de la aplicación
func GetConfig() *Config {
	configOnce.Do(func() {
		// Intentar cargar .env
		if err := godotenv.Load(); err != nil {
			log.Println("No se pudo cargar el archivo .env. Usando variables del entorno.")
		}

		config = &Config{}

		// Configuración de base de datos
		config.Database.URL = os.Getenv("DATABASE_URL")
		config.Database.User = os.Getenv("SQL_USER")
		config.Database.Password = os.Getenv("GENERATED_PASSWORD")
		config.Database.Host = os.Getenv("CLUSTER_HOST")
		config.Database.Port = os.Getenv("CLUSTER_PORT")
		config.Database.Name = os.Getenv("CLUSTER_NAME")

		// Valores por defecto para la configuración del pool de conexiones
		config.Database.MaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", 20)
		config.Database.MaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
		config.Database.ConnMaxLifetime = getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute)

		// Configuración de la API
		config.API.URL = os.Getenv("API_URL")
		config.API.Token = os.Getenv("API_TOKEN")

		// Configuración del servidor
		config.Server.Host = os.Getenv("HOST")
		if config.Server.Host == "" {
			config.Server.Host = "0.0.0.0" // Por defecto escucha en todas las interfaces
		}
		// Asegúrate de que el puerto no tenga prefijos como "http://"
		config.Server.Port = os.Getenv("PORT")
		if config.Server.Port == "" {
			config.Server.Port = "10000" // Valor predeterminado
		}
		config.Server.ReadTimeout = getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second)
		config.Server.WriteTimeout = getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second)
	})

	return config
}

// Helper para obtener variables de entorno como enteros
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// Helper para obtener variables de entorno como duración
func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	valueStr := os.Getenv(name)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultVal
}
