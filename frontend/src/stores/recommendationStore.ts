import { defineStore } from 'pinia';
import type { Recommendation } from '@/types/recommendation';
import { fetchRecommendations } from '@/services/stockService'; // Importar la nueva función

interface RecommendationState {
  recommendations: Recommendation[];
  loading: boolean;
  error: string | null;
}

export const useRecommendationStore = defineStore('recommendation', {
  state: (): RecommendationState => ({
    recommendations: [],
    loading: false,
    error: null,
  }),
  actions: {
    async fetchRecommendations(limit: number = 5) {
      this.loading = true;
      this.error = null;
      try {
        this.recommendations = await fetchRecommendations(limit);
      } catch (error: any) {
        this.error = error.message || 'Failed to fetch recommendations';
        this.recommendations = []; // Limpiar en caso de error
      } finally {
        this.loading = false;
      }
    },
  },
});