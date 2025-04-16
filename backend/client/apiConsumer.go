// Project: go-api-consumer
// File: backend/client/apiConsumer.go
// Description: This file contains the APIConsumer struct and methods for making HTTP requests to an external API.
package client

import (
	"backend/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Item representa un elemento de acción individual devuelto por la API
type Item struct {
	Ticker     string    `json:"ticker"`
	Company    string    `json:"company"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

// Response representa la estructura de la respuesta completa de la API
type Response struct {
	Items    []Item `json:"items"`
	NextPage string `json:"next_page"`
}

// Aseguramos que APIConsumer implementa la interfaz APIClient
var _ interfaces.APIClient = (*APIConsumer)(nil)

// APIConsumer es un cliente para consumir la API de stocks
type APIConsumer struct {
	BaseURL string
	Client  HTTPClient
	Token   string
}

// HTTPClient define una interfaz para operaciones HTTP
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewAPIConsumer crea una nueva instancia del consumidor de API
func NewAPIConsumer(baseURL string) *APIConsumer {
	return &APIConsumer{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
		Token:   os.Getenv("API_TOKEN"),
	}
}

// Get realiza una solicitud GET a un endpoint específico de la API
func (api *APIConsumer) Get(endpoint string, nextPage string) (*interfaces.APIResponse, error) {
	url := fmt.Sprintf("%s/%s", api.BaseURL, endpoint)
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Añadir el token de autorización si está disponible
	if api.Token != "" {
		req.Header.Add("Authorization", "Bearer "+api.Token)
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	var apiResp Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Convertir la respuesta al formato de la interfaz
	items := make([]interfaces.StockItem, len(apiResp.Items))
	for i, item := range apiResp.Items {
		items[i] = interfaces.StockItem{
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

	return &interfaces.APIResponse{
		Items:    items,
		NextPage: apiResp.NextPage,
	}, nil
}
