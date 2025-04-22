<template>
  <BaseTable>
    <div class="recommendations-container">
      <h2 class="section-title">Recomendaciones</h2>

      <div v-if="recommendationStore.loading" class="loading-message">
        Cargando recomendaciones...
      </div>
      <div v-else-if="recommendationStore.error" class="error-message">
        Error: {{ recommendationStore.error }}
      </div>
      <div v-else-if="recommendations.length === 0" class="no-data-message">
        No hay recomendaciones disponibles.
      </div>

      <!-- Mostrar las recomendaciones del backend -->
      <div v-else class="recommendation-cards">
        <div v-for="rec in recommendations" :key="rec.ticker" class="card" :class="getCardClass(rec.score)">
          <h3>{{ rec.ticker }}</h3>
          <div class="company">{{ rec.company }}</div>
          <div class="details">
            <div class="score">Puntuación: {{ rec.score.toFixed(2) }}</div>
            <div class="last-update">Últ. Señal: {{ formatDate(rec.last_update) }}</div>
          </div>
          <div class="reason">{{ rec.reason }}</div>
        </div>
      </div>
    </div>
  </BaseTable>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import BaseTable from './base/BaseTable.vue'
// Importar el nuevo store de recomendaciones
import { useRecommendationStore } from '@/stores/recommendationStore'
import type { Recommendation } from '@/types/recommendation' // Importar el tipo

const recommendationStore = useRecommendationStore()

// Obtener las recomendaciones del store
const recommendations = computed(() => recommendationStore.recommendations)

// Cargar recomendaciones al montar el componente
onMounted(() => {
  recommendationStore.fetchRecommendations()
})

// Función para determinar la clase CSS basada en la puntuación
function getCardClass(score: number): string {
  if (score >= 2.0) return 'top-pick' // Puntuación alta -> verde
  if (score <= 0.5) return 'avoid'    // Puntuación baja -> rojo
  return 'neutral'                   // Puntuación media -> gris/neutro
}

// Función para formatear la fecha
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  if (isNaN(date.getTime())) return dateString; // Devolver original si no es fecha válida
  return date.toLocaleDateString('es-ES', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  });
}
</script>

<style scoped>
.section-title {
  text-align: center;
  margin-bottom: 10px;
  color: #fff;
  font-weight: bold;
  text-transform: uppercase;
}

.recommendation-cards {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.card {
  border: 1px solid black;
  border-radius: 8px;
  padding: 8px;
  background-color: rgba(255, 255, 255, 0.7);
  -webkit-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  -moz-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
}

.card:hover {
  background-color: #646464;
  color: #fff;
}

.card:hover .company,
.card:hover .details,
.card:hover .reason {
  border-color: #fff;
}

.company,
.details,
.reason {
  margin: 4px 0;
  padding: 4px 0;
  border-bottom: 1px solid #646464;
}

@media (min-width:992px) {
  .recommendation-cards {
    flex-direction: row;
  }

  .card {
    flex: 1 1 20%;
  }
}
</style>