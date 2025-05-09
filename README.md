# StockAnalyzer - Sistema de Análisis y Recomendaciones para Inversiones

![App](./docs/AppStockAnalyzer.png)

## Descripción del Proyecto

StockAnalyzer es un sistema completo para analizar acciones bursátiles y generar recomendaciones de inversión basadas en datos actualizados del mercado. La aplicación permite a los usuarios explorar tendencias, filtrar recomendaciones según diversos criterios y visualizar datos de acciones de manera intuitiva.

## Arquitectura

El proyecto implementa una **Arquitectura Limpia (Clean Architecture)** con una clara separación de responsabilidades en capas:

### Capas de la arquitectura:

1. **Capa de Presentación**: 
   - Chi Router para manejo de rutas HTTP
   - Handlers para procesamiento de solicitudes/respuestas

2. **Capa de Aplicación**: 
   - Controladores (StockController)
   - Factory para creación de servicios con inyección de dependencias

3. **Capa de Dominio**: 
   - Servicios con lógica de negocio
   - Interfaces que definen contratos entre capas
   - Modelos de entidades

4. **Capa de Infraestructura**: 
   - Repositorios para acceso a datos
   - Adaptadores para base de datos (GORM/CockroachDB)
   - Cliente API para servicios externos
   - Configuración de la aplicación

## Tecnologías Utilizadas

### Backend
- **Go 1.23**: Lenguaje principal para la API REST
- **Chi Router**: Framework HTTP ligero para Go
- **GORM**: ORM para operaciones con base de datos
- **CockroachDB**: Base de datos distribuida SQL
- **Testify**: Framework para pruebas unitarias
- **go-sqlmock**: Herramienta para mock de base de datos en pruebas

### Frontend
- **Vue 3**: Framework progresivo para JavaScript
- **Vite**: Herramienta de construcción frontend
- **TailwindCSS**: Framework CSS de utilidades
- **Pinia**: Gestión de estado centralizado

## Características Principales

### API REST

- **Sincronización de Datos**: Obtiene datos actualizados de fuentes externas
- **Filtrado Avanzado**: Por ticker, empresa, precio, rating, bróker, etc.
- **Estadísticas**: Análisis de distribución por tipo de recomendación
- **Búsqueda Global**: Implementación de búsqueda en múltiples campos
- **Rango de Precios**: Filtrado por precio objetivo

### Análisis y Visualización

- **Recomendaciones Inteligentes**: Sistema de puntuación basado en múltiples factores
- **Visualización de Tendencias**: Identificación de patrones en las recomendaciones
- **Interfaz Intuitiva**: Dashboard con componentes interactivos y responsivos
- **Filtros Combinados**: Posibilidad de aplicar múltiples filtros simultáneamente

## Testing

El proyecto incluye pruebas unitarias exhaustivas para todos los componentes:

- **Tests de Cliente API**: Validan la comunicación con servicios externos
- **Tests de Servicios**: Verifican la lógica de negocio
- **Tests de Repositorios**: Comprueban las operaciones con base de datos
- **Tests de Controladores**: Validan los endpoints HTTP

Para ejecutar las pruebas:

```bash
cd backend
./run_tests.bat
```


### Arquitectura del Proyecto

### Arquitectura del Proyecto

![Diagrama de Arquitectura de Stock Analyzer](./docs/Arquitectura.png)


## Diagrama de Clases

![Diagrama de Clases de Stock Analyzer](./docs/StockApp.png)


# Desarrollado por

Bryan Cartagena Hincapie - Abril 2025