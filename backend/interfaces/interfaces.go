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

// StockRepository define la interfaz para el acceso a datos
type StockRepository interface {
	// --- Métodos Paginados ---
	GetAll(page, pageSize int) ([]models.Stock, int, int, error) // Devuelve: items, totalItems, totalPages, error
	GetByTicker(ticker string, page, pageSize int) ([]models.Stock, int, int, error)
	GetByAction(action string, page, pageSize int) ([]models.Stock, int, int, error)
	GetByRatingTo(rating string, page, pageSize int) ([]models.Stock, int, int, error)
	GetByRatingFrom(rating string, page, pageSize int) ([]models.Stock, int, int, error)
	GetByBrokerage(brokerage string, page, pageSize int) ([]models.Stock, int, int, error)
	GetByDateRange(startDate, endDate time.Time, page, pageSize int) ([]models.Stock, int, int, error)
	GetByCompany(company string, page, pageSize int) ([]models.Stock, int, int, error)
	SearchStocks(query string, page, pageSize int) ([]models.Stock, int, int, error)
	GetByPriceRange(minPrice, maxPrice string, page, pageSize int) ([]models.Stock, int, int, error)
	GetSortedStocks(sortBy, sortOrder, search string, page, pageSize int) ([]models.Stock, int, int, error)

	// --- Método para Recomendaciones ---
	GetRecentRecommendations(since time.Time) ([]models.Stock, error) // Fetch recent stock data

	// --- Otros Métodos ---
	GetCount() (int, error) // Devuelve int
	UpsertStocksParallel(stocks []models.Stock, syncDate time.Time) (int, map[string]string, error)
	GetActionCounts() (map[string]int, error)
	GetRatingCounts() (map[string]int, error)
	InsertStock(stock models.Stock) error // (Considera si estos métodos Insert* son necesarios si ya tienes Upsert)
	InsertStocks(stocks []models.Stock) (int, map[string]string, error)
	InsertStocksParallel(stocks []models.Stock) (int, map[string]string, error)
	TruncateTable() error
}

// StockService define la interfaz para la lógica de negocio
type StockService interface {
	GetAllStocks(page, pageSize int) ([]models.Stock, int, int, error) // Devuelve: items, totalItems, totalPages, error
	GetTotalCount() (int, error)
	GetStockByTicker(ticker string, page, pageSize int) ([]models.Stock, int, int, error)
	GetStocksByAction(action string, page, pageSize int) ([]models.Stock, int, int, error)
	GetStocksByRatingTo(rating string, page, pageSize int) ([]models.Stock, int, int, error)
	GetStocksByRatingFrom(rating string, page, pageSize int) ([]models.Stock, int, int, error)
	GetStocksByCompany(company string, page, pageSize int) ([]models.Stock, int, int, error)
	GetActionStats() (map[string]int, error)
	GetRatingStats() (map[string]int, error)
	SyncStockData() (SyncResult, error)
	SearchStocks(query string, page, pageSize int) ([]models.Stock, int, int, error)
	GetStocksByBrokerage(brokerage string, page, pageSize int) ([]models.Stock, int, int, error)
	GetStocksByDateRange(startDate, endDate time.Time, page, pageSize int) ([]models.Stock, int, int, error)
	GetSortedStocks(sortBy, search string, page, pageSize int) ([]models.Stock, int, int, error)
	GetStocksByPriceRange(minPrice, maxPrice string, page, pageSize int) ([]models.Stock, int, int, error)

	// --- Método para Recomendaciones ---
	RecommendStocks(limit int) ([]models.Recommendation, error) // Generate recommendations
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
