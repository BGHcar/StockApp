package repositories

import (
	"backend/interfaces"
	"backend/models"
	"database/sql"
	"fmt"
	"log"
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

// GetByTicker obtiene un stock específico por su ticker
func (r *StockRepository) GetByTicker(ticker string) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.QueryRow(`
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        WHERE ticker = $1
    `, ticker).Scan(
		&stock.Ticker, &stock.Company,
		&stock.TargetFrom, &stock.TargetTo,
		&stock.Action, &stock.Brokerage,
		&stock.RatingFrom, &stock.RatingTo, &stock.Time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró el stock, no es un error
		}
		return nil, err
	}
	return &stock, nil
}

// GetByAction obtiene stocks filtrados por tipo de acción
func (r *StockRepository) GetByAction(action string) ([]models.Stock, error) {
	rows, err := r.db.Query(`
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        WHERE action = $1
        ORDER BY time DESC
    `, action)
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
	rows, err := r.db.Query(`
        SELECT ticker, company, target_from, target_to, 
               action, brokerage, rating_from, rating_to, time
        FROM stocks
        WHERE rating_to = $1
        ORDER BY time DESC
    `, rating)
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
		_, err := r.db.Exec(`
            INSERT INTO stocks (
                ticker, company, target_from, target_to, 
                action, brokerage, rating_from, rating_to, time
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            ON CONFLICT (ticker) DO NOTHING
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
