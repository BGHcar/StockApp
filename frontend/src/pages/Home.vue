<template>
  <div class="container">
    <div class="title">
      <h1 class="title-text font-bold text-2xl text-white text-center">
        <span class="text-[#ef4444]">STOCKS </span>
        <span class="text-[#60a5fa]">SEARCH </span>
        <span class="text-white">APP</span> <!-- Cambiado de verde a blanco -->
      </h1>
    </div>

    <div class="content">
      <!-- Escuchar eventos y llamar acciones del store -->
      <StockSearch @search="handleSearch" @reset="handleReset" class="search-container" />

      <div v-if="stockStore.loading" class="status-message">
        Cargando...
      </div>
      <div v-else-if="stockStore.error" class="status-message error">
        {{ stockStore.error }}
      </div>
      <div v-else class="components-wrapper">
        <!-- Pasar los stocks de la página actual a StockTable -->
        <StockTable :stocks="stockStore.stocks" :headers="tableHeaders" class="table-container" />
        <!-- StockRecommendations ahora obtiene datos de su propio store -->
      </div>
      <StockRecommendations class="recommendations-container" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '@/stores/stockStore'
import StockTable from '@/components/StockTable.vue'
import StockSearch from '@/components/StockSearch.vue'
import StockRecommendations from '@/components/StockRecommendations.vue' // Asegúrate que la ruta sea correcta
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
  // { key: 'time', label : 'Fecha', class: 'w-[10%]' }, // Descomentar si quieres mostrarla
]

onMounted(() => {
  // Cargar la primera página al montar
  stockStore.loadStocks(1, stockStore.pagination.pageSize);
  // No necesitas cargar recomendaciones aquí, el componente lo hace
})

// Llamar a la acción de búsqueda del store
async function handleSearch(type: string, query: string) {
  await stockStore.searchStocks(type, query); // El store maneja el loading y error
}

// Llamar a la acción de reseteo del store
function handleReset() {
  stockStore.resetStocks();
}
</script>

<style scoped>

.container {
  width: 100%;
  max-width: 1366px;
  margin: 0 auto;
  background-color: #292f36;
  border-radius: 8px;
}

.content {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.title {
  width: 100%;
  padding: 20px;
  text-align: center;
  font: bold 2rem 'Arial', sans-serif;
  color: #fff;
}

.title-text {
  font-size: 2rem;
  background-color: #937666;
  border-radius: 8px;
  padding: 8px;
  -webkit-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  -moz-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
}

.components-wrapper {
  width: 100%;
}

.search-container {
  border-radius: 8px;
  width: 100%;
  background-color: #937666;
  padding: 8px;
  -webkit-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  -moz-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
}

.table-container {
  border: 1px solid;
  border-radius: 8px;
  margin-bottom: 10px;
  height: 70vh;
  width: 100%;
  overflow: auto;
  background-color: #937666;
  -webkit-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  -moz-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
}

.recommendations-container {
  border: 1px solid;
  border-radius: 8px;
  background-color: #937666;
  padding: 8px;
  -webkit-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  -moz-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
}
</style>