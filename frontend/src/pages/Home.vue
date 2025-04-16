<template>
  <div class="container">
    <h1 class="title font-bold text-2xl text-white text-center mb-2">
      <span class="text-[#ef4444]">STOCKS </span> 
      <span class="text-[#60a5fa]">SEARCH </span>
      <span class="text-[#4ade80]">APP</span>
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
        <StockTable 
          :stocks="stockStore.stocks" 
          :headers="tableHeaders" 
          class="table-container" 
        />
        <StockRecommendations 
          v-if="stockStore.stocks.length > 0" 
          :stocks="stockStore.stocks" 
          class="recommendations-container"
        />
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
  max-width: 95vw;
  margin: 0 auto; /* Cambiado de '0 0' a '0 auto' para centrar horizontalmente */
  padding: 0.5rem;
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.title {
  font-family: Arial, Helvetica, sans-serif;
  font-size: 1.75rem; /* Reducido aún más */
  font-weight: 900; /* Aumentar a 900 para una negrita más fuerte */
  color: white;
  text-align: center;
  margin: 0 0 0.25rem 0; /* Eliminar margen superior */
  flex: 0 0 auto; /* No crecer */
  letter-spacing: 0.05em; /* Añadir espaciado entre letras */
  width: 100%; /* Asegurar que ocupe todo el ancho disponible */
  text-transform: uppercase; /* Opcional: texto en mayúsculas */
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3); /* Sombra sutil para mejor legibilidad */
}

/* Específicamente para los spans dentro del título */
.title span {
  font-weight: 900; /* Asegurar que los spans también están en negrita */
  padding: 0 0.1em; /* Añadir un poco de espacio horizontal */
}

.content {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  width: 100%;
  flex: 1;
  overflow: hidden;
}

.search-container {
  width: 100%;
  border-radius: 8px;
  background: rgba(60, 16, 83, 0.95);
  margin-bottom: 0.25rem;
  flex: 0 0 auto; /* No crecer */
}

.components-wrapper {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
  overflow: hidden;
  min-height: 0;
  max-height: calc(100vh - 150px); /* Limitar altura máxima */
  height: 100%; /* Asegurar que ocupe todo el espacio disponible */
}

.table-container {
  width: 100%;
  border-radius: 8px;
  background: rgba(60, 16, 83, 0.95);
  flex: 1 1 auto; /* Cambiar a flex-grow: 1 para que ocupe el espacio restante */
  min-height: 0; /* Permitir que se encoja si es necesario */
  overflow: auto; /* Esta tendrá scroll */
}

.recommendations-container {
  width: 100%;
  border-radius: 8px;
  background: rgba(60, 16, 83, 0.95);
  flex: 0 0 auto; /* No crecer ni encoger, ajustarse al contenido */
  max-height: none; /* Quitar la restricción de altura máxima */
  overflow: visible; /* Sin scroll */
  padding-bottom: 0.25rem;
}

.status-message {
  text-align: center;
  padding: 0.5rem;
  color: white;
  flex: 0 0 auto; /* No crecer */
}

.error {
  color: #ef4444;
}
</style>