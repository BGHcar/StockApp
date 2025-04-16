<template>
  <div class="container">
    <h1 class="title">Lista de Stocks</h1>
    
    <div class="content">
      <StockSearch @search="handleSearch" @reset="handleReset" class="component-container" />
      
      <div v-if="stockStore.loading" class="status-message">
        Cargando...
      </div>
      <div v-else-if="stockStore.error" class="status-message error">
        {{ stockStore.error }}
      </div>
      <StockTable 
        v-else 
        :stocks="stockStore.stocks" 
        :headers="tableHeaders" 
        class="component-container" 
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '@/stores/stockStore'
import StockTable from '@/components/StockTable.vue'
import StockSearch from '@/components/StockSearch.vue'
import { searchStocks } from '@/services/stockService'
import type { TableHeader } from '@/types/table'

const stockStore = useStockStore()

const tableHeaders: TableHeader[] = [
  { key: 'ticker', label: 'Ticker', class: 'w-[8%]' },
  { key: 'company', label: 'Compañía', class: 'w-[15%]' },
  { key: 'target_from', label: 'Precio Desde', class: 'w-[10%]' },
  { key: 'target_to', label: 'Precio Hasta', class: 'w-[10%]' },
  { key: 'action', label: 'Acción', class: 'w-[12%]' },
  { key: 'brokerage', label: 'Brokerage', class: 'w-[15%]' },
  { key: 'rating_from', label: 'Rating Ant.', class: 'w-[10%]' },
  { key: 'rating_to', label: 'Rating Act.', class: 'w-[10%]' },
  { key: 'time', label: 'Fecha', class: 'w-[10%]' }
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
.container {
  width: 100%;
  max-width: 95vw; /* Reducimos un poco para tener margen mínimo */
  margin: 0 auto; /* Centramos el contenedor */
  padding: 0.5rem;
}

.title {
  font-size: 2rem;
  font-weight: bold;
  color: white;
  text-align: center;
  margin-bottom: 0.5rem;
}

.content {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  width: 100%; /* Aseguramos que ocupe todo el ancho */
}

.component-container {
  width: 100%; /* Aseguramos que los componentes ocupen todo el ancho */
  border-radius: 8px;
  background: rgba(60, 16, 83, 0.95);
}

.status-message {
  text-align: center;
  padding: 1rem;
  color: white;
}

.error {
  color: #ef4444;
}
</style>