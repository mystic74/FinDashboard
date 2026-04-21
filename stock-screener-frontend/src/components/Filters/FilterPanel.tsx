import { useState } from 'react';
import { ChevronDown, ChevronUp, X, Plus } from 'lucide-react';
import type { Filter, FilterDefinition, FilterOperator } from '../../types';
import { useFilters } from '../../hooks/useStocks';

interface FilterPanelProps {
  filters: Filter[];
  onChange: (filters: Filter[]) => void;
}

const operatorLabels: Record<FilterOperator, string> = {
  eq: 'equals',
  ne: 'not equals',
  gt: 'greater than',
  gte: 'at least',
  lt: 'less than',
  lte: 'at most',
  between: 'between',
  in: 'in',
  notIn: 'not in',
  contains: 'contains',
};

const categoryLabels: Record<string, string> = {
  price_volume: 'Price & Volume',
  valuation: 'Valuation',
  dividends: 'Dividends',
  financial_health: 'Financial Health',
  profitability: 'Profitability',
  growth: 'Growth',
  technical: 'Technical',
  profile: 'Profile',
};

export function FilterPanel({ filters, onChange }: FilterPanelProps) {
  const { data: filterData } = useFilters();
  const [expandedCategories, setExpandedCategories] = useState<string[]>(['price_volume', 'valuation']);
  const [showAddFilter, setShowAddFilter] = useState(false);

  const toggleCategory = (category: string) => {
    setExpandedCategories((prev) =>
      prev.includes(category) ? prev.filter((c) => c !== category) : [...prev, category]
    );
  };

  const addFilter = (definition: FilterDefinition) => {
    const newFilter: Filter = {
      field: definition.field,
      operator: definition.operators[0],
      value: definition.type === 'string' ? '' : 0,
    };
    onChange([...filters, newFilter]);
    setShowAddFilter(false);
  };

  const updateFilter = (index: number, updates: Partial<Filter>) => {
    const newFilters = [...filters];
    newFilters[index] = { ...newFilters[index], ...updates };
    onChange(newFilters);
  };

  const removeFilter = (index: number) => {
    onChange(filters.filter((_, i) => i !== index));
  };

  const getDefinition = (field: string): FilterDefinition | undefined => {
    return filterData?.filters.find((f) => f.field === field);
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
      {/* Active Filters */}
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between mb-3">
          <h3 className="font-semibold text-gray-900 dark:text-white">Active Filters</h3>
          <button
            onClick={() => setShowAddFilter(!showAddFilter)}
            className="flex items-center gap-1 px-3 py-1.5 text-sm font-medium text-primary-600 dark:text-primary-400 hover:bg-primary-50 dark:hover:bg-primary-900/20 rounded-lg transition-colors"
          >
            <Plus className="w-4 h-4" />
            Add Filter
          </button>
        </div>

        {filters.length === 0 ? (
          <p className="text-sm text-gray-500 dark:text-gray-400">
            No filters applied. Click "Add Filter" to start screening.
          </p>
        ) : (
          <div className="space-y-2">
            {filters.map((filter, index) => {
              const definition = getDefinition(filter.field);
              return (
                <div
                  key={index}
                  className="flex items-center gap-2 p-2 bg-gray-50 dark:bg-gray-700/50 rounded-lg"
                >
                  <div className="flex-1 grid grid-cols-3 gap-2 min-w-0">
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300 truncate">
                      {definition?.label || filter.field}
                    </span>
                    <select
                      value={filter.operator}
                      onChange={(e) => updateFilter(index, { operator: e.target.value as FilterOperator })}
                      className="text-sm bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded px-2 py-1"
                    >
                      {definition?.operators.map((op) => (
                        <option key={op} value={op}>
                          {operatorLabels[op]}
                        </option>
                      ))}
                    </select>
                    <div className="flex gap-1">
                      <input
                        type={definition?.type === 'string' ? 'text' : 'number'}
                        value={filter.value}
                        onChange={(e) =>
                          updateFilter(index, {
                            value: definition?.type === 'string' ? e.target.value : Number(e.target.value),
                          })
                        }
                        className="w-full text-sm bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded px-2 py-1"
                        placeholder={definition?.type === 'percent' ? '%' : ''}
                      />
                      {filter.operator === 'between' && (
                        <input
                          type="number"
                          value={filter.value2 || ''}
                          onChange={(e) => updateFilter(index, { value2: Number(e.target.value) })}
                          className="w-full text-sm bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded px-2 py-1"
                        />
                      )}
                    </div>
                  </div>
                  <button
                    onClick={() => removeFilter(index)}
                    className="p-1 text-gray-400 hover:text-red-500 transition-colors"
                  >
                    <X className="w-4 h-4" />
                  </button>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Add Filter Modal/Dropdown */}
      {showAddFilter && filterData && (
        <div className="p-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50 max-h-80 overflow-y-auto">
          {Object.entries(filterData.grouped).map(([category, categoryFilters]) => (
            <div key={category} className="mb-2">
              <button
                onClick={() => toggleCategory(category)}
                className="w-full flex items-center justify-between p-2 text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 rounded"
              >
                {categoryLabels[category] || category}
                {expandedCategories.includes(category) ? (
                  <ChevronUp className="w-4 h-4" />
                ) : (
                  <ChevronDown className="w-4 h-4" />
                )}
              </button>
              {expandedCategories.includes(category) && (
                <div className="ml-2 space-y-1">
                  {categoryFilters.map((def) => (
                    <button
                      key={def.field}
                      onClick={() => addFilter(def)}
                      disabled={filters.some((f) => f.field === def.field)}
                      className="w-full text-left p-2 text-sm text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <div className="font-medium">{def.label}</div>
                      <div className="text-xs text-gray-400">{def.description}</div>
                    </button>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Quick Filters */}
      <div className="p-4">
        <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Quick Filters</h4>
        <div className="flex flex-wrap gap-2">
          <button
            onClick={() => onChange([
              { field: 'peRatio', operator: 'lt', value: 15 },
              { field: 'peRatio', operator: 'gt', value: 0 },
            ])}
            className="px-3 py-1 text-xs font-medium text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
          >
            Low P/E (&lt;15)
          </button>
          <button
            onClick={() => onChange([
              { field: 'dividendYield', operator: 'gt', value: 3 },
            ])}
            className="px-3 py-1 text-xs font-medium text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
          >
            High Dividend (&gt;3%)
          </button>
          <button
            onClick={() => onChange([
              { field: 'marketCap', operator: 'gt', value: 10000000000 },
            ])}
            className="px-3 py-1 text-xs font-medium text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
          >
            Large Cap (&gt;$10B)
          </button>
          <button
            onClick={() => onChange([
              { field: 'beta', operator: 'lt', value: 1 },
            ])}
            className="px-3 py-1 text-xs font-medium text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 rounded-full hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
          >
            Low Beta (&lt;1)
          </button>
          <button
            onClick={() => onChange([])}
            className="px-3 py-1 text-xs font-medium text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 rounded-full hover:bg-red-100 dark:hover:bg-red-900/40 transition-colors"
          >
            Clear All
          </button>
        </div>
      </div>
    </div>
  );
}
