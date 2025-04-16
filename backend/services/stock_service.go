package services

import (
	"backend/interfaces"
	"backend/models"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

// StockService implementa la interfaz Service
type StockService struct {
	repo      interfaces.StockRepository
	apiClient interfaces.APIClient
}

// Verificamos que StockService implementa la interfaz
var _ interfaces.StockService = (*StockService)(nil)

// NewStockService crea una nueva instancia del servicio de stocks
func NewStockService(repo interfaces.StockRepository, apiClient interfaces.APIClient) *StockService {
	return &StockService{
		repo:      repo,
		apiClient: apiClient,
	}
}

// GetAllStocks obtiene todos los stocks
func (s *StockService) GetAllStocks() ([]models.Stock, error) {
	return s.repo.GetAll()
}

// GetTotalCount obtiene el número total de stocks
func (s *StockService) GetTotalCount() (int, error) {
	return s.repo.GetCount()
}

// GetStockByTicker obtiene un stock específico por su ticker
func (s *StockService) GetStockByTicker(ticker string) (*models.Stock, error) {
	return s.repo.GetByTicker(ticker)
}

// GetStocksByAction obtiene stocks filtrados por tipo de acción
func (s *StockService) GetStocksByAction(action string) ([]models.Stock, error) {
	return s.repo.GetByAction(action)
}

// GetStocksByRating obtiene stocks filtrados por rating
func (s *StockService) GetStocksByRating(rating string) ([]models.Stock, error) {
	return s.repo.GetByRating(rating)
}

func (s *StockService) GetStocksByBrokerage(brokerage string) ([]models.Stock, error) {
	return s.repo.GetByBrokerage(brokerage)
}

// GetActionStats obtiene estadísticas por tipo de acción
func (s *StockService) GetActionStats() (map[string]int, error) {
	return s.repo.GetActionCounts()
}

// GetRatingStats obtiene estadísticas por tipo de rating
func (s *StockService) GetRatingStats() (map[string]int, error) {
	return s.repo.GetRatingCounts()
}

// GetStocksByDateRange obtiene stocks por rango de fechas
func (s *StockService) GetStocksByDateRange(startDate, endDate time.Time) ([]models.Stock, error) {
	return s.repo.GetByDateRange(startDate, endDate)
}


// SyncStockData sincroniza datos de stocks desde la API
func (s *StockService) SyncStockData() (interfaces.SyncResult, error) {
	result := interfaces.SyncResult{
		FailedInsertDetails: make(map[string]string),
	}

	// Mostrar mensaje de inicio con timestamp
	startTime := time.Now()
	log.Printf("Iniciando sincronización de datos a las %s", 
		startTime.Format("15:04:05"))

	// Limpiar la tabla antes de sincronizar
	if err := s.repo.TruncateTable(); err != nil {
		return result, fmt.Errorf("error truncating table: %w", err)
	}
	log.Println("Tabla stocks limpiada correctamente")

	// Recolectar datos de la API
	stockItems, duplicates, err := s.collectAPIData()
	if err != nil {
		return result, fmt.Errorf("error collecting API data: %w", err)
	}

	// Preparar resultado con datos de duplicados
	result.TotalProcessed = len(stockItems)
	result.UniqueTickersAPI = len(duplicates.tickerCounts)
	result.DuplicateTickers = duplicates.count
	result.DuplicatesList = duplicates.items

	// Convertir a modelos de dominio
	stocks := s.convertToModels(stockItems)

	log.Printf("Iniciando inserción paralela de %d registros...", len(stocks))
	insertStartTime := time.Now()
	
	// Usar la nueva inserción paralela
	inserted, failedInserts, err := s.repo.InsertStocksParallel(stocks)
	
	insertDuration := time.Since(insertStartTime)
	log.Printf("Inserción completada en %.2f segundos", 
		insertDuration.Seconds())
	
	if err != nil {
		return result, fmt.Errorf("error inserting stocks: %w", err)
	}

	// Completar el resultado
	result.TotalInserted = inserted
	result.FailedInserts = len(failedInserts)
	result.UniqueTickersDB = inserted
	result.FailedInsertDetails = failedInserts

	// Calcular tiempo total
	totalDuration := time.Since(startTime)
	log.Printf("Sincronización completada en %.2f segundos", 
		totalDuration.Seconds())
	
	// Guardar log de errores
	s.logFailedInserts(failedInserts)

	// Resumen en log
	s.logSyncSummary(result)

	return result, nil
}

// --- Métodos privados para descomponer la responsabilidad ---

type duplicateInfo struct {
	count        int
	tickerCounts map[string]int
	items        []interfaces.TickerDuplicate
}

// collectAPIData obtiene y analiza datos de la API
func (s *StockService) collectAPIData() ([]interfaces.StockItem, duplicateInfo, error) {
	var allItems []interfaces.StockItem
	tickerCounts := make(map[string]int)
	nextPage := ""
	duplicateFound := false

	log.Println("Recolectando datos de la API...")

	for {
		resp, err := s.apiClient.Get("list", nextPage)
		if err != nil {
			return nil, duplicateInfo{}, fmt.Errorf("error fetching data: %w", err)
		}

		// Acumular items y contar tickers
		for _, item := range resp.Items {
			allItems = append(allItems, item)
			tickerCounts[item.Ticker]++

			if tickerCounts[item.Ticker] > 1 && !duplicateFound {
				duplicateFound = true
				log.Printf("¡ATENCIÓN! Se detectó un ticker duplicado: %s", item.Ticker)
				log.Println("Los tickers duplicados pueden causar problemas si la tabla tiene ticker como clave primaria")
			}
		}

		// Verificar si hay más páginas
		if resp.NextPage == "" {
			break
		}
		nextPage = resp.NextPage
	}

	// Contar duplicados y prepararlos para el informe
	duplicatedCount := 0
	var duplicatesList []interfaces.TickerDuplicate

	for ticker, count := range tickerCounts {
		if count > 1 {
			duplicatedCount++
			duplicatesList = append(duplicatesList, interfaces.TickerDuplicate{
				Ticker: ticker,
				Count:  count,
			})
		}
	}

	// Ordenar la lista de duplicados por conteo
	sort.Slice(duplicatesList, func(i, j int) bool {
		return duplicatesList[i].Count > duplicatesList[j].Count
	})

	duplicates := duplicateInfo{
		count:        duplicatedCount,
		tickerCounts: tickerCounts,
		items:        duplicatesList,
	}

	log.Printf("Total de registros recogidos de la API: %d", len(allItems))
	log.Printf("Total de tickers únicos: %d", len(tickerCounts))

	return allItems, duplicates, nil
}

// convertToModels convierte items de API a modelos de dominio
func (s *StockService) convertToModels(items []interfaces.StockItem) []models.Stock {
	stocks := make([]models.Stock, len(items))
	for i, item := range items {
		stocks[i] = models.Stock{
			Ticker:     item.Ticker,
			Company:    item.Company,
			TargetFrom: item.TargetFrom,
			TargetTo:   item.TargetTo,
			Action:     item.Action,
			Brokerage:  item.Brokerage,
			RatingFrom: item.RatingFrom,
			RatingTo:   item.RatingTo,
			Time:       item.Time,
		}
	}
	return stocks
}

// logFailedInserts guarda los errores de inserción en un archivo de log
func (s *StockService) logFailedInserts(failedInserts map[string]string) {
	if len(failedInserts) == 0 {
		return
	}

	// Crear o abrir el archivo de log
	f, err := os.OpenFile("failed_inserts.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("Error creating log file: %v", err)
		return
	}
	defer f.Close()

	// Escribir una cabecera con timestamp
	f.WriteString("\n===== ERROR LOG - " + time.Now().Format("2006-01-02 15:04:05") + " =====\n")
	f.WriteString(fmt.Sprintf("Total elementos fallidos: %d\n\n", len(failedInserts)))

	// Ordenar los tickers para que el log sea más legible
	var failedKeys []string
	for k := range failedInserts {
		failedKeys = append(failedKeys, k)
	}
	sort.Strings(failedKeys)

	// Escribir cada error
	for _, k := range failedKeys {
		f.WriteString(fmt.Sprintf("Ticker: %s - Motivo: %s\n", k, failedInserts[k]))
	}

	log.Printf("Se ha guardado el detalle de %d elementos no insertados en 'failed_inserts.log'", len(failedInserts))
}

// logSyncSummary muestra un resumen de la sincronización
func (s *StockService) logSyncSummary(result interfaces.SyncResult) {
	log.Println("\n===== INFORME DE SINCRONIZACIÓN =====")
	log.Printf("Total elementos procesados: %d", result.TotalProcessed)
	log.Printf("Total elementos insertados: %d", result.TotalInserted)
	log.Printf("Errores de inserción: %d", result.FailedInserts)
	log.Printf("Tickers únicos en API: %d", result.UniqueTickersAPI)
	log.Printf("Tickers únicos insertados: %d", result.UniqueTickersDB)

	if result.DuplicateTickers > 0 {
		log.Printf("\n===== INFORME DE DUPLICADOS =====")
		log.Printf("Se encontraron %d tickers duplicados en la API", result.DuplicateTickers)
		log.Printf("Esto representa %.2f%% del total de tickers", float64(result.DuplicateTickers)/float64(result.UniqueTickersAPI)*100)

		// Mostrar los 10 tickers con más duplicados
		log.Printf("\nTop 10 tickers duplicados:")

		// Mostrar top 10 o menos si hay menos
		limit := 10
		if len(result.DuplicatesList) < limit {
			limit = len(result.DuplicatesList)
		}

		for i := 0; i < limit; i++ {
			log.Printf("%d. %s - %d apariciones", i+1,
				result.DuplicatesList[i].Ticker,
				result.DuplicatesList[i].Count)
		}

		log.Printf("\n===== RECOMENDACIÓN TÉCNICA =====")
		if result.TotalInserted < result.TotalProcessed {
			log.Println("⚠️ ATENCIÓN: No todos los registros pudieron ser insertados debido a restricciones de clave primaria.")
			log.Println("Se recomienda modificar la tabla para usar una clave primaria compuesta (ticker, time):")
			log.Println("ALTER TABLE stocks DROP CONSTRAINT IF EXISTS stocks_pkey;")
			log.Println("ALTER TABLE stocks ADD PRIMARY KEY (ticker, time);")
		}
	} else {
		log.Println("\nNo se detectaron tickers duplicados en la API. La estructura actual de la tabla es adecuada.")
	}
}

// SearchStocks realiza una búsqueda general por ticker, compañía o brokerage
func (s *StockService) SearchStocks(query string) ([]models.Stock, error) {
	return s.repo.SearchStocks(query)
}
