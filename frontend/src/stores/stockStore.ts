import { defineStore } from "pinia";
import type { Stock } from "@/types/Stock";
import { fetchStocks } from "@/services/stockService";

export const useStockStore = defineStore("stock", {
  state: () => ({
    stocks: [] as Stock[],
    loading: false,
    error : null as string | null,
  }),
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