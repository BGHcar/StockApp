// frontend/src/services/analysisService.ts
import type { Stock } from '@/types/stock';

export interface StockAnalysis {
  recommendation: string;
  confidence: number;
  metrics: {
    trendDirection: 'up' | 'down' | 'neutral';
    riskLevel: 'low' | 'medium' | 'high';
    potentialReturn: number;
  };
}

export function analyzeStock(stock: Stock): StockAnalysis {
  // Implementar lógica de análisis
  return {
    recommendation: 'buy',
    confidence: 0.85,
    metrics: {
      trendDirection: 'up',
      riskLevel: 'medium',
      potentialReturn: 15.5
    }
  };
}

export function getTopRecommendations(stocks: Stock[]): Stock[] {
  // Ordenar por rating y fecha más reciente
  return stocks
    .filter(stock => stock.rating_to.toLowerCase().includes('buy') || 
                     stock.rating_to.toLowerCase().includes('outperform'))
    .sort((a, b) => {
      // Primero por rating
      const ratingA = getRatingScore(a.rating_to);
      const ratingB = getRatingScore(b.rating_to);
      if (ratingB !== ratingA) {
        return ratingB - ratingA;
      }
      // Luego por fecha más reciente
      return new Date(b.time).getTime() - new Date(a.time).getTime();
    })
    .slice(0, 5); // Top 5 recomendaciones
}

function getRatingScore(rating: string): number {
  const lowercaseRating = rating.toLowerCase();
  if (lowercaseRating.includes('strong buy')) return 5;
  if (lowercaseRating.includes('buy')) return 4;
  if (lowercaseRating.includes('outperform')) return 3;
  if (lowercaseRating.includes('neutral')) return 2;
  if (lowercaseRating.includes('hold')) return 1;
  return 0;
}