package routes

import (
	"backend/factory"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Controlador para los endpoints de la API
type StockController struct {
	factory *factory.ServiceFactory
}

// NewStockController crea un nuevo controlador
func NewStockController(factory *factory.ServiceFactory) *StockController {
	return &StockController{
		factory: factory,
	}
}

// SetupRoutes configura las rutas para la API
func SetupRoutes(r chi.Router) {
	controller := NewStockController(factory.NewServiceFactory())

	// Define the routes for the API
	r.Get("/api/sync", controller.SyncHandler)
	r.Get("/api/stocks", controller.GetAllStocks)
	r.Get("/api/stocks/{ticker}", controller.GetStockByTicker)
	r.Get("/api/stocks/action/{action}", controller.GetStocksByAction)
	r.Get("/api/stocks/ratingto/{rating}", controller.GetStocksByRatingTo)
	r.Get("/api/stocks/ratingfrom/{rating}", controller.GetStocksByRatingFrom)
	r.Get("/api/stocks/company/{company}", controller.GetStocksByCompany)
	r.Get("/api/stats/actions", controller.GetActionStats)
	r.Get("/api/stats/ratings", controller.GetRatingStats)
	r.Get("/api/stats/count", controller.GetTotalCount)
	r.Get("/api/stocks/brokerage/{brokerage}", controller.GetStocksByBrokerage)
	r.Get("/api/stocks/price-range/{min}/{max}", controller.GetStocksByPriceRange) // Nueva ruta
	r.Get("/api/stocks/search/{query}", controller.SearchStocks)
}

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

	// Devolver un resumen estructurado
	summary := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"processed":        result.TotalProcessed,
			"inserted":         result.TotalInserted,
			"failed":           result.FailedInserts,
			"duplicates_found": result.DuplicateTickers > 0,
		},
	}

	respondJSON(w, summary)
}

func (c *StockController) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stocks, err := service.GetAllStocks()
	if err != nil {
		http.Error(w, "Error fetching stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stocks)
}

func (c *StockController) GetStockByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "Ticker parameter is required", http.StatusBadRequest)
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stock, err := service.GetStockByTicker(ticker)
	if err != nil {
		http.Error(w, "Error fetching stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if stock == nil {
		http.Error(w, "Stock not found", http.StatusNotFound)
		return
	}

	respondJSON(w, stock)
}

func (c *StockController) GetStocksByAction(w http.ResponseWriter, r *http.Request) {
	action := chi.URLParam(r, "action")
	if action == "" {
		http.Error(w, "Action parameter is required", http.StatusBadRequest)
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stocks, err := service.GetStocksByAction(action)
	if err != nil {
		http.Error(w, "Error fetching stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stocks)
}

func (c *StockController) GetStocksByRatingTo(w http.ResponseWriter, r *http.Request) {
	rating := chi.URLParam(r, "rating")
	if rating == "" {
		http.Error(w, "Rating parameter is required", http.StatusBadRequest)
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stocks, err := service.GetStocksByRatingTo(rating)
	if err != nil {
		http.Error(w, "Error fetching stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	respondJSON(w, stocks)
}

func (c *StockController) GetStocksByRatingFrom(w http.ResponseWriter, r *http.Request){
	rating := chi.URLParam(r, "rating")
	if rating == "" {
		http.Error(w, "Rating parameter is required", http.StatusBadRequest)
		return
	}
	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	stocks, err := service.GetStocksByRatingFrom(rating)
	if err != nil {
		http.Error(w, "Error fetching stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stocks)
}

func (c *StockController) GetStocksByCompany(w http.ResponseWriter, r *http.Request) {
	company := chi.URLParam(r, "company")
	if company == "" {
		http.Error(w, "Company parameter is required", http.StatusBadRequest)
		return
	}
	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	stocks, err := service.GetStocksByCompany(company)
	if err != nil {
		http.Error(w, "Error fetching stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stocks)
}

func (c *StockController) GetStocksByBrokerage(w http.ResponseWriter, r *http.Request) {
	brokerage := chi.URLParam(r, "brokerage")
	if brokerage == "" {
		http.Error(w, "Brokerage parameter is required", http.StatusBadRequest)
		return
	}
	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}
	stocks, err := service.GetStocksByBrokerage(brokerage)
	if err != nil {
		http.Error(w, "Error fetching stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stocks)
}

func (c *StockController) GetActionStats(w http.ResponseWriter, r *http.Request) {
	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stats, err := service.GetActionStats()
	if err != nil {
		http.Error(w, "Error fetching statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stats)
}

func (c *StockController) GetRatingStats(w http.ResponseWriter, r *http.Request) {
	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stats, err := service.GetRatingStats()
	if err != nil {
		http.Error(w, "Error fetching statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stats)
}

func (c *StockController) GetTotalCount(w http.ResponseWriter, r *http.Request) {
	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := service.GetTotalCount()
	if err != nil {
		http.Error(w, "Error fetching count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]int{"count": count})
}

// Helper function to respond with JSON
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *StockController) SearchStocks(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")

	// Si la consulta está vacía, devolver todos los stocks
	if query == "" {
		service, err := c.factory.CreateStockService()
		if err != nil {
			http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
			return
		}

		stocks, err := service.GetAllStocks()
		if err != nil {
			http.Error(w, "Error fetching stocks: "+err.Error(), http.StatusInternalServerError)
			return
		}

		respondJSON(w, stocks)
		return
	}

	// Resto del código para búsquedas regulares...
}

func (c *StockController) GetStocksByPriceRange(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la URL
	minPrice := chi.URLParam(r, "min")
	maxPrice := chi.URLParam(r, "max")

	if minPrice == "" || maxPrice == "" {
		http.Error(w, "Se requieren los parámetros min y max", http.StatusBadRequest)
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir precios a float para validación (aunque el servicio maneje strings)
	_, errMin := parsePrice(minPrice)
	_, errMax := parsePrice(maxPrice)

	if errMin != nil || errMax != nil {
		http.Error(w, "Los precios deben ser valores numéricos válidos", http.StatusBadRequest)
		return
	}

	stocks, err := service.GetStocksByPriceRange(minPrice, maxPrice)
	if err != nil {
		http.Error(w, "Error fetching stocks by price range: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stocks)
}

// Helper para parsear precios
func parsePrice(price string) (float64, error) {
	// Eliminar el símbolo de moneda si existe
	priceStr := price
	if len(priceStr) > 0 && (priceStr[0] == '$') {
		priceStr = priceStr[1:]
	}

	// Intentar convertir a float
	return strconv.ParseFloat(priceStr, 64)
}
