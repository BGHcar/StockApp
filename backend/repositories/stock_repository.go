package repositories

import (
	"backend/interfaces"
	"backend/models"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm/clause" // Importar el paquete clause
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

// TruncateTable elimina todos los registros de la tabla stocks
func (r *StockRepository) TruncateTable() error {
	return r.db.DB().Exec("TRUNCATE TABLE stocks").Error
}

// GetAll obtiene todos los stocks de la base de datos
func (r *StockRepository) GetAll() ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.DB().Order("time DESC").Find(&stocks)
	return stocks, result.Error
}

// GetCount obtiene el número total de registros en la tabla stocks
func (r *StockRepository) GetCount() (int64, error) {
	var count int64
	result := r.db.DB().Model(&models.Stock{}).Count(&count)
	return count, result.Error
}

// GetByTicker obtiene stocks que coincidan con el ticker (búsqueda parcial)
func (r *StockRepository) GetByTicker(ticker string) ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.DB().Where("ticker ILIKE ?", "%"+ticker+"%").
		Order("time DESC").
		Find(&stocks)
	return stocks, result.Error
}

// GetByAction obtiene stocks filtrados por tipo de acción
func (r *StockRepository) GetByAction(action string) ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.DB().Where("action ILIKE ?", "%"+action+"%").
		Order("time DESC").
		Find(&stocks)
	return stocks, result.Error
}

// GetByRating obtiene stocks filtrados por rating
func (r *StockRepository) GetByRating(rating string) ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.DB().Where("rating_to ILIKE ?", "%"+rating+"%").
		Order("time DESC").
		Find(&stocks)
	return stocks, result.Error
}

// GetActionCounts obtiene conteos por tipo de acción
func (r *StockRepository) GetActionCounts() (map[string]int, error) {
	type Result struct {
		Action string
		Count  int
	}
	var results []Result

	err := r.db.DB().Model(&models.Stock{}).
		Select("action, count(*) as count").
		Group("action").
		Order("count DESC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, result := range results {
		counts[result.Action] = result.Count
	}

	return counts, nil
}

// GetByCompany obtiene stocks filtrados por compañía
func (r *StockRepository) GetByCompany(company string) ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.DB().Where("company ILIKE ?", "%"+company+"%").
		Order("time DESC").
		Find(&stocks)
	return stocks, result.Error
}

// GetRatingCounts obtiene conteos por tipo de rating
func (r *StockRepository) GetRatingCounts() (map[string]int, error) {
	type Result struct {
		Rating string `gorm:"column:rating_to"`
		Count  int
	}
	var results []Result

	err := r.db.DB().Model(&models.Stock{}).
		Select("rating_to, count(*) as count").
		Group("rating_to").
		Order("count DESC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, result := range results {
		counts[result.Rating] = result.Count
	}

	return counts, nil
}

// GetByBrokerage obtiene stocks filtrados por brokerage
func (r *StockRepository) GetByBrokerage(brokerage string) ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.DB().Where("brokerage ILIKE ?", "%"+brokerage+"%").
		Order("time DESC").
		Find(&stocks)

	// Añadir log para depuración
	log.Printf("Búsqueda por brokerage '%s': %d resultados encontrados", brokerage, len(stocks))

	return stocks, result.Error
}

// InsertStock inserta un único registro de stock
func (r *StockRepository) InsertStock(stock models.Stock) error {
	return r.db.DB().Create(&stock).Error
}

// InsertStocks inserta múltiples registros de stock
func (r *StockRepository) InsertStocks(stocks []models.Stock) (int, map[string]string, error) {
	inserted := 0
	failedTickers := make(map[string]string)

	// Crear un mapa para rastrear tickers insertados
	insertedTickers := make(map[string]bool)

	// Insertar registros uno por uno para mejor control
	for _, stock := range stocks {
		// Usar directamente clause.Columns y clause.OnConflict
		err := r.db.DB().Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "ticker"}, {Name: "time"}},
			DoNothing: true,
		}).Create(&stock).Error

		if err != nil {
			failedTickers[stock.Ticker] = err.Error()
		} else if _, isInserted := insertedTickers[stock.Ticker]; !isInserted {
			insertedTickers[stock.Ticker] = true
			inserted++

			// Log cada 100 inserciones
			if inserted%100 == 0 {
				log.Printf("Progreso: %d elementos insertados", inserted)
			}
		} else {
			// Duplicado
			tickerWithTime := fmt.Sprintf("%s-%s", stock.Ticker, stock.Time.Format("20060102150405"))
			failedTickers[tickerWithTime] = "Ticker duplicado (ya insertado previamente)"
		}
	}

	return inserted, failedTickers, nil
}

// InsertStocksParallel inserta múltiples registros en paralelo
func (r *StockRepository) InsertStocksParallel(stocks []models.Stock) (int, map[string]string, error) {
	inserted := int32(0)
	insertedTickers := sync.Map{}
	failedTickers := sync.Map{}

	// Número de workers y tamaño de lotes
	numWorkers := 10
	batchSize := 50

	// Dividir los stocks en lotes
	var stockBatches [][]models.Stock
	for i := 0; i < len(stocks); i += batchSize {
		end := i + batchSize
		if end > len(stocks) {
			end = len(stocks)
		}
		stockBatches = append(stockBatches, stocks[i:end])
	}

	log.Printf("Procesando %d registros en %d lotes con %d workers",
		len(stocks), len(stockBatches), numWorkers)

	// Canal para distribuir lotes a workers
	batchChan := make(chan []models.Stock, len(stockBatches))
	var wg sync.WaitGroup

	// Iniciar workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			// Cada worker procesa lotes del canal
			for batch := range batchChan {
				if len(batch) == 0 {
					continue
				}

				// Crear transacción para el lote
				tx := r.db.DB().Begin()

				for _, stock := range batch {
					// Usar directamente clause.Columns y clause.OnConflict
					err := tx.Clauses(clause.OnConflict{
						Columns:   []clause.Column{{Name: "ticker"}, {Name: "time"}},
						DoNothing: true,
					}).Create(&stock).Error

					if err != nil {
						failedTickers.Store(stock.Ticker, err.Error())
					} else {
						ticker := stock.Ticker
						if _, loaded := insertedTickers.LoadOrStore(ticker, true); !loaded {
							atomic.AddInt32(&inserted, 1)
						}
					}
				}

				// Commit transacción
				if err := tx.Commit().Error; err != nil {
					log.Printf("Error en worker %d al hacer commit del lote: %v", workerId, err)
					for _, stock := range batch {
						failedTickers.Store(stock.Ticker, err.Error())
					}
				}
			}
		}(i)
	}

	// Enviar lotes a los workers
	for _, batch := range stockBatches {
		batchChan <- batch
	}

	close(batchChan)
	wg.Wait()

	// Convertir resultados de sync.Map a map estándar
	resultInserted := int(inserted)
	resultFailed := make(map[string]string)

	failedTickers.Range(func(key, value interface{}) bool {
		resultFailed[key.(string)] = value.(string)
		return true
	})

	log.Printf("Inserción paralela completada: %d elementos insertados, %d errores",
		resultInserted, len(resultFailed))

	return resultInserted, resultFailed, nil
}

// GetByDateRange obtiene stocks dentro de un rango de fechas
func (r *StockRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.DB().Where("time BETWEEN ? AND ?", startDate, endDate).
		Order("time DESC").
		Find(&stocks)
	return stocks, result.Error
}

// SearchStocks búsqueda general por compañía o ticker
func (r *StockRepository) SearchStocks(query string) ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.DB().Where(
		"ticker ILIKE ? OR company ILIKE ? OR brokerage ILIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Order("time DESC").
		Limit(100).
		Find(&stocks)
	return stocks, result.Error
}

// GetByPriceRange obtiene stocks dentro de un rango de precios objetivo
func (r *StockRepository) GetByPriceRange(minPrice, maxPrice string) ([]models.Stock, error) {
	// Limpiar valores de precio
	minPriceClean := minPrice
	maxPriceClean := maxPrice

	if len(minPrice) > 0 && minPrice[0] == '$' {
		minPriceClean = minPrice[1:]
	}

	if len(maxPrice) > 0 && maxPrice[0] == '$' {
		maxPriceClean = maxPrice[1:]
	}

	// Usar el método Where de GORM con SQL personalizado para el filtrado
	var stocks []models.Stock
	result := r.db.DB().Where(`
        (CAST(REGEXP_REPLACE(target_from, '[^0-9.]', '', 'g') AS DECIMAL(10,2)) 
         BETWEEN ?::DECIMAL(10,2) AND ?::DECIMAL(10,2))
        AND
        (CAST(REGEXP_REPLACE(target_to, '[^0-9.]', '', 'g') AS DECIMAL(10,2)) 
         BETWEEN ?::DECIMAL(10,2) AND ?::DECIMAL(10,2))
    `, minPriceClean, maxPriceClean, minPriceClean, maxPriceClean).
		Order("time DESC").
		Find(&stocks)

	// Registrar para depuración
	log.Printf("Búsqueda por rango de precio %s-%s: %d resultados encontrados",
		minPrice, maxPrice, len(stocks))

	return stocks, result.Error
}

// Función auxiliar para extraer precio numérico
func extractNumericPrice(priceStr string) float64 {
	if len(priceStr) == 0 {
		return 0
	}

	// Eliminar símbolo de dólar
	if priceStr[0] == '$' {
		priceStr = priceStr[1:]
	}

	// Eliminar caracteres no numéricos
	numStr := ""
	for _, r := range priceStr {
		if (r >= '0' && r <= '9') || r == '.' {
			numStr += string(r)
		}
	}

	// Convertir a float
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}
	return val
}
