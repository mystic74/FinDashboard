import axios, { AxiosError } from 'axios';
import type {
  Stock,
  Screener,
  ScreenerResult,
  FilterRequest,
  FilterResponse,
  FilterDefinition,
  SectorPerformance,
  ScreenerSummary,
} from '../types';

const API_BASE_URL = '/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError<{ error?: string }>) => {
    const message = error.response?.data?.error || error.message || 'An error occurred';
    console.error('API Error:', message);
    return Promise.reject(new Error(message));
  }
);

// Screener APIs
export const screenerApi = {
  getAll: async (): Promise<Screener[]> => {
    const { data } = await api.get('/screeners');
    return data.screeners;
  },

  getSummaries: async (): Promise<ScreenerSummary[]> => {
    const { data } = await api.get('/screeners/summary');
    return data.summaries;
  },

  run: async (screenerId: string, country?: string, sector?: string): Promise<ScreenerResult> => {
    const params: Record<string, string> = {};
    if (country) params.country = country;
    if (sector) params.sector = sector;
    const { data } = await api.get(`/screeners/${screenerId}`, { params });
    return data.result;
  },

  runCustom: async (request: FilterRequest): Promise<FilterResponse> => {
    const { data } = await api.post('/screeners/custom', request);
    return data.result;
  },
};

// Stock APIs
export const stockApi = {
  get: async (symbol: string): Promise<Stock> => {
    const { data } = await api.get(`/stocks/${symbol}`);
    return data.stock;
  },

  getFundamentals: async (symbol: string): Promise<Stock> => {
    const { data } = await api.get(`/stocks/${symbol}/fundamentals`);
    return data.stock;
  },

  getMultiple: async (symbols: string[]): Promise<Stock[]> => {
    const { data } = await api.get('/stocks', {
      params: { symbols: symbols.join(',') },
    });
    return data.stocks;
  },

  getQuote: async (symbol: string): Promise<Stock> => {
    const { data } = await api.get(`/stocks/${symbol}/quote`);
    return data.quote;
  },

  getHistory: async (
    symbol: string,
    period: string = '1y'
  ): Promise<{ date: string; close: number; volume: number }[]> => {
    const { data } = await api.get(`/stocks/${symbol}/history`, {
      params: { period },
    });
    return data.prices;
  },
};

// Filter APIs
export const filterApi = {
  getAll: async (): Promise<{ filters: FilterDefinition[]; grouped: Record<string, FilterDefinition[]> }> => {
    const { data } = await api.get('/filters');
    return { filters: data.filters, grouped: data.grouped };
  },

  getCategories: async (): Promise<{ id: string; name: string; description: string }[]> => {
    const { data } = await api.get('/filters/categories');
    return data.categories;
  },

  getMarketCapRanges: async (): Promise<{ id: string; name: string; min: number; max: number }[]> => {
    const { data } = await api.get('/filters/marketcap-ranges');
    return data.ranges;
  },
};

// Sector APIs
export const sectorApi = {
  getPerformance: async (): Promise<SectorPerformance[]> => {
    const { data } = await api.get('/sectors');
    return data.sectors;
  },

  getList: async (): Promise<string[]> => {
    const { data } = await api.get('/sectors/list');
    return data.sectors;
  },
};

export default api;
