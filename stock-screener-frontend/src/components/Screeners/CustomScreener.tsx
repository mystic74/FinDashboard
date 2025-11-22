import { useState, useCallback } from 'react';
import { Play, Loader2, Clock, Globe, Building2 } from 'lucide-react';
import { useCustomScreener } from '../../hooks/useScreener';
import { FilterPanel } from '../Filters/FilterPanel';
import { StockTable } from '../StockTable/StockTable';
import type { Filter } from '../../types';

// Available countries and sectors for quick filtering
const COUNTRIES = [
  'All', 'USA', 'UK', 'Israel', 'Germany', 'Japan', 'China', 'India',
  'Brazil', 'Canada', 'France', 'Switzerland', 'Australia', 'Netherlands',
  'South Korea', 'Taiwan'
];

const SECTORS = [
  'All', 'Technology', 'Financial Services', 'Healthcare', 'Energy',
  'Consumer Cyclical', 'Consumer Defensive', 'Industrials', 'Basic Materials',
  'Communication Services', 'Real Estate'
];

export function CustomScreener() {
  const [filters, setFilters] = useState<Filter[]>([]);
  const [sortBy, setSortBy] = useState('marketCap');
  const [sortOrder, setSortOrder] = useState('desc');
  const [selectedCountry, setSelectedCountry] = useState('All');
  const [selectedSector, setSelectedSector] = useState('All');

  const { mutate: runScreener, data: result, isPending, error } = useCustomScreener();

  // Build filters including country/sector
  const buildFilters = useCallback(() => {
    const allFilters = [...filters];
    if (selectedCountry !== 'All') {
      allFilters.push({ field: 'country', operator: 'eq', value: selectedCountry });
    }
    if (selectedSector !== 'All') {
      allFilters.push({ field: 'sector', operator: 'eq', value: selectedSector });
    }
    return allFilters;
  }, [filters, selectedCountry, selectedSector]);

  const handleRunScreener = () => {
    const allFilters = buildFilters();
    runScreener({
      filters: allFilters,
      sortBy,
      sortOrder,
      limit: 100,
    });
  };

  // Quick filter handlers for StockTable clicks
  const handleSectorClick = (sector: string) => {
    setSelectedSector(sector);
  };

  const handleCountryClick = (country: string) => {
    setSelectedCountry(country);
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

      {/* Quick Filters - Country & Sector */}
      <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h3 className="font-semibold text-gray-900 dark:text-white mb-3">Quick Filters</h3>
        <div className="flex flex-wrap gap-4">
          <div>
            <label className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
              <Globe className="w-4 h-4" />
              Country
            </label>
            <select
              value={selectedCountry}
              onChange={(e) => setSelectedCountry(e.target.value)}
              className="bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm min-w-[150px]"
            >
              {COUNTRIES.map((country) => (
                <option key={country} value={country}>{country}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 mb-1">
              <Building2 className="w-4 h-4" />
              Sector
            </label>
            <select
              value={selectedSector}
              onChange={(e) => setSelectedSector(e.target.value)}
              className="bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm min-w-[180px]"
            >
              {SECTORS.map((sector) => (
                <option key={sector} value={sector}>{sector}</option>
              ))}
            </select>
          </div>
        </div>
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
              disabled={isPending}
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
          <StockTable
            stocks={result.stocks}
            onSectorClick={handleSectorClick}
            onCountryClick={handleCountryClick}
          />
        </div>
      )}

      {/* Empty State */}
      {!result && !isPending && (
        <div className="text-center py-12 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
          <p className="text-gray-500 dark:text-gray-400">
            {filters.length > 0 || selectedCountry !== 'All' || selectedSector !== 'All'
              ? 'Click "Run Screener" to see results'
              : 'Select country, sector, or add filters to start screening stocks'}
          </p>
        </div>
      )}
    </div>
  );
}
