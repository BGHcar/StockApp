import type { Stock } from '@/types/stock'

const API_URL = import.meta.env.VITE_API_URL

export async function fetchStocks(): Promise<Stock[]> {
  try {
    const response = await fetch(`${API_URL}/stocks`)
    if (!response.ok) throw new Error(`Error al obtener stocks: ${response.status}`)
    const data = await response.json()
    return data || []
  } catch (error) {
    console.error('Error al obtener stocks:', error)
    return []
  }
}

export async function searchStocks(type: string, query: string): Promise<Stock[]> {
  let url = ''
  try {
    switch (type) {
      case 'ticker':
        url = `${API_URL}/stocks/${encodeURIComponent(query)}`
        break
      case 'action':
        url = `${API_URL}/stocks/action/${encodeURIComponent(query)}`
        break
      case 'ratingTo':
        url = `${API_URL}/stocks/ratingto/${encodeURIComponent(query)}`
        break
      case 'ratingFrom':
        url = `${API_URL}/stocks/ratingfrom/${encodeURIComponent(query)}`
        break
      case 'brokerage':
        url = `${API_URL}/stocks/brokerage/${encodeURIComponent(query)}`
        break
      case 'company':
        url = `${API_URL}/stocks/company/${encodeURIComponent(query)}`
        break
      case 'price':
        const [min, max] = query.split('-').map(Number)
        url = `${API_URL}/stocks/price-range/${min}/${max}`
        break
      case 'general':
      default:
        url = `${API_URL}/stocks/search/${encodeURIComponent(query)}`
        break
    }
    
    const response = await fetch(url)
    if (!response.ok) {
      console.warn(`Búsqueda por ${type} (${query}) no devolvió resultados`)
      return []
    }
    
    const data = await response.json()
    
    // Asegurarse de manejar null/undefined correctamente
    if (!data) {
      console.warn(`Búsqueda por ${type} (${query}) devolvió datos nulos`)
      return []
    }
    
    // Convertir a array si no lo es
    const stocksArray = Array.isArray(data) ? data : [data]
    console.log(`Búsqueda por ${type} (${query}) encontró ${stocksArray.length} resultados`)
    
    return stocksArray
  } catch (error) {
    console.error(`Error en búsqueda por ${type} (${query}):`, error)
    return []
  }
}

