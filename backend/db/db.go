// File: db.go
// File: backend/db/db.go

package db

import (
	"backend/interfaces"
	"database/sql"
	"log"
	"os"

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

// InitDB inicializa la conexión a la base de datos
func InitDB() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	log.Println("Connected to database successfully")

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
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	log.Println("Database connection closed")
}
