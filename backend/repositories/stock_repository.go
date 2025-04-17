package repositories

import (
	"backend/interfaces"
	"backend/models"
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
	_, err := r.db.Exec("TRUNCATE TABLE stocks")
	return err
}

// GetAll obtiene todos los stocks de la base de datos
func (r *StockRepository) GetAll() ([]models.Stock, error) {
	rows, err := r.db.Query(`
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        ORDER BY time DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}

// GetCount obtiene el número total de registros en la tabla stocks
func (r *StockRepository) GetCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM stocks").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetByTicker obtiene stocks que coincidan con el ticker (búsqueda parcial)
func (r *StockRepository) GetByTicker(ticker string) ([]models.Stock, error) {
	// Modificamos la consulta para usar ILIKE con comodines
	rows, err := r.db.Query(`
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        WHERE ticker ILIKE $1
        ORDER BY time DESC
    `, "%"+ticker+"%") // Agregamos comodines antes y después

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetByAction obtiene stocks filtrados por tipo de acción
func (r *StockRepository) GetByAction(action string) ([]models.Stock, error) {
	// Modificamos la consulta para usar ILIKE con comodines
	rows, err := r.db.Query(`
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        WHERE action ILIKE $1
        ORDER BY time DESC
    `, "%"+action+"%") // Agregamos comodines antes y después

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetByRating obtiene stocks filtrados por rating
func (r *StockRepository) GetByRating(rating string) ([]models.Stock, error) {
	// Modificamos la consulta para usar ILIKE con comodines
	rows, err := r.db.Query(`
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        WHERE rating_to ILIKE $1
        ORDER BY time DESC
    `, "%"+rating+"%") // Agregamos comodines antes y después

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetActionCounts obtiene conteos por tipo de acción
func (r *StockRepository) GetActionCounts() (map[string]int, error) {
	rows, err := r.db.Query(`
        SELECT action, COUNT(*) as count
        FROM stocks
        GROUP BY action
        ORDER BY count DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var action string
		var count int
		if err := rows.Scan(&action, &count); err != nil {
			return nil, err
		}
		counts[action] = count
	}

	return counts, nil
}

func (r *StockRepository) GetByCompany(company string) ([]models.Stock, error) {
	// Modificamos la consulta para usar ILIKE con comodines
	rows, err := r.db.Query(`
		SELECT ticker, company, target_from, target_to,
			action, brokerage, rating_from, rating_to, time
			FROM stocks
			WHERE company ILIKE $1
			ORDER BY time DESC
	`, "%"+company+"%") // Agregamos comodines antes y después
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

// GetRatingCounts obtiene conteos por tipo de rating
func (r *StockRepository) GetRatingCounts() (map[string]int, error) {
	rows, err := r.db.Query(`
        SELECT rating_to, COUNT(*) as count
        FROM stocks
        GROUP BY rating_to
        ORDER BY count DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var rating string
		var count int
		if err := rows.Scan(&rating, &count); err != nil {
			return nil, err
		}
		counts[rating] = count
	}

	return counts, nil
}

func (r *StockRepository) GetByBrokerage(brokerage string) ([]models.Stock, error) {
	// Modificamos la consulta para usar ILIKE con comodines
	rows, err := r.db.Query(`
		SELECT ticker, company, target_from, target_to,
			action, brokerage, rating_from, rating_to, time
			FROM stocks 
			WHERE brokerage ILIKE $1
			ORDER BY time DESC
	`, "%"+brokerage+"%") // Agregamos comodines antes y después

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

// InsertStock inserta un único registro de stock
func (r *StockRepository) InsertStock(stock models.Stock) error {
	_, err := r.db.Exec(`
        INSERT INTO stocks (
            ticker, company, target_from, target_to, 
            action, brokerage, rating_from, rating_to, time
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `,
		stock.Ticker, stock.Company,
		stock.TargetFrom, stock.TargetTo,
		stock.Action, stock.Brokerage,
		stock.RatingFrom, stock.RatingTo, stock.Time,
	)
	return err
}

// InsertStocks inserta múltiples registros de stock y devuelve detalles sobre éxitos y fallos
func (r *StockRepository) InsertStocks(stocks []models.Stock) (int, map[string]string, error) {
	inserted := 0
	insertedTickers := make(map[string]bool)
	failedTickers := make(map[string]string)

	for _, stock := range stocks {
		// Modificar para usar la clave primaria compuesta (ticker, time)
		_, err := r.db.Exec(`
            INSERT INTO stocks (
                ticker, company, target_from, target_to, 
                action, brokerage, rating_from, rating_to, time
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            ON CONFLICT (ticker, time) DO NOTHING
        `,
			stock.Ticker, stock.Company,
			stock.TargetFrom, stock.TargetTo,
			stock.Action, stock.Brokerage,
			stock.RatingFrom, stock.RatingTo, stock.Time,
		)

		if err != nil {
			log.Printf("Error inserting item %s: %v", stock.Ticker, err)
			failedTickers[stock.Ticker] = err.Error()
		} else if _, isInserted := insertedTickers[stock.Ticker]; !isInserted {
			// Si es la primera vez que vemos este ticker
			insertedTickers[stock.Ticker] = true
			inserted++

			// Para reducir el spam de logs, solo mostramos cada 100 inserciones
			if inserted%100 == 0 {
				log.Printf("Progreso: %d elementos insertados", inserted)
			}
		} else {
			// Este ticker ya fue insertado previamente (duplicado)
			tickerWithTime := fmt.Sprintf("%s-%s", stock.Ticker, stock.Time.Format("20060102150405"))
			failedTickers[tickerWithTime] = "Ticker duplicado (ya insertado previamente)"
		}
	}

	return inserted, failedTickers, nil
}

// InsertStocksParallel inserta múltiples registros de stock en paralelo usando workers
func (r *StockRepository) InsertStocksParallel(stocks []models.Stock) (int, map[string]string, error) {
	inserted := int32(0)
	insertedTickers := sync.Map{}
	failedTickers := sync.Map{}

	// Número de workers - ajustar según la capacidad del cluster y latencia
	numWorkers := 10
	batchSize := 50 // Tamaño de cada lote para inserción por lotes

	// Preparar la consulta para inserción por lotes
	// Esto reduce drásticamente el número de viajes de red
	insertQuery := `
        INSERT INTO stocks (
            ticker, company, target_from, target_to, 
            action, brokerage, rating_from, rating_to, time
        ) VALUES 
    `

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

				// Construir la consulta para este lote
				query := insertQuery
				var values []interface{}
				placeholders := []string{}

				for i, stock := range batch {
					// Crear placeholders para consulta preparada
					baseIdx := i * 9
					phRow := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
						baseIdx+1, baseIdx+2, baseIdx+3, baseIdx+4,
						baseIdx+5, baseIdx+6, baseIdx+7, baseIdx+8, baseIdx+9)
					placeholders = append(placeholders, phRow)

					// Agregar valores para esta fila
					values = append(values,
						stock.Ticker, stock.Company,
						stock.TargetFrom, stock.TargetTo,
						stock.Action, stock.Brokerage,
						stock.RatingFrom, stock.RatingTo, stock.Time)

					// Registrar este ticker como procesado
					ticker := stock.Ticker
					if _, loaded := insertedTickers.LoadOrStore(ticker, true); !loaded {
						atomic.AddInt32(&inserted, 1)
					}
				}

				// Finalizar la consulta con ON CONFLICT
				query += strings.Join(placeholders, ", ")
				query += " ON CONFLICT (ticker, time) DO NOTHING"

				// Ejecutar la inserción por lotes
				_, err := r.db.Exec(query, values...)
				if err != nil {
					log.Printf("Error en worker %d insertando lote: %v", workerId, err)
					// Registrar error para cada ticker en este lote
					for _, stock := range batch {
						failedTickers.Store(stock.Ticker, err.Error())
					}
				}
			}
		}(i)
	}

	// Enviar todos los lotes a los workers
	for _, batch := range stockBatches {
		batchChan <- batch
	}

	// Cerrar el canal y esperar a que terminen todos los workers
	close(batchChan)
	wg.Wait()

	// Convertir resultados de sync.Map a map estándar
	resultInserted := int(inserted)
	resultFailed := make(map[string]string)

	failedTickers.Range(func(key, value interface{}) bool {
		resultFailed[key.(string)] = value.(string)
		return true
	})

	// Registrar progreso final
	log.Printf("Inserción paralela completada: %d elementos insertados, %d errores",
		resultInserted, len(resultFailed))

	return resultInserted, resultFailed, nil
}

// GetByDateRange obtiene stocks dentro de un rango de fechas
func (r *StockRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Stock, error) {
	rows, err := r.db.Query(`
		SELECT ticker, company, target_from, target_to,
			action, brokerage, rating_from, rating_to, time
			FROM stocks
			WHERE time BETWEEN $1 AND $2
			ORDER BY time DESC
	`, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// Añadimos un nuevo método para búsqueda general por compañía o ticker
func (r *StockRepository) SearchStocks(query string) ([]models.Stock, error) {
	rows, err := r.db.Query(`
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        WHERE 
            ticker ILIKE $1 OR 
            company ILIKE $1 OR 
            brokerage ILIKE $1
        ORDER BY time DESC
        LIMIT 100
    `, "%"+query+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil
}

// GetByPriceRange obtiene stocks dentro de un rango de precios objetivo
func (r *StockRepository) GetByPriceRange(minPrice, maxPrice string) ([]models.Stock, error) {
	// Limpiar los valores de precio
	minPriceClean := minPrice
	maxPriceClean := maxPrice

	if len(minPrice) > 0 && minPrice[0] == '$' {
		minPriceClean = minPrice[1:]
	}

	if len(maxPrice) > 0 && maxPrice[0] == '$' {
		maxPriceClean = maxPrice[1:]
	}

	// Usar una consulta SQL directa con extracción numérica
	query := `
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        WHERE 
            (
                CAST(REGEXP_REPLACE(target_from, '[^0-9.]', '', 'g') AS DECIMAL(10,2)) 
                BETWEEN $1::DECIMAL(10,2) AND $2::DECIMAL(10,2)
            )
            AND
            (
                CAST(REGEXP_REPLACE(target_to, '[^0-9.]', '', 'g') AS DECIMAL(10,2)) 
                BETWEEN $1::DECIMAL(10,2) AND $2::DECIMAL(10,2)
            )
        ORDER BY time DESC
    `

	rows, err := r.db.Query(query, minPriceClean, maxPriceClean)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		if err := rows.Scan(
			&stock.Ticker, &stock.Company,
			&stock.TargetFrom, &stock.TargetTo,
			&stock.Action, &stock.Brokerage,
			&stock.RatingFrom, &stock.RatingTo, &stock.Time,
		); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}
