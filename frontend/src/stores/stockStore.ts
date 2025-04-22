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

    // Acción para cambiar de página
    async goToPage(page: number) {
      if (page < 1 || page > this.pagination.totalPages || page === this.pagination.currentPage) {
        return; // No hacer nada si la página es inválida o es la actual
      }
      // Determinar si estamos en modo búsqueda o carga general
      if (this.currentSearchType && this.currentSearchType !== 'loadAll' && this.currentSearchQuery) {
        // Recargar la búsqueda en la nueva página
        await this.searchStocks(this.currentSearchType, this.currentSearchQuery, page);
      } else {
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