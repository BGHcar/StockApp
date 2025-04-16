<!-- filepath: c:\Users\User\Documents\stock-analyzer\frontend\src\components\StockSearch.vue -->
<template>
  <BaseTable>
    <div class="search-container">
      <form @submit.prevent="onSearch" class="search-form">
        <select v-model="type" class="search-select">
          <option v-for="option in searchOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>

        <!-- Input normal para otros tipos de búsqueda -->
        <input 
          v-if="type !== 'price'" 
          v-model="query" 
          type="text" 
          class="search-input" 
          :placeholder="placeholder" 
        />

        <!-- Inputs para rango de precios -->
        <div v-else class="price-range-container">
          <input 
            v-model="minPrice" 
            type="number" 
            class="price-input" 
            placeholder="Mínimo" 
            min="0"
            step="0.01"
          />
          <span class="price-separator">a</span>
          <input 
            v-model="maxPrice" 
            type="number" 
            class="price-input" 
            placeholder="Máximo" 
            min="0"
            step="0.01"
          />
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
  { value: 'rating', label: 'Rating' },
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
.search-container {
  padding: 0.5rem;
}

.search-form {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  justify-content: space-between;
}

.search-select,
.search-input,
.btn {
  padding: 0.35rem 0.75rem;
  border-radius: 0.5rem;
  border: 1px solid rgba(173, 83, 137, 0.3);
  background: rgba(60, 16, 83, 0.9);
  color: white;
}

.search-input {
  flex: 1;
  min-width: 200px;
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
  border: 1px solid rgba(173, 83, 137, 0.3);
  background: rgba(60, 16, 83, 0.9);
  color: white;
}

.price-separator {
  color: white;
  font-weight: bold;
}

.button-group {
  display: flex;
  gap: 0.35rem;
}

.btn {
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary {
  background: rgba(173, 83, 137, 0.8);
}

.btn-primary:hover {
  background: rgba(173, 83, 137, 1);
}

.btn-secondary {
  background: transparent;
  border: 1px solid rgba(173, 83, 137, 0.3);
}

.btn-secondary:hover {
  background: rgba(173, 83, 137, 0.2);
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
</style>