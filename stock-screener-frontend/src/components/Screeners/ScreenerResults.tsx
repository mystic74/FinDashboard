import { ArrowLeft, Clock, Filter } from 'lucide-react';
import { useRunScreener } from '../../hooks/useScreener';
import { StockTable } from '../StockTable/StockTable';

interface ScreenerResultsProps {
  screenerId: string;
  onBack: () => void;
}

export function ScreenerResults({ screenerId, onBack }: ScreenerResultsProps) {
  const { data: result, isLoading, error } = useRunScreener(screenerId);

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
          <div className="bg-white dark:bg-gray-800 rounded-lg px-4 py-2 border border-gray-200 dark:border-gray-700 flex items-center gap-2">
            <Filter className="w-4 h-4 text-gray-400" />
            <span className="text-sm text-gray-600 dark:text-gray-400">
              {result.screener.filters.length} filters applied
            </span>
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
