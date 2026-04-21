import { useState } from 'react';
import { ArrowLeft, Clock, Filter, Info } from 'lucide-react';
import { useRunScreener } from '../../hooks/useScreener';
import { StockTable } from '../StockTable/StockTable';
import type { Filter as FilterType } from '../../types';

// Helper to format filter for display
function formatFilter(filter: FilterType): string {
  const operatorMap: Record<string, string> = {
    eq: '=',
    ne: '≠',
    gt: '>',
    gte: '≥',
    lt: '<',
    lte: '≤',
    between: 'between',
    in: 'in',
    notIn: 'not in',
    contains: 'contains',
  };
  const op = operatorMap[filter.operator] || filter.operator;
  if (filter.operator === 'between' && filter.value2 !== undefined) {
    return `${filter.field} ${op} ${filter.value} and ${filter.value2}`;
  }
  return `${filter.field} ${op} ${filter.value}`;
}

interface ScreenerResultsProps {
  screenerId: string;
  country?: string;
  onBack: () => void;
}

export function ScreenerResults({ screenerId, country, onBack }: ScreenerResultsProps) {
  const { data: result, isLoading, error } = useRunScreener(screenerId, country);
  const [showFilters, setShowFilters] = useState(false);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <button
          onClick={onBack}
          className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
        >
          <ArrowLeft className="w-5 h-5 text-gray-600 dark:text-gray-400" />
        </button>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {result?.screener?.name || 'Loading...'}
          </h1>
          <p className="text-gray-500 dark:text-gray-400">
            {result?.screener?.description}
          </p>
        </div>
      </div>

      {/* Stats */}
      {result && (
        <div className="flex flex-wrap gap-4">
          {country && (
            <div className="bg-primary-100 dark:bg-primary-900/30 rounded-lg px-4 py-2 border border-primary-200 dark:border-primary-700 flex items-center gap-2">
              <span className="text-sm font-medium text-primary-700 dark:text-primary-300">
                Filtered by: {country}
              </span>
            </div>
          )}
          <div className="relative">
            <button
              onMouseEnter={() => setShowFilters(true)}
              onMouseLeave={() => setShowFilters(false)}
              onClick={() => setShowFilters(!showFilters)}
              className="bg-white dark:bg-gray-800 rounded-lg px-4 py-2 border border-gray-200 dark:border-gray-700 flex items-center gap-2 hover:border-primary-300 dark:hover:border-primary-600 transition-colors"
            >
              <Filter className="w-4 h-4 text-gray-400" />
              <span className="text-sm text-gray-600 dark:text-gray-400">
                {result.screener.filters.length} filters applied
              </span>
              <Info className="w-3 h-3 text-gray-400" />
            </button>
            {/* Filter Tooltip */}
            {showFilters && result.screener.filters.length > 0 && (
              <div className="absolute top-full left-0 mt-2 z-50 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 shadow-lg p-3 min-w-[280px]">
                <div className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">
                  Applied Filters
                </div>
                <ul className="space-y-1">
                  {result.screener.filters.map((filter, idx) => (
                    <li key={idx} className="text-sm text-gray-700 dark:text-gray-300 font-mono bg-gray-50 dark:bg-gray-900 rounded px-2 py-1">
                      {formatFilter(filter)}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-lg px-4 py-2 border border-gray-200 dark:border-gray-700 flex items-center gap-2">
            <span className="text-sm font-medium text-primary-600 dark:text-primary-400">
              {result.total} stocks found
            </span>
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-lg px-4 py-2 border border-gray-200 dark:border-gray-700 flex items-center gap-2">
            <Clock className="w-4 h-4 text-gray-400" />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {result.executionMs}ms
            </span>
          </div>
        </div>
      )}

      {/* Error */}
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <p className="text-red-600 dark:text-red-400">{error.message}</p>
        </div>
      )}

      {/* Results Table */}
      <StockTable stocks={result?.stocks || []} isLoading={isLoading} />
    </div>
  );
}
