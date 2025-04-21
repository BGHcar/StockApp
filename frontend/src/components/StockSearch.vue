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
.search-container {
  padding: 0;
  background: #D1CEC8;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(173, 169, 150, 0.25);
  border: 2px solid #646464; /* Borde más oscuro y consistente */
}

.search-form {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem; /* Mover el padding aquí */
}

.search-select,
.search-input,
.btn {
  padding: 0.35rem 0.75rem;
  border-radius: 0.5rem;
  border: 1px solid rgba(173, 169, 150, 0.5);
  background: #646464;
  color: white; /* Asegurar que el texto sea blanco en todos estos elementos */
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
  gap: 0.35rem;
}

.btn {
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: bold;
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
.search-select:focus,
.search-input:focus,
.price-input:focus,
.btn:focus {
  outline: none;
  border-color: #646464;
  box-shadow: 0 0 0 2px rgba(173, 169, 150, 0.3); /* Sombra de foco acorde al tema */
}

@media (max-width: 768px) {
  .search-form {
    flex-direction: column; /* Cambiar a diseño vertical */
    gap: 0.75rem;
  }

  .search-select,
  .search-input,
  .price-input,
  .btn {
    width: 100%; /* Asegurar que ocupen todo el ancho */
    font-size: 0.9rem; /* Reducir el tamaño de fuente */
  }

  .price-range-container {
    flex-direction: column; /* Cambiar a diseño vertical */
    gap: 0.5rem;
  }

  .button-group {
    flex-direction: column; /* Botones en columna */
    gap: 0.5rem;
  }

  .search-container {
    padding: 0.5rem !important;
    width: 100% !important;
  }

  .search-form {
    display: grid !important;
    grid-template-columns: 1fr 1fr !important;
    gap: 0.5rem !important;
    width: 100% !important;
  }

  .search-select {
    grid-column: 1 / -1 !important; /* Span en toda la fila */
    width: 100% !important;
    padding: 0.25rem !important;
    font-size: 0.8rem !important;
  }

  .search-input, 
  .price-input {
    padding: 0.25rem !important;
    font-size: 0.8rem !important;
    width: 100% !important;
    min-width: 0 !important;
  }

  .price-range-container {
    grid-column: 1 / -1 !important;
    display: grid !important;
    grid-template-columns: 1fr auto 1fr !important;
    gap: 0.5rem !important;
    align-items: center !important;
    width: 100% !important;
  }

  .button-group {
    grid-column: 1 / -1 !important;
    display: grid !important;
    grid-template-columns: 1fr 1fr !important;
    gap: 0.5rem !important;
    width: 100% !important;
  }

  .btn {
    padding: 0.25rem !important;
    font-size: 0.8rem !important;
    width: 100% !important;
  }

  /* ... mantener estilos existentes ... */
  
  .search-form {
    display: grid !important;
    grid-template-columns: 1fr !important; /* Cambiar a una columna */
    gap: 0.5rem !important;
    width: 100% !important;
    overflow-x: hidden !important;
    max-width: 100% !important; /* Importante para evitar desbordamiento */
  }

  .search-select {
    grid-column: 1 !important; /* Ocupar solo una columna */
    width: 100% !important;
    padding: 0.25rem !important;
    font-size: 0.8rem !important;
    max-width: 100% !important; /* Evitar desbordamiento */
    text-overflow: ellipsis !important; /* Cortar texto si es necesario */
  }

  .search-input, 
  .price-input {
    grid-column: 1 !important;
    padding: 0.25rem !important;
    font-size: 0.8rem !important;
    width: 100% !important;
    min-width: 0 !important;
    max-width: 100% !important; /* Evitar desbordamiento */
  }

  .price-range-container {
    grid-column: 1 !important;
    display: grid !important;
    grid-template-columns: 1fr auto 1fr !important;
    gap: 0.5rem !important;
    align-items: center !important;
    width: 100% !important;
    max-width: 100% !important; /* Evitar desbordamiento */
  }

  .button-group {
    grid-column: 1 !important;
    display: grid !important;
    grid-template-columns: 1fr 1fr !important;
    gap: 0.5rem !important;
    width: 100% !important;
    max-width: 100% !important; /* Evitar desbordamiento */
  }
}
</style>