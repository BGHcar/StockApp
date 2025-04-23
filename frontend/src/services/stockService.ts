import type { Stock } from '@/types/stock'
// Importar el tipo Recommendation
import type { Recommendation } from '@/types/recommendation'

const API_URL = import.meta.env.VITE_API_URL

// Define la estructura de la respuesta paginada que esperas del backend
export interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    page: number;
    pageSize: number;
    totalItems: number;
    totalPages: number;
  };
}

// Función para construir URLs con paginación
function buildUrlWithPagination(baseUrl: string, page: number, pageSize: number): string {
  const url = new URL(baseUrl);
  url.searchParams.append('page', page.toString());
  url.searchParams.append('pageSize', pageSize.toString());
  return url.toString();
}

export async function fetchStocks(page: number, pageSize: number): Promise<PaginatedResponse<Stock>> {
  try {
    const url = buildUrlWithPagination(`${API_URL}/stocks`, page, pageSize);
    const response = await fetch(url);
    if (!response.ok) throw new Error(`Error al obtener stocks: ${response.status}`);
    const data: PaginatedResponse<Stock> = await response.json();
    // Asegurarse de que la respuesta tenga la estructura esperada
    return data || { items: [], pagination: { page: 1, pageSize: pageSize, totalItems: 0, totalPages: 0 } };
  } catch (error) {
    console.error('Error al obtener stocks:', error);
    // Devolver estructura vacía en caso de error
    return { items: [], pagination: { page: 1, pageSize: pageSize, totalItems: 0, totalPages: 0 } };
  }
}

export async function searchStocks(type: string, query: string, page: number, pageSize: number, activeSearchQuery?: string): Promise<PaginatedResponse<Stock>> {
  let baseUrl = '';
  try {
    switch (type) {
      case 'ticker':
        baseUrl = `${API_URL}/stocks/ticker/${encodeURIComponent(query)}`;
        break;
      case 'action':
        baseUrl = `${API_URL}/stocks/action/${encodeURIComponent(query)}`;
        break;
      case 'ratingTo':
        baseUrl = `${API_URL}/stocks/ratingto/${encodeURIComponent(query)}`;
        break;
      case 'ratingFrom':
        baseUrl = `${API_URL}/stocks/ratingfrom/${encodeURIComponent(query)}`;
        break;
      case 'brokerage':
        baseUrl = `${API_URL}/stocks/brokerage/${encodeURIComponent(query)}`;
        break;
      case 'company':
        baseUrl = `${API_URL}/stocks/company/${encodeURIComponent(query)}`;
        break;
      case 'sortBy':
        // Para ordenamiento siempre usamos el endpoint /stocks/sorted/{field}
        baseUrl = `${API_URL}/stocks/sorted/${encodeURIComponent(query)}`;
        
        // Si hay una búsqueda activa proporcionada, añadirla como parámetro query
        if (activeSearchQuery && activeSearchQuery.trim() !== '') {
          // Construir URL con el parámetro query primero, luego los parámetros de paginación
          const url = new URL(baseUrl);
          url.searchParams.append('query', activeSearchQuery);
          url.searchParams.append('page', page.toString());
          url.searchParams.append('pageSize', pageSize.toString());

          console.log (url.toString())
          
          const response = await fetch(url.toString());
          if (!response.ok) {
            return { items: [], pagination: { page, pageSize, totalItems: 0, totalPages: 0 } };
          }
          
          return await response.json();
        }
        break;
      case 'price':
        const [min, max] = query.split('-').map(Number);
        baseUrl = `${API_URL}/stocks/price-range/${min}/${max}`;
        break;
      case 'general':
      default:
        baseUrl = `${API_URL}/stocks/search/${encodeURIComponent(query)}`;
        break;
    }

    // Para todos los demás casos, usar la URL base con paginación
    const url = buildUrlWithPagination(baseUrl, page, pageSize);
    const response = await fetch(url);

    // Si no hay resultados, devolver vacío
    if (!response.ok) {
      return { items: [], pagination: { page, pageSize, totalItems: 0, totalPages: 0 } };
    }

    const data: PaginatedResponse<Stock> = await response.json();
    return data || { items: [], pagination: { page, pageSize, totalItems: 0, totalPages: 0 } };

  } catch (error) {
    console.error(`Error en búsqueda por ${type} (${query}) pag ${page}:`, error);
    return { items: [], pagination: { page, pageSize, totalItems: 0, totalPages: 0 } };
  }
}

// --- Nueva función para obtener recomendaciones ---
export async function fetchRecommendations(limit: number = 5): Promise<Recommendation[]> {
  try {
    const url = new URL(`${API_URL}/recommendations`);
    url.searchParams.append('limit', limit.toString());
    const response = await fetch(url.toString());
    if (!response.ok) {
      // Intentar leer el mensaje de error del backend si existe
      let errorMsg = `Error fetching recommendations: ${response.status}`;
      try {
        const errorBody = await response.json();
        if (errorBody && errorBody.error) {
          errorMsg += ` - ${errorBody.error}`;
        }
      } catch (e) { /* Ignorar si el cuerpo no es JSON */ }
      throw new Error(errorMsg);
    }
    const data: Recommendation[] = await response.json();
    return data || []; // Devolver array vacío si la respuesta es null/undefined
  } catch (error) {
    console.error('Error fetching recommendations:', error);
    // Relanzar el error para que el store lo maneje
    throw error;
  }
}

