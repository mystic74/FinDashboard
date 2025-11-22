import { useState, useEffect, useCallback, useMemo } from 'react';
import { ArrowLeft, Sliders, Eye, RotateCcw, Sparkles, TrendingUp, DollarSign, Percent, Clock } from 'lucide-react';
import { useScreeners, useCustomScreener } from '../../hooks/useScreener';
import { StockTable } from '../StockTable/StockTable';
import type { Filter } from '../../types';

// Filter metadata for display and slider configuration
const FILTER_CONFIG: Record<string, {
  label: string;
  unit: string;
  min: number;
  max: number;
  step: number;
  format: (v: number) => string;
  icon: React.ReactNode;
}> = {
  marketCap: {
    label: 'Market Cap',
    unit: 'B',
    min: 0,
    max: 500,
    step: 0.5,
    format: (v) => `$${v >= 1000 ? (v/1000).toFixed(1) + 'T' : v.toFixed(1) + 'B'}`,
    icon: <DollarSign className="w-4 h-4" />,
  },
  peRatio: {
    label: 'P/E Ratio',
    unit: 'x',
    min: 0,
    max: 100,
    step: 0.5,
    format: (v) => v.toFixed(1) + 'x',
    icon: <TrendingUp className="w-4 h-4" />,
  },
  forwardPE: {
    label: 'Forward P/E',
    unit: 'x',
    min: 0,
    max: 100,
    step: 0.5,
    format: (v) => v.toFixed(1) + 'x',
    icon: <TrendingUp className="w-4 h-4" />,
  },
  pegRatio: {
    label: 'PEG Ratio',
    unit: '',
    min: 0,
    max: 5,
    step: 0.1,
    format: (v) => v.toFixed(2),
    icon: <TrendingUp className="w-4 h-4" />,
  },
  pbRatio: {
    label: 'P/B Ratio',
    unit: 'x',
    min: 0,
    max: 20,
    step: 0.1,
    format: (v) => v.toFixed(1) + 'x',
    icon: <TrendingUp className="w-4 h-4" />,
  },
  psRatio: {
    label: 'P/S Ratio',
    unit: 'x',
    min: 0,
    max: 50,
    step: 0.5,
    format: (v) => v.toFixed(1) + 'x',
    icon: <TrendingUp className="w-4 h-4" />,
  },
  evToEbitda: {
    label: 'EV/EBITDA',
    unit: 'x',
    min: 0,
    max: 50,
    step: 0.5,
    format: (v) => v.toFixed(1) + 'x',
    icon: <TrendingUp className="w-4 h-4" />,
  },
  dividendYield: {
    label: 'Dividend Yield',
    unit: '%',
    min: 0,
    max: 15,
    step: 0.1,
    format: (v) => v.toFixed(1) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  beta: {
    label: 'Beta',
    unit: '',
    min: 0,
    max: 3,
    step: 0.05,
    format: (v) => v.toFixed(2),
    icon: <TrendingUp className="w-4 h-4" />,
  },
  roe: {
    label: 'ROE',
    unit: '%',
    min: -50,
    max: 100,
    step: 1,
    format: (v) => v.toFixed(0) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  roa: {
    label: 'ROA',
    unit: '%',
    min: -20,
    max: 50,
    step: 0.5,
    format: (v) => v.toFixed(1) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  roic: {
    label: 'ROIC',
    unit: '%',
    min: -20,
    max: 100,
    step: 1,
    format: (v) => v.toFixed(0) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  grossMargin: {
    label: 'Gross Margin',
    unit: '%',
    min: 0,
    max: 100,
    step: 1,
    format: (v) => v.toFixed(0) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  operatingMargin: {
    label: 'Operating Margin',
    unit: '%',
    min: -50,
    max: 100,
    step: 1,
    format: (v) => v.toFixed(0) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  netMargin: {
    label: 'Net Margin',
    unit: '%',
    min: -50,
    max: 100,
    step: 1,
    format: (v) => v.toFixed(0) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  currentRatio: {
    label: 'Current Ratio',
    unit: 'x',
    min: 0,
    max: 10,
    step: 0.1,
    format: (v) => v.toFixed(1) + 'x',
    icon: <TrendingUp className="w-4 h-4" />,
  },
  quickRatio: {
    label: 'Quick Ratio',
    unit: 'x',
    min: 0,
    max: 10,
    step: 0.1,
    format: (v) => v.toFixed(1) + 'x',
    icon: <TrendingUp className="w-4 h-4" />,
  },
  debtToEquity: {
    label: 'Debt/Equity',
    unit: 'x',
    min: 0,
    max: 5,
    step: 0.1,
    format: (v) => v.toFixed(1) + 'x',
    icon: <TrendingUp className="w-4 h-4" />,
  },
  revenueGrowth: {
    label: 'Revenue Growth',
    unit: '%',
    min: -50,
    max: 200,
    step: 1,
    format: (v) => v.toFixed(0) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  epsGrowth: {
    label: 'EPS Growth',
    unit: '%',
    min: -100,
    max: 300,
    step: 1,
    format: (v) => v.toFixed(0) + '%',
    icon: <Percent className="w-4 h-4" />,
  },
  volume: {
    label: 'Volume',
    unit: 'M',
    min: 0,
    max: 100,
    step: 0.5,
    format: (v) => (v >= 1 ? v.toFixed(1) + 'M' : (v * 1000).toFixed(0) + 'K'),
    icon: <TrendingUp className="w-4 h-4" />,
  },
  avgVolume: {
    label: 'Avg Volume',
    unit: 'M',
    min: 0,
    max: 100,
    step: 0.5,
    format: (v) => (v >= 1 ? v.toFixed(1) + 'M' : (v * 1000).toFixed(0) + 'K'),
    icon: <TrendingUp className="w-4 h-4" />,
  },
  piotroskiFScore: {
    label: 'F-Score',
    unit: '',
    min: 0,
    max: 9,
    step: 1,
    format: (v) => v.toFixed(0),
    icon: <Sparkles className="w-4 h-4" />,
  },
  altmanZScore: {
    label: 'Z-Score',
    unit: '',
    min: 0,
    max: 10,
    step: 0.1,
    format: (v) => v.toFixed(1),
    icon: <TrendingUp className="w-4 h-4" />,
  },
  rsi14: {
    label: 'RSI (14)',
    unit: '',
    min: 0,
    max: 100,
    step: 1,
    format: (v) => v.toFixed(0),
    icon: <TrendingUp className="w-4 h-4" />,
  },
  consecutiveDivYears: {
    label: 'Consecutive Div Years',
    unit: 'yrs',
    min: 0,
    max: 50,
    step: 1,
    format: (v) => v.toFixed(0) + ' yrs',
    icon: <Clock className="w-4 h-4" />,
  },
  dividendGrowthYears: {
    label: 'Div Growth Years',
    unit: 'yrs',
    min: 0,
    max: 50,
    step: 1,
    format: (v) => v.toFixed(0) + ' yrs',
    icon: <Clock className="w-4 h-4" />,
  },
  cashToDebt: {
    label: 'Cash/Debt',
    unit: 'x',
    min: 0,
    max: 10,
    step: 0.1,
    format: (v) => v.toFixed(1) + 'x',
    icon: <DollarSign className="w-4 h-4" />,
  },
};

// Convert marketCap from raw to billions for slider
const toSliderValue = (field: string, value: number): number => {
  if (field === 'marketCap') return value / 1e9;
  if (field === 'volume' || field === 'avgVolume') return value / 1e6;
  return value;
};

// Convert slider value back to raw value
const fromSliderValue = (field: string, value: number): number => {
  if (field === 'marketCap') return value * 1e9;
  if (field === 'volume' || field === 'avgVolume') return value * 1e6;
  return value;
};

interface ScreenerTweakerProps {
  screenerId: string;
  country?: string;
  onBack: () => void;
}

export function ScreenerTweaker({ screenerId, country, onBack }: ScreenerTweakerProps) {
  const { data: screeners } = useScreeners();
  const { mutate: runCustom, data: previewResult, isPending } = useCustomScreener();

  const [filters, setFilters] = useState<Filter[]>([]);
  const [originalFilters, setOriginalFilters] = useState<Filter[]>([]);
  const [autoPreview, setAutoPreview] = useState(true);
  const [previewDebounce, setPreviewDebounce] = useState<ReturnType<typeof setTimeout> | null>(null);

  // Find the base screener
  const baseScreener = useMemo(() =>
    screeners?.find(s => s.id === screenerId),
    [screeners, screenerId]
  );

  // Initialize filters from screener
  useEffect(() => {
    if (baseScreener) {
      const clonedFilters = JSON.parse(JSON.stringify(baseScreener.filters));
      setFilters(clonedFilters);
      setOriginalFilters(JSON.parse(JSON.stringify(baseScreener.filters)));
    }
  }, [baseScreener]);

  // Build filters for API including country
  const buildApiFilters = useCallback(() => {
    const apiFilters = [...filters];
    if (country) {
      apiFilters.push({ field: 'country', operator: 'eq', value: country });
    }
    return apiFilters;
  }, [filters, country]);

  // Run preview
  const runPreview = useCallback(() => {
    const apiFilters = buildApiFilters();
    runCustom({
      filters: apiFilters,
      sortBy: baseScreener?.sortBy || 'marketCap',
      sortOrder: baseScreener?.sortOrder || 'desc',
      limit: 50,
    });
  }, [buildApiFilters, runCustom, baseScreener]);

  // Auto preview with debounce
  useEffect(() => {
    if (autoPreview && filters.length > 0) {
      if (previewDebounce) clearTimeout(previewDebounce);
      const timeout = setTimeout(runPreview, 500);
      setPreviewDebounce(timeout);
      return () => clearTimeout(timeout);
    }
  }, [filters, autoPreview]); // eslint-disable-line react-hooks/exhaustive-deps

  // Update a filter value
  const updateFilter = (index: number, field: 'value' | 'value2', newValue: number) => {
    setFilters(prev => {
      const updated = [...prev];
      const filterField = updated[index].field;
      updated[index] = {
        ...updated[index],
        [field]: fromSliderValue(filterField, newValue),
      };
      return updated;
    });
  };

  // Reset to original
  const resetFilters = () => {
    setFilters(JSON.parse(JSON.stringify(originalFilters)));
  };

  // Check if filter has changed
  const hasChanges = useMemo(() => {
    return JSON.stringify(filters) !== JSON.stringify(originalFilters);
  }, [filters, originalFilters]);

  if (!baseScreener) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-pulse text-gray-500">Loading screener...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button
            onClick={onBack}
            className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          >
            <ArrowLeft className="w-5 h-5 text-gray-600 dark:text-gray-400" />
          </button>
          <div>
            <div className="flex items-center gap-2">
              <Sliders className="w-5 h-5 text-primary-600" />
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                Customize Screener
              </h1>
            </div>
            <p className="text-gray-500 dark:text-gray-400">
              Based on: {baseScreener.name}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          {country && (
            <span className="px-3 py-1.5 bg-primary-100 dark:bg-primary-900/30 text-primary-700 dark:text-primary-300 rounded-full text-sm font-medium">
              {country}
            </span>
          )}
          <button
            onClick={resetFilters}
            disabled={!hasChanges}
            className="flex items-center gap-2 px-3 py-2 text-sm rounded-lg border border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            <RotateCcw className="w-4 h-4" />
            Reset
          </button>
        </div>
      </div>

      {/* Main Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Filters Panel */}
        <div className="lg:col-span-1 space-y-4">
          <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
            <div className="flex items-center justify-between mb-4">
              <h2 className="font-semibold text-gray-900 dark:text-white flex items-center gap-2">
                <Sliders className="w-4 h-4" />
                Filter Adjustments
              </h2>
              <label className="flex items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={autoPreview}
                  onChange={(e) => setAutoPreview(e.target.checked)}
                  className="rounded border-gray-300 dark:border-gray-600 text-primary-600 focus:ring-primary-500"
                />
                <span className="text-gray-600 dark:text-gray-400">Live preview</span>
              </label>
            </div>

            <div className="space-y-6">
              {filters.map((filter, index) => {
                const config = FILTER_CONFIG[filter.field];
                if (!config) return null; // Skip unknown filters

                const isNumeric = typeof filter.value === 'number';
                if (!isNumeric) return null; // Skip non-numeric filters

                const isBetween = filter.operator === 'between';
                const value1 = toSliderValue(filter.field, filter.value as number);
                const value2 = filter.value2 !== undefined ? toSliderValue(filter.field, filter.value2 as number) : undefined;

                return (
                  <div key={index} className="space-y-2">
                    <div className="flex items-center justify-between">
                      <label className="flex items-center gap-2 text-sm font-medium text-gray-700 dark:text-gray-300">
                        {config.icon}
                        {config.label}
                      </label>
                      <span className="text-sm font-mono text-primary-600 dark:text-primary-400">
                        {isBetween && value2 !== undefined
                          ? `${config.format(value1)} - ${config.format(value2)}`
                          : config.format(value1)}
                      </span>
                    </div>

                    {isBetween ? (
                      <div className="space-y-3">
                        <div>
                          <div className="text-xs text-gray-500 mb-1">Min: {config.format(value1)}</div>
                          <input
                            type="range"
                            min={config.min}
                            max={config.max}
                            step={config.step}
                            value={value1}
                            onChange={(e) => updateFilter(index, 'value', parseFloat(e.target.value))}
                            className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer slider-thumb"
                          />
                        </div>
                        <div>
                          <div className="text-xs text-gray-500 mb-1">Max: {config.format(value2 || config.max)}</div>
                          <input
                            type="range"
                            min={config.min}
                            max={config.max}
                            step={config.step}
                            value={value2 || config.max}
                            onChange={(e) => updateFilter(index, 'value2', parseFloat(e.target.value))}
                            className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer slider-thumb"
                          />
                        </div>
                      </div>
                    ) : (
                      <input
                        type="range"
                        min={config.min}
                        max={config.max}
                        step={config.step}
                        value={value1}
                        onChange={(e) => updateFilter(index, 'value', parseFloat(e.target.value))}
                        className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer slider-thumb"
                      />
                    )}

                    {/* Show operator hint */}
                    <div className="text-xs text-gray-400">
                      {filter.operator === 'gt' && 'Greater than'}
                      {filter.operator === 'gte' && 'Greater or equal'}
                      {filter.operator === 'lt' && 'Less than'}
                      {filter.operator === 'lte' && 'Less or equal'}
                      {filter.operator === 'between' && 'Between range'}
                      {filter.operator === 'eq' && 'Equals'}
                    </div>
                  </div>
                );
              })}
            </div>

            {!autoPreview && (
              <button
                onClick={runPreview}
                disabled={isPending}
                className="w-full mt-4 flex items-center justify-center gap-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 transition-colors"
              >
                <Eye className="w-4 h-4" />
                {isPending ? 'Loading...' : 'Preview Results'}
              </button>
            )}
          </div>

          {/* Quick Actions */}
          <div className="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 space-y-2">
            <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Quick Adjustments</h3>
            <button
              onClick={() => {
                setFilters(prev => prev.map(f => {
                  if (f.field === 'marketCap' && typeof f.value === 'number') {
                    return { ...f, value: f.value * 0.5, value2: f.value2 ? (f.value2 as number) * 0.5 : undefined };
                  }
                  return f;
                }));
              }}
              className="w-full text-left px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              Scale down market cap (smaller companies)
            </button>
            <button
              onClick={() => {
                setFilters(prev => prev.map(f => {
                  if (f.field === 'marketCap' && typeof f.value === 'number') {
                    return { ...f, value: f.value * 2, value2: f.value2 ? (f.value2 as number) * 2 : undefined };
                  }
                  return f;
                }));
              }}
              className="w-full text-left px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              Scale up market cap (larger companies)
            </button>
            <button
              onClick={() => {
                setFilters(prev => prev.map(f => {
                  if ((f.field === 'peRatio' || f.field === 'pbRatio' || f.field === 'psRatio') && typeof f.value === 'number') {
                    return { ...f, value: f.value * 0.8 };
                  }
                  return f;
                }));
              }}
              className="w-full text-left px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              Tighter value criteria (-20%)
            </button>
            <button
              onClick={() => {
                setFilters(prev => prev.map(f => {
                  if ((f.field === 'peRatio' || f.field === 'pbRatio' || f.field === 'psRatio') && typeof f.value === 'number') {
                    return { ...f, value: f.value * 1.2 };
                  }
                  return f;
                }));
              }}
              className="w-full text-left px-3 py-2 text-sm rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
            >
              Looser value criteria (+20%)
            </button>
          </div>
        </div>

        {/* Results Panel */}
        <div className="lg:col-span-2 space-y-4">
          {/* Stats Bar */}
          <div className="flex items-center gap-4 p-4 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
            <div className="flex items-center gap-2">
              <Sparkles className="w-5 h-5 text-primary-600" />
              <span className="font-semibold text-gray-900 dark:text-white">
                {previewResult?.total ?? '-'} stocks match
              </span>
            </div>
            {isPending && (
              <div className="flex items-center gap-2 text-gray-500">
                <div className="w-4 h-4 border-2 border-primary-600 border-t-transparent rounded-full animate-spin" />
                Updating...
              </div>
            )}
            {hasChanges && (
              <span className="px-2 py-0.5 text-xs font-medium bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-300 rounded">
                Modified
              </span>
            )}
          </div>

          {/* Results Table */}
          <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
            <StockTable
              stocks={previewResult?.stocks || []}
              isLoading={isPending}
            />
          </div>

          {/* Empty State */}
          {!isPending && previewResult?.total === 0 && (
            <div className="text-center py-8 bg-gray-50 dark:bg-gray-800/50 rounded-lg">
              <p className="text-gray-500 dark:text-gray-400">
                No stocks match your current criteria. Try loosening some filters.
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Custom CSS for slider */}
      <style>{`
        .slider-thumb::-webkit-slider-thumb {
          -webkit-appearance: none;
          appearance: none;
          width: 18px;
          height: 18px;
          border-radius: 50%;
          background: #2563eb;
          cursor: pointer;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
          transition: transform 0.1s, box-shadow 0.1s;
        }
        .slider-thumb::-webkit-slider-thumb:hover {
          transform: scale(1.1);
          box-shadow: 0 3px 6px rgba(0, 0, 0, 0.3);
        }
        .slider-thumb::-moz-range-thumb {
          width: 18px;
          height: 18px;
          border-radius: 50%;
          background: #2563eb;
          cursor: pointer;
          border: none;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }
      `}</style>
    </div>
  );
}
