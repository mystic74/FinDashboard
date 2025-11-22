import { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ThemeProvider } from './context/ThemeContext';
import { Header } from './components/Layout/Header';
import { Dashboard } from './components/Dashboard/Dashboard';
import { ScreenerResults } from './components/Screeners/ScreenerResults';
import { CustomScreener } from './components/Screeners/CustomScreener';
import { ScreenerTweaker } from './components/Screeners/ScreenerTweaker';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 2,
      refetchOnWindowFocus: false,
    },
  },
});

type Page = 'dashboard' | 'screeners' | 'custom' | 'screener-results' | 'screener-tweaker';

function AppContent() {
  const [currentPage, setCurrentPage] = useState<Page>('dashboard');
  const [selectedScreener, setSelectedScreener] = useState<string | null>(null);
  const [selectedCountry, setSelectedCountry] = useState<string | undefined>(undefined);

  const handleSelectScreener = (screenerId: string, country?: string) => {
    setSelectedScreener(screenerId);
    setSelectedCountry(country);
    setCurrentPage('screener-results');
  };

  const handleCustomizeScreener = (screenerId: string, country?: string) => {
    setSelectedScreener(screenerId);
    setSelectedCountry(country);
    setCurrentPage('screener-tweaker');
  };

  const handleBack = () => {
    setSelectedScreener(null);
    setSelectedCountry(undefined);
    setCurrentPage('dashboard');
  };

  const handleNavigate = (page: string) => {
    setCurrentPage(page as Page);
    if (page !== 'screener-results' && page !== 'screener-tweaker') {
      setSelectedScreener(null);
      setSelectedCountry(undefined);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors">
      <Header
        currentPage={currentPage === 'screener-results' ? 'dashboard' : currentPage}
        onNavigate={handleNavigate}
      />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {currentPage === 'dashboard' && (
          <Dashboard onSelectScreener={handleSelectScreener} onCustomizeScreener={handleCustomizeScreener} />
        )}
        {currentPage === 'screeners' && (
          <Dashboard onSelectScreener={handleSelectScreener} onCustomizeScreener={handleCustomizeScreener} />
        )}
        {currentPage === 'custom' && <CustomScreener />}
        {currentPage === 'screener-results' && selectedScreener && (
          <ScreenerResults screenerId={selectedScreener} country={selectedCountry} onBack={handleBack} />
        )}
        {currentPage === 'screener-tweaker' && selectedScreener && (
          <ScreenerTweaker screenerId={selectedScreener} country={selectedCountry} onBack={handleBack} />
        )}
      </main>
      <footer className="border-t border-gray-200 dark:border-gray-800 py-6 mt-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center text-sm text-gray-500 dark:text-gray-400">
          Stock Screener Dashboard - Data provided by Yahoo Finance
        </div>
      </footer>
    </div>
  );
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <AppContent />
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;
