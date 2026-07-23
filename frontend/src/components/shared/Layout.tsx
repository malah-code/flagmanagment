import { Link, Outlet } from 'react-router-dom';
import { Flag } from 'lucide-react';

export const Layout = () => {
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
