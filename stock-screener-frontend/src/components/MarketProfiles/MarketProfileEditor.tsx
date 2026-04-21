import { useState, useEffect, useMemo } from 'react';
import {
  Globe2,
  TrendingUp,
  Volume2,
  Coins,
  Sprout,
  RotateCcw,
  Save,
  Eye,
  ChevronDown,
  ChevronUp,
  Sparkles,
  Info,
  Check,
} from 'lucide-react';

// Market profile type
interface MarketProfile {
  country: string;
  flag: string;
  name: string;
  marketCapMultiplier: number;
  volumeMultiplier: number;
  dividendMultiplier: number;
  growthMultiplier: number;
  description: string;
  marketSize: 'large' | 'medium' | 'small';
}

// Default profiles - these match the backend
const DEFAULT_PROFILES: MarketProfile[] = [
  {
    country: 'USA',
    flag: '🇺🇸',
    name: 'United States',
    marketCapMultiplier: 1.0,
    volumeMultiplier: 1.0,
    dividendMultiplier: 1.0,
    growthMultiplier: 1.0,
    description: 'Baseline - largest and most liquid market',
    marketSize: 'large',
  },
  {
    country: 'Israel',
    flag: '🇮🇱',
    name: 'Israel',
    marketCapMultiplier: 0.1,
    volumeMultiplier: 0.3,
    dividendMultiplier: 0.8,
    growthMultiplier: 0.8,
    description: 'Smaller market with strong tech sector',
    marketSize: 'small',
  },
  {
    country: 'UK',
    flag: '🇬🇧',
    name: 'United Kingdom',
    marketCapMultiplier: 0.5,
    volumeMultiplier: 0.6,
    dividendMultiplier: 1.2,
    growthMultiplier: 0.9,
    description: 'Mature market with dividend focus',
    marketSize: 'large',
  },
  {
    country: 'Germany',
    flag: '🇩🇪',
    name: 'Germany',
    marketCapMultiplier: 0.5,
    volumeMultiplier: 0.5,
    dividendMultiplier: 1.0,
    growthMultiplier: 0.9,
    description: 'Industrial powerhouse of Europe',
    marketSize: 'large',
  },
  {
    country: 'Japan',
    flag: '🇯🇵',
    name: 'Japan',
    marketCapMultiplier: 0.6,
    volumeMultiplier: 0.7,
    dividendMultiplier: 0.7,
    growthMultiplier: 0.8,
    description: 'Established market with unique characteristics',
    marketSize: 'large',
  },
  {
    country: 'China',
    flag: '🇨🇳',
    name: 'China',
    marketCapMultiplier: 0.8,
    volumeMultiplier: 1.2,
    dividendMultiplier: 0.5,
    growthMultiplier: 1.2,
    description: 'High growth, high volume emerging market',
    marketSize: 'large',
  },
  {
    country: 'India',
    flag: '🇮🇳',
    name: 'India',
    marketCapMultiplier: 0.2,
    volumeMultiplier: 0.4,
    dividendMultiplier: 0.6,
    growthMultiplier: 1.3,
    description: 'Fast-growing emerging market',
    marketSize: 'medium',
  },
  {
    country: 'Brazil',
    flag: '🇧🇷',
    name: 'Brazil',
    marketCapMultiplier: 0.3,
    volumeMultiplier: 0.4,
    dividendMultiplier: 1.5,
    growthMultiplier: 0.9,
    description: 'Latin America\'s largest market',
    marketSize: 'medium',
  },
  {
    country: 'Canada',
    flag: '🇨🇦',
    name: 'Canada',
    marketCapMultiplier: 0.4,
    volumeMultiplier: 0.5,
    dividendMultiplier: 1.1,
    growthMultiplier: 0.9,
    description: 'Resource-rich developed market',
    marketSize: 'medium',
  },
  {
    country: 'France',
    flag: '🇫🇷',
    name: 'France',
    marketCapMultiplier: 0.5,
    volumeMultiplier: 0.5,
    dividendMultiplier: 1.1,
    growthMultiplier: 0.9,
    description: 'Major European market',
    marketSize: 'large',
  },
  {
    country: 'Switzerland',
    flag: '🇨🇭',
    name: 'Switzerland',
    marketCapMultiplier: 0.5,
    volumeMultiplier: 0.4,
    dividendMultiplier: 1.3,
    growthMultiplier: 0.8,
    description: 'Quality-focused defensive market',
    marketSize: 'medium',
  },
  {
    country: 'Australia',
    flag: '🇦🇺',
    name: 'Australia',
    marketCapMultiplier: 0.3,
    volumeMultiplier: 0.4,
    dividendMultiplier: 1.4,
    growthMultiplier: 0.85,
    description: 'Resource and banking focused',
    marketSize: 'medium',
  },
];

// Multiplier config for UI
const MULTIPLIER_CONFIG = {
  marketCapMultiplier: {
    label: 'Market Cap',
    icon: TrendingUp,
    color: 'blue',
    description: 'Scales market cap thresholds relative to US market',
    gradient: 'from-blue-500 to-blue-600',
    bgLight: 'bg-blue-50',
    bgDark: 'dark:bg-blue-900/20',
    textColor: 'text-blue-600 dark:text-blue-400',
  },
  volumeMultiplier: {
    label: 'Volume',
    icon: Volume2,
    color: 'purple',
    description: 'Adjusts trading volume expectations',
    gradient: 'from-purple-500 to-purple-600',
    bgLight: 'bg-purple-50',
    bgDark: 'dark:bg-purple-900/20',
    textColor: 'text-purple-600 dark:text-purple-400',
  },
  dividendMultiplier: {
    label: 'Dividend',
    icon: Coins,
    color: 'emerald',
    description: 'Scales dividend yield expectations',
    gradient: 'from-emerald-500 to-emerald-600',
    bgLight: 'bg-emerald-50',
    bgDark: 'dark:bg-emerald-900/20',
    textColor: 'text-emerald-600 dark:text-emerald-400',
  },
  growthMultiplier: {
    label: 'Growth',
    icon: Sprout,
    color: 'amber',
    description: 'Adjusts growth rate expectations',
    gradient: 'from-amber-500 to-amber-600',
    bgLight: 'bg-amber-50',
    bgDark: 'dark:bg-amber-900/20',
    textColor: 'text-amber-600 dark:text-amber-400',
  },
};

type MultiplierKey = keyof typeof MULTIPLIER_CONFIG;

// Beautiful slider component
function ProfileSlider({
  value,
  onChange,
  config,
  disabled = false,
}: {
  value: number;
  onChange: (value: number) => void;
  config: typeof MULTIPLIER_CONFIG[MultiplierKey];
  disabled?: boolean;
}) {
  const percentage = Math.min(value * 100, 200);
  const Icon = config.icon;

  return (
    <div className="group">
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <div className={`p-1.5 rounded-lg ${config.bgLight} ${config.bgDark}`}>
            <Icon className={`w-4 h-4 ${config.textColor}`} />
          </div>
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {config.label}
          </span>
        </div>
        <div className="flex items-center gap-2">
          <span className={`text-lg font-bold tabular-nums ${config.textColor}`}>
            {value.toFixed(2)}x
          </span>
        </div>
      </div>

      <div className="relative h-3 bg-gray-100 dark:bg-gray-700 rounded-full overflow-hidden">
        {/* Background track */}
        <div className="absolute inset-0 opacity-20">
          <div
            className={`h-full bg-gradient-to-r ${config.gradient}`}
            style={{ width: '100%' }}
          />
        </div>

        {/* Filled track */}
        <div
          className={`absolute inset-y-0 left-0 bg-gradient-to-r ${config.gradient} rounded-full transition-all duration-200`}
          style={{ width: `${Math.min(percentage, 100)}%` }}
        />

        {/* 1.0x marker */}
        <div
          className="absolute top-0 bottom-0 w-0.5 bg-gray-400 dark:bg-gray-500"
          style={{ left: '50%' }}
        />

        {/* Input slider */}
        <input
          type="range"
          min="0.05"
          max="2"
          step="0.05"
          value={value}
          onChange={(e) => onChange(parseFloat(e.target.value))}
          disabled={disabled}
          className="absolute inset-0 w-full h-full opacity-0 cursor-pointer disabled:cursor-not-allowed"
        />

        {/* Custom thumb indicator */}
        <div
          className={`absolute top-1/2 -translate-y-1/2 w-5 h-5 rounded-full bg-white border-2 shadow-lg transition-all duration-200 group-hover:scale-110 ${
            disabled ? 'border-gray-300' : `border-current ${config.textColor}`
          }`}
          style={{ left: `calc(${Math.min(percentage / 2, 100)}% - 10px)` }}
        />
      </div>

      <div className="flex justify-between mt-1 text-xs text-gray-400">
        <span>0.05x</span>
        <span className="font-medium">1.0x (baseline)</span>
        <span>2.0x</span>
      </div>
    </div>
  );
}

// Country profile card
function ProfileCard({
  profile,
  isExpanded,
  onToggle,
  onChange,
  onReset,
  isModified,
  isUSA,
}: {
  profile: MarketProfile;
  isExpanded: boolean;
  onToggle: () => void;
  onChange: (key: MultiplierKey, value: number) => void;
  onReset: () => void;
  isModified: boolean;
  isUSA: boolean;
}) {
  const sizeColors = {
    large: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    medium: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
    small: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
  };

  return (
    <div
      className={`bg-white dark:bg-gray-800 rounded-xl border transition-all duration-300 ${
        isExpanded
          ? 'border-primary-500 shadow-lg shadow-primary-500/10'
          : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
      }`}
    >
      {/* Header */}
      <button
        onClick={onToggle}
        className="w-full p-4 flex items-center gap-4 text-left"
      >
        {/* Flag */}
        <div className="text-4xl">{profile.flag}</div>

        {/* Info */}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <h3 className="font-semibold text-gray-900 dark:text-white">
              {profile.name}
            </h3>
            {isUSA && (
              <span className="px-2 py-0.5 text-xs font-medium bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400 rounded-full">
                Baseline
              </span>
            )}
            {isModified && !isUSA && (
              <span className="px-2 py-0.5 text-xs font-medium bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400 rounded-full">
                Modified
              </span>
            )}
          </div>
          <p className="text-sm text-gray-500 dark:text-gray-400 truncate">
            {profile.description}
          </p>
        </div>

        {/* Quick stats */}
        <div className="hidden sm:flex items-center gap-3">
          <span className={`px-2 py-1 text-xs font-medium rounded-lg ${sizeColors[profile.marketSize]}`}>
            {profile.marketSize.charAt(0).toUpperCase() + profile.marketSize.slice(1)} Market
          </span>
          <div className="flex items-center gap-1 text-sm text-gray-500">
            <TrendingUp className="w-4 h-4" />
            <span className="tabular-nums">{profile.marketCapMultiplier.toFixed(1)}x</span>
          </div>
        </div>

        {/* Expand icon */}
        <div className={`p-2 rounded-lg transition-colors ${isExpanded ? 'bg-primary-100 dark:bg-primary-900/30' : 'bg-gray-100 dark:bg-gray-700'}`}>
          {isExpanded ? (
            <ChevronUp className={`w-5 h-5 ${isExpanded ? 'text-primary-600' : 'text-gray-500'}`} />
          ) : (
            <ChevronDown className="w-5 h-5 text-gray-500" />
          )}
        </div>
      </button>

      {/* Expanded content */}
      {isExpanded && (
        <div className="px-4 pb-4 border-t border-gray-100 dark:border-gray-700">
          <div className="pt-4 space-y-5">
            {/* Multiplier sliders */}
            {(Object.keys(MULTIPLIER_CONFIG) as MultiplierKey[]).map((key) => (
              <ProfileSlider
                key={key}
                value={profile[key]}
                onChange={(value) => onChange(key, value)}
                config={MULTIPLIER_CONFIG[key]}
                disabled={isUSA}
              />
            ))}

            {/* Actions */}
            {!isUSA && (
              <div className="flex items-center justify-between pt-4 border-t border-gray-100 dark:border-gray-700">
                <p className="text-xs text-gray-400 flex items-center gap-1">
                  <Info className="w-3 h-3" />
                  Changes apply to all screeners for this market
                </p>
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    onReset();
                  }}
                  disabled={!isModified}
                  className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <RotateCcw className="w-4 h-4" />
                  Reset to Default
                </button>
              </div>
            )}

            {isUSA && (
              <div className="flex items-center gap-2 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg text-sm text-blue-700 dark:text-blue-300">
                <Info className="w-4 h-4 flex-shrink-0" />
                <span>USA is the baseline market. All other markets are scaled relative to it.</span>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

// Comparison visualization
function ProfileComparison({ profiles }: { profiles: MarketProfile[] }) {
  const sortedByMarketCap = [...profiles].sort((a, b) => b.marketCapMultiplier - a.marketCapMultiplier);

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
      <h3 className="font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
        <Eye className="w-5 h-5 text-primary-600" />
        Market Comparison
      </h3>

      <div className="space-y-3">
        {sortedByMarketCap.map((profile) => (
          <div key={profile.country} className="flex items-center gap-3">
            <span className="text-2xl w-10">{profile.flag}</span>
            <span className="w-24 text-sm font-medium text-gray-700 dark:text-gray-300 truncate">
              {profile.country}
            </span>
            <div className="flex-1 h-6 bg-gray-100 dark:bg-gray-700 rounded-full overflow-hidden relative">
              {/* Market Cap bar */}
              <div
                className="absolute inset-y-0 left-0 bg-gradient-to-r from-blue-500 to-blue-400 rounded-full transition-all duration-500"
                style={{ width: `${Math.min(profile.marketCapMultiplier * 50, 100)}%` }}
              />
              {/* Volume overlay */}
              <div
                className="absolute inset-y-0 left-0 bg-purple-500/30 rounded-full transition-all duration-500"
                style={{ width: `${Math.min(profile.volumeMultiplier * 50, 100)}%` }}
              />
            </div>
            <span className="w-16 text-right text-sm tabular-nums font-medium text-gray-600 dark:text-gray-400">
              {profile.marketCapMultiplier.toFixed(2)}x
            </span>
          </div>
        ))}
      </div>

      <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700 flex items-center gap-4 text-xs text-gray-500">
        <div className="flex items-center gap-1.5">
          <div className="w-3 h-3 rounded-full bg-gradient-to-r from-blue-500 to-blue-400" />
          <span>Market Cap</span>
        </div>
        <div className="flex items-center gap-1.5">
          <div className="w-3 h-3 rounded-full bg-purple-500/50" />
          <span>Volume (overlay)</span>
        </div>
      </div>
    </div>
  );
}

// Main component
export function MarketProfileEditor() {
  const [profiles, setProfiles] = useState<MarketProfile[]>(DEFAULT_PROFILES);
  const [expandedCountry, setExpandedCountry] = useState<string | null>('Israel');
  const [saveStatus, setSaveStatus] = useState<'idle' | 'saving' | 'saved'>('idle');

  // Load from localStorage
  useEffect(() => {
    const saved = localStorage.getItem('marketProfiles');
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        setProfiles(DEFAULT_PROFILES.map(dp => {
          const savedProfile = parsed.find((p: MarketProfile) => p.country === dp.country);
          return savedProfile ? { ...dp, ...savedProfile } : dp;
        }));
      } catch (e) {
        console.error('Failed to load profiles:', e);
      }
    }
  }, []);

  // Check if profile is modified
  const isModified = (country: string) => {
    const current = profiles.find(p => p.country === country);
    const original = DEFAULT_PROFILES.find(p => p.country === country);
    if (!current || !original) return false;
    return (
      current.marketCapMultiplier !== original.marketCapMultiplier ||
      current.volumeMultiplier !== original.volumeMultiplier ||
      current.dividendMultiplier !== original.dividendMultiplier ||
      current.growthMultiplier !== original.growthMultiplier
    );
  };

  // Update a profile
  const updateProfile = (country: string, key: MultiplierKey, value: number) => {
    setProfiles(prev => prev.map(p =>
      p.country === country ? { ...p, [key]: value } : p
    ));
    setSaveStatus('idle');
  };

  // Reset a profile
  const resetProfile = (country: string) => {
    const original = DEFAULT_PROFILES.find(p => p.country === country);
    if (original) {
      setProfiles(prev => prev.map(p =>
        p.country === country ? { ...original } : p
      ));
    }
    setSaveStatus('idle');
  };

  // Reset all profiles
  const resetAll = () => {
    setProfiles([...DEFAULT_PROFILES]);
    localStorage.removeItem('marketProfiles');
    setSaveStatus('idle');
  };

  // Save profiles
  const saveProfiles = () => {
    setSaveStatus('saving');
    localStorage.setItem('marketProfiles', JSON.stringify(profiles));
    setTimeout(() => setSaveStatus('saved'), 500);
    setTimeout(() => setSaveStatus('idle'), 2000);
  };

  // Count modified profiles
  const modifiedCount = useMemo(() =>
    profiles.filter(p => p.country !== 'USA' && isModified(p.country)).length,
    [profiles]
  );

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <div className="p-2 bg-gradient-to-br from-primary-500 to-primary-600 rounded-xl text-white shadow-lg shadow-primary-500/25">
              <Globe2 className="w-6 h-6" />
            </div>
            <div>
              <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
                Market Profiles
              </h1>
              <p className="text-gray-500 dark:text-gray-400">
                Customize how screeners adapt to different markets
              </p>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-3">
          {modifiedCount > 0 && (
            <span className="text-sm text-amber-600 dark:text-amber-400">
              {modifiedCount} profile{modifiedCount > 1 ? 's' : ''} modified
            </span>
          )}
          <button
            onClick={resetAll}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
          >
            <RotateCcw className="w-4 h-4" />
            Reset All
          </button>
          <button
            onClick={saveProfiles}
            disabled={saveStatus === 'saving'}
            className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-gradient-to-r from-primary-600 to-primary-500 rounded-lg hover:from-primary-700 hover:to-primary-600 shadow-lg shadow-primary-500/25 transition-all disabled:opacity-50"
          >
            {saveStatus === 'saving' ? (
              <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
            ) : saveStatus === 'saved' ? (
              <Check className="w-4 h-4" />
            ) : (
              <Save className="w-4 h-4" />
            )}
            {saveStatus === 'saved' ? 'Saved!' : 'Save Changes'}
          </button>
        </div>
      </div>

      {/* Info banner */}
      <div className="flex items-start gap-4 p-4 bg-gradient-to-r from-primary-50 to-blue-50 dark:from-primary-900/20 dark:to-blue-900/20 rounded-xl border border-primary-100 dark:border-primary-800/50">
        <div className="p-2 bg-white dark:bg-gray-800 rounded-lg shadow-sm">
          <Sparkles className="w-5 h-5 text-primary-600" />
        </div>
        <div>
          <h3 className="font-medium text-gray-900 dark:text-white">How it works</h3>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            Market profiles automatically adjust screener thresholds when filtering by country.
            For example, Israel's 0.1x market cap multiplier means a "small cap" stock in Israel
            would be 10x smaller than a US small cap. This ensures your screening criteria makes
            sense for each market's unique characteristics.
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Profile list */}
        <div className="lg:col-span-2 space-y-3">
          {profiles.map((profile) => (
            <ProfileCard
              key={profile.country}
              profile={profile}
              isExpanded={expandedCountry === profile.country}
              onToggle={() => setExpandedCountry(
                expandedCountry === profile.country ? null : profile.country
              )}
              onChange={(key, value) => updateProfile(profile.country, key, value)}
              onReset={() => resetProfile(profile.country)}
              isModified={isModified(profile.country)}
              isUSA={profile.country === 'USA'}
            />
          ))}
        </div>

        {/* Comparison sidebar */}
        <div className="space-y-6">
          <ProfileComparison profiles={profiles} />

          {/* Quick reference */}
          <div className="bg-gray-50 dark:bg-gray-800/50 rounded-xl p-4">
            <h4 className="font-medium text-gray-900 dark:text-white mb-3">
              Multiplier Guide
            </h4>
            <div className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <div className="flex items-center gap-2">
                <span className="w-12 font-mono text-xs bg-gray-200 dark:bg-gray-700 px-1.5 py-0.5 rounded">0.1x</span>
                <span>Much smaller (e.g., Israel tech)</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-12 font-mono text-xs bg-gray-200 dark:bg-gray-700 px-1.5 py-0.5 rounded">0.5x</span>
                <span>Half scale (e.g., UK, Germany)</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-12 font-mono text-xs bg-gray-200 dark:bg-gray-700 px-1.5 py-0.5 rounded">1.0x</span>
                <span>Baseline (USA)</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-12 font-mono text-xs bg-gray-200 dark:bg-gray-700 px-1.5 py-0.5 rounded">1.5x</span>
                <span>Higher expectations</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
