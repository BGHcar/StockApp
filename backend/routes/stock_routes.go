package routes

import (
	"backend/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)


func StockRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET"},
		MaxAge: 300,
	}))
	

	r.Get("/api/sync", handlers.FetchAndStoreStock)
	r.Get("/api/stocks/all", handlers.GetAllStoreData)
	r.Get("/api/stocks/ticker/{ticker}", handlers.GetStoreByTicker)
	r.Get("/api/stocks/company/{company}", handlers.GetStoreByCompany)
	r.Get("/api/stocks/brokerage/{brokerage}", handlers.GetStoreByBrokerage)
	r.Get("/api/stocks/action/{action}", handlers.GetStoreByAction)
	r.Get("/api/stocks/rating-to/{rating}", handlers.GetStoreByRatingTo)
	r.Get("/api/stocks/rating-from/{rating}", handlers.GetStoreByRatingFrom)
	r.Get("/api/stocks/price-range/{min}/{max}", handlers.GetStoreByPrice)
	r.Get("/api/recommendations", handlers.GetStoreByRecommendation)
	r.Get("/api/sync/del",  handlers.DeleteTable)

	return r
}