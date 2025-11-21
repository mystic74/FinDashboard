export interface Stock {
  symbol: string;
  name: string;
  exchange: string;
  currency: string;
  country: string;
  sector: string;
  industry: string;
  price: number;
  open: number;
  high: number;
  low: number;
  previousClose: number;
  change: number;
  changePercent: number;
  volume: number;
  avgVolume: number;
  marketCap: number;
  sharesOutstanding: number;
  peRatio: number;
  forwardPE: number;
  pegRatio: number;
  pbRatio: number;
  psRatio: number;
  evToEbitda: number;
  dividendYield: number;
  dividendPerShare: number;
  payoutRatio: number;
  consecutiveDivYears: number;
  dividendGrowthYears: number;
  beta: number;
  roe: number;
  roa: number;
  roic: number;
  grossMargin: number;
  operatingMargin: number;
  netMargin: number;
  currentRatio: number;
  quickRatio: number;
  debtToEquity: number;
  interestCoverage: number;
  altmanZScore: number;
  cashToDebt: number;
  revenueGrowth: number;
  epsGrowth: number;
  bookValueGrowth: number;
  fcfGrowth: number;
  eps: number;
  bookValuePerShare: number;
  week52High: number;
  week52Low: number;
  ma50: number;
  ma200: number;
  rsi14: number;
  macd: number;
  return1W: number;
  return1M: number;
  return3M: number;
  return6M: number;
  return1Y: number;
  piotroskiFScore: number;
  lastUpdated: string;
}

export interface Screener {
  id: string;
  name: string;
  description: string;
  category: string;
  filters: Filter[];
  sortBy?: string;
  sortOrder?: string;
  icon?: string;
  isCustom: boolean;
}

export interface Filter {
  field: string;
  operator: FilterOperator;
  value: number | string;
  value2?: number | string;
}

export type FilterOperator =
  | 'eq'
  | 'ne'
  | 'gt'
  | 'gte'
  | 'lt'
  | 'lte'
  | 'between'
  | 'in'
  | 'notIn'
  | 'contains';

export interface FilterDefinition {
  field: string;
  label: string;
  description: string;
  category: string;
  type: 'number' | 'percent' | 'currency' | 'string' | 'boolean';
  unit?: string;
  min?: number;
  max?: number;
  options?: string[];
  operators: FilterOperator[];
}

export interface ScreenerResult {
  screener: Screener;
  stocks: Stock[];
  total: number;
  executionMs: number;
  lastUpdated: string;
}

export interface FilterRequest {
  filters: Filter[];
  sortBy?: string;
  sortOrder?: string;
  limit?: number;
  offset?: number;
}

export interface FilterResponse {
  stocks: Stock[];
  total: number;
  page: number;
  pageSize: number;
  appliedFilters: Filter[];
}

export interface SectorPerformance {
  sector: string;
  change1D: number;
  change1W: number;
  change1M: number;
  change3M: number;
  changeYtd: number;
  change1Y: number;
  stockCount: number;
  marketCap: number;
  topPerformer: string;
  worstPerformer: string;
}

export interface ApiResponse<T> {
  success: boolean;
  error?: string;
  [key: string]: T | boolean | string | undefined;
}

export interface ScreenerSummary {
  id: string;
  name: string;
  description: string;
  category: string;
  icon: string;
  matchCount?: number;
  topStock?: string;
}

export type SortDirection = 'asc' | 'desc';

export interface TableColumn {
  key: keyof Stock;
  label: string;
  sortable?: boolean;
  format?: 'number' | 'percent' | 'currency' | 'compact' | 'date';
  width?: string;
  align?: 'left' | 'center' | 'right';
}
