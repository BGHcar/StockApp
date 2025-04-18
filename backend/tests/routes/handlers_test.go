package routes_test

import (
	"backend/interfaces"
	"backend/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStockService es un mock para StockService
type MockStockService struct {
	mock.Mock
}

func (m *MockStockService) GetAllStocks() ([]models.Stock, error) {
	args := m.Called()
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockService) GetTotalCount() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *MockStockService) GetStockByTicker(ticker string) ([]models.Stock, error) {
	args := m.Called(ticker)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockService) GetStocksByAction(action string) ([]models.Stock, error) {
	args := m.Called(action)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockService) GetStocksByRating(rating string) ([]models.Stock, error) {
	args := m.Called(rating)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockService) GetStocksByCompany(company string) ([]models.Stock, error) {
	args := m.Called(company)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockService) GetActionStats() (map[string]int, error) {
	args := m.Called()
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockStockService) GetRatingStats() (map[string]int, error) {
	args := m.Called()
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockStockService) SyncStockData() (interfaces.SyncResult, error) {
	args := m.Called()
	return args.Get(0).(interfaces.SyncResult), args.Error(1)
}

func (m *MockStockService) SearchStocks(query string) ([]models.Stock, error) {
	args := m.Called(query)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockService) GetStocksByBrokerage(brokerage string) ([]models.Stock, error) {
	args := m.Called(brokerage)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockService) GetStocksByDateRange(startDate, endDate time.Time) ([]models.Stock, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockService) GetStocksByPriceRange(minPrice, maxPrice string) ([]models.Stock, error) {
	args := m.Called(minPrice, maxPrice)
	return args.Get(0).([]models.Stock), args.Error(1)
}

// MockServiceFactory es un mock para ServiceFactory
type MockServiceFactory struct {
	mock.Mock
}

func (m *MockServiceFactory) CreateStockService() (interfaces.StockService, error) {
	args := m.Called()
	return args.Get(0).(interfaces.StockService), args.Error(1)
}

// Añadir esta nueva estructura que envuelve tu mock
type MockFactoryWrapper struct {
	Mock *MockServiceFactory
}

// Implementa el mismo método que factory.ServiceFactory para hacerlo compatible
func (m *MockFactoryWrapper) CreateStockService() (interfaces.StockService, error) {
	return m.Mock.CreateStockService()
}

func TestGetAllStocksHandler(t *testing.T) {
	// Arrange
	mockService := new(MockStockService)

	// Crear datos de prueba
	expectedStocks := []models.Stock{
		{Ticker: "AAPL", Company: "Apple Inc."},
		{Ticker: "MSFT", Company: "Microsoft Corporation"},
	}

	// Configurar expectativas
	mockService.On("GetAllStocks").Return(expectedStocks, nil)

	// Crear un handler para la prueba
	handler := func(w http.ResponseWriter, r *http.Request) {
		stocks, err := mockService.GetAllStocks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stocks)
	}

	// Crear request y recorder
	req := httptest.NewRequest("GET", "/api/stocks", nil)
	rr := httptest.NewRecorder()

	// Act
	http.HandlerFunc(handler).ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseStocks []models.Stock
	err := json.Unmarshal(rr.Body.Bytes(), &responseStocks)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(responseStocks))
	assert.Equal(t, "AAPL", responseStocks[0].Ticker)
	assert.Equal(t, "MSFT", responseStocks[1].Ticker)

	mockService.AssertExpectations(t)
}

func TestGetStockByTickerHandler(t *testing.T) {
	// Arrange
	mockService := new(MockStockService)

	// Declarar las variables primero
	ticker := "AAPL"
	stocks := []models.Stock{
		{Ticker: ticker, Company: "Apple Inc."},
	}

	// Set expectations
	mockService.On("GetStockByTicker", ticker).Return(stocks, nil)

	// Crear un handler para la prueba - similar a TestGetAllStocksHandler
	handler := func(w http.ResponseWriter, r *http.Request) {
		tickerParam := chi.URLParam(r, "ticker")
		stocks, err := mockService.GetStockByTicker(tickerParam)
		if err != nil {
			http.Error(w, "Error fetching stock: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stocks)
	}

	// Setup router with chi para obtener los parámetros de URL
	r := chi.NewRouter()
	r.Get("/{ticker}", handler)

	// Create request
	req, err := http.NewRequest("GET", "/"+ticker, nil)
	assert.NoError(t, err)

	// Create recorder
	rr := httptest.NewRecorder()

	// Act: call the router
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify response body
	var responseStocks []models.Stock
	err = json.Unmarshal(rr.Body.Bytes(), &responseStocks)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(responseStocks))
	assert.Equal(t, ticker, responseStocks[0].Ticker)

	mockService.AssertExpectations(t)
}

func TestGetStockByTickerHandler_Error(t *testing.T) {
	// Arrange
	mockService := new(MockStockService)

	// Mock data
	ticker := "AAPL"
	expectedError := errors.New("database error")

	// Set expectations
	mockService.On("GetStockByTicker", ticker).Return([]models.Stock{}, expectedError)

	// Usar el mismo enfoque que en TestGetStockByTickerHandler
	handler := func(w http.ResponseWriter, r *http.Request) {
		tickerParam := chi.URLParam(r, "ticker")
		stocks, err := mockService.GetStockByTicker(tickerParam)
		if err != nil {
			http.Error(w, "Error fetching stock: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stocks)
	}

	// Setup router con chi
	r := chi.NewRouter()
	r.Get("/{ticker}", handler)

	// Create request
	req, err := http.NewRequest("GET", "/"+ticker, nil)
	assert.NoError(t, err)

	// Create recorder
	rr := httptest.NewRecorder()

	// Act: call the router
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error fetching stock")

	mockService.AssertExpectations(t)
}

func TestGetStocksByPriceRangeHandler(t *testing.T) {
	// Arrange
	mockService := new(MockStockService)

	// Mock data
	minPrice := "100"
	maxPrice := "200"
	stocks := []models.Stock{
		{Ticker: "AAPL", Company: "Apple Inc.", TargetFrom: "$150", TargetTo: "$180"},
	}

	// Set expectations
	mockService.On("GetStocksByPriceRange", minPrice, maxPrice).Return(stocks, nil)

	// Crear un handler para la prueba
	handler := func(w http.ResponseWriter, r *http.Request) {
		minParam := chi.URLParam(r, "min")
		maxParam := chi.URLParam(r, "max")
		stocks, err := mockService.GetStocksByPriceRange(minParam, maxParam)
		if err != nil {
			http.Error(w, "Error fetching stocks: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stocks)
	}

	// Setup router with chi
	r := chi.NewRouter()
	r.Get("/{min}/{max}", handler)

	// Create request
	req, err := http.NewRequest("GET", "/"+minPrice+"/"+maxPrice, nil)
	assert.NoError(t, err)

	// Create recorder
	rr := httptest.NewRecorder()

	// Act: call the router
	r.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify response body
	var responseStocks []models.Stock
	err = json.Unmarshal(rr.Body.Bytes(), &responseStocks)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(responseStocks))
	assert.Equal(t, "AAPL", responseStocks[0].Ticker)
	assert.Equal(t, "$150", responseStocks[0].TargetFrom)
	assert.Equal(t, "$180", responseStocks[0].TargetTo)

	mockService.AssertExpectations(t)
}
