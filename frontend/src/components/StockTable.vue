<template>
  <BaseTable>
    <BaseScroll class="table-container">
      <table class="stock-table">
        <thead>
          <tr>
            <th v-for="header in headers" 
                :key="header.key" 
                :class="header.class">
              {{ header.label }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="stock in stocks" 
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
import { defineProps } from 'vue'
import type { Stock } from '@/types/stock'
import type { TableHeader } from '@/types/table'
import BaseTable from './base/BaseTable.vue'
import BaseScroll from './base/BaseScroll.vue'

defineProps<{ 
  stocks: Stock[],
  headers: TableHeader[] 
}>()
</script>

<style scoped>
.table-container {
  width: 100%;
  height: calc(100vh - 200px);
  overflow-y: auto;
}

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

th:nth-child(1), td:nth-child(1) { width: 10%; } /* Ticker */
th:nth-child(2), td:nth-child(2) { width: 25%; } /* Compañía */
th:nth-child(3), td:nth-child(3) { width: 15%; } /* Acción */
th:nth-child(4), td:nth-child(4) { width: 20%; } /* Brokerage */
th:nth-child(5), td:nth-child(5) { width: 15%; } /* Rating */
th:nth-child(6), td:nth-child(6) { width: 15%; } /* Fecha */

th {
  position: sticky;
  top: 0;
  background: rgba(60, 16, 83, 0.98);
  font-family: Arial, Helvetica, sans-serif;
  font-weight: 700;
  color: white;
  text-transform: uppercase;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  border-bottom: 2px solid rgba(173, 83, 137, 0.4);
  z-index: 10;
}

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