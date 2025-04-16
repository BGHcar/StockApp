// frontend/src/components/StockRecommendations.vue
<template>
  <BaseTable>
    <div class="recommendations-container">
      <h2 class="title">Recomendaciones de Inversión</h2>
      <div v-if="loading" class="loading">
        Analizando stocks...
      </div>
      <div v-else-if="error" class="error">
        {{ error }}
      </div>
      <div v-else class="recommendations-list">
        <div v-for="stock in topStocks" 
             :key="stock.ticker" 
             class="recommendation-item">
          <div class="stock-info">
            <h3>{{ stock.company }}</h3>
            <span class="ticker">{{ stock.ticker }}</span>
          </div>
          <div class="recommendation-details">
            <p class="rating">Rating: {{ stock.rating_to }}</p>
            <p class="action">{{ stock.action }}</p>
            <p class="brokerage">Por: {{ stock.brokerage }}</p>
          </div>
        </div>
      </div>
    </div>
  </BaseTable>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import type { Stock } from '@/types/stock';
import { getTopRecommendations } from '@/services/analysisService';
import { useStockStore } from '@/stores/stockStore';

const store = useStockStore();
const topStocks = ref<Stock[]>([]);
const loading = ref(false);
const error = ref<string | null>(null);

async function loadRecommendations() {
  loading.value = true;
  error.value = null;
  try {
    topStocks.value = getTopRecommendations(store.stocks);
  } catch (e) {
    error.value = 'Error al cargar recomendaciones';
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadRecommendations();
});
</script>

<style scoped>
.recommendations-container {
  padding: 1rem;
}

.title {
  color: white;
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.recommendations-list {
  display: grid;
  gap: 1rem;
}

.recommendation-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: rgba(60, 16, 83, 0.85);
  border: 1px solid rgba(173, 83, 137, 0.3);
  border-radius: 8px;
}

.stock-info {
  flex: 1;
}

.stock-info h3 {
  color: white;
  margin: 0;
  font-size: 1.1rem;
}

.ticker {
  color: rgba(173, 83, 137, 0.8);
  font-size: 0.9rem;
}

.recommendation-details {
  text-align: right;
}

.rating {
  color: #10B981;
  font-weight: bold;
}

.action {
  color: white;
  font-size: 0.9rem;
}

.brokerage {
  color: rgba(255, 255, 255, 0.7);
}
</style>