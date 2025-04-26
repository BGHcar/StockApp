<template>
  <div class="container">
    <div class="title">
      <h1 class="title-text font-bold text-2xl text-white text-center">
        STOCKS SEARCH APP
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
        <StockTable :stocks="stockStore.stocks" :headers="tableHeaders" class="table-container" />
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
import StockRecommendations from '@/components/StockRecommendations.vue'
import type { TableHeader } from '@/types/table'

const stockStore = useStockStore()

const tableHeaders: TableHeader[] = [
  { key: 'Ticker', label: 'Ticker', class: 'w-[8%]' },
  { key: 'Company', label: 'Compañía', class: 'w-[15%]' },
  { key: 'Brokerage', label: 'Brokerage', class: 'w-[15%]' },
  { key: 'Action', label: 'Acción', class: 'w-[12%]' },
  { key: 'RatingFrom', label: 'Rating Ant.', class: 'w-[10%]' }, 
  { key: 'RatingTo', label: 'Rating Act.', class: 'w-[10%]' },
  { key: 'TargetFrom', label: 'Precio Desde', class: 'w-[10%]' },
  { key: 'TargetTo', label: 'Precio Hasta', class: 'w-[10%]' },
  { key: 'Time', label : 'Fecha', class: 'w-[10%]' }, 
]

onMounted(() => {
  // Cargar la primera página al montar
  stockStore.loadStocks(1, stockStore.pagination.pageSize);
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
  background-color: #292F36;
  border-radius: 8px;
}

.status-message {
  width: 100%;
  padding: 20px;
  text-align: center;
  font: bold 2rem 'Arial', sans-serif;
  color: #fff;
  background-color: #9dcbba;
  border-radius: 8px;
  border: 1px solid;
  -webkit-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  -moz-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
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
  font-size: 2.5rem;
  font-weight: bold;
  background-color: #9dcbba;
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
  background-color: #9dcbba;
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
  background-color: #9dcbba;
  -webkit-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  -moz-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
}

.recommendations-container {
  border: 1px solid;
  border-radius: 8px;
  background-color: #9dcbba;
  padding: 8px;
  -webkit-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  -moz-box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
  box-shadow: 2px 4px 16px 0px rgba(0, 0, 0, 0.75);
}


@media (min-width: 768px) {
  .table-container {
  height: 40vh;
}

  .content{
    height: 90vh;
  }
}

</style>