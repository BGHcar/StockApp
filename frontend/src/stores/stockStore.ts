import { defineStore } from "pinia";
import type { Stock } from "@/types/stock"; 
import { fetchStocks } from "@/services/stockService";

interface StockState {
  stocks: Stock[];
  loading: boolean;
  error: string | null;
  selectedStock: Stock | null;
  filters: {
    sortBy: keyof Stock | null;
    sortOrder: 'asc' | 'desc';
    searchQuery: string;
    filterType: string;
  };
}

export const useStockStore = defineStore("stock", {
  state: (): StockState => ({
    stocks: [],
    loading: false,
    error: null,
    selectedStock: null,
    filters: {
      sortBy: null,
      sortOrder: 'asc',
      searchQuery: '',
      filterType: 'general'
    }
  }),
  
  getters: {
    filteredStocks: (state) => {
      // Lógica de filtrado
      return state.stocks;
    },
    
    stockAnalytics: (state) => {
      // Análisis de stocks
      return {
        totalStocks: state.stocks.length,
        // Más métricas...
      };
    }
  },
  
  actions: {
    async loadStocks() {
      this.loading = true;
      this.error = null;
      try {
        const stocks = await fetchStocks();
        this.stocks = stocks;
      } catch (error) {
        this.error = "Failed to load stocks";
      } finally {
        this.loading = false;
      }
    },
  },
});