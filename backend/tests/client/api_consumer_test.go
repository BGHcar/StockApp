package client_test

import (
	"backend/client"
	"bytes"
	"encoding/json"
	"errors"
	"io" // Importar io en lugar de ioutil
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient es un mock para http.Client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGet_Success(t *testing.T) {
	// Arrange
	mockClient := new(MockHTTPClient)
	consumer := &client.APIConsumer{
		BaseURL: "https://api.example.com",
		Client:  mockClient,
		Token:   "test-token",
	}

	// Mock API response
	now := time.Now()
	apiResp := client.Response{
		Items: []client.Item{
			{
				Ticker:     "AAPL",
				Company:    "Apple Inc.",
				TargetFrom: "$150",
				TargetTo:   "$180",
				Action:     "Buy",
				Brokerage:  "TestBrokerage",
				RatingFrom: "Hold",
				RatingTo:   "Buy",
				Time:       now,
			},
		},
		NextPage: "",
	}

	// Serialize expected response
	jsonData, _ := json.Marshal(apiResp)

	// Create mock response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(jsonData)), // Usar io.NopCloser en lugar de ioutil.NopCloser
	}

	// Set expectations
	mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResp, nil)

	// Act
	resp, err := consumer.Get("list", "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 1, len(resp.Items))
	assert.Equal(t, "AAPL", resp.Items[0].Ticker)
	assert.Equal(t, "Apple Inc.", resp.Items[0].Company)
	assert.Equal(t, now.Unix(), resp.Items[0].Time.Unix()) // Compare Unix timestamps to avoid time zone issues
	mockClient.AssertExpectations(t)
}

func TestGet_WithNextPage(t *testing.T) {
	// Arrange
	mockClient := new(MockHTTPClient)
	consumer := &client.APIConsumer{
		BaseURL: "https://api.example.com",
		Client:  mockClient,
		Token:   "test-token",
	}

	// Mock empty API response
	apiResp := client.Response{
		Items:    []client.Item{},
		NextPage: "",
	}

	// Serialize expected response
	jsonData, _ := json.Marshal(apiResp)

	// Create mock response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(jsonData)), // Cambiar aquí también
	}

	// Set expectations for request with next_page parameter
	mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == "https://api.example.com/list?next_page=page2"
	})).Return(mockResp, nil)

	// Act
	resp, err := consumer.Get("list", "page2")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 0, len(resp.Items))
	mockClient.AssertExpectations(t)
}

func TestGet_HTTPError(t *testing.T) {
	// Arrange
	mockClient := new(MockHTTPClient)
	consumer := &client.APIConsumer{
		BaseURL: "https://api.example.com",
		Client:  mockClient,
		Token:   "test-token",
	}

	// Set expectations
	mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(nil, errors.New("connection error"))

	// Act
	resp, err := consumer.Get("list", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "error making request")
	mockClient.AssertExpectations(t)
}

func TestGet_NonOkStatus(t *testing.T) {
	// Arrange
	mockClient := new(MockHTTPClient)
	consumer := &client.APIConsumer{
		BaseURL: "https://api.example.com",
		Client:  mockClient,
		Token:   "test-token",
	}

	// Create mock response with non-200 status
	mockResp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"error": "unauthorized"}`))), // Cambiar aquí también
	}

	// Set expectations
	mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResp, nil)

	// Act
	resp, err := consumer.Get("list", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "API returned non-200 status code")
	mockClient.AssertExpectations(t)
}

func TestGet_InvalidJSON(t *testing.T) {
	// Arrange
	mockClient := new(MockHTTPClient)
	consumer := &client.APIConsumer{
		BaseURL: "https://api.example.com",
		Client:  mockClient,
		Token:   "test-token",
	}

	// Create mock response with invalid JSON
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{invalid json}`))), // Cambiar aquí también
	}

	// Set expectations
	mockClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResp, nil)

	// Act
	resp, err := consumer.Get("list", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "error decoding response")
	mockClient.AssertExpectations(t)
}
