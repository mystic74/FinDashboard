import { useState, useMemo } from 'react';
import {
  ArrowUpDown,
  ArrowUp,
  ArrowDown,
  Download,
  ChevronLeft,
  ChevronRight,
} from 'lucide-react';
import type { Stock, TableColumn, SortDirection } from '../../types';
import {
  formatCurrency,
  formatPercent,
  formatCompact,
  formatVolume,
  formatRatio,
  getChangeColor,
} from '../../utils/formatters';

interface StockTableProps {
  stocks: Stock[];
  isLoading?: boolean;
  onStockClick?: (symbol: string) => void;
}

const defaultColumns: TableColumn[] = [
  { key: 'symbol', label: 'Symbol', sortable: true, width: '100px' },
  { key: 'name', label: 'Name', sortable: true, width: '180px' },
  { key: 'price', label: 'Price', sortable: true, format: 'currency', align: 'right' },
  { key: 'changePercent', label: 'Change %', sortable: true, format: 'percent', align: 'right' },
  { key: 'volume', label: 'Volume', sortable: true, format: 'compact', align: 'right' },
  { key: 'marketCap', label: 'Market Cap', sortable: true, format: 'compact', align: 'right' },
  { key: 'peRatio', label: 'P/E', sortable: true, format: 'number', align: 'right' },
  { key: 'dividendYield', label: 'Div Yield', sortable: true, format: 'percent', align: 'right' },
  { key: 'roe', label: 'ROE', sortable: true, format: 'percent', align: 'right' },
  { key: 'sector', label: 'Sector', sortable: true, width: '140px' },
];

export function StockTable({ stocks, isLoading, onStockClick }: StockTableProps) {
  const [sortKey, setSortKey] = useState<keyof Stock>('marketCap');
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc');
  const [page, setPage] = useState(0);
  const [pageSize] = useState(20);

  const sortedStocks = useMemo(() => {
    return [...stocks].sort((a, b) => {
      const aVal = a[sortKey];
      const bVal = b[sortKey];

      if (aVal === null || aVal === undefined) return 1;
      if (bVal === null || bVal === undefined) return -1;

      if (typeof aVal === 'string' && typeof bVal === 'string') {
        return sortDirection === 'asc'
          ? aVal.localeCompare(bVal)
          : bVal.localeCompare(aVal);
      }

      if (typeof aVal === 'number' && typeof bVal === 'number') {
        return sortDirection === 'asc' ? aVal - bVal : bVal - aVal;
      }

      return 0;
    });
  }, [stocks, sortKey, sortDirection]);

  const paginatedStocks = useMemo(() => {
    const start = page * pageSize;
    return sortedStocks.slice(start, start + pageSize);
  }, [sortedStocks, page, pageSize]);

  const totalPages = Math.ceil(stocks.length / pageSize);

  const handleSort = (key: keyof Stock) => {
    if (sortKey === key) {
      setSortDirection((prev) => (prev === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortKey(key);
      setSortDirection('desc');
    }
    setPage(0);
  };

  const formatValue = (value: unknown, format?: string): string => {
    if (value === null || value === undefined) return '-';

    switch (format) {
      case 'currency':
        return formatCurrency(value as number);
      case 'percent':
        return formatPercent(value as number);
      case 'compact':
        return typeof value === 'number' && value > 1e6
          ? formatCompact(value)
          : formatVolume(value as number);
      case 'number':
        return formatRatio(value as number);
      default:
        return String(value);
    }
  };

  const exportToCSV = () => {
    const headers = defaultColumns.map((col) => col.label).join(',');
    const rows = sortedStocks.map((stock) =>
      defaultColumns.map((col) => {
        const val = stock[col.key];
        if (typeof val === 'string' && val.includes(',')) {
          return `"${val}"`;
        }
        return val ?? '';
      }).join(',')
    );
    const csv = [headers, ...rows].join('\n');
    const blob = new Blob([csv], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `stocks-${new Date().toISOString().split('T')[0]}.csv`;
    a.click();
    URL.revokeObjectURL(url);
  };

  if (isLoading) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div className="animate-pulse">
          <div className="h-12 bg-gray-100 dark:bg-gray-700" />
          {[...Array(10)].map((_, i) => (
            <div key={i} className="h-14 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800" />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
      {/* Table Header Actions */}
      <div className="flex items-center justify-between px-4 py-3 border-b border-gray-200 dark:border-gray-700">
        <div className="text-sm text-gray-600 dark:text-gray-400">
          Showing {paginatedStocks.length} of {stocks.length} stocks
        </div>
        <button
          onClick={exportToCSV}
          className="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
        >
          <Download className="w-4 h-4" />
          Export CSV
        </button>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="bg-gray-50 dark:bg-gray-900/50 border-b border-gray-200 dark:border-gray-700">
              {defaultColumns.map((column) => (
                <th
                  key={column.key}
                  className={`px-4 py-3 text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider ${
                    column.align === 'right' ? 'text-right' : 'text-left'
                  } ${column.sortable ? 'cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800' : ''}`}
                  style={{ width: column.width }}
                  onClick={() => column.sortable && handleSort(column.key)}
                >
                  <div className={`flex items-center gap-1 ${column.align === 'right' ? 'justify-end' : ''}`}>
                    {column.label}
                    {column.sortable && (
                      <span className="text-gray-400">
                        {sortKey === column.key ? (
                          sortDirection === 'asc' ? (
                            <ArrowUp className="w-3 h-3" />
                          ) : (
                            <ArrowDown className="w-3 h-3" />
                          )
                        ) : (
                          <ArrowUpDown className="w-3 h-3" />
                        )}
                      </span>
                    )}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {paginatedStocks.map((stock, index) => (
              <tr
                key={stock.symbol}
                onClick={() => onStockClick?.(stock.symbol)}
                className={`border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700/50 cursor-pointer transition-colors ${
                  index % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50/50 dark:bg-gray-800/50'
                }`}
              >
                {defaultColumns.map((column) => (
                  <td
                    key={column.key}
                    className={`px-4 py-3 text-sm ${
                      column.align === 'right' ? 'text-right' : 'text-left'
                    }`}
                  >
                    {column.key === 'symbol' ? (
                      <span className="font-semibold text-primary-600 dark:text-primary-400">
                        {stock.symbol}
                      </span>
                    ) : column.key === 'name' ? (
                      <span className="text-gray-900 dark:text-white truncate max-w-[180px] block">
                        {stock.name}
                      </span>
                    ) : column.key === 'changePercent' ? (
                      <span className={getChangeColor(stock.changePercent)}>
                        {formatValue(stock[column.key], column.format)}
                      </span>
                    ) : column.key === 'sector' ? (
                      <span className="text-gray-600 dark:text-gray-400 truncate max-w-[140px] block">
                        {stock.sector || '-'}
                      </span>
                    ) : (
                      <span className="text-gray-900 dark:text-white">
                        {formatValue(stock[column.key], column.format)}
                      </span>
                    )}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-700">
          <div className="text-sm text-gray-600 dark:text-gray-400">
            Page {page + 1} of {totalPages}
          </div>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage((p) => Math.max(0, p - 1))}
              disabled={page === 0}
              className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <ChevronLeft className="w-4 h-4" />
            </button>
            <button
              onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
              disabled={page === totalPages - 1}
              className="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
