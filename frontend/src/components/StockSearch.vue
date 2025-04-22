<!-- filepath: c:\Users\User\Documents\stock-analyzer\frontend\src\components\StockSearch.vue -->
<template>
  <BaseTable>
    <div class="search-container">
      <form @submit.prevent="onSearch" class="search-form">
        <div class="search-select">
          <select v-model="type" class="select">
            <option v-for="option in searchOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>

          <!-- Input normal para otros tipos de búsqueda -->
          <input v-if="type !== 'price'" v-model="query" type="text" class="search-input" :placeholder="placeholder" />

          <!-- Inputs para rango de precios -->
          <div v-else class="price-range-container">
            <input v-model="minPrice" type="number" class="price-input" placeholder="Mínimo" min="0" step="0.01" />
            <span class="price-separator">a</span>
            <input v-model="maxPrice" type="number" class="price-input" placeholder="Máximo" min="0" step="0.01" />
          </div>
        </div>

        <div class="button-group">
          <button type="submit" class="btn btn-primary">
            Buscar
          </button>
          <button type="button" @click="onReset" class="btn btn-secondary">
            Limpiar
          </button>
        </div>
      </form>
    </div>
  </BaseTable>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import BaseTable from './base/BaseTable.vue'

const searchOptions = [
  { value: 'general', label: 'General' },
  { value: 'ticker', label: 'Ticker' },
  { value: 'company', label: 'Compañía' },
  { value: 'brokerage', label: 'Brokerage' },
  { value: 'action', label: 'Acción' },
  { value: 'ratingFrom', label: 'Rating Ant' },
  { value: 'ratingTo', label: 'Rating Act' },
  { value: 'price', label: 'Precio' }
]

const emit = defineEmits<{
  (e: 'search', type: string, query: string): void
  (e: 'reset'): void
}>()

const type = ref('general')
const query = ref('')
const minPrice = ref('')
const maxPrice = ref('')

const placeholder = computed(() => `Buscar por ${searchOptions.find(opt => opt.value === type.value)?.label.toLowerCase()}...`)

// Limpia los campos de precio cuando cambia el tipo de búsqueda
watch(type, (newType) => {
  if (newType !== 'price') {
    minPrice.value = ''
    maxPrice.value = ''
  } else {
    query.value = ''
  }
})

function onSearch() {
  if (type.value === 'price') {
    // Validar que ambos campos tengan valores
    if (!minPrice.value || !maxPrice.value) {
      return
    }

    // Para búsqueda por rango de precio, concatenamos los valores con un guion
    emit('search', 'price', `${minPrice.value}-${maxPrice.value}`)
  } else {
    // Si la consulta está vacía para otros tipos, emitir reset
    if (query.value.trim() === '') {
      onReset()
      return
    }
    emit('search', type.value, query.value)
  }
}

function onReset() {
  query.value = ''
  minPrice.value = ''
  maxPrice.value = ''
  type.value = 'general'
  emit('reset')
}
</script>

<style scoped>
.search-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: space-between;
}

.search-select {
  display: flex;
  flex: auto;
  gap: 8px;
}

.select,
.search-input,
.btn {
  padding: 0.35rem 0.75rem;
  border-radius: 0.5rem;
  border: 1px solid rgba(173, 169, 150, 0.5);
  background: #646464;
  color: white;
  box-shadow: 0 1px 3px rgba(173, 169, 150, 0.3);
}

.search-input {
  flex: 1;
  min-width: 200px;
  color: white; /* Asegurar que el color del texto sea blanco */
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.7); /* Placeholder más claro pero aún visible */
}

/* Estilos para el contenedor de rango de precios */
.price-range-container {
  display: flex;
  flex: 1;
  align-items: center;
  gap: 0.5rem;
}

.price-input {
  flex: 1;
  min-width: 80px;
  padding: 0.35rem 0.75rem;
  border-radius: 0.5rem;
  border: 1px solid rgba(173, 169, 150, 0.5); /* Borde acorde al tema */
  background: #646464;
  color: white; /* Añadir color blanco explícitamente */
  box-shadow: 0 1px 3px rgba(173, 169, 150, 0.3); /* Sombra acorde al fondo */
}

.price-input::placeholder {
  color: rgba(255, 255, 255, 0.7); /* Placeholder más claro pero aún visible */
}

.price-separator {
  color: #333; /* Color que contrasta con el fondo del gradiente */
  font-weight: bold;
}

.button-group {
  display: flex;
  justify-content: space-between;
  gap: 0.35rem;
}

.btn {
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: bold;
  flex: 50%;
}

/* Unificar el color de ambos botones */
.btn-primary,
.btn-secondary {
  background: #4a4a4a; /* Mismo color para ambos botones */
}

/* Mantener el mismo comportamiento de hover para ambos */
.btn-primary:hover,
.btn-secondary:hover {
  background: #333; /* Mismo color de hover */
}

/* Remover las flechas de los inputs numéricos */
input[type=number]::-webkit-inner-spin-button, 
input[type=number]::-webkit-outer-spin-button { 
  -webkit-appearance: none; 
  margin: 0; 
}
input[type=number] {
  -moz-appearance: textfield;
  appearance: textfield;
}

/* Añadir un estilo de foco para mejorar la accesibilidad */
.select:focus,
.search-input:focus,
.price-input:focus,
.btn:focus {
  outline: none;
  border-color: #646464;
  box-shadow: 0 0 0 2px rgba(173, 169, 150, 0.3); /* Sombra de foco acorde al tema */
}

@media (min-width: 640px) {
  .search-form {
    flex-direction: row;
    justify-content: space-between;
  }
}
</style>