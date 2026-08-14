import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { ArrowLeft, Server, Flag, Settings, Folder, Activity, Plus } from 'lucide-react';
import type { Environment } from '../../types';

interface EnvironmentSidebarProps {
  projectId: string;
  projectName: string;
  environments: Environment[];
  onAddEnvironment?: () => void;
}

export const EnvironmentSidebar: React.FC<EnvironmentSidebarProps> = ({
  projectId,
  projectName,
  environments,
  onAddEnvironment,
}) => {
  const location = useLocation();

  return (
    <div className="w-full lg:w-64 flex-shrink-0 bg-slate-900 text-slate-200 rounded-2xl p-5 shadow-xl border border-slate-800 space-y-6">
      {/* Header */}
      <div>
        <Link
          to="/projects"
          className="inline-flex items-center gap-1.5 text-xs font-semibold text-slate-400 hover:text-slate-200 transition-colors mb-4 group"
        >
          <ArrowLeft className="w-3.5 h-3.5 group-hover:-translate-x-0.5 transition-transform" />
          Back to Projects
        </Link>
        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 rounded-xl">
            <Folder className="w-5 h-5" />
          </div>
          <div className="min-w-0">
            <h1 className="text-base font-bold text-white leading-tight truncate">{projectName}</h1>
            <span className="text-[11px] font-medium text-emerald-400 flex items-center gap-1 mt-0.5">
              <Activity className="w-3 h-3" /> Live Workspace
            </span>
          </div>
        </div>
      </div>

      <nav className="space-y-6">
        {/* Flag Definitions */}
        <div>
          <div className="text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2 px-2">
            Overview
          </div>
          <Link
            to={`/projects/${projectId}/flags`}
            className={`flex items-center gap-2.5 px-3 py-2.5 rounded-xl text-xs font-semibold transition-all ${
              location.pathname === `/projects/${projectId}/flags`
                ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30'
                : 'text-slate-300 hover:bg-slate-800/80 hover:text-white'
            }`}
          >
            <Flag className="w-4 h-4" />
            <span>All Flag Definitions</span>
          </Link>
        </div>

        {/* Environments */}
        <div>
          <div className="flex items-center justify-between px-2 mb-2">
            <span className="text-[11px] font-bold text-slate-400 uppercase tracking-wider">
              Environments
            </span>
            {onAddEnvironment && (
              <button
                onClick={onAddEnvironment}
                className="text-slate-400 hover:text-indigo-400 p-0.5 rounded transition-colors"
                title="Add Environment"
              >
                <Plus className="w-3.5 h-3.5" />
              </button>
            )}
          </div>
          <div className="space-y-1">
            {environments.length === 0 ? (
              <div className="text-xs text-slate-500 px-3 py-2 italic">No environments yet</div>
            ) : (
              environments.map((env) => {
                const envRoute = `/projects/${projectId}/env/${env.id}`;
                const isEnvActive = location.pathname.startsWith(envRoute);
                return (
                  <Link
                    key={env.id}
                    to={envRoute}
                    className={`flex items-center justify-between px-3 py-2.5 rounded-xl text-xs font-semibold transition-all ${
                      isEnvActive
                        ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30'
                        : 'text-slate-300 hover:bg-slate-800/80 hover:text-white'
                    }`}
                  >
                    <div className="flex items-center gap-2.5 truncate">
                      <Server className={`w-4 h-4 shrink-0 ${isEnvActive ? 'text-white' : 'text-indigo-400'}`} />
                      <span className="truncate">{env.name}</span>
                    </div>
                    <span
                      className={`w-2 h-2 rounded-full shrink-0 ${
                        env.name.toLowerCase().includes('prod')
                          ? 'bg-rose-400'
                          : env.name.toLowerCase().includes('stage')
                          ? 'bg-amber-400'
                          : 'bg-emerald-400'
                      }`}
                    />
                  </Link>
                );
              })
            )}
          </div>
        </div>

        {/* Settings */}
        <div>
          <div className="text-[11px] font-bold text-slate-400 uppercase tracking-wider mb-2 px-2">
            Management
          </div>
          <Link
            to={`/projects/${projectId}/settings`}
            className={`flex items-center gap-2.5 px-3 py-2.5 rounded-xl text-xs font-semibold transition-all ${
              location.pathname.startsWith(`/projects/${projectId}/settings`)
                ? 'bg-indigo-600 text-white shadow-md shadow-indigo-600/30'
                : 'text-slate-300 hover:bg-slate-800/80 hover:text-white'
            }`}
          >
            <Settings className="w-4 h-4" />
            <span>Environment Settings</span>
          </Link>
        </div>
      </nav>
    </div>
  );
};
