package db

import (
	"backend/interfaces"
	"backend/models"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB es la conexión global a la base de datos
var DB *gorm.DB  // el *gorm.DB es el objeto de conexión a la base de datos, * es para que sea un puntero a la conexión


// DBAdapter adapta gorm.DB a la interfaz DatabaseHandler
type DBAdapter struct {
	db *gorm.DB
}

// Verificamos que DBAdapter implementa interfaces.DatabaseHandler
var _ interfaces.DatabaseHandler = (*DBAdapter)(nil)

// DB devuelve el objeto *gorm.DB subyacente
func (d *DBAdapter) DB() *gorm.DB {
	return d.db
}

// NewDatabaseHandler crea un nuevo handler de base de datos
func NewDatabaseHandler() (*DBAdapter, error) {
	return &DBAdapter{db: DB}, nil
}

// InitDB inicializa la conexión a la base de datos con GORM
func InitDB() {
	var err error

	// Obtenemos la configuración de conexión
	connStr := getDatabaseConfig()
	if connStr == "" {
		log.Fatal("No se pudo obtener una cadena de conexión válida. Verifique las variables de entorno.")
	}

	// Configurar el logger de GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 200 * time.Second,
			LogLevel:      logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:      true,
		},
	)

	// Abrimos la conexión con GORM
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC() // Usar UTC para timestamps
		},
	})

	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	}

	log.Println("Conexión a base de datos establecida correctamente")

	// Configuración del pool de conexiones
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Error al obtener conexión SQL subyacente:", err)
	}

	// Configuración del pool desde variables de entorno
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 10)
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetimeMin := getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 5)
	connMaxIdleTimeMin := getEnvAsInt("DB_CONN_MAX_IDLE_TIME_MIN", 3)

	// Aplicar configuración
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetimeMin) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(connMaxIdleTimeMin) * time.Minute)

	// Auto-migración para crear/actualizar la estructura de la tabla
	err = DB.AutoMigrate(&models.Stock{})
	if err != nil {
		log.Fatal("Error durante la migración automática:", err)
	}

	log.Println("Estructura de datos verificada correctamente mediante migración")
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
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

// Helper para obtener variables de entorno como enteros
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
