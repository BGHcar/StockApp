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
<!--                 <template v-if = "header.key != 'time'">
                  {{ stock[header.key] }}
                </template>
                <template v-else>
                </template> -->
                {{ formatDate(stock[header.key]) }}
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


function formatDate (dateString: string) : string{
  const date = new Date(dateString);
  
  if (isNaN(date.getTime())) return dateString;

  return date.toLocaleDateString('es-ES', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });

}
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
  border: 2px solid #646464;
  border-radius: 8px;
  max-height: calc(100vh - 260px);
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
</style>