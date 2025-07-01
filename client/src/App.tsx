import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Dashboard } from './components/Dashboard';
import { ScheduleDetails } from './components/ScheduleDetails';
import { ClockOut } from './components/ClockOut';

// Create a client with optimized settings to prevent redundant calls
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      refetchOnWindowFocus: false,
      refetchOnMount: false, // Prevent refetch on mount if data is fresh
      refetchOnReconnect: false, // Prevent refetch on reconnect if data is fresh
      retry: 1, // Reduce retry attempts
      retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <div className="min-h-screen bg-background">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/schedule/:id" element={<ScheduleDetails />} />
            <Route path="/schedule/:id/clock-out" element={<ClockOut />} />
          </Routes>
        </div>
      </Router>
    </QueryClientProvider>
  );
}

export default App;
