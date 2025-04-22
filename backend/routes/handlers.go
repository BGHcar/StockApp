package routes

import (
	"backend/factory"
	"backend/models"
	"encoding/json"
	"log" // Import log
	"net/http"
	"strconv"
	"time" // Necesario para GetStocksByDateRange

	"github.com/go-chi/chi/v5"
)

// Controlador para los endpoints de la API
type StockController struct {
	factory *factory.ServiceFactory // Usar la interfaz si defines una, si no, el tipo concreto está bien
}

// NewStockController crea un nuevo controlador
func NewStockController(factory *factory.ServiceFactory) *StockController {
	return &StockController{
		factory: factory,
	}
}

// SetupRoutes configura las rutas para la API
func SetupRoutes(r chi.Router) {
	// Considera inyectar la factoría si es necesario para testing o flexibilidad
	controller := NewStockController(factory.NewServiceFactory())

	// Define the routes for the API
	r.Get("/api/sync", controller.SyncHandler)
	r.Get("/api/stocks", controller.GetAllStocks)
	r.Get("/api/stocks/ticker/{ticker}", controller.GetStockByTicker) // Cambiado para claridad
	r.Get("/api/stocks/action/{action}", controller.GetStocksByAction)
	r.Get("/api/stocks/rating-to/{rating}", controller.GetStocksByRatingTo)     // Cambiado
	r.Get("/api/stocks/rating-from/{rating}", controller.GetStocksByRatingFrom) // Cambiado
	r.Get("/api/stocks/company/{company}", controller.GetStocksByCompany)
	r.Get("/api/stocks/brokerage/{brokerage}", controller.GetStocksByBrokerage)
	r.Get("/api/stocks/date-range", controller.GetStocksByDateRange) // Usar query params
	r.Get("/api/stocks/price-range/{min}/{max}", controller.GetStocksByPriceRange)
	r.Get("/api/stocks/search", controller.SearchStocks) // Usar query params

	r.Get("/api/stats/actions", controller.GetActionStats)
	r.Get("/api/stats/ratings", controller.GetRatingStats)
	r.Get("/api/stats/count", controller.GetTotalCount)

	// --- Add Recommendation Route ---
	r.Get("/api/recommendations", controller.GetRecommendations) // New route
}

// --- Sync Handler (sin cambios necesarios para paginación) ---
func (c *StockController) SyncHandler(w http.ResponseWriter, r *http.Request) {
	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := service.SyncStockData()
	if err != nil {
		http.Error(w, "Error syncing data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	summary := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"processed_api":  result.TotalProcessed,
			"processed_db":   result.TotalInserted, // Renombrado para claridad
			"failed":         result.FailedInserts,
			"duplicates_api": result.DuplicateTickers,
		},
	}
	respondJSON(w, http.StatusOK, summary) // Añadir código de estado
}

// --- Estructura de Respuesta Paginada ---
type PaginatedResponse struct {
	Items      interface{} `json:"items"` // Usar interface{} para flexibilidad (o tipo específico si prefieres)
	Pagination struct {
		Page       int `json:"page"`
		PageSize   int `json:"pageSize"`
		TotalItems int `json:"totalItems"`
		TotalPages int `json:"totalPages"`
	} `json:"pagination"`
}

// --- Helpers ---
func getPaginationParams(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	defaultPageSize := 20
	maxPageSize := 100
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	} else if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil { // Evitar escribir cuerpo vacío si no hay datos
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// Log el error, pero es difícil enviar otro error al cliente aquí
			http.Error(w, `{"error":"Error encoding response"}`, http.StatusInternalServerError)
		}
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// --- Handlers Paginados ---

func (c *StockController) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	page, pageSize := getPaginationParams(r)
	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}

	stocks, totalItems, totalPages, err := service.GetAllStocks(page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response)
}

func (c *StockController) GetStockByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	page, pageSize := getPaginationParams(r)
	if ticker == "" {
		respondError(w, http.StatusBadRequest, "Ticker parameter is required")
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}

	stocks, totalItems, totalPages, err := service.GetStockByTicker(ticker, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks by ticker: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response)
}

func (c *StockController) GetStocksByAction(w http.ResponseWriter, r *http.Request) {
	action := chi.URLParam(r, "action")
	page, pageSize := getPaginationParams(r)
	if action == "" {
		respondError(w, http.StatusBadRequest, "Action parameter is required")
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}

	stocks, totalItems, totalPages, err := service.GetStocksByAction(action, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks by action: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response)
}

func (c *StockController) GetStocksByRatingTo(w http.ResponseWriter, r *http.Request) {
	rating := chi.URLParam(r, "rating")
	page, pageSize := getPaginationParams(r)
	if rating == "" {
		respondError(w, http.StatusBadRequest, "Rating parameter is required")
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}

	stocks, totalItems, totalPages, err := service.GetStocksByRatingTo(rating, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks by rating_to: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response)
}

func (c *StockController) GetStocksByRatingFrom(w http.ResponseWriter, r *http.Request) {
	rating := chi.URLParam(r, "rating")
	page, pageSize := getPaginationParams(r)
	if rating == "" {
		respondError(w, http.StatusBadRequest, "Rating parameter is required")
		return
	}
	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}
	stocks, totalItems, totalPages, err := service.GetStocksByRatingFrom(rating, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks by rating_from: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response)
}

func (c *StockController) GetStocksByCompany(w http.ResponseWriter, r *http.Request) {
	company := chi.URLParam(r, "company")
	page, pageSize := getPaginationParams(r)
	if company == "" {
		respondError(w, http.StatusBadRequest, "Company parameter is required")
		return
	}
	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}
	stocks, totalItems, totalPages, err := service.GetStocksByCompany(company, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks by company: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response)
}

func (c *StockController) GetStocksByBrokerage(w http.ResponseWriter, r *http.Request) {
	brokerage := chi.URLParam(r, "brokerage")
	page, pageSize := getPaginationParams(r)
	if brokerage == "" {
		respondError(w, http.StatusBadRequest, "Brokerage parameter is required")
		return
	}
	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}
	stocks, totalItems, totalPages, err := service.GetStocksByBrokerage(brokerage, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks by brokerage: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response)
}

func (c *StockController) GetStocksByDateRange(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("startDate") // e.g., 2023-10-26
	endDateStr := r.URL.Query().Get("endDate")     // e.g., 2023-10-27
	page, pageSize := getPaginationParams(r)

	if startDateStr == "" || endDateStr == "" {
		respondError(w, http.StatusBadRequest, "startDate and endDate query parameters are required (YYYY-MM-DD)")
		return
	}

	// Parsear fechas (considera añadir validación más robusta)
	layout := "2006-01-02"
	startDate, errStart := time.Parse(layout, startDateStr)
	endDate, errEnd := time.Parse(layout, endDateStr)

	if errStart != nil || errEnd != nil {
		respondError(w, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD.")
		return
	}
	// Ajustar endDate para incluir todo el día
	endDate = endDate.Add(24*time.Hour - 1*time.Nanosecond)

	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}

	stocks, totalItems, totalPages, err := service.GetStocksByDateRange(startDate, endDate, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks by date range: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response)
}

func (c *StockController) GetStocksByPriceRange(w http.ResponseWriter, r *http.Request) {
	minPrice := chi.URLParam(r, "min")
	maxPrice := chi.URLParam(r, "max")
	page, pageSize := getPaginationParams(r)

	if minPrice == "" || maxPrice == "" {
		respondError(w, http.StatusBadRequest, "min and max path parameters are required")
		return
	}

	// Validar si son números (opcional, depende de cómo los maneje el repo)
	// _, errMin := strconv.ParseFloat(minPrice, 64)
	// _, errMax := strconv.ParseFloat(maxPrice, 64)
	// if errMin != nil || errMax != nil {
	// 	respondError(w, http.StatusBadRequest, "min and max must be valid numbers")
	// 	return
	// }

	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}

	stocks, totalItems, totalPages, err := service.GetStocksByPriceRange(minPrice, maxPrice, page, pageSize)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching stocks by price range: "+err.Error())
		return
	}

	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response) // Devolver la estructura completa
}

func (c *StockController) SearchStocks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query") // Obtener de query param
	page, pageSize := getPaginationParams(r)

	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}

	var stocks []models.Stock
	var totalItems, totalPages int

	if query == "" {
		// Si la consulta está vacía, devolver todos los stocks paginados
		stocks, totalItems, totalPages, err = service.GetAllStocks(page, pageSize)
	} else {
		// Realizar la búsqueda
		stocks, totalItems, totalPages, err = service.SearchStocks(query, page, pageSize)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error searching stocks: "+err.Error())
		return
	}

	// Construir y devolver la respuesta paginada en ambos casos
	response := PaginatedResponse{Items: stocks}
	response.Pagination.Page = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalItems = totalItems
	response.Pagination.TotalPages = totalPages

	respondJSON(w, http.StatusOK, response) // Devolver la estructura completa
}

// --- Handlers de Estadísticas (sin cambios necesarios para paginación) ---

func (c *StockController) GetActionStats(w http.ResponseWriter, r *http.Request) {
	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}
	stats, err := service.GetActionStats()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching action statistics: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, stats)
}

func (c *StockController) GetRatingStats(w http.ResponseWriter, r *http.Request) {
	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}
	stats, err := service.GetRatingStats()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching rating statistics: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, stats)
}

func (c *StockController) GetTotalCount(w http.ResponseWriter, r *http.Request) {
	service, err := c.factory.CreateStockService()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}
	count, err := service.GetTotalCount()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error fetching total count: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]int{"count": count})
}

// --- Recommendation Handler ---

func (c *StockController) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// Get optional limit from query params
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5 // Default limit if not specified or invalid
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		log.Printf("Error creating service for recommendations: %v", err) // Log error
		respondError(w, http.StatusInternalServerError, "Error creating service: "+err.Error())
		return
	}

	recommendations, err := service.RecommendStocks(limit)
	if err != nil {
		log.Printf("Error generating recommendations: %v", err) // Log error
		respondError(w, http.StatusInternalServerError, "Error generating recommendations: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, recommendations)
}

// --- Helper parsePrice (si es necesario, sin cambios) ---
func parsePrice(price string) (float64, error) {
	// Eliminar el símbolo de moneda si existe
	priceStr := price
	if len(priceStr) > 0 && (priceStr[0] == '$') {
		priceStr = priceStr[1:]
	}

	// Intentar convertir a float
	return strconv.ParseFloat(priceStr, 64)
}
