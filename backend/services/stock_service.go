package services

import (
	"backend/interfaces"
	"backend/models"
	"fmt"
	"log"
	"os"
	"sort"
	"strings" // Import strings
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

// --- Helper para normalizar paginación ---
func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	// Definir límites por defecto y máximo
	defaultPageSize := 20
	maxPageSize := 100
	if pageSize < 1 || pageSize > maxPageSize {
		pageSize = defaultPageSize
	}
	return page, pageSize
}

// GetAllStocks obtiene todos los stocks paginados
func (s *StockService) GetAllStocks(page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	// Llama al repositorio que ahora devuelve los 4 valores
	stocks, totalItems, totalPages, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, 0, err // Devuelve ceros para int
	}
	return stocks, totalItems, totalPages, nil
}

// GetTotalCount obtiene el número total de stocks
func (s *StockService) GetTotalCount() (int, error) {
	// El repositorio ya devuelve int
	return s.repo.GetCount()
}

// GetStockByTicker obtiene stocks por ticker con paginación
func (s *StockService) GetStockByTicker(ticker string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.GetByTicker(ticker, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

// GetStocksByAction obtiene stocks por acción con paginación
func (s *StockService) GetStocksByAction(action string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.GetByAction(action, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

// GetStocksByRatingTo obtiene stocks filtrados por rating_to paginados
func (s *StockService) GetStocksByRatingTo(rating string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.GetByRatingTo(rating, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

func (s *StockService) GetStocksByRatingFrom(rating string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.GetByRatingFrom(rating, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

// Variables para almacenar el último sortBy y sortOrder
var lastSortBy string = "time"
var currentSortOrder string = "DESC"

func (s *StockService) GetSortedStocks(sortBy string, search string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)

	// Si el sortBy es el mismo que el anterior, invertimos el orden
	if sortBy == lastSortBy {
		if currentSortOrder == "DESC" {
			currentSortOrder = "ASC"
		} else {
			currentSortOrder = "DESC"
		}
	} else {
		// Si es un nuevo campo de ordenación, usamos DESC por defecto
		currentSortOrder = "DESC"
		lastSortBy = sortBy
	}

	// Validar que sortBy sea un campo válido para evitar SQL injection
	validFields := map[string]bool{
		"ticker": true, "company": true, "target_from": true, "target_to": true,
		"action": true, "brokerage": true, "rating_from": true, "rating_to": true, "time": true,
	}

	if !validFields[sortBy] {
		sortBy = "time" // Valor predeterminado seguro si el campo no es válido
		lastSortBy = sortBy
	}

	stocks, totalItems, totalPages, err := s.repo.GetSortedStocks(sortBy, currentSortOrder, search, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

// GetStocksByBrokerage obtiene stocks por brokerage con paginación
func (s *StockService) GetStocksByBrokerage(brokerage string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.GetByBrokerage(brokerage, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

// GetStocksByDateRange obtiene stocks por rango de fechas paginados
func (s *StockService) GetStocksByDateRange(startDate, endDate time.Time, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.GetByDateRange(startDate, endDate, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

// GetStocksByCompany obtiene stocks por compañía con paginación
func (s *StockService) GetStocksByCompany(company string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.GetByCompany(company, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
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
	newStocks := s.convertToModels(stockItems)

	// Verificar si hay cambios reales antes de actualizar la base de datos
	log.Printf("Verificando cambios en %d registros...", len(newStocks))

	// Crear un mapa de los nuevos stocks para facilitar la comparación
	newStocksMap := make(map[string]models.Stock)
	for _, stock := range newStocks {
		// Usar ticker+tiempo como clave única
		key := fmt.Sprintf("%s_%s", stock.Ticker, stock.Time.Format(time.RFC3339))
		newStocksMap[key] = stock
	}

	// Obtener todos los stocks existentes en los últimos X días
	// Usamos un período más largo para asegurarnos de capturar todos los posibles registros a comparar
	comparisonPeriod := 30 * 24 * time.Hour // 30 días
	since := time.Now().Add(-comparisonPeriod)
	existingStocks, err := s.repo.GetRecentRecommendations(since)
	if err != nil {
		return result, fmt.Errorf("error obteniendo datos existentes para comparación: %w", err)
	}

	// Crear mapa de stocks existentes
	existingStocksMap := make(map[string]models.Stock)
	for _, stock := range existingStocks {
		key := fmt.Sprintf("%s_%s", stock.Ticker, stock.Time.Format(time.RFC3339))
		existingStocksMap[key] = stock
	}

	// Identificar registros nuevos o modificados
	var stocksToUpsert []models.Stock
	for key, newStock := range newStocksMap {
		existingStock, exists := existingStocksMap[key]
		if !exists {
			// El registro es completamente nuevo
			stocksToUpsert = append(stocksToUpsert, newStock)
		} else {
			// El registro existe, verificar si hubo cambios en campos relevantes
			if s.stockHasChanged(existingStock, newStock) {
				stocksToUpsert = append(stocksToUpsert, newStock)
			}
		}
	}

	syncDate := time.Now()
	if len(stocksToUpsert) == 0 {
		log.Printf("No se detectaron cambios en los datos. No se realizarán actualizaciones.")

		// Completar el resultado
		result.TotalInserted = 0
		result.FailedInserts = 0
		result.UniqueTickersDB = len(existingStocksMap)

		// Calcular tiempo total
		totalDuration := time.Since(startTime)
		log.Printf("Sincronización completada en %.2f segundos. No hubo cambios.",
			totalDuration.Seconds())

		// Resumen en log
		s.logSyncSummary(result)

		return result, nil
	}

	// Solo actualizamos los registros que han cambiado o son nuevos
	log.Printf("Se detectaron %d registros nuevos o modificados de un total de %d",
		len(stocksToUpsert), len(newStocks))
	log.Printf("Iniciando actualización de %d registros...", len(stocksToUpsert))

	insertStartTime := time.Now()
	inserted, failedInserts, err := s.repo.UpsertStocksParallel(stocksToUpsert, syncDate)

	insertDuration := time.Since(insertStartTime)
	log.Printf("Operación de actualización completada en %.2f segundos",
		insertDuration.Seconds())

	if err != nil {
		return result, fmt.Errorf("error inserting/updating stocks: %w", err)
	}

	// Completar el resultado
	result.TotalInserted = inserted
	result.FailedInserts = len(failedInserts)
	result.UniqueTickersDB = len(existingStocksMap) + inserted - len(failedInserts)
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

// stockHasChanged compara dos stocks y determina si hay cambios relevantes
func (s *StockService) stockHasChanged(existingStock, newStock models.Stock) bool {
	// Comparar campos que pueden cambiar
	return existingStock.TargetFrom != newStock.TargetFrom ||
		existingStock.TargetTo != newStock.TargetTo ||
		existingStock.Action != newStock.Action ||
		existingStock.Brokerage != newStock.Brokerage ||
		existingStock.RatingFrom != newStock.RatingFrom ||
		existingStock.RatingTo != newStock.RatingTo ||
		existingStock.Company != newStock.Company
	// No comparamos ID, CreatedAt, UpdatedAt o DeletedAt
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
	log.Printf("Total elementos obtenidos de la API: %d", result.TotalProcessed)

	if result.TotalInserted > 0 {
		log.Printf("Total elementos actualizados/insertados: %d", result.TotalInserted)
	} else {
		log.Printf("No se detectaron cambios que requieran actualización")
	}

	if result.FailedInserts > 0 {
		log.Printf("Errores de inserción: %d", result.FailedInserts)
	} else {
		log.Printf("No hubo errores de inserción")
	}

	log.Printf("Tickers únicos en API: %d", result.UniqueTickersAPI)
	log.Printf("Total registros en base de datos: %d", result.UniqueTickersDB)

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
func (s *StockService) SearchStocks(query string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.SearchStocks(query, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

// GetStocksByPriceRange obtiene stocks por rango de precios paginados
func (s *StockService) GetStocksByPriceRange(minPrice, maxPrice string, page, pageSize int) ([]models.Stock, int, int, error) {
	page, pageSize = normalizePagination(page, pageSize)
	stocks, totalItems, totalPages, err := s.repo.GetByPriceRange(minPrice, maxPrice, page, pageSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return stocks, totalItems, totalPages, nil
}

// GetActionStats obtiene estadísticas por tipo de acción
func (s *StockService) GetActionStats() (map[string]int, error) {
	// Llama al método correspondiente del repositorio
	return s.repo.GetActionCounts()
}

// GetRatingStats obtiene estadísticas por tipo de rating
func (s *StockService) GetRatingStats() (map[string]int, error) {
	// Llama al método correspondiente del repositorio
	return s.repo.GetRatingCounts()
}

// --- Recommendation Logic ---

// isPositiveRating checks if a rating string is generally positive.
func isPositiveRating(rating string) bool {
	ratingLower := strings.ToLower(rating)
	positiveRatings := []string{"buy", "strong buy", "outperform", "overweight", "accumulate"}
	for _, pr := range positiveRatings {
		if strings.Contains(ratingLower, pr) {
			return true
		}
	}
	return false
}

// isNegativeRating checks if a rating string is generally negative.
func isNegativeRating(rating string) bool {
	ratingLower := strings.ToLower(rating)
	negativeRatings := []string{"sell", "strong sell", "underperform", "underweight", "reduce"}
	for _, nr := range negativeRatings {
		if strings.Contains(ratingLower, nr) {
			return true
		}
	}
	return false
}

// isUpgrade checks if the rating change represents an upgrade.
func isUpgrade(from, to string) bool {
	// Simple logic: Upgrade if 'to' is positive and 'from' is not positive (or is negative/neutral)
	return isPositiveRating(to) && !isPositiveRating(from)
}

// RecommendStocks generates stock recommendations based on recent data.
func (s *StockService) RecommendStocks(limit int) ([]models.Recommendation, error) {
	// 1. Define time window for recent data (e.g., last 7 days)
	recommendationPeriod := 7 * 24 * time.Hour
	since := time.Now().Add(-recommendationPeriod)

	// 2. Fetch recent stock data from repository
	recentStocks, err := s.repo.GetRecentRecommendations(since)
	if err != nil {
		return nil, fmt.Errorf("error fetching recent data for recommendations: %w", err)
	}

	if len(recentStocks) == 0 {
		log.Println("No recent stock data found for recommendations.")
		return []models.Recommendation{}, nil // Return empty slice, not an error
	}

	// 3. Calculate scores for each ticker
	scores := make(map[string]*models.Recommendation)
	reasons := make(map[string][]string) // Store reasons per ticker

	for _, stock := range recentStocks {
		if stock.Ticker == "" {
			continue // Skip entries without a ticker
		}

		currentScore := 0.0
		reasonFragments := []string{}

		// Score based on Action
		if strings.EqualFold(stock.Action, "Buy") {
			currentScore += 1.0
			reasonFragments = append(reasonFragments, "Buy Action")
		}

		// Score based on Rating To
		if isPositiveRating(stock.RatingTo) {
			currentScore += 1.0
			reasonFragments = append(reasonFragments, fmt.Sprintf("Positive Rating (%s)", stock.RatingTo))
		}

		// Score based on Upgrade
		if isUpgrade(stock.RatingFrom, stock.RatingTo) {
			currentScore += 0.5 // Bonus for upgrade
			reasonFragments = append(reasonFragments, fmt.Sprintf("Upgrade (%s -> %s)", stock.RatingFrom, stock.RatingTo))
		}

		// --- Potential Enhancements (Optional) ---
		// - Weight by brokerage reputation
		// - Weight by recency (newer data gets higher score)
		// - Factor in target price upside (requires current price data)
		// -----------------------------------------

		if currentScore > 0 {
			ticker := stock.Ticker
			if existing, ok := scores[ticker]; ok {
				// Update existing score and latest timestamp
				existing.Score += currentScore
				if stock.Time.After(existing.LastUpdate) {
					existing.LastUpdate = stock.Time
					existing.Company = stock.Company // Update company name too
				}
				reasons[ticker] = append(reasons[ticker], reasonFragments...)
			} else {
				// Create new recommendation entry
				scores[ticker] = &models.Recommendation{
					Ticker:     ticker,
					Company:    stock.Company,
					Score:      currentScore,
					LastUpdate: stock.Time,
					// Reason will be built later
				}
				reasons[ticker] = reasonFragments
			}
		}
	}

	// 4. Convert map to slice and build reason string
	recommendations := make([]models.Recommendation, 0, len(scores))
	for ticker, rec := range scores {
		// Build a concise reason string from unique fragments
		uniqueReasons := make(map[string]bool)
		for _, frag := range reasons[ticker] {
			uniqueReasons[frag] = true
		}
		var reasonParts []string
		for part := range uniqueReasons {
			reasonParts = append(reasonParts, part)
		}
		sort.Strings(reasonParts) // Consistent order
		rec.Reason = strings.Join(reasonParts, ", ")

		recommendations = append(recommendations, *rec)
	}

	// 5. Sort recommendations by score (descending)
	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].Score != recommendations[j].Score {
			return recommendations[i].Score > recommendations[j].Score
		}
		// Tie-breaker: more recent update first
		return recommendations[i].LastUpdate.After(recommendations[j].LastUpdate)
	})

	// 6. Apply limit
	if limit <= 0 {
		limit = 5 // Default limit
	}
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	log.Printf("Generated %d stock recommendations.", len(recommendations))
	return recommendations, nil
}
