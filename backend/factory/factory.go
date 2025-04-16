package factory

import (
	"backend/client"
	"backend/db"
	"backend/interfaces"
	"backend/repositories"
	"backend/services"
	"os"
)

// ServiceFactory crea instancias de servicios
type ServiceFactory struct{}

// NewServiceFactory crea una nueva fábrica de servicios
func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{}
}

// CreateStockService crea y devuelve un servicio de stocks
func (f *ServiceFactory) CreateStockService() (interfaces.StockService, error) {
	// Crear handler de base de datos
	dbHandler, err := db.NewDatabaseHandler()
	if err != nil {
		return nil, err
	}

	// Crear repositorio
	repo := repositories.NewStockRepository(dbHandler)

	// Crear cliente API
	apiClient := client.NewAPIConsumer(os.Getenv("API_URL"))

	// Crear y devolver el servicio
	return services.NewStockService(repo, apiClient), nil
}
