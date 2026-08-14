import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { authService } from '../services/auth';
import { AlertCircle } from 'lucide-react';

export const Login: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      await authService.login(email, password);
      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Invalid email or password');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-950 text-slate-50">
      <div className="w-full max-w-md space-y-6 rounded-xl bg-slate-900 p-8 shadow-2xl border border-slate-800">
        <div className="space-y-2 text-center">
          <h1 className="text-3xl font-bold tracking-tight">FlagManagment</h1>
          <p className="text-sm text-slate-400">Sign in to your account</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-300">Email</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => {
                setEmail(e.target.value);
                setError(null);
              }}
              className="mt-1 block w-full rounded-md bg-slate-950 border border-slate-700 px-3 py-2 text-slate-100 focus:border-indigo-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-300">Password</label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => {
                setPassword(e.target.value);
                setError(null);
              }}
              className="mt-1 block w-full rounded-md bg-slate-950 border border-slate-700 px-3 py-2 text-slate-100 focus:border-indigo-500 focus:outline-none"
            />
            {error && (
              <div className="mt-3 p-3 bg-red-950/60 border border-red-800/80 rounded-lg flex items-center gap-2 text-xs font-medium text-red-300">
                <AlertCircle className="w-4 h-4 text-red-400 flex-shrink-0" />
                <span>{error === 'API Error' ? 'Invalid email or password. Please try again.' : error}</span>
              </div>
            )}
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-md bg-indigo-600 px-4 py-2 font-medium text-white hover:bg-indigo-500 disabled:opacity-50"
          >
            {loading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>

        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-slate-700"></div>
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="bg-slate-900 px-2 text-slate-400">Or continue with</span>
          </div>
        </div>

        <button
          onClick={() => authService.ssoLogin('oidc')}
          disabled={loading}
          className="w-full rounded-md bg-slate-800 px-4 py-2 font-medium text-slate-100 border border-slate-700 hover:bg-slate-700 disabled:opacity-50"
        >
          Log in with SSO (OIDC)
        </button>
      </div>
    </div>
  );
};
