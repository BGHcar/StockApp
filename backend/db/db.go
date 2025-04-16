// File: db.go
// File: backend/db/db.go

package db

import (
	"backend/interfaces"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB es la conexión global a la base de datos
var DB *sql.DB

// DBAdapter adapta sql.DB a la interfaz DatabaseHandler
type DBAdapter struct {
	DB *sql.DB
}

// Verificamos que DBAdapter implementa interfaces.DatabaseHandler
var _ interfaces.DatabaseHandler = (*DBAdapter)(nil)

// Close cierra la conexión a la base de datos
func (d *DBAdapter) Close() error {
	return d.DB.Close()
}

// Exec ejecuta una consulta que no devuelve filas
func (d *DBAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

// Query ejecuta una consulta que devuelve filas
func (d *DBAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

// QueryRow ejecuta una consulta que devuelve una sola fila
func (d *DBAdapter) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}

// NewDatabaseHandler crea un nuevo handler de base de datos
func NewDatabaseHandler() (*DBAdapter, error) {
	return &DBAdapter{DB: DB}, nil
}

// getDatabaseConfig obtiene la configuración de la base de datos desde variables de entorno
func getDatabaseConfig() string {
	// Primero verificamos si hay una URL completa de base de datos
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" && !strings.Contains(dbURL, "localhost") {
		log.Println("Usando DATABASE_URL completa para conexión")
		return dbURL
	}

	// Si no, construimos la URL con componentes individuales
	user := os.Getenv("SQL_USER")
	password := os.Getenv("GENERATED_PASSWORD")
	host := os.Getenv("CLUSTER_HOST")
	port := os.Getenv("CLUSTER_PORT")
	dbName := os.Getenv("CLUSTER_NAME")

	if user != "" && password != "" && host != "" {
		if port == "" {
			port = "26257" // Puerto por defecto de CockroachDB
		}
		if dbName == "" {
			dbName = "stock_data"
		}

		// Construir URL de conexión
		connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require",
			user, password, host, port, dbName)
		log.Println("URL de conexión a cluster construida desde variables de entorno")
		return connStr
	}

	// Si no hay configuración para cluster, usamos conexión local
	log.Println("ADVERTENCIA: No se encontró configuración válida para cluster, usando conexión local")
	return "postgresql://root@localhost:26257/stock_data?sslmode=disable"
}

// InitDB inicializa la conexión a la base de datos
func InitDB() {
	var err error

	// Obtenemos la configuración de conexión
	connStr := getDatabaseConfig()

	log.Printf("Conectando a base de datos: %s", sanitizeConnectionString(connStr))

	// Abrimos la conexión
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Configuración optimizada del pool para conexiones remotas
	// Aumentamos el número de conexiones para soportar operaciones paralelas
	if strings.Contains(connStr, "cockroachlabs.cloud") {
		// Configuración para cluster remoto
		DB.SetMaxOpenConns(25) // Más conexiones para paralelismo
		DB.SetMaxIdleConns(10) // Mantener más conexiones inactivas
		DB.SetConnMaxLifetime(10 * time.Minute)
		DB.SetConnMaxIdleTime(5 * time.Minute)
	} else {
		// Configuración para base de datos local
		DB.SetMaxOpenConns(10)
		DB.SetMaxIdleConns(5)
		DB.SetConnMaxLifetime(5 * time.Minute)
		DB.SetConnMaxIdleTime(1 * time.Minute)
	}

	// Verificamos la conexión
	if err = DB.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	log.Println("Conectado a la base de datos exitosamente")

	// Crear la tabla si no existe
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS stocks (
			ticker TEXT NOT NULL,
			company TEXT,
			target_from TEXT,
			target_to TEXT,
			action TEXT,
			brokerage TEXT,
			rating_from TEXT,
			rating_to TEXT,
			time TIMESTAMP NOT NULL,
			PRIMARY KEY (ticker, time)
		)
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Println("Estructura de tabla verificada/creada correctamente")
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	log.Println("Database connection closed")
}

// sanitizeConnectionString oculta completamente la información sensible
func sanitizeConnectionString(connStr string) string {
	if connStr == "" {
		return "[cadena de conexión vacía]"
	}

	// Detectar el tipo de conexión sin mostrar detalles
	if strings.Contains(connStr, "cockroachlabs.cloud") {
		return "[conexión segura a CockroachDB Cloud]"
	} else if strings.Contains(connStr, "localhost") || strings.Contains(connStr, "127.0.0.1") {
		return "[conexión local a CockroachDB]"
	}

	// Para cualquier otro caso
	return "[conexión de base de datos configurada]"
}
