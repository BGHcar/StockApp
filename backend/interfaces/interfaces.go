package interfaces

import (
	"backend/models"
	"time"

	"gorm.io/gorm"
)

// DatabaseHandler define la interfaz para operaciones de base de datos con GORM
type DatabaseHandler interface {
	DB() *gorm.DB
}

// APIClient define la interfaz para consumir APIs externas
type APIClient interface {
	Get(endpoint string, nextPage string) (*APIResponse, error)
}

// APIResponse define la estructura de respuesta de una API
type APIResponse struct {
	Items    []StockItem
	NextPage string
}

// StockItem define la estructura de un elemento de stock devuelto por la API
type StockItem struct {
	Ticker     string
	Company    string
	TargetFrom string
	TargetTo   string
	Action     string
	Brokerage  string
	RatingFrom string
	RatingTo   string
	Time       time.Time
}

// Repository define la interfaz para el acceso a datos
type StockRepository interface {
	GetAll() ([]models.Stock, error)
	GetCount() (int64, error)
	GetByTicker(ticker string) ([]models.Stock, error) // Ya devuelve un slice
	GetByAction(action string) ([]models.Stock, error)
	GetByRating(rating string) ([]models.Stock, error)
	GetByBrokerage(brokerage string) ([]models.Stock, error)
	GetByDateRange(startDate, endDate time.Time) ([]models.Stock, error)
	GetByCompany(company string) ([]models.Stock, error)
	GetActionCounts() (map[string]int, error)
	GetRatingCounts() (map[string]int, error)
	InsertStock(stock models.Stock) error
	InsertStocks(stocks []models.Stock) (int, map[string]string, error)
	InsertStocksParallel(stocks []models.Stock) (int, map[string]string, error)
	TruncateTable() error
	SearchStocks(query string) ([]models.Stock, error)
	GetByPriceRange(minPrice, maxPrice string) ([]models.Stock, error) // Nueva función
}

// Service define la interfaz para la lógica de negocio
type StockService interface {
	GetAllStocks() ([]models.Stock, error)
	GetTotalCount() (int, error)
	GetStockByTicker(ticker string) ([]models.Stock, error)
	GetStocksByAction(action string) ([]models.Stock, error)
	GetStocksByRating(rating string) ([]models.Stock, error)
	GetStocksByCompany(company string) ([]models.Stock, error)
	GetActionStats() (map[string]int, error)
	GetRatingStats() (map[string]int, error)
	SyncStockData() (SyncResult, error)
	SearchStocks(query string) ([]models.Stock, error)
	GetStocksByBrokerage(brokerage string) ([]models.Stock, error)
	GetStocksByDateRange(startDate, endDate time.Time) ([]models.Stock, error)
	GetStocksByPriceRange(minPrice, maxPrice string) ([]models.Stock, error) // Nueva función
}

// SyncResult define el resultado de un proceso de sincronización
type SyncResult struct {
	TotalProcessed      int
	TotalInserted       int
	FailedInserts       int
	UniqueTickersAPI    int
	UniqueTickersDB     int
	DuplicateTickers    int
	DuplicatesList      []TickerDuplicate
	FailedInsertDetails map[string]string
}

// TickerDuplicate representa información sobre tickers duplicados
type TickerDuplicate struct {
	Ticker string
	Count  int
}
