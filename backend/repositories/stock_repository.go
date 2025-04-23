package repositories

import (
	"backend/interfaces"
	"backend/models"
	"log"
	"math" // Necesario para Ceiling
	"strings"
	"sync"
	"sync/atomic"
	"time"

	// Asegúrate de que gorm esté importado
	"gorm.io/gorm/clause"
)

// StockRepository implementa la interfaz Repository
type StockRepository struct {
	db interfaces.DatabaseHandler
}

// Verificamos que StockRepository implementa la interfaz
var _ interfaces.StockRepository = (*StockRepository)(nil)

// NewStockRepository crea una nueva instancia del repositorio
func NewStockRepository(db interfaces.DatabaseHandler) *StockRepository {
	return &StockRepository{
		db: db,
	}
}

// --- Helper para calcular totalPages ---
func calculateTotalPages(totalItems int64, pageSize int) int {
	if pageSize <= 0 {
		return 0 // O 1 si prefieres que siempre haya al menos una página
	}
	return int(math.Ceil(float64(totalItems) / float64(pageSize)))
}

// GetAll obtiene todos los stocks de la base de datos paginados
func (r *StockRepository) GetAll(page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64

	// Contar total de items (GORM devuelve int64)
	if err := r.db.DB().Model(&models.Stock{}).Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	// Calcular offset
	offset := (page - 1) * pageSize

	// Obtener items paginados
	result := r.db.DB().
		Order("time DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&stocks)

	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	// Calcular totalPages
	totalPages := calculateTotalPages(totalCount, pageSize)

	return stocks, int(totalCount), totalPages, nil // Devolver int, int, error
}

// GetCount obtiene el número total de registros en la tabla stocks
func (r *StockRepository) GetCount() (int, error) {
	var count int64
	result := r.db.DB().Model(&models.Stock{}).Count(&count)
	return int(count), result.Error // Convertir a int
}

// GetByTicker obtiene stocks que coincidan con el ticker (búsqueda parcial) paginados
func (r *StockRepository) GetByTicker(ticker string, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64
	query := r.db.DB().Model(&models.Stock{}).Where("ticker ILIKE ?", "%"+ticker+"%")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// GetByAction obtiene stocks filtrados por tipo de acción paginados
func (r *StockRepository) GetByAction(action string, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64
	query := r.db.DB().Model(&models.Stock{}).Where("action ILIKE ?", "%"+action+"%")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// GetByRatingTo obtiene stocks filtrados por rating_to paginados
func (r *StockRepository) GetByRatingTo(rating string, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64
	query := r.db.DB().Model(&models.Stock{}).Where("rating_to ILIKE ?", "%"+rating+"%")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// GetByRatingFrom obtiene stocks filtrados por rating_from paginados
func (r *StockRepository) GetByRatingFrom(rating string, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64
	query := r.db.DB().Model(&models.Stock{}).Where("rating_from ILIKE ?", "%"+rating+"%")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// GetByBrokerage obtiene stocks filtrados por brokerage paginados
func (r *StockRepository) GetByBrokerage(brokerage string, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64
	query := r.db.DB().Model(&models.Stock{}).Where("brokerage ILIKE ?", "%"+brokerage+"%")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// GetByDateRange obtiene stocks dentro de un rango de fechas paginados
func (r *StockRepository) GetByDateRange(startDate, endDate time.Time, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64
	query := r.db.DB().Model(&models.Stock{}).Where("time BETWEEN ? AND ?", startDate, endDate)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// GetByCompany obtiene stocks filtrados por compañía paginados
func (r *StockRepository) GetByCompany(company string, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64
	query := r.db.DB().Model(&models.Stock{}).Where("company ILIKE ?", "%"+company+"%")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// SearchStocks búsqueda general por compañía, ticker o brokerage paginada
func (r *StockRepository) SearchStocks(query string, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64
	dbQuery := r.db.DB().Model(&models.Stock{}).
		Where("ticker ILIKE ? OR company ILIKE ? OR brokerage ILIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	if err := dbQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := dbQuery.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// GetByPriceRange obtiene stocks dentro de un rango de precios objetivo paginados
func (r *StockRepository) GetByPriceRange(minPrice, maxPrice string, page, pageSize int) ([]models.Stock, int, int, error) {
	var stocks []models.Stock
	var totalCount int64

	// Convertir los valores de minPrice y maxPrice eliminando comas y el símbolo $
	minPrice = strings.ReplaceAll(minPrice, ",", "")
	minPrice = strings.ReplaceAll(minPrice, "$", "")
	maxPrice = strings.ReplaceAll(maxPrice, ",", "")
	maxPrice = strings.ReplaceAll(maxPrice, "$", "")

	// Ajustar la consulta para considerar solo target_from dentro del rango
	query := r.db.DB().Model(&models.Stock{}).
		Where("CAST(REPLACE(REPLACE(target_from, '$', ''), ',', '') AS DECIMAL) BETWEEN ? AND ?", minPrice, maxPrice)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, 0, err
	}

	offset := (page - 1) * pageSize
	result := query.Limit(pageSize).Offset(offset).Order("time DESC").Find(&stocks)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	totalPages := calculateTotalPages(totalCount, pageSize)
	return stocks, int(totalCount), totalPages, nil
}

// --- Métodos No Paginados (Revisados) ---

// GetActionCounts obtiene conteos por tipo de acción
func (r *StockRepository) GetActionCounts() (map[string]int, error) {
	type Result struct {
		Action string
		Count  int
	}
	var results []Result
	err := r.db.DB().Model(&models.Stock{}).Select("action, count(*) as count").Group("action").Order("count DESC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, result := range results {
		counts[result.Action] = result.Count
	}
	return counts, nil
}

// GetRatingCounts obtiene conteos por tipo de rating_to
func (r *StockRepository) GetRatingCounts() (map[string]int, error) {
	type Result struct {
		Rating string `gorm:"column:rating_to"`
		Count  int
	}
	var results []Result
	err := r.db.DB().Model(&models.Stock{}).Select("rating_to, count(*) as count").Group("rating_to").Order("count DESC").Find(&results).Error
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, result := range results {
		counts[result.Rating] = result.Count
	}
	return counts, nil
}

// TruncateTable elimina todos los registros de la tabla stocks
func (r *StockRepository) TruncateTable() error {
	// ¡CUIDADO! Esto elimina todos los datos. Asegúrate de que es lo que quieres.
	return r.db.DB().Exec("TRUNCATE TABLE stocks").Error
}

// InsertStock inserta un único registro de stock
func (r *StockRepository) InsertStock(stock models.Stock) error {
	return r.db.DB().Create(&stock).Error
}

// InsertStocks (Considera eliminar si UpsertStocksParallel es suficiente)
func (r *StockRepository) InsertStocks(stocks []models.Stock) (int, map[string]string, error) {
	// Implementación simplificada, puede ser ineficiente para muchos registros
	inserted := 0
	failedTickers := make(map[string]string)
	for _, stock := range stocks {
		err := r.InsertStock(stock)
		if err != nil {
			failedTickers[stock.Ticker+"_"+stock.Time.String()] = err.Error() // Usar clave única
		} else {
			inserted++
		}
	}
	return inserted, failedTickers, nil
}

// InsertStocksParallel (Considera eliminar si UpsertStocksParallel es suficiente)
func (r *StockRepository) InsertStocksParallel(stocks []models.Stock) (int, map[string]string, error) {
	// Esta implementación necesita ser revisada/completada si se va a usar
	log.Println("ADVERTENCIA: InsertStocksParallel no está completamente implementada.")
	return r.InsertStocks(stocks) // Llama a la versión secuencial por ahora
}

// UpsertStocksParallel inserta o actualiza múltiples registros en paralelo
func (r *StockRepository) UpsertStocksParallel(stocks []models.Stock, syncDate time.Time) (int, map[string]string, error) {
	inserted := int32(0)
	failedTickers := sync.Map{} // Usar sync.Map para concurrencia

	numWorkers := 10
	batchSize := 100 // Aumentar tamaño de lote para Upsert

	var stockBatches [][]models.Stock
	for i := 0; i < len(stocks); i += batchSize {
		end := i + batchSize
		if end > len(stocks) {
			end = len(stocks)
		}
		stockBatches = append(stockBatches, stocks[i:end])
	}

	log.Printf("Upsert: Procesando %d registros en %d lotes con %d workers", len(stocks), len(stockBatches), numWorkers)

	batchChan := make(chan []models.Stock, len(stockBatches))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for batch := range batchChan {
				// Usar transacción para el lote
				tx := r.db.DB().Begin()
				if tx.Error != nil {
					log.Printf("Worker %d: Error iniciando transacción: %v", workerId, tx.Error)
					for _, stock := range batch {
						failedTickers.Store(stock.Ticker+"_"+stock.Time.String(), tx.Error.Error())
					}
					continue // Saltar este lote
				}

				// Preparar la cláusula ON CONFLICT
				// Asumiendo clave primaria (ticker, time)
				conflictClause := clause.OnConflict{
					Columns: []clause.Column{{Name: "ticker"}, {Name: "time"}},
					DoUpdates: clause.AssignmentColumns([]string{
						"company", "target_from", "target_to", "action",
						"brokerage", "rating_from", "rating_to", "updated_at",
					}), // Actualizar todos los campos excepto claves primarias y created_at
				}

				// Ejecutar Upsert para el lote
				result := tx.Clauses(conflictClause).Create(&batch)

				if result.Error != nil {
					tx.Rollback()
					log.Printf("Worker %d: Error en Upsert del lote: %v", workerId, result.Error)
					for _, stock := range batch {
						failedTickers.Store(stock.Ticker+"_"+stock.Time.String(), result.Error.Error())
					}
					continue
				}

				// Contar resultados (GORM no facilita esto directamente en Upsert masivo)
				// Estimación basada en RowsAffected (puede no ser precisa para skipped)
				// Una mejor aproximación requeriría consultar antes/después o usar RETURNING si la DB lo soporta bien con GORM
				atomic.AddInt32(&inserted, int32(result.RowsAffected)) // Asumimos que RowsAffected son inserciones/actualizaciones

				// Commit transacción
				if err := tx.Commit().Error; err != nil {
					log.Printf("Worker %d: Error haciendo commit: %v", workerId, err)
					// Marcar como fallidos si el commit falla
					for _, stock := range batch {
						failedTickers.Store(stock.Ticker+"_"+stock.Time.String(), err.Error())
					}
				}
			}
		}(i)
	}

	for _, batch := range stockBatches {
		batchChan <- batch
	}
	close(batchChan)
	wg.Wait()

	resultInserted := int(inserted) // O una lógica más precisa si se implementa
	resultFailed := make(map[string]string)
	failedTickers.Range(func(key, value interface{}) bool {
		resultFailed[key.(string)] = value.(string)
		return true
	})

	log.Printf("Upsert paralelo completado: %d procesados (aprox.), %d errores", resultInserted, len(resultFailed))
	return resultInserted, resultFailed, nil
}

func (r *StockRepository) GetRecentRecommendations(since time.Time) ([]models.Stock, error) {
	var stocks []models.Stock
	err := r.db.DB().Model(&models.Stock{}).
		Where("time >= ?", since).
		Order("time DESC").
		Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

// --- Funciones Auxiliares (si son necesarias) ---
// ... (como extractNumericPrice si se usa) ...
