import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Layout } from './components/shared/Layout';

import { ProjectsList } from './pages/ProjectsList';
import { ProjectDetail } from './pages/ProjectDetail';
import { Login } from './pages/Login';
import { SSOSuccess } from './pages/SSOSuccess';
import { FlagDetail } from './pages/FlagDetail';
import { UsersManagement } from './pages/UsersManagement';
import { SystemSettings } from './pages/SystemSettings';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

import { Toaster } from 'react-hot-toast';

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Toaster position="bottom-right" />
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/sso-success" element={<SSOSuccess />} />
          <Route path="/" element={<Layout />}>
            <Route index element={<Navigate to="/projects" replace />} />
            <Route path="projects" element={<ProjectsList />} />
            <Route path="projects/:projectId/*" element={<ProjectDetail />} />
            <Route path="projects/:projectId/flags/:flagId" element={<FlagDetail />} />
            <Route path="settings/users" element={<UsersManagement />} />
            <Route path="settings/system" element={<SystemSettings />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
