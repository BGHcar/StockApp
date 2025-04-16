// File: db.go
// File: backend/db/db.go

package db

import (
	"backend/interfaces"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
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
	if dbURL != "" {
		return dbURL
	}

	// Si no, construimos la URL con componentes individuales
	user := os.Getenv("SQL_USER")
	password := os.Getenv("GENERATED_PASSWORD")
	host := os.Getenv("CLUSTER_HOST")
	port := os.Getenv("CLUSTER_PORT")
	dbName := os.Getenv("CLUSTER_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")

	// Verificar que todos los parámetros necesarios están presentes
	if user == "" || password == "" || host == "" {
		log.Println("ADVERTENCIA: Configuración de base de datos incompleta en variables de entorno")
		return ""
	}

	// Valores por defecto solo si no están definidos en variables de entorno
	if port == "" {
		port = os.Getenv("DB_DEFAULT_PORT")
	}

	if dbName == "" {
		dbName = os.Getenv("DB_DEFAULT_NAME")
	}

	if sslMode == "" {
		sslMode = os.Getenv("DB_DEFAULT_SSL_MODE")
	}

	// Si falta algún componente crítico después de intentar usar valores por defecto
	if port == "" || dbName == "" || sslMode == "" {
		log.Println("ADVERTENCIA: Faltan valores críticos para la conexión a base de datos")
		return ""
	}

	// Construir URL de conexión
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)
}

// InitDB inicializa la conexión a la base de datos
func InitDB() {
	var err error

	// Obtenemos la configuración de conexión
	connStr := getDatabaseConfig()
	if connStr == "" {
		log.Fatal("No se pudo obtener una cadena de conexión válida. Verifique las variables de entorno.")
	}

	// No mostramos ninguna información sobre la cadena de conexión
	log.Println("Iniciando conexión a base de datos...")

	// Abrimos la conexión
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error al conectar con la base de datos. Verifique las credenciales y la conectividad.")
	}

	// Configuración del pool desde variables de entorno
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 10)
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetimeMin := getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 5)
	connMaxIdleTimeMin := getEnvAsInt("DB_CONN_MAX_IDLE_TIME_MIN", 3)

	// Aplicar configuración
	DB.SetMaxOpenConns(maxOpenConns)
	DB.SetMaxIdleConns(maxIdleConns)
	DB.SetConnMaxLifetime(time.Duration(connMaxLifetimeMin) * time.Minute)
	DB.SetConnMaxIdleTime(time.Duration(connMaxIdleTimeMin) * time.Minute)

	// Verificamos la conexión
	if err = DB.Ping(); err != nil {
		log.Fatal("No se pudo establecer comunicación con la base de datos.")
	}

	log.Println("Conexión a base de datos establecida correctamente")

	// Crear la tabla si no existe, usando una consulta parametrizada desde un archivo o variable de entorno
	createTableSQL := os.Getenv("DB_CREATE_TABLE_SQL")
	if createTableSQL == "" {
		// Si no se proporciona, usar una consulta básica
		createTableSQL = `
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
        )`
	}

	if _, err = DB.Exec(createTableSQL); err != nil {
		log.Printf("Error al preparar estructura de datos: %v", err)
	} else {
		log.Println("Estructura de datos verificada correctamente")
	}
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

// Helper para obtener variables de entorno como enteros
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
