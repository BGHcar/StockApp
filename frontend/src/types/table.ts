import type { Stock } from '@/types/stock'
export interface TableHeader {
  key: keyof Stock; // Esto asegura que solo uses propiedades válidas de Stock
  label: string;
  class?: string;
}