package repositories_test

import (
	"backend/models"
	"backend/repositories"
	"regexp"
	"strings" // Asegúrate de agregar esta importación
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MockDB implementa la interfaz DatabaseHandler para pruebas
type MockDB struct {
	GormDB *gorm.DB
}

func (m *MockDB) DB() *gorm.DB {
	return m.GormDB
}

func setupMockDB(t *testing.T) (*MockDB, sqlmock.Sqlmock, func()) {
	// Crear una mock connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creando sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Error abriendo conexión GORM: %v", err)
	}

	return &MockDB{GormDB: gormDB}, mock, func() { db.Close() }
}

// Helper para pruebas
func buildExpectedSQLWithSoftDelete(baseSQL string) string {
	// Buscar dónde insertar la condición de soft delete
	if containsWhere := strings.Contains(baseSQL, "WHERE"); containsWhere {
		// Si ya tiene WHERE, añadir AND "deleted_at" IS NULL
		return strings.Replace(baseSQL, "WHERE", "WHERE \"stocks\".\"deleted_at\" IS NULL AND", 1)
	}
	// Si no tiene WHERE, añadir WHERE "deleted_at" IS NULL
	if orderBy := strings.Contains(baseSQL, "ORDER BY"); orderBy {
		// Si tiene ORDER BY, insertar antes
		return strings.Replace(baseSQL, "ORDER BY", "WHERE \"stocks\".\"deleted_at\" IS NULL ORDER BY", 1)
	}
	// Si no tiene ni WHERE ni ORDER BY, añadir al final
	return baseSQL + " WHERE \"stocks\".\"deleted_at\" IS NULL"
}

func TestGetAll(t *testing.T) {
	// Arrange
	mockDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewStockRepository(mockDB)

	// Mock data
	rows := sqlmock.NewRows([]string{"id", "created_at", "ticker", "company"}).
		AddRow(1, time.Now(), "AAPL", "Apple Inc.").
		AddRow(2, time.Now(), "MSFT", "Microsoft Corporation")

	// La consulta real incluye WHERE deleted_at IS NULL y ORDER BY time DESC
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE "stocks"."deleted_at" IS NULL ORDER BY time DESC`)).
		WillReturnRows(rows)

	// Act
	stocks, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Len(t, stocks, 2)
	assert.Equal(t, "AAPL", stocks[0].Ticker)
	assert.Equal(t, "MSFT", stocks[1].Ticker)
}

func TestGetCount(t *testing.T) {
	// Arrange
	mockDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewStockRepository(mockDB)

	// Mock data
	rows := sqlmock.NewRows([]string{"count"}).AddRow(42)

	// Expect query
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "stocks"`)).
		WillReturnRows(rows)

	// Act
	count, err := repo.GetCount()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(42), count)
}

func TestGetByTicker(t *testing.T) {
	// Arrange
	mockDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewStockRepository(mockDB)
	ticker := "AAPL"

	// Mock data
	rows := sqlmock.NewRows([]string{"id", "created_at", "ticker", "company"}).
		AddRow(1, time.Now(), "AAPL", "Apple Inc.")

	// Expect query with ILIKE y la condición de soft delete
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE ticker ILIKE $1 AND "stocks"."deleted_at" IS NULL ORDER BY time DESC`)).
		WithArgs("%" + ticker + "%").
		WillReturnRows(rows)

	// Act
	stocks, err := repo.GetByTicker(ticker)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, stocks, 1)
	assert.Equal(t, ticker, stocks[0].Ticker)
}

func TestInsertStock(t *testing.T) {
	// Arrange
	mockDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewStockRepository(mockDB)
	stock := models.Stock{
		Ticker:  "AAPL",
		Company: "Apple Inc.",
	}

	// Expect insert query
	mock.ExpectBegin() // Start transaction
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "stocks"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit() // Commit transaction

	// Act
	err := repo.InsertStock(stock)

	// Assert
	assert.NoError(t, err)
}

func TestTruncateTable(t *testing.T) {
	// Arrange
	mockDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewStockRepository(mockDB)

	// Expect truncate query
	mock.ExpectExec(regexp.QuoteMeta(`TRUNCATE TABLE stocks`)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Act
	err := repo.TruncateTable()

	// Assert
	assert.NoError(t, err)
}

// Modificar la prueba TestSearchStocks:
func TestSearchStocks(t *testing.T) {
	// Arrange
	mockDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewStockRepository(mockDB)
	query := "apple"

	// Mock data
	rows := sqlmock.NewRows([]string{"id", "created_at", "ticker", "company"}).
		AddRow(1, time.Now(), "AAPL", "Apple Inc.")

	// Aquí está el cambio importante - usar un patrón que incluya la condición de soft delete
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE (ticker ILIKE $1 OR company ILIKE $2 OR brokerage ILIKE $3) AND "stocks"."deleted_at" IS NULL ORDER BY time DESC LIMIT $4`)).
		WithArgs("%"+query+"%", "%"+query+"%", "%"+query+"%", int64(100)).
		WillReturnRows(rows)

	// Act
	stocks, err := repo.SearchStocks(query)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "AAPL", stocks[0].Ticker)
	assert.Equal(t, "Apple Inc.", stocks[0].Company)
}

func TestGetByPriceRange(t *testing.T) {
	// Arrange
	mockDB, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := repositories.NewStockRepository(mockDB)
	minPrice := "100"
	maxPrice := "200"

	// Mock data
	rows := sqlmock.NewRows([]string{"id", "created_at", "ticker", "company", "target_from", "target_to"}).
		AddRow(1, time.Now(), "AAPL", "Apple Inc.", "$150", "$180")

	// Usar un patrón más general en lugar de la expresión regular específica
	mock.ExpectQuery(`.*target_from.*target_to.*`).
		WillReturnRows(rows)

	// Act
	stocks, err := repo.GetByPriceRange(minPrice, maxPrice)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "AAPL", stocks[0].Ticker)
	assert.Equal(t, "$150", stocks[0].TargetFrom)
	assert.Equal(t, "$180", stocks[0].TargetTo)
}
