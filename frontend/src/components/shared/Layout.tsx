import { Navigate, Link, Outlet, useNavigate } from 'react-router-dom';
import { Flag, User, LogOut, Settings, Server } from 'lucide-react';
import { authService } from '../../services/auth';
import { useState, useRef, useEffect } from 'react';

export const Layout = () => {
  const navigate = useNavigate();
  const [isProfileOpen, setIsProfileOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsProfileOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  if (!authService.isAuthenticated()) {
    return <Navigate to="/login" replace />;
  }

  const handleLogout = () => {
    authService.logout();
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-slate-50 text-slate-900 flex flex-col font-sans">
      <header className="bg-white border-b border-slate-200 px-6 py-4 flex items-center justify-between sticky top-0 z-10">
        <div className="flex items-center gap-2">
          <div className="bg-indigo-600 p-1.5 rounded-md">
            <Flag className="w-5 h-5 text-white" />
          </div>
          <Link to="/" className="text-xl font-bold tracking-tight text-slate-900 hover:text-indigo-600 transition-colors">
            Flag<span className="text-indigo-600">Management</span>
          </Link>
        </div>
        <nav className="flex items-center gap-6 text-sm font-medium">
          <Link to="/projects" className="text-slate-600 hover:text-indigo-600 transition-colors">
            Projects
          </Link>
          
          <div className="relative" ref={dropdownRef}>
            <button 
              onClick={() => setIsProfileOpen(!isProfileOpen)}
              className="flex items-center justify-center w-8 h-8 rounded-full bg-slate-100 hover:bg-slate-200 text-slate-600 transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
            >
              <User className="w-4 h-4" />
            </button>

            {isProfileOpen && (
              <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-slate-200 py-1 z-50">
                <div className="px-4 py-2 border-b border-slate-100">
                  <p className="text-sm font-medium text-slate-900">Administrator</p>
                </div>
                <Link 
                  to="/settings/users" 
                  onClick={() => setIsProfileOpen(false)}
                  className="flex items-center gap-2 w-full px-4 py-2 text-sm text-slate-600 hover:bg-slate-50 transition-colors"
                >
                  <Settings className="w-4 h-4" />
                  <span>Team Settings</span>
                </Link>
                <Link 
                  to="/settings/system" 
                  onClick={() => setIsProfileOpen(false)}
                  className="flex items-center gap-2 w-full px-4 py-2 text-sm text-slate-600 hover:bg-slate-50 transition-colors"
                >
                  <Server className="w-4 h-4" />
                  <span>System Settings</span>
                </Link>
                <button 
                  onClick={handleLogout}
                  className="flex items-center gap-2 w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors text-left"
                >
                  <LogOut className="w-4 h-4" />
                  <span>Sign out</span>
                </button>
              </div>
            )}
          </div>
        </nav>
      </header>

      <main className="flex-1 max-w-7xl w-full mx-auto p-6 md:p-8">
        <Outlet />
      </main>

      <footer className="bg-white border-t border-slate-200 py-6 text-center text-sm text-slate-500 mt-auto">
        <p>FlagManagement Dashboard &copy; {new Date().getFullYear()}</p>
      </footer>
    </div>
  );
};
