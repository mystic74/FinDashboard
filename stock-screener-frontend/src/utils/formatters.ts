export function formatNumber(value: number, decimals: number = 2): string {
  if (value === null || value === undefined || isNaN(value)) return '-';
  return value.toLocaleString('en-US', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  });
}

export function formatCurrency(value: number, currency: string = 'USD'): string {
  if (value === null || value === undefined || isNaN(value)) return '-';
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value);
}

export function formatPercent(value: number, decimals: number = 2): string {
  if (value === null || value === undefined || isNaN(value)) return '-';
  return `${value >= 0 ? '+' : ''}${value.toFixed(decimals)}%`;
}

export function formatCompact(value: number): string {
  if (value === null || value === undefined || isNaN(value)) return '-';

  const absValue = Math.abs(value);
  const sign = value < 0 ? '-' : '';

  if (absValue >= 1e12) {
    return `${sign}$${(absValue / 1e12).toFixed(2)}T`;
  }
  if (absValue >= 1e9) {
    return `${sign}$${(absValue / 1e9).toFixed(2)}B`;
  }
  if (absValue >= 1e6) {
    return `${sign}$${(absValue / 1e6).toFixed(2)}M`;
  }
  if (absValue >= 1e3) {
    return `${sign}$${(absValue / 1e3).toFixed(2)}K`;
  }
  return `${sign}$${absValue.toFixed(2)}`;
}

export function formatVolume(value: number): string {
  if (value === null || value === undefined || isNaN(value)) return '-';

  if (value >= 1e9) {
    return `${(value / 1e9).toFixed(2)}B`;
  }
  if (value >= 1e6) {
    return `${(value / 1e6).toFixed(2)}M`;
  }
  if (value >= 1e3) {
    return `${(value / 1e3).toFixed(2)}K`;
  }
  return value.toLocaleString();
}

export function formatMarketCap(value: number): string {
  return formatCompact(value);
}

export function formatRatio(value: number, decimals: number = 2): string {
  if (value === null || value === undefined || isNaN(value)) return '-';
  if (value === 0) return '0';
  return value.toFixed(decimals);
}

export function formatDate(date: string | Date): string {
  if (!date) return '-';
  const d = new Date(date);
  return d.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

export function formatDateTime(date: string | Date): string {
  if (!date) return '-';
  const d = new Date(date);
  return d.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function getChangeColor(value: number): string {
  if (value > 0) return 'text-green-600 dark:text-green-400';
  if (value < 0) return 'text-red-600 dark:text-red-400';
  return 'text-gray-600 dark:text-gray-400';
}

export function getChangeBgColor(value: number): string {
  if (value > 0) return 'bg-green-100 dark:bg-green-900/30';
  if (value < 0) return 'bg-red-100 dark:bg-red-900/30';
  return 'bg-gray-100 dark:bg-gray-800';
}

export function getRSIColor(value: number): string {
  if (value >= 70) return 'text-red-600 dark:text-red-400'; // Overbought
  if (value <= 30) return 'text-green-600 dark:text-green-400'; // Oversold
  return 'text-gray-600 dark:text-gray-400';
}

export function getPERatioColor(value: number): string {
  if (value <= 0) return 'text-gray-400';
  if (value < 15) return 'text-green-600 dark:text-green-400';
  if (value < 25) return 'text-yellow-600 dark:text-yellow-400';
  return 'text-red-600 dark:text-red-400';
}

export function getScoreColor(score: number, max: number = 9): string {
  const ratio = score / max;
  if (ratio >= 0.8) return 'text-green-600 dark:text-green-400';
  if (ratio >= 0.5) return 'text-yellow-600 dark:text-yellow-400';
  return 'text-red-600 dark:text-red-400';
}

export function truncateText(text: string, maxLength: number): string {
  if (!text) return '';
  if (text.length <= maxLength) return text;
  return `${text.slice(0, maxLength)}...`;
}
