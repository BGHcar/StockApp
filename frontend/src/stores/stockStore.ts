import { defineStore } from "pinia";
import type { Stock } from "@/types/stock";
// Importar PaginatedResponse y las funciones actualizadas
import { fetchStocks, searchStocks, type PaginatedResponse } from "@/services/stockService";

interface PaginationState {
  currentPage: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
}

interface StockState {
  stocks: Stock[];
  loading: boolean;
  error: string | null;
  selectedStock: Stock | null;
  pagination: PaginationState; // Añadir estado de paginación
  // Mantener filtros si los usas para búsquedas
  currentSearchType: string;
  currentSearchQuery: string;
  currentSortField: string | null; // Campo por el que se está ordenando actualmente
}

export const useStockStore = defineStore("stock", {
  state: (): StockState => ({
    stocks: [],
    loading: false,
    error: null,
    selectedStock: null,
    pagination: { // Inicializar paginación
      currentPage: 1,
      pageSize: 20, // Tamaño de página por defecto
      totalItems: 0,
      totalPages: 0,
    },
    currentSearchType: 'general', // Guardar el tipo de búsqueda actual
    currentSearchQuery: '',      // Guardar la query actual
    currentSortField: null,      // Guardar el campo de ordenamiento actual
  }),

  getters: {
    // Los getters existentes pueden permanecer igual o adaptarse si es necesario
    filteredStocks: (state) => {
      // Esta lógica podría eliminarse si la paginación/filtrado se hace en backend
      return state.stocks;
    },
    stockAnalytics: (state) => {
      return {
        // totalStocks ahora es el total de la página actual, usar totalItems para el global
        totalStocksOnPage: state.stocks.length,
        totalStocksInDB: state.pagination.totalItems,
        // Más métricas...
      };
    }
  },

  actions: {
    // Acción para cargar stocks paginados
    async loadStocks(page?: number, pageSize?: number) {
      this.loading = true;
      this.error = null;
      const targetPage = page ?? this.pagination.currentPage;
      const targetPageSize = pageSize ?? this.pagination.pageSize;

      // Resetear búsqueda actual al cargar todos
      this.currentSearchType = 'loadAll';
      this.currentSearchQuery = '';

      try {
        const response: PaginatedResponse<Stock> = await fetchStocks(targetPage, targetPageSize);
        this.stocks = response.items;
        this.pagination.currentPage = response.pagination.page;
        this.pagination.pageSize = response.pagination.pageSize;
        this.pagination.totalItems = response.pagination.totalItems;
        this.pagination.totalPages = response.pagination.totalPages;
      } catch (error) {
        this.error = "Failed to load stocks";
        this.stocks = []; // Limpiar stocks en caso de error
        this.pagination.totalItems = 0;
        this.pagination.totalPages = 0;
      } finally {
        this.loading = false;
      }
    },

    // Acción para realizar búsquedas paginadas
    async searchStocks(type: string, query: string, page?: number, pageSize?: number) {
      this.loading = true;
      this.error = null;
      const targetPage = page ?? 1; // Siempre empezar en la página 1 al buscar
      const targetPageSize = pageSize ?? this.pagination.pageSize;

      // Guardar búsqueda actual
      this.currentSearchType = type;
      this.currentSearchQuery = query;

      try {
        const response: PaginatedResponse<Stock> = await searchStocks(type, query, targetPage, targetPageSize);
        this.stocks = response.items;
        this.pagination.currentPage = response.pagination.page;
        this.pagination.pageSize = response.pagination.pageSize;
        this.pagination.totalItems = response.pagination.totalItems;
        this.pagination.totalPages = response.pagination.totalPages;
      } catch (e: any) {
        this.error = e.message || "Failed to search stocks";
        this.stocks = [];
        this.pagination.totalItems = 0;
        this.pagination.totalPages = 0;
      } finally {
        this.loading = false;
      }
    },

    // Método simplificado para cargar stocks ordenados
    async loadSortedStocks(sortBy: string) {
      this.loading = true;
      this.error = null;
      const targetPage = 1; // Al cambiar ordenamiento, volvemos a la primera página
      
      // Actualizar el campo de ordenamiento actual
      this.currentSortField = sortBy;
      
      try {
        // Determinar si hay un término de búsqueda activo que debemos preservar
        let activeSearchQuery = '';
        if (this.currentSearchType !== 'sortBy' && this.currentSearchType !== 'loadAll') {
          activeSearchQuery = this.currentSearchQuery;
        }
        
        // Llamar al servicio pasando el término de búsqueda activo si existe
        const response = await searchStocks('sortBy', sortBy, targetPage, this.pagination.pageSize, activeSearchQuery);
        
        this.stocks = response.items;
        this.pagination.currentPage = response.pagination.page;
        this.pagination.pageSize = response.pagination.pageSize;
        this.pagination.totalItems = response.pagination.totalItems;
        this.pagination.totalPages = response.pagination.totalPages;
        
        // No actualizamos el tipo de búsqueda si ya había una búsqueda activa
        if (this.currentSearchType === 'loadAll' || this.currentSearchType === 'sortBy') {
          this.currentSearchType = 'sortBy';
          this.currentSearchQuery = sortBy;
        }
      }
      catch (error) {
        this.error = "Error al cargar stocks ordenados";
        this.stocks = [];
        this.pagination.totalItems = 0;
        this.pagination.totalPages = 0;
      }
      finally {
        this.loading = false;
      }
    },

    // Acción para cambiar de página que soporte ordenamiento
    async goToPage(page: number) {
      if (page < 1 || page > this.pagination.totalPages || page === this.pagination.currentPage) {
        return; // No hacer nada si la página es inválida o es la actual
      }
      
      // Determinar si estamos en modo ordenamiento, búsqueda o carga general
      if (this.currentSearchType === 'sortBy' && this.currentSortField) {
        // Determinar si hay un término de búsqueda activo que debemos preservar
        let activeSearchQuery = '';
        if (this.currentSearchQuery && this.currentSearchQuery !== this.currentSortField) {
          activeSearchQuery = this.currentSearchQuery;
        }
        
        // Usar searchStocks con el parámetro de búsqueda activa
        await searchStocks('sortBy', this.currentSortField, page, this.pagination.pageSize, activeSearchQuery)
          .then(response => {
            this.stocks = response.items;
            this.pagination.currentPage = response.pagination.page;
            this.pagination.pageSize = response.pagination.pageSize;
            this.pagination.totalItems = response.pagination.totalItems;
            this.pagination.totalPages = response.pagination.totalPages;
          })
          .catch(() => {
            this.error = "Error al cambiar página con ordenamiento";
            this.stocks = [];
          });
      } 
      else if (this.currentSearchType && this.currentSearchType !== 'loadAll' && this.currentSearchQuery) {
        // Recargar la búsqueda en la nueva página
        await this.searchStocks(this.currentSearchType, this.currentSearchQuery, page);
      } 
      else {
        // Recargar la lista general en la nueva página
        await this.loadStocks(page);
      }
    },

    // Acción para resetear (volver a cargar la primera página de todos)
    async resetStocks() {
      await this.loadStocks(1); // Carga la primera página
    }
  },
});