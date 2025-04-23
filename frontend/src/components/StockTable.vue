<template>
  <BaseTable>
    <BaseScroll class="table-container">
      <table class="stock-table">
        <thead>
          <tr>
            <th v-for="header in headers" :key="header.key" :class="[header.class, 'sortable']"
              @click="handleSort(header.key)">
              {{ header.label }}
              <span class="sort-icon" v-if="currentSortField === header.key">
                {{ currentSortIndicator }}
              </span>
            </th>
          </tr>
        </thead>
        <tbody>
          <!-- Mostrar mensaje si no hay stocks en la página actual -->
          <tr v-if="props.stocks.length === 0">
            <td :colspan="headers.length" class="no-data-cell">No hay datos para mostrar en esta página.</td>
          </tr>
          <!-- Renderizar filas de stocks -->
          <tr v-for="stock in props.stocks" :key="stock.ticker + stock.time" v-else>
            <td v-for="header in headers" :key="header.key" :class="header.class">
              {{ formatValue(stock, header.key) }}
            </td>
          </tr>
        </tbody>
      </table>
      <div v-if="pagination.totalPages > 0" class="pagination-container">
        <div class="pagination-controls">
          <button @click="goToPage(pagination.currentPage - 1)" :disabled="pagination.currentPage <= 1"
            class="pagination-button">
            &lt; Anterior
          </button>
          <span class="pagination-info">
            &nbsp Página {{ pagination.currentPage }} de {{ pagination.totalPages }} &nbsp
          </span>
          <button @click="goToPage(pagination.currentPage + 1)"
            :disabled="pagination.currentPage >= pagination.totalPages" class="pagination-button">
            Siguiente &gt;
          </button>
        </div>
      </div>
    </BaseScroll>
  </BaseTable>
</template>

<script setup lang="ts">
import { defineProps, computed } from 'vue'
import type { Stock } from '@/types/stock'
import type { TableHeader } from '@/types/table'
import BaseTable from './base/BaseTable.vue'
import BaseScroll from './base/BaseScroll.vue'
import { useStockStore } from '@/stores/stockStore' // Importar el store

const props = defineProps<{
  stocks: Stock[], // Estos son los stocks de la página actual
  headers: TableHeader[]
}>()

const stockStore = useStockStore() // Usar el store
const pagination = computed(() => stockStore.pagination) // Acceder al estado de paginación

// Obtenemos el campo de ordenamiento actual del store
const currentSortField = computed(() => stockStore.currentSortField);

// Indicador visual de ordenamiento
const currentSortIndicator = computed(() => '▼'); // Siempre mostramos triángulo hacia abajo ya que el backend maneja la dirección

// Función para formatear valores (incluyendo fechas)
function formatValue(stock: Stock, key: keyof Stock): string {
  const value = stock[key];
  if (key === 'time') {
    const date = new Date(value as string);
    if (isNaN(date.getTime())) return value as string; // Devolver original si no es fecha válida
    return date.toLocaleDateString('es-ES', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit'
    });
  }
  // Devolver el valor como string para otras claves
  return String(value);
}

// Función para manejar el ordenamiento desde el backend
async function handleSort(field: keyof Stock) {
  // Llamamos a la acción del store para ordenar
  // Pasamos el campo de ordenamiento y preservamos la búsqueda actual si existe
  if (stockStore.currentSearchType !== 'sortBy' && stockStore.currentSearchType !== 'loadAll') {
    // Si hay una búsqueda activa que no es de ordenamiento, enviar la consulta actual
    await stockStore.loadSortedStocks(field as string, stockStore.currentSearchQuery);
  } else {
    // Si no hay búsqueda o ya estamos en modo ordenamiento, solo enviar el campo
    await stockStore.loadSortedStocks(field as string);
  }
}

// Función para llamar a la acción del store para cambiar de página
function goToPage(page: number) {
  stockStore.goToPage(page);
}
</script>

<style scoped>
.stock-table {
  width: 100%;
  min-width: 800px;
  table-layout: fixed;
  border-collapse: separate;
  /* Cambiado de collapse a separate para mejor control de bordes */
  border-spacing: 0;
  /* Eliminar espacios entre celdas */
  background: transparent;
  /* Cambiado para que use el gradiente del contenedor */
}


.no-data-cell {
  text-align: center;
  padding: 1rem;
  font-size: 1.2rem;
  color: #646464;
  /* Color de texto acorde al tema */
  background: rgba(255, 255, 255, 0.5);
}

/* Definimos anchos específicos para cada columna */
th,
td {
  padding: 0.35rem 0.5rem;
  text-align: left;
  border: 1px solid rgba(173, 169, 150, 0.3);
  /* Color de borde que complementa el gradiente */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Asegurar que todas las celdas de cabecera tengan bordes */
th {
  position: sticky;
  top: 0;
  background: #646464;
  /* Color sólido más oscuro que armoniza con el gradiente */
  font-family: Arial, Helvetica, sans-serif;
  font-weight: 700;
  color: white;
  text-transform: uppercase;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  border: 1px solid #4a4a4a;
  border-bottom: 2px solid #4a4a4a;
  /* Borde inferior más visible */
  z-index: 10;
  box-shadow: 0 2px 4px rgba(173, 169, 150, 0.4);
  /* Sombra acorde al fondo */
}

/* Específicamente para la columna Rating Act. y Precio Desde */
th:nth-child(6),
th:nth-child(7) {
  border-right: 2px solid #4a4a4a !important;
  /* Línea más gruesa y visible */
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
  background: #8a8a8a;
  /* Un tono más claro cuando se hace hover */
  z-index: -1;
  /* Coloca detrás del texto pero delante de cualquier contenido inferior */
  pointer-events: none;
  /* Para que no interfiera con los eventos del mouse */
}

.sort-icon {
  position: relative;
  /* Asegura que esté encima del fondo */
  z-index: 1;
  margin-left: 4px;
  font-size: 0.7em;
}


td {
  background: rgba(255, 255, 255, 0.7);
  /* Fondo claro semi-transparente */
  color: #333;
  /* Texto oscuro para contraste */
  border: 1px solid rgba(173, 169, 150, 0.3);
  /* Borde más acorde al tema */
}

tr:hover td {
  background: rgba(173, 169, 150, 0.2);
  /* Hover sutil acorde al fondo */
  color: #000;
  /* Texto más oscuro en hover */
  transition: background-color 0.3s ease;
  box-shadow: inset 0 0 5px rgba(173, 169, 150, 0.2);
  /* Sombra interna sutil */
}

.pagination-container {
  width: 100%;
  display: flex;
  justify-content: center;
  padding: 8px 0;
  background-color: #646464;
  font-family: 15px Arial, Helvetica, sans-serif;
  border: 1px solid rgba(173, 169, 150, 0.3);
  text-transform: uppercase;
}

.pagination-info {
  font-family: Arial, Helvetica, sans-serif;
  font-weight: 700;
  color: white;
}

.pagination-button {
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: bold;
  padding: 0.35rem 0.75rem;
  border-radius: 0.5rem;
  border: 1px solid rgba(173, 169, 150, 0.5);
  background: #4a4a4a;
  color: white;
  box-shadow: 0 1px 3px rgba(173, 169, 150, 0.3);
}

.pagination-button:hover {
  background: #333;
  border: 1px solid #646464;
}

.pagination-button:focus {
  outline: none;
  border-color: #646464;
  box-shadow: 0 0 0 2px rgba(173, 169, 150, 0.3);
}
</style>