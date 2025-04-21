<template>
  <div class="container">
    <h1 class="title font-bold text-2xl text-white text-center mb-2">
      <span class="text-[#ef4444]">STOCKS </span>
      <span class="text-[#60a5fa]">SEARCH </span>
      <span class="text-white">APP</span> <!-- Cambiado de verde a blanco -->
    </h1>

    <div class="content">
      <StockSearch @search="handleSearch" @reset="handleReset" class="search-container" />

      <div v-if="stockStore.loading" class="status-message">
        Cargando...
      </div>
      <div v-else-if="stockStore.error" class="status-message error">
        {{ stockStore.error }}
      </div>
      <div v-else class="components-wrapper">
        <StockTable v-if="stockStore.stocks.length > 0" :stocks="stockStore.stocks" :headers="tableHeaders"
          class="table-container" />
        <div v-else class="no-results">
          No se encontraron resultados
        </div>
        <StockRecommendations v-if="stockStore.stocks.length > 0" :stocks="stockStore.stocks"
          class="recommendations-container" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '@/stores/stockStore'
import StockTable from '@/components/StockTable.vue'
import StockSearch from '@/components/StockSearch.vue'
import { searchStocks } from '@/services/stockService'
import StockRecommendations from '@/components/StockRecommendations.vue'
import type { TableHeader } from '@/types/table'

const stockStore = useStockStore()

const tableHeaders: TableHeader[] = [
  { key: 'ticker', label: 'Ticker', class: 'w-[8%]' },
  { key: 'company', label: 'Compañía', class: 'w-[15%]' },
  { key: 'brokerage', label: 'Brokerage', class: 'w-[15%]' },
  { key: 'action', label: 'Acción', class: 'w-[12%]' },
  { key: 'rating_from', label: 'Rating Ant.', class: 'w-[10%]' }, // Había un label vacío aquí
  { key: 'rating_to', label: 'Rating Act.', class: 'w-[10%]' },
  { key: 'target_from', label: 'Precio Desde', class: 'w-[10%]' },
  { key: 'target_to', label: 'Precio Hasta', class: 'w-[10%]' },
  /* { key: 'time', label : 'Fecha', class: 'w-[10%]' },   */
]

onMounted(() => {
  stockStore.loadStocks()
})

async function handleSearch(type: string, query: string) {
  stockStore.loading = true
  stockStore.error = null
  try {
    const results = await searchStocks(type, query)
    stockStore.stocks = results
  } catch (e: any) {
    stockStore.error = e.message
    stockStore.stocks = []
  } finally {
    stockStore.loading = false
  }
}

function handleReset() {
  stockStore.loadStocks()
}
</script>

<style scoped>
.title {
  
  font-size: 1.75rem;
  /* Reducido aún más */
  font-weight: 900;
  /* Aumentar a 900 para una negrita más fuerte */
  /* Cambiar color a blanco por defecto */
  color: white;
  text-align: center;
  margin: 0 0 0.25rem 0;
  /* Eliminar margen superior */
  flex: 0 0 auto;
  /* No crecer */
  
  /* Añadir espaciado entre letras */
  width: 100%;
  /* Asegurar que ocupe todo el ancho disponible */
  text-transform: uppercase;
  /* Opcional: texto en mayúsculas */
  /* Actualizar sombra para que complemente el nuevo tema */
  
}

/* Específicamente para los spans dentro del título */
.title span {
  font-weight: 900;
  /* Asegurar que los spans también están en negrita */
  padding: 0 0.1em;
  /* Añadir un poco de espacio horizontal */
}

.content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
  flex: 1;
  overflow: hidden;
}

.search-container {
  width: 100%;
  border-radius: 8px;
  /* Cambiar del fondo morado al transparente */
  background: transparent;
  flex: 0 0 auto;
  /* No crecer */
}

.components-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  /* Aumentar ligeramente el espacio entre componentes */
  flex: 1;
  overflow: hidden;
  min-height: 0;
  /* Limitar altura máxima */
  height: 100%;
  /* Asegurar que ocupe todo el espacio disponible */
}

.recommendations-container {
  width: 100%;
  border-radius: 8px;
  /* Cambiar del fondo morado al transparente */
  background: transparent;
  flex: 0 0 auto;
  /* No crecer ni encoger, ajustarse al contenido */
  max-height: none;
  /* Quitar la restricción de altura máxima */
  overflow: visible;
  /* Sin scroll */
  padding-bottom: 0;
  /* Eliminar el padding-bottom adicional */
}

.status-message {
  text-align: center;
  padding: 0.5rem;
  color: #333;
  flex: 0 0 auto;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 8px;
  margin: 1rem 0;
  border: 2px solid #646464;
  /* Borde más oscuro y consistente */
}

.error {
  color: #d32f2f;
  border: 2px solid #d32f2f;
  /* Borde rojo para errores */
}

/* Estilo para mensaje de no resultados */
.no-results {
  text-align: center;
  padding: 1rem;
  color: #333;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 8px;
  margin: 1rem 0;
  border: 2px solid #646464;
  /* Borde más oscuro y consistente */
}
</style>