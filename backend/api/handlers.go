package api

import (
	"backend/factory"
	"encoding/json"
	"net/http"

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
	r.Get("/api/stocks/rating/{rating}", controller.GetStocksByRating)
	r.Get("/api/stats/actions", controller.GetActionStats)
	r.Get("/api/stats/ratings", controller.GetRatingStats)
	r.Get("/api/stats/count", controller.GetTotalCount)

	// Nueva ruta de búsqueda
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

func (c *StockController) GetStocksByRating(w http.ResponseWriter, r *http.Request) {
	rating := chi.URLParam(r, "rating")
	if rating == "" {
		// example: http://localhost:8080/api/stocks/rating/Buy
		http.Error(w, "Rating parameter is required", http.StatusBadRequest)
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stocks, err := service.GetStocksByRating(rating)
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
	if query == "" {
		http.Error(w, "Query parameter is required", http.StatusBadRequest)
		return
	}

	service, err := c.factory.CreateStockService()
	if err != nil {
		http.Error(w, "Error creating service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	stocks, err := service.SearchStocks(query)
	if err != nil {
		http.Error(w, "Error searching stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, stocks)
}
