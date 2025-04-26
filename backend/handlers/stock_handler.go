package handlers

import (
	"backend/api"
	"backend/db"
	"backend/repositories"
	"backend/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func FetchAndStoreStock(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/sync")
	total, err := api.FetchData()
	if err != nil{
		http.Error(w, "failed to fetch data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message": "data fetched and stored",
		"total_items": total,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}


func GetAllStoreData(w http.ResponseWriter, r *http.Request){	

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
    if page <= 0 {
      page = 1
    }

    pageSize, _ := strconv.Atoi(q.Get("page_size"))
    switch {
    case pageSize > 100:
      pageSize = 100
    case pageSize <= 0:
      pageSize = 20
    }

	items, newpage, newpageSize, totalItems, err := repositories.GetAll(page, pageSize)
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface {}{
		"items": items,
		"pagination": map[string]interface{}{
			"page": newpage,
			"pageSize": newpageSize,
			"totalItems": totalItems,
			"totalPages": (totalItems/pageSize),

		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}

func DeleteTable(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/sync/del")

	if err := db.Drop(); err != nil{
		http.Error(w, "failed to delete de table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface {}{
		"message": "table has been deleted",
		"data": "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}

func GetStoreByTicker(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/stocks/ticker")

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
    if page <= 0 {
      page = 1
    }

    pageSize, _ := strconv.Atoi(q.Get("page_size"))
    switch {
    case pageSize > 100:
      pageSize = 100
    case pageSize <= 0:
      pageSize = 20
    }

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker parameter is required",  http.StatusInternalServerError)
	}

	items, newpage, newpageSize, totalItems, err := repositories.GetByTicker(ticker, page, pageSize)
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := float64(totalItems/pageSize)

	if totalItems%pageSize != 0 {
		totalPages += 1
	}


	resp := map[string]interface {}{
		"items": items,
		"pagination": map[string]interface{}{
			"page": newpage,
			"pageSize": newpageSize,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}


func GetStoreByCompany(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/stocks/company")

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
	  page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
	  pageSize = 100
	case pageSize <= 0:
	  pageSize = 20
	}

	company := chi.URLParam(r, "company")
	if company == "" {
		http.Error(w, "company parameter is required",  http.StatusInternalServerError)
	}

	items, newpage, newpageSize, totalItems, err := repositories.GetByCompany(company, page, pageSize)
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := float64(totalItems/pageSize)

	if totalItems%pageSize != 0 {
		totalPages += 1
	}

	resp := map[string]interface {}{
		"items": items,
		"pagination": map[string]interface{}{
			"page": newpage,
			"pageSize": newpageSize,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}
/* 
r.Get("/api/stocks/brokerage/{brokerage}", handlers.GetStoreByBrokerage)
r.Get("/api/stocks/action/{action}", handlers.GetStoreByAction)
r.Get("/api/stocks/rating-to/{rating}", handlers.GetStoreByRatingTo)
r.Get("/api/stocks/rating-from/{rating}", handlers.GetStoreByRatingFrom)
r.Get("/api/stocks/price-to/{price}", handlers.GetStoreByPriceTo) */

func GetStoreByBrokerage(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/stocks/brokerage")

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
	  page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
	  pageSize = 100
	case pageSize <= 0:
	  pageSize = 20
	}

	brokerage := chi.URLParam(r, "brokerage")
	if brokerage == "" {
		http.Error(w, "brokerage parameter is required",  http.StatusInternalServerError)
	}

	items, newpage, newpageSize, totalItems, err := repositories.GetByBrokerage(brokerage, page, pageSize)
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := float64(totalItems/pageSize)

	if totalItems%pageSize != 0 {
		totalPages += 1
	}

	resp := map[string]interface {}{
		"items": items,
		"pagination": map[string]interface{}{
			"page": newpage,
			"pageSize": newpageSize,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}

func GetStoreByAction(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/stocks/action")

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
	  page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
	  pageSize = 100
	case pageSize <= 0:
	  pageSize = 20
	}

	action := chi.URLParam(r, "action")
	if action == "" {
		http.Error(w, "action parameter is required",  http.StatusInternalServerError)
	}

	items, newpage, newpageSize, totalItems, err := repositories.GetByAction(action, page, pageSize)
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := float64(totalItems/pageSize)

	if totalItems%pageSize != 0 {
		totalPages += 1
	}

	resp := map[string]interface {}{
		"items": items,
		"pagination": map[string]interface{}{
			"page": newpage,
			"pageSize": newpageSize,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}

func GetStoreByRatingTo(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/stocks/rating-to")

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
	  page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
	  pageSize = 100
	case pageSize <= 0:
	  pageSize = 20
	}

	rating := chi.URLParam(r, "rating")
	if rating == "" {
		http.Error(w, "rating parameter is required",  http.StatusInternalServerError)
	}

	items, newpage, newpageSize, totalItems, err := repositories.GetByRatingTo(rating, page, pageSize)
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := float64(totalItems/pageSize)

	if totalItems%pageSize != 0 {
		totalPages += 1
	}

	resp := map[string]interface {}{
		"items": items,
		"pagination": map[string]interface{}{
			"page": newpage,
			"pageSize": newpageSize,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}

func GetStoreByRatingFrom(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/stocks/rating-from")

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
	  page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
	  pageSize = 100
	case pageSize <= 0:
	  pageSize = 20
	}

	rating := chi.URLParam(r, "rating")
	if rating == "" {
		http.Error(w, "rating parameter is required",  http.StatusInternalServerError)
	}

	items, newpage, newpageSize, totalItems, err := repositories.GetByRatingFrom(rating, page, pageSize)
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := float64(totalItems/pageSize)

	if totalItems%pageSize != 0 {
		totalPages += 1
	}

	resp := map[string]interface {}{
		"items": items,
		"pagination": map[string]interface{}{
			"page": newpage,
			"pageSize": newpageSize,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}

func GetStoreByPrice(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/stocks/price-range")

	q := r.URL.Query()

	minPrice := chi.URLParam(r, "min")
	maxPrice := chi.URLParam(r, "max")
	if minPrice == "" || maxPrice == "" {
		http.Error(w, "min and max price parameters are required",  http.StatusInternalServerError)
		return
	}
	min, err := strconv.ParseFloat(minPrice, 64)
	if err != nil {
		http.Error(w, "invalid min price: "+ err.Error(), http.StatusInternalServerError)
		return
	}
	max, err := strconv.ParseFloat(maxPrice, 64)
	if err != nil {
		http.Error(w, "invalid max price: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
	  page = 1
	}

	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	switch {
	case pageSize > 100:
	  pageSize = 100
	case pageSize <= 0:
	  pageSize = 20
	}


	items, newpage, newpageSize, totalItems, err := repositories.GetByPrice(min, max, page, pageSize)
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := float64(totalItems/pageSize)

	if totalItems%pageSize != 0 {
		totalPages += 1
	}

	resp := map[string]interface {}{
		"items": items,
		"pagination": map[string]interface{}{
			"page": newpage,
			"pageSize": newpageSize,
			"totalItems": totalItems,
			"totalPages": totalPages,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}

//recommendations always return 5 items and these are not paginated
func GetStoreByRecommendation(w http.ResponseWriter, r *http.Request){
	fmt.Println("received request for /api/recommendations")

	items, err := services.GetRecommendationsService()
	if err != nil {
		http.Error(w, "failed to get data: "+ err.Error(), http.StatusInternalServerError)
		return
	}

	resp := items
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	fmt.Println(r.Body)
}

