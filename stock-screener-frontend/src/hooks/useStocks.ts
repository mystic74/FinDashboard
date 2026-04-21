import { useQuery } from '@tanstack/react-query';
import { stockApi, filterApi, sectorApi } from '../services/api';
import type { Stock, FilterDefinition, SectorPerformance } from '../types';

export function useStock(symbol: string, enabled: boolean = true) {
  return useQuery<Stock>({
    queryKey: ['stock', symbol],
    queryFn: () => stockApi.get(symbol),
    enabled: enabled && !!symbol,
    staleTime: 60 * 1000,
  });
}

export function useStockFundamentals(symbol: string, enabled: boolean = true) {
  return useQuery<Stock>({
    queryKey: ['stock-fundamentals', symbol],
    queryFn: () => stockApi.getFundamentals(symbol),
    enabled: enabled && !!symbol,
    staleTime: 5 * 60 * 1000,
  });
}

export function useMultipleStocks(symbols: string[], enabled: boolean = true) {
  return useQuery<Stock[]>({
    queryKey: ['stocks', symbols],
    queryFn: () => stockApi.getMultiple(symbols),
    enabled: enabled && symbols.length > 0,
    staleTime: 60 * 1000,
  });
}

export function useStockHistory(symbol: string, period: string = '1y', enabled: boolean = true) {
  return useQuery({
    queryKey: ['stock-history', symbol, period],
    queryFn: () => stockApi.getHistory(symbol, period),
    enabled: enabled && !!symbol,
    staleTime: 5 * 60 * 1000,
  });
}

export function useFilters() {
  return useQuery<{ filters: FilterDefinition[]; grouped: Record<string, FilterDefinition[]> }>({
    queryKey: ['filters'],
    queryFn: filterApi.getAll,
    staleTime: 30 * 60 * 1000,
  });
}

export function useFilterCategories() {
  return useQuery({
    queryKey: ['filter-categories'],
    queryFn: filterApi.getCategories,
    staleTime: 30 * 60 * 1000,
  });
}

export function useSectorPerformance() {
  return useQuery<SectorPerformance[]>({
    queryKey: ['sector-performance'],
    queryFn: sectorApi.getPerformance,
    staleTime: 2 * 60 * 1000,
  });
}

export function useSectors() {
  return useQuery<string[]>({
    queryKey: ['sectors'],
    queryFn: sectorApi.getList,
    staleTime: 30 * 60 * 1000,
  });
}
