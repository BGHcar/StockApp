import type { Stock } from '@/types/stock'
export interface TableHeader {
  // creo un key que filtre todo menos la fecha
  key: keyof Stock; // Esto asegura que solo uses propiedades válidas de Stock eliminando la fecha
  label: string;
  class?: string;
}