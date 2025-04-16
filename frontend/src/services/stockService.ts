import type { Stock } from '@/types/stock'

const API_URL = import.meta.env.VITE_API_URL

export async function fetchStocks(): Promise<Stock[]> {
  const response = await fetch(`${API_URL}/stocks`)
  if (!response.ok) throw new Error('Failed to fetch stocks')
  return response.json()
}

export async function searchStocks(type: string, query: string): Promise<Stock[]> {
  let url = ''
  switch (type) {
    case 'ticker':
      url = `${API_URL}/stocks/${encodeURIComponent(query)}`
      break
    case 'action':
      url = `${API_URL}/stocks/action/${encodeURIComponent(query)}`
      break
    case 'rating':
      url = `${API_URL}/stocks/rating/${encodeURIComponent(query)}`
      break
    case 'brokerage':
      url = `${API_URL}/stocks/brokerage/${encodeURIComponent(query)}`
      break
    case 'general':
    default:
      url = `${API_URL}/stocks/search/${encodeURIComponent(query)}`
      break
  }
  const response = await fetch(url)
  if (!response.ok) throw new Error('No se encontraron resultados')
  const data = await response.json()
  return Array.isArray(data) ? data : [data]
}

