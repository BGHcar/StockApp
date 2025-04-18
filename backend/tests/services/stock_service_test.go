package services_test

import (
	"backend/interfaces"
	"backend/models"
	"backend/services"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStockRepository es un mock para StockRepository
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) GetAll() ([]models.Stock, error) {
	args := m.Called()
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStockRepository) GetByTicker(ticker string) ([]models.Stock, error) {
	args := m.Called(ticker)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetByAction(action string) ([]models.Stock, error) {
	args := m.Called(action)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetByRating(rating string) ([]models.Stock, error) {
	args := m.Called(rating)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetByBrokerage(brokerage string) ([]models.Stock, error) {
	args := m.Called(brokerage)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Stock, error) {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]models.Stock), args.Error(1) // Añadir esta línea que falta
}

func (m *MockStockRepository) GetByCompany(company string) ([]models.Stock, error) {
	args := m.Called(company)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetActionCounts() (map[string]int, error) {
	args := m.Called()
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockStockRepository) GetRatingCounts() (map[string]int, error) {
	args := m.Called()
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockStockRepository) InsertStock(stock models.Stock) error {
	args := m.Called(stock)
	return args.Error(0)
}

func (m *MockStockRepository) InsertStocks(stocks []models.Stock) (int, map[string]string, error) {
	args := m.Called(stocks)
	return args.Int(0), args.Get(1).(map[string]string), args.Error(2)
}

func (m *MockStockRepository) InsertStocksParallel(stocks []models.Stock) (int, map[string]string, error) {
	args := m.Called(stocks)
	return args.Int(0), args.Get(1).(map[string]string), args.Error(2)
}

func (m *MockStockRepository) TruncateTable() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStockRepository) SearchStocks(query string) ([]models.Stock, error) {
	args := m.Called(query)
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetByPriceRange(minPrice, maxPrice string) ([]models.Stock, error) {
	args := m.Called(minPrice, maxPrice)
	return args.Get(0).([]models.Stock), args.Error(1)
}

// MockAPIClient es un mock para APIClient
type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) Get(endpoint string, nextPage string) (*interfaces.APIResponse, error) {
	args := m.Called(endpoint, nextPage)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.APIResponse), args.Error(1)
}

func TestGetAllStocks(t *testing.T) {
	// Arrange
	mockRepo := new(MockStockRepository)
	mockAPI := new(MockAPIClient)
	service := services.NewStockService(mockRepo, mockAPI)

	expectedStocks := []models.Stock{
		{Ticker: "AAPL", Company: "Apple Inc."},
		{Ticker: "MSFT", Company: "Microsoft Corporation"},
	}

	mockRepo.On("GetAll").Return(expectedStocks, nil)

	// Act
	stocks, err := service.GetAllStocks()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStocks, stocks)
	mockRepo.AssertExpectations(t)
}

func TestGetStockByTicker(t *testing.T) {
	// Arrange
	mockRepo := new(MockStockRepository)
	mockAPI := new(MockAPIClient)
	service := services.NewStockService(mockRepo, mockAPI)

	ticker := "AAPL"
	expectedStocks := []models.Stock{{Ticker: ticker, Company: "Apple Inc."}}

	mockRepo.On("GetByTicker", ticker).Return(expectedStocks, nil)

	// Act
	stocks, err := service.GetStockByTicker(ticker)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStocks, stocks)
	mockRepo.AssertExpectations(t)
}

func TestGetTotalCount(t *testing.T) {
	// Arrange
	mockRepo := new(MockStockRepository)
	mockAPI := new(MockAPIClient)
	service := services.NewStockService(mockRepo, mockAPI)

	expectedCount := int64(42)

	mockRepo.On("GetCount").Return(expectedCount, nil)

	// Act
	count, err := service.GetTotalCount()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int(expectedCount), count)
	mockRepo.AssertExpectations(t)
}

func TestGetStocksByPriceRange(t *testing.T) {
	// Arrange
	mockRepo := new(MockStockRepository)
	mockAPI := new(MockAPIClient)
	service := services.NewStockService(mockRepo, mockAPI)

	minPrice := "100"
	maxPrice := "200"
	expectedStocks := []models.Stock{
		{Ticker: "AAPL", Company: "Apple Inc.", TargetFrom: "$150", TargetTo: "$180"},
		{Ticker: "MSFT", Company: "Microsoft Corporation", TargetFrom: "$120", TargetTo: "$190"},
	}

	mockRepo.On("GetByPriceRange", minPrice, maxPrice).Return(expectedStocks, nil)

	// Act
	stocks, err := service.GetStocksByPriceRange(minPrice, maxPrice)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedStocks, stocks)
	mockRepo.AssertExpectations(t)
}

func TestSyncStockData_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockStockRepository)
	mockAPI := new(MockAPIClient)
	service := services.NewStockService(mockRepo, mockAPI)

	// Mock API response
	mockResponse := &interfaces.APIResponse{
		Items: []interfaces.StockItem{
			{
				Ticker:  "AAPL",
				Company: "Apple Inc.",
				Time:    time.Now(),
			},
		},
		NextPage: "",
	}

	// Mock repository behavior
	mockRepo.On("TruncateTable").Return(nil)
	mockRepo.On("InsertStocksParallel", mock.AnythingOfType("[]models.Stock")).Return(1, map[string]string{}, nil)
	mockAPI.On("Get", "list", "").Return(mockResponse, nil)

	// Act
	result, err := service.SyncStockData()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, result.TotalProcessed)
	assert.Equal(t, 1, result.TotalInserted)
	mockRepo.AssertExpectations(t)
	mockAPI.AssertExpectations(t)
}

func TestSyncStockData_APIError(t *testing.T) {
	// Arrange
	mockRepo := new(MockStockRepository)
	mockAPI := new(MockAPIClient)
	service := services.NewStockService(mockRepo, mockAPI)

	// Mock repository behavior
	mockRepo.On("TruncateTable").Return(nil)
	mockAPI.On("Get", "list", "").Return(nil, errors.New("API error"))

	// Act
	_, err := service.SyncStockData()

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error collecting API data")
	mockRepo.AssertExpectations(t)
	mockAPI.AssertExpectations(t)
}
