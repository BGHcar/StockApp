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

        <input v-model="query" type="text" class="search-input" :placeholder="placeholder" />

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
import { ref, computed } from 'vue'
import BaseTable from './base/BaseTable.vue'

const searchOptions = [
  { value: 'general', label: 'General' },
  { value: 'ticker', label: 'Ticker' },
  { value: 'company', label: 'Compañía' },
  { value: 'brokerage', label: 'Brokerage' },
  { value: 'action', label: 'Acción' },
  { value: 'rating', label: 'Rating' }
]

const emit = defineEmits<{
  (e: 'search', type: string, query: string): void
  (e: 'reset'): void
}>()

const type = ref('general')
const query = ref('')

const placeholder = computed(() => `Buscar por ${searchOptions.find(opt => opt.value === type.value)?.label.toLowerCase()}...`)

function onSearch() {
  // Si la consulta está vacía, emitir reset en lugar de search
  if (query.value.trim() === '') {
    onReset();
    return;
  }
  emit('search', type.value, query.value)
}

function onReset() {
  query.value = ''
  type.value = 'general'
  emit('reset')
}
</script>

<style scoped>
.search-container {
  padding: 0.5rem;
  /* Reducido de 1rem a 0.5rem */
}

.search-form {
  display: flex;
  gap: 0.5rem;
  /* Reducido de 1rem a 0.5rem */
  align-items: center;
  justify-content: space-between;
}

.search-select,
.search-input,
.btn {
  padding: 0.35rem 0.75rem;
  /* Reducido de 0.5rem 1rem */
  border-radius: 0.5rem;
  border: 1px solid rgba(173, 83, 137, 0.3);
  background: rgba(60, 16, 83, 0.9);
  color: white;
}

.search-input {
  flex: 1;
  min-width: 200px;
}

.button-group {
  display: flex;
  gap: 0.35rem;
  /* Reducido de 0.5rem a 0.35rem */
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
</style>