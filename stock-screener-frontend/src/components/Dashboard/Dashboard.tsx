import { useState } from 'react';
import {
  TrendingUp,
  TrendingDown,
  DollarSign,
  BarChart3,
  Zap,
  Shield,
  Award,
  Search,
  Cpu,
  RefreshCw,
  Star,
  Sprout,
  Globe,
  Sliders,
} from 'lucide-react';
import { useScreeners } from '../../hooks/useScreener';
import { useSectorPerformance } from '../../hooks/useStocks';
import { formatPercent, getChangeColor } from '../../utils/formatters';

// Available countries with their flags and market context
const COUNTRIES = [
  { code: 'All', name: 'All Markets', flag: '🌍' },
  { code: 'USA', name: 'United States', flag: '🇺🇸' },
  { code: 'UK', name: 'United Kingdom', flag: '🇬🇧' },
  { code: 'Israel', name: 'Israel', flag: '🇮🇱' },
  { code: 'Germany', name: 'Germany', flag: '🇩🇪' },
  { code: 'Japan', name: 'Japan', flag: '🇯🇵' },
  { code: 'China', name: 'China', flag: '🇨🇳' },
  { code: 'India', name: 'India', flag: '🇮🇳' },
  { code: 'Brazil', name: 'Brazil', flag: '🇧🇷' },
  { code: 'Canada', name: 'Canada', flag: '🇨🇦' },
  { code: 'France', name: 'France', flag: '🇫🇷' },
  { code: 'Switzerland', name: 'Switzerland', flag: '🇨🇭' },
  { code: 'Australia', name: 'Australia', flag: '🇦🇺' },
  { code: 'Netherlands', name: 'Netherlands', flag: '🇳🇱' },
  { code: 'South Korea', name: 'South Korea', flag: '🇰🇷' },
  { code: 'Taiwan', name: 'Taiwan', flag: '🇹🇼' },
];

const iconMap: Record<string, React.ReactNode> = {
  'trending-up': <TrendingUp className="w-5 h-5" />,
  'dollar-sign': <DollarSign className="w-5 h-5" />,
  'search': <Search className="w-5 h-5" />,
  'zap': <Zap className="w-5 h-5" />,
  'shield': <Shield className="w-5 h-5" />,
  'award': <Award className="w-5 h-5" />,
  'sprout': <Sprout className="w-5 h-5" />,
  'cpu': <Cpu className="w-5 h-5" />,
  'star': <Star className="w-5 h-5" />,
  'refresh-cw': <RefreshCw className="w-5 h-5" />,
};

const categoryColors: Record<string, string> = {
  'Momentum': 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300',
  'Value': 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300',
  'Income': 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300',
  'Growth': 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-300',
  'Quality': 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-300',
  'Financial Health': 'bg-teal-100 text-teal-700 dark:bg-teal-900/30 dark:text-teal-300',
  'Defensive': 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
};

interface DashboardProps {
  onSelectScreener: (screenerId: string, country?: string) => void;
  onCustomizeScreener?: (screenerId: string, country?: string) => void;
}

export function Dashboard({ onSelectScreener, onCustomizeScreener }: DashboardProps) {
  const [selectedCountry, setSelectedCountry] = useState('All');
  const { data: screeners, isLoading: screenersLoading } = useScreeners();
  const { data: sectors, isLoading: sectorsLoading } = useSectorPerformance();

  const handleScreenerClick = (screenerId: string) => {
    // Pass the selected country along with the screener ID
    onSelectScreener(screenerId, selectedCountry !== 'All' ? selectedCountry : undefined);
  };

  const handleCustomizeClick = (e: React.MouseEvent, screenerId: string) => {
    e.stopPropagation(); // Prevent triggering the card click
    onCustomizeScreener?.(screenerId, selectedCountry !== 'All' ? selectedCountry : undefined);
  };

  return (
    <div className="space-y-8">
      {/* Country Selector */}
      <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div className="flex items-center gap-3 mb-3">
          <Globe className="w-5 h-5 text-primary-600 dark:text-primary-400" />
          <h2 className="font-semibold text-gray-900 dark:text-white">Market Region</h2>
          {selectedCountry !== 'All' && (
            <span className="text-sm text-primary-600 dark:text-primary-400">
              Filtering by {COUNTRIES.find(c => c.code === selectedCountry)?.name}
            </span>
          )}
        </div>
        <div className="flex flex-wrap gap-2">
          {COUNTRIES.map((country) => (
            <button
              key={country.code}
              onClick={() => setSelectedCountry(country.code)}
              className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
                selectedCountry === country.code
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              }`}
            >
              <span className="mr-1">{country.flag}</span>
              {country.code === 'All' ? 'All' : country.code}
            </button>
          ))}
        </div>
        {selectedCountry !== 'All' && (
          <p className="mt-3 text-sm text-gray-500 dark:text-gray-400">
            Screener results will be filtered to show only {COUNTRIES.find(c => c.code === selectedCountry)?.name} stocks.
            Market cap and other metrics are relative to the local market context.
          </p>
        )}
      </div>

      {/* Hero Section */}
      <div className="bg-gradient-to-r from-primary-600 to-primary-800 rounded-2xl p-8 text-white">
        <h1 className="text-3xl font-bold mb-2">Stock Screener Dashboard</h1>
        <p className="text-primary-100 text-lg">
          Find your next investment with powerful screening tools
        </p>
        <div className="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-white/10 rounded-lg p-4">
            <div className="text-2xl font-bold">{screeners?.length || 0}</div>
            <div className="text-sm text-primary-200">Predefined Screeners</div>
          </div>
          <div className="bg-white/10 rounded-lg p-4">
            <div className="text-2xl font-bold">50+</div>
            <div className="text-sm text-primary-200">Filter Options</div>
          </div>
          <div className="bg-white/10 rounded-lg p-4">
            <div className="text-2xl font-bold">100+</div>
            <div className="text-sm text-primary-200">Stocks Tracked</div>
          </div>
          <div className="bg-white/10 rounded-lg p-4">
            <div className="text-2xl font-bold">Real-time</div>
            <div className="text-sm text-primary-200">Data Updates</div>
          </div>
        </div>
      </div>

      {/* Sector Performance */}
      <section>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          Sector Performance
        </h2>
        {sectorsLoading ? (
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="animate-pulse bg-gray-200 dark:bg-gray-700 rounded-lg h-20" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {sectors?.slice(0, 6).map((sector) => (
              <div
                key={sector.sector}
                className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700"
              >
                <div className="text-sm font-medium text-gray-600 dark:text-gray-400 truncate">
                  {sector.sector}
                </div>
                <div className={`text-lg font-bold ${getChangeColor(sector.change1D)}`}>
                  {formatPercent(sector.change1D)}
                </div>
                <div className="flex items-center gap-1 text-xs text-gray-500">
                  {sector.change1D >= 0 ? (
                    <TrendingUp className="w-3 h-3 text-green-500" />
                  ) : (
                    <TrendingDown className="w-3 h-3 text-red-500" />
                  )}
                  <span>{sector.stockCount} stocks</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Predefined Screeners */}
      <section>
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          Popular Screeners
        </h2>
        {screenersLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="animate-pulse bg-gray-200 dark:bg-gray-700 rounded-lg h-32" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {screeners?.map((screener) => (
              <div
                key={screener.id}
                className="bg-white dark:bg-gray-800 rounded-lg p-5 border border-gray-200 dark:border-gray-700 hover:border-primary-500 dark:hover:border-primary-500 hover:shadow-lg transition-all group"
              >
                <button
                  onClick={() => handleScreenerClick(screener.id)}
                  className="text-left w-full"
                >
                  <div className="flex items-start gap-3">
                    <div className="p-2 bg-primary-100 dark:bg-primary-900/30 rounded-lg text-primary-600 dark:text-primary-400 group-hover:bg-primary-200 dark:group-hover:bg-primary-800/30 transition-colors">
                      {iconMap[screener.icon || 'search'] || <BarChart3 className="w-5 h-5" />}
                    </div>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400">
                        {screener.name}
                      </h3>
                      <p className="text-sm text-gray-500 dark:text-gray-400 mt-1 line-clamp-2">
                        {screener.description}
                      </p>
                    </div>
                  </div>
                </button>
                <div className="mt-3 flex items-center justify-between">
                  <span
                    className={`inline-block px-2 py-0.5 text-xs font-medium rounded ${
                      categoryColors[screener.category] ||
                      'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
                    }`}
                  >
                    {screener.category}
                  </span>
                  {onCustomizeScreener && (
                    <button
                      onClick={(e) => handleCustomizeClick(e, screener.id)}
                      className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors"
                      title="Customize this screener"
                    >
                      <Sliders className="w-3.5 h-3.5" />
                      Customize
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      {/* Quick Stats */}
      <section className="bg-gray-50 dark:bg-gray-800/50 rounded-xl p-6">
        <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          Market Overview
        </h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div className="text-sm text-gray-500 dark:text-gray-400">Top Sector</div>
            <div className="text-lg font-bold text-gray-900 dark:text-white">
              {sectors?.[0]?.sector || '-'}
            </div>
            <div className={`text-sm ${getChangeColor(sectors?.[0]?.change1D || 0)}`}>
              {formatPercent(sectors?.[0]?.change1D || 0)}
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div className="text-sm text-gray-500 dark:text-gray-400">Worst Sector</div>
            <div className="text-lg font-bold text-gray-900 dark:text-white">
              {sectors?.[sectors.length - 1]?.sector || '-'}
            </div>
            <div className={`text-sm ${getChangeColor(sectors?.[sectors.length - 1]?.change1D || 0)}`}>
              {formatPercent(sectors?.[sectors.length - 1]?.change1D || 0)}
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div className="text-sm text-gray-500 dark:text-gray-400">Total Sectors</div>
            <div className="text-lg font-bold text-gray-900 dark:text-white">
              {sectors?.length || 0}
            </div>
          </div>
          <div className="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div className="text-sm text-gray-500 dark:text-gray-400">Screeners Available</div>
            <div className="text-lg font-bold text-gray-900 dark:text-white">
              {screeners?.length || 0}
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
