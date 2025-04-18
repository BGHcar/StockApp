// frontend/src/components/StockRecommendations.vue
<template>
  <BaseTable>
    <div class="recommendations-container">
      <h2 class="section-title">Recomendaciones de Inversión</h2>
      
      <div class="recommendation-cards">
        <div class="card top-pick" v-if="topPick">
          <h3>Mejor Oportunidad</h3>
          <div class="ticker">{{ topPick.ticker }}</div>
          <div class="company">{{ topPick.company }}</div>
          <div class="details">
            <div class="rating">Rating: {{ topPick.rating_to || 'N/A' }}</div>
            <div class="target">Precio Objetivo: {{ topPick.target_to }}</div>
            <div class="score">Puntuación: {{ topPick.score.toFixed(2) }}</div>
          </div>
          <div class="reason">{{ topPick.reason }}</div>
        </div>
        
        <div class="card avoid" v-if="avoid">
          <h3>Evitar</h3>
          <div class="ticker">{{ avoid.ticker }}</div>
          <div class="company">{{ avoid.company }}</div>
          <div class="details">
            <div class="rating">Rating: {{ avoid.rating_to || 'N/A' }}</div>
            <div class="target">Precio Objetivo: {{ avoid.target_to }}</div>
            <div class="score">Puntuación: {{ avoid.score.toFixed(2) }}</div>
          </div>
          <div class="reason">{{ avoid.reason }}</div>
        </div>
      </div>
    </div>
  </BaseTable>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Stock } from '@/types/stock'
import BaseTable from './base/BaseTable.vue'

const props = defineProps<{
  stocks: Stock[]
}>()

// Sistema de puntuación para evaluar acciones
function calculateScore(stock: Stock): { score: number, reason: string } {
  let score = 0
  let reason = ''
  
  // 1. Evaluación por Rating
  const ratingScores: Record<string, number> = {
    'Buy': 2,
    'Overweight': 1.5,
    'Equal Weight': 0,
    'Hold': 0,
    'Underweight': -1.5,
    'Sell': -2
  }
  
  if (stock.rating_to && ratingScores[stock.rating_to]) {
    score += ratingScores[stock.rating_to]
    
    if (ratingScores[stock.rating_to] > 0) {
      reason += `Calificación positiva (${stock.rating_to}) de ${stock.brokerage}. `
    } else if (ratingScores[stock.rating_to] < 0) {
      reason += `Calificación negativa (${stock.rating_to}) de ${stock.brokerage}. `
    }
  }
  
  // 2. Evaluación por cambio en el precio objetivo
  if (stock.target_from && stock.target_to) {
    const targetFrom = parseFloat(stock.target_from.replace(/[^\d.-]/g, ''))
    const targetTo = parseFloat(stock.target_to.replace(/[^\d.-]/g, ''))
    
    if (!isNaN(targetFrom) && !isNaN(targetTo)) {
      const percentChange = ((targetTo - targetFrom) / targetFrom) * 100
      
      // Asignar puntos basados en el cambio porcentual
      if (percentChange > 10) {
        score += 2
        reason += `Aumento significativo del precio objetivo (+${percentChange.toFixed(1)}%). `
      } else if (percentChange > 0) {
        score += 1
        reason += `Ligero aumento del precio objetivo (+${percentChange.toFixed(1)}%). `
      } else if (percentChange < -10) {
        score -= 2
        reason += `Reducción significativa del precio objetivo (${percentChange.toFixed(1)}%). `
      } else if (percentChange < 0) {
        score -= 1
        reason += `Ligera reducción del precio objetivo (${percentChange.toFixed(1)}%). `
      }
    }
  }
  
  // 3. Evaluación por acción realizada
  if (stock.action) {
    if (stock.action.includes('raised')) {
      score += 1.5
      reason += 'Aumento del precio objetivo. '
    } else if (stock.action.includes('lowered')) {
      score -= 1.5
      reason += 'Reducción del precio objetivo. '
    } else if (stock.action.includes('reiterated')) {
      score += 0.5
      reason += 'Reiteración del precio objetivo. '
    }
  }
  
  // 4. Factor de reputación de la correduría
  const brokerageReputation: Record<string, number> = {
    'Morgan Stanley': 0.5,
    'JPMorgan Chase & Co.': 0.5,
    'Goldman Sachs': 0.5,
    'Wells Fargo & Company': 0.3,
    'Needham & Company LLC': 0.3,
    'Truist Financial': 0.2
  }
  
  if (stock.brokerage && brokerageReputation[stock.brokerage]) {
    score += brokerageReputation[stock.brokerage]
    reason += `Análisis de una correduría de alta reputación (${stock.brokerage}). `
  }
  
  // 5. Actualidad de la recomendación (priorizando las más recientes)
  const today = new Date()
  const stockDate = new Date(stock.time)
  const daysDifference = Math.floor((today.getTime() - stockDate.getTime()) / (1000 * 3600 * 24))
  
  if (daysDifference <= 7) {
    score += 0.5
    reason += 'Recomendación reciente (última semana). '
  }
  
  return { score, reason }
}

// Procesamiento de los stocks para añadir puntuaciones
const scoredStocks = computed(() => {
  return props.stocks.map(stock => {
    const { score, reason } = calculateScore(stock)
    return {
      ...stock,
      score,
      reason
    }
  }).sort((a, b) => b.score - a.score) // Ordenar por puntuación descendente
})

// Mejores y peores recomendaciones
const topPick = computed(() => scoredStocks.value.length > 0 ? scoredStocks.value[0] : null)
const avoid = computed(() => scoredStocks.value.length > 0 ? scoredStocks.value[scoredStocks.value.length - 1] : null)
</script>

<style scoped>
.recommendations-container {
  padding: 0;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 8px;
  box-shadow: 0 2px 6px rgba(173, 169, 150, 0.3);
  border: 2px solid #646464; /* Borde más oscuro y consistente */
  /* Eliminar desbordamiento para evitar problemas con esquinas redondeadas */
  overflow: hidden;
}

.section-title {
  padding: 0.4rem 0.5rem;
  /* Eliminar border-radius para evitar inconsistencias con el contenedor padre */
  border-radius: 0;
  font-size: 1.4rem;
  color: white;
  /* Eliminar margin-bottom para evitar espacios innecesarios */
  margin: 0;
  text-align: center;
  font-weight: 900;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.4);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  background: #646464;
  /* Añadir borde inferior para separar visualmente del contenido */
  border-bottom: 2px solid #646464;
}

.recommendation-cards {
  display: flex;
  justify-content: space-between;
  /* Ajustar padding para eliminar espacio al fondo */
  padding: 0.5rem 0.5rem 0; /* Quitar padding inferior */
  align-items: stretch;
  /* Evitar desbordamiento */
  overflow: hidden;
}

/* Ajustar márgenes de las tarjetas para alinear correctamente */
.card {
  flex: 1;
  margin: 0 0.25rem 0.5rem; /* Añadir margen inferior para compensar */
  border-radius: 8px;
  padding: 0.75rem;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
}

/* Ajustar la primera tarjeta para eliminar margen izquierdo */
.card:first-child {
  margin-left: 0;
}

/* Ajustar la última tarjeta para eliminar margen derecho */
.card:last-child {
  margin-right: 0;
}

/* El resto del código permanece igual */
.top-pick {
  background: rgba(77, 174, 128, 0.7); /* Verde más intenso */
  border: 2px solid #4DAE80; /* Borde verde más visible */
}

.avoid {
  background: rgba(239, 68, 68, 0.7); /* Rojo más intenso */
  border: 2px solid #EF4444; /* Borde rojo más visible */
}

/* El resto del código se mantiene igual */
.ticker {
  font-size: 1.6rem;
  font-weight: 700;
  margin-bottom: 0.2rem;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
  text-align: left; /* Texto alineado a la izquierda */
  letter-spacing: 0.05em;
}

.company {
  font-size: 0.95rem;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 0.5rem;
  text-align: left; /* Texto alineado a la izquierda */
  font-style: italic;
}

.details {
  margin-bottom: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.3);
  text-align: left; /* Asegurar que todo está alineado a la izquierda */
}

.rating, .target, .score {
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.95);
  margin-bottom: 0.2rem;
  font-weight: 500;
  text-align: left; /* Texto alineado a la izquierda */
}

.reason {
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.9);
  line-height: 1.4;
  padding-top: 0.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.3);
  margin-top: auto;
  text-align: left; /* Texto alineado a la izquierda */
}
</style>