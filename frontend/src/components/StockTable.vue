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
  height: 100%;
  position: relative;
  background: #D1CEC8;
  box-shadow: 0 3px 5px rgba(173, 169, 150, 0.2);
  border: 2px solid #646464; /* Borde más oscuro y consistente */
  border-radius: 8px;
}

/* El resto del CSS se mantiene igual */
.stock-table {
  width: 100%;
  min-width: 800px;
  table-layout: fixed;
  border-collapse: separate; /* Cambiado de collapse a separate para mejor control de bordes */
  border-spacing: 0; /* Eliminar espacios entre celdas */
  background: transparent; /* Cambiado para que use el gradiente del contenedor */
}

/* Definimos anchos específicos para cada columna */
th, td {
  padding: 0.35rem 0.5rem;
  text-align: left;
  border: 1px solid rgba(173, 169, 150, 0.3); /* Color de borde que complementa el gradiente */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Asegurar que todas las celdas de cabecera tengan bordes */
th {
  position: sticky;
  top: 0;
  background: #646464; /* Color sólido más oscuro que armoniza con el gradiente */
  font-family: Arial, Helvetica, sans-serif;
  font-weight: 700;
  color: white;
  text-transform: uppercase;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  border: 1px solid #4a4a4a;
  border-bottom: 2px solid #4a4a4a; /* Borde inferior más visible */
  z-index: 10;
  box-shadow: 0 2px 4px rgba(173, 169, 150, 0.4); /* Sombra acorde al fondo */
}

/* Específicamente para la columna Rating Act. y Precio Desde */
th:nth-child(6), th:nth-child(7) {
  border-right: 2px solid #4a4a4a !important; /* Línea más gruesa y visible */
}

.sortable {
  cursor: pointer;
  user-select: none;
  position: relative; /* Para crear un nuevo contexto de apilamiento */
}

.sortable:hover::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: #8a8a8a; /* Un tono más claro cuando se hace hover */
  z-index: -1; /* Coloca detrás del texto pero delante de cualquier contenido inferior */
  pointer-events: none; /* Para que no interfiera con los eventos del mouse */
}

.sort-icon {
  position: relative; /* Asegura que esté encima del fondo */
  z-index: 1;
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
  background: rgba(255, 255, 255, 0.7); /* Fondo claro semi-transparente */
  color: #333; /* Texto oscuro para contraste */
  border: 1px solid rgba(173, 169, 150, 0.3); /* Borde más acorde al tema */
}

tr:hover td {
  background: rgba(173, 169, 150, 0.2); /* Hover sutil acorde al fondo */
  color: #000; /* Texto más oscuro en hover */
  transition: background-color 0.3s ease;
  box-shadow: inset 0 0 5px rgba(173, 169, 150, 0.2); /* Sombra interna sutil */
}

@media (max-width: 768px) {
  .table-container {
    overflow-x: auto !important; /* Permitir scroll horizontal */
    overflow-y: hidden !important; /* NO permitir scroll vertical dentro de la tabla */
    width: 100% !important;
    max-width: 100% !important;
    height: calc(100vh - 120px) !important; /* Ocupar el resto de la pantalla móvil */
    padding: 0 !important;
    margin: 0 !important;
    border-radius: 4px !important;
    background: #D1CEC8 !important;
  }

  .stock-table {
    width: 600px !important; 
    min-width: 600px !important;
    table-layout: fixed !important;
  }

  th {
    position: sticky !important;
    top: 0 !important;
    font-size: 0.7rem !important;
    padding: 0.25rem !important;
    background: #646464 !important;
    z-index: 10 !important;
  }

  td {
    font-size: 0.7rem !important;
    padding: 0.25rem !important;
    white-space: nowrap !important;
    overflow: hidden !important;
    text-overflow: ellipsis !important;
    background: rgba(255, 255, 255, 0.7) !important;
  }
  
  /* Ajustar anchos para móvil */
  th:nth-child(1), td:nth-child(1) { width: 10% !important; } /* Ticker */
  th:nth-child(2), td:nth-child(2) { width: 15% !important; } /* Compañía */
  th:nth-child(3), td:nth-child(3) { width: 15% !important; } /* Brokerage */
  th:nth-child(4), td:nth-child(4) { width: 12% !important; } /* Acción */
  th:nth-child(5), td:nth-child(5) { width: 12% !important; } /* Rating Ant. */
  th:nth-child(6), td:nth-child(6) { width: 12% !important; } /* Rating Act. */
  th:nth-child(7), td:nth-child(7) { width: 8% !important; } /* Precio Desde */
  th:nth-child(8), td:nth-child(8) { width: 8% !important; } /* Precio Hasta */
  th:nth-child(9), td:nth-child(9) { width: 8% !important; } /* Fecha */
}
</style>