export interface Recommendation {
  ticker: string;
  company: string;
  score: number;
  reason: string;
  last_update: string; // Mantener como string, se formateará en el componente
}