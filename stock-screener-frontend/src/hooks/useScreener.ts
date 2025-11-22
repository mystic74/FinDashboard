import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { screenerApi } from '../services/api';
import type { FilterRequest, Screener, ScreenerResult, FilterResponse, ScreenerSummary } from '../types';

export function useScreeners() {
  return useQuery<Screener[]>({
    queryKey: ['screeners'],
    queryFn: screenerApi.getAll,
    staleTime: 5 * 60 * 1000,
  });
}

export function useScreenerSummaries() {
  return useQuery<ScreenerSummary[]>({
    queryKey: ['screener-summaries'],
    queryFn: screenerApi.getSummaries,
    staleTime: 5 * 60 * 1000,
  });
}

export function useRunScreener(
  screenerId: string,
  country?: string,
  sector?: string,
  enabled: boolean = true
) {
  return useQuery<ScreenerResult>({
    queryKey: ['screener-result', screenerId, country, sector],
    queryFn: () => screenerApi.run(screenerId, country, sector),
    enabled: enabled && !!screenerId,
    staleTime: 2 * 60 * 1000,
  });
}

export function useCustomScreener() {
  const queryClient = useQueryClient();

  return useMutation<FilterResponse, Error, FilterRequest>({
    mutationFn: screenerApi.runCustom,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['custom-screener'] });
    },
  });
}
