import { useState } from 'react';
import { Play, Loader2, Clock } from 'lucide-react';
import { useCustomScreener } from '../../hooks/useScreener';
import { FilterPanel } from '../Filters/FilterPanel';
import { StockTable } from '../StockTable/StockTable';
import type { Filter } from '../../types';

export function CustomScreener() {
  const [filters, setFilters] = useState<Filter[]>([]);
  const [sortBy, setSortBy] = useState('marketCap');
  const [sortOrder, setSortOrder] = useState('desc');

  const { mutate: runScreener, data: result, isPending, error } = useCustomScreener();

  const handleRunScreener = () => {
    if (filters.length === 0) return;
    runScreener({
      filters,
      sortBy,
      sortOrder,
      limit: 100,
    });
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Custom Screener</h1>
        <p className="text-gray-500 dark:text-gray-400">
          Build your own stock screener with custom filters
        </p>
      </div>

      {/* Filter Panel */}
      <FilterPanel filters={filters} onChange={setFilters} />

      {/* Sort Options */}
      <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h3 className="font-semibold text-gray-900 dark:text-white mb-3">Sort Results</h3>
        <div className="flex flex-wrap gap-4">
          <div>
            <label className="block text-sm text-gray-600 dark:text-gray-400 mb-1">Sort By</label>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm"
            >
              <option value="marketCap">Market Cap</option>
              <option value="price">Price</option>
              <option value="changePercent">Change %</option>
              <option value="volume">Volume</option>
              <option value="peRatio">P/E Ratio</option>
              <option value="dividendYield">Dividend Yield</option>
              <option value="roe">ROE</option>
              <option value="revenueGrowth">Revenue Growth</option>
              <option value="piotroskiFScore">F-Score</option>
            </select>
          </div>
          <div>
            <label className="block text-sm text-gray-600 dark:text-gray-400 mb-1">Order</label>
            <select
              value={sortOrder}
              onChange={(e) => setSortOrder(e.target.value)}
              className="bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm"
            >
              <option value="desc">Descending</option>
              <option value="asc">Ascending</option>
            </select>
          </div>
          <div className="flex items-end">
            <button
              onClick={handleRunScreener}
              disabled={filters.length === 0 || isPending}
              className="flex items-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isPending ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <Play className="w-4 h-4" />
              )}
              Run Screener
            </button>
          </div>
        </div>
      </div>

      {/* Error */}
      {error && (
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
          <p className="text-red-600 dark:text-red-400">{error.message}</p>
        </div>
      )}

      {/* Results */}
      {result && (
        <div className="space-y-4">
          <div className="flex items-center gap-4">
            <span className="text-lg font-semibold text-gray-900 dark:text-white">
              Results: {result.total} stocks
            </span>
            <div className="flex items-center gap-1 text-sm text-gray-500 dark:text-gray-400">
              <Clock className="w-4 h-4" />
              Page {result.page}
            </div>
          </div>
          <StockTable stocks={result.stocks} />
        </div>
      )}

      {/* Empty State */}
      {!result && !isPending && filters.length > 0 && (
        <div className="text-center py-12 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
          <p className="text-gray-500 dark:text-gray-400">
            Click "Run Screener" to see results
          </p>
        </div>
      )}

      {filters.length === 0 && (
        <div className="text-center py-12 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
          <p className="text-gray-500 dark:text-gray-400">
            Add filters above to start screening stocks
          </p>
        </div>
      )}
    </div>
  );
}
