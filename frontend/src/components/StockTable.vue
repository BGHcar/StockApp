<template>
  <BaseTable>
    <BaseScroll class="table-container">
      <table class="stock-table">
        <thead>
          <tr>
            <th 
              v-for="header in headers" 
              :key="header.key" 
              :class="[header.class, 'sortable']"
              @click="sort(header.key)"
            >
              {{ header.label }}
              <span class="sort-icon" v-if="sortKey === header.key">
                {{ sortOrder === 'asc' ? '▲' : '▼' }}
              </span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="stock in sortedStocks" 
              :key="stock.ticker + stock.time">
            <td v-for="header in headers" 
                :key="header.key"
                :class="header.class">
              {{ stock[header.key] }}
            </td>
          </tr>
        </tbody>
      </table>
    </BaseScroll>
  </BaseTable>
</template>

<script setup lang="ts">
import { defineProps, ref, computed } from 'vue'
import type { Stock } from '@/types/stock'
import type { TableHeader } from '@/types/table'
import BaseTable from './base/BaseTable.vue'
import BaseScroll from './base/BaseScroll.vue'

const props = defineProps<{ 
  stocks: Stock[],
  headers: TableHeader[] 
}>()

// Estado para el ordenamiento
const sortKey = ref<keyof Stock | null>(null);
const sortOrder = ref<'asc' | 'desc'>('asc');

// Función para ordenar
function sort(key: keyof Stock) {
  if (sortKey.value === key) {
    // Si ya está ordenado por esta columna, invertir el orden
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
  } else {
    // Si no, ordenar ascendente
    sortKey.value = key;
    sortOrder.value = 'asc';
  }
}

// Ordenar los stocks
const sortedStocks = computed(() => {
  if (!sortKey.value) return props.stocks;
  
  return [...props.stocks].sort((a, b) => {
    const aValue = a[sortKey.value as keyof Stock];
    const bValue = b[sortKey.value as keyof Stock];
    
    // Para fechas
    if (sortKey.value === 'time') {
      return sortOrder.value === 'asc' 
        ? new Date(aValue as string).getTime() - new Date(bValue as string).getTime()
        : new Date(bValue as string).getTime() - new Date(aValue as string).getTime();
    }
    
    // Para precios
    if (sortKey.value === 'target_from' || sortKey.value === 'target_to') {
      const numA = parseFloat((aValue as string).replace(/[^\d.-]/g, ''));
      const numB = parseFloat((bValue as string).replace(/[^\d.-]/g, ''));
      return sortOrder.value === 'asc' ? numA - numB : numB - numA;
    }
    
    // Para strings
    return sortOrder.value === 'asc'
      ? String(aValue).localeCompare(String(bValue))
      : String(bValue).localeCompare(String(aValue));
  });
});
</script>

<style scoped>
.table-container {
  width: 100%;
  max-height: calc(100vh - 200px); /* Usamos max-height en lugar de height */
  min-height: 150px; /* Establecemos una altura mínima */
  overflow-y: auto;
}

/* El resto del CSS se mantiene igual */
.stock-table {
  width: 100%;
  min-width: 800px; /* Asegura un ancho mínimo */
  table-layout: fixed; /* Importante para controlar el ancho de columnas */
  border-collapse: collapse;
  background: linear-gradient(180deg, rgba(60, 16, 83, 0.95) 0%, rgba(60, 16, 83, 0.85) 100%);
}

/* Definimos anchos específicos para cada columna */
th, td {
  padding: 0.35rem 0.5rem;
  text-align: left;
  border: 1px solid rgba(173, 83, 137, 0.3);
  white-space: nowrap; /* Evita que el texto se rompa */
  overflow: hidden;
  text-overflow: ellipsis; /* Muestra ... si el texto es muy largo */
}

/* Asegurar que todas las celdas de cabecera tengan bordes */
th {
  position: sticky;
  top: 0;
  background: rgba(60, 16, 83, 0.98);
  font-family: Arial, Helvetica, sans-serif;
  font-weight: 700;
  color: white;
  text-transform: uppercase;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(173, 83, 137, 0.3) !important; /* Forzar bordes con !important */
  border-bottom: 2px solid rgba(173, 83, 137, 0.4) !important;
  z-index: 10;
}

/* Específicamente para la columna Rating Act. y Precio Desde */
th:nth-child(6), th:nth-child(7) {
  border-right: 2px solid rgba(173, 83, 137, 0.4) !important; /* Línea más gruesa y visible */
}

.sortable {
  cursor: pointer;
  user-select: none;
}

.sortable:hover {
  background: rgba(173, 83, 137, 0.2);
}

.sort-icon {
  margin-left: 4px;
  font-size: 0.7em;
}

th:nth-child(1), td:nth-child(1) { width: 8%; } /* Ticker */
th:nth-child(2), td:nth-child(2) { width: 15%; } /* Compañía */
th:nth-child(3), td:nth-child(3) { width: 15%; } /* Brokerage */
th:nth-child(4), td:nth-child(4) { width: 12%; } /* Acción */
th:nth-child(5), td:nth-child(5) { width: 10%; } /* Rating Ant. */
th:nth-child(6), td:nth-child(6) { width: 10%; } /* Rating Act. */
th:nth-child(7), td:nth-child(7) { width: 10%; } /* Precio Desde */
th:nth-child(8), td:nth-child(8) { width: 10%; } /* Precio Hasta */
th:nth-child(9), td:nth-child(9) { width: 10%; } /* Fecha */

td {
  background: rgba(60, 16, 83, 0.85);
  color: #f3f4f6;
}

tr:hover td {
  background: rgba(173, 83, 137, 0.15);
  transition: background-color 0.3s ease;
}

/* Estilizar el scroll */
.table-container::-webkit-scrollbar {
  width: 10px;
}

.table-container::-webkit-scrollbar-track {
  background: rgba(60, 16, 83, 0.8);
  border-radius: 15px;
}

.table-container::-webkit-scrollbar-thumb {
  background: rgba(173, 83, 137, 0.3);
  border-radius: 15px;
}

.table-container::-webkit-scrollbar-thumb:hover {
  background: rgba(173, 83, 137, 0.5);
}
</style>