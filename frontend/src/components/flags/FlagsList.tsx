import { useState } from 'react';
import { useFlags, useDeleteFlag } from '../../hooks/useFlags';
import { useFlagStates, useUpdateFlagState } from '../../hooks/useFlagStates';
import { CreateFlagDialog } from './CreateFlagDialog';
import { StaleBadge } from './StaleBadge';
import { FlagActions } from './FlagActions';
import type { LifecycleState } from '../../types';
import { Link } from 'react-router-dom';
import { Plus, Flag, Trash2, Loader2, Search, Tag, Filter, Link2, Sliders, Edit3 } from 'lucide-react';
import { Switch } from '../ui/Switch';
import { TagFilter } from './TagFilter';
import toast from 'react-hot-toast';

interface FlagsListProps {
  projectId: string;
  environmentId?: string;
  onNavigateToTargeting?: () => void;
}

export const FlagsList = ({ projectId, environmentId, onNavigateToTargeting }: FlagsListProps) => {
  const { data: flags = [], isLoading, isError, error, refetch } = useFlags(projectId);
  const deleteMutation = useDeleteFlag(projectId);

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [lifecycleFilter, setLifecycleFilter] = useState<LifecycleState | 'ALL'>('ALL');
  const [tagFilter, setTagFilter] = useState('');
  
  const { data: flagStates = [] } = useFlagStates(environmentId || '');
  const updateMutation = useUpdateFlagState(environmentId || '');
  const [togglingStateId, setTogglingStateId] = useState<string | null>(null);

  const handleToggle = async (flagId: string, currentEnabled: boolean) => {
    const state = flagStates.find((s) => s.flagId === flagId);
    if (!state) return;
    
    setTogglingStateId(flagId);
    try {
      await updateMutation.mutateAsync({
        flagId,
        payload: { 
          isEnabled: !currentEnabled,
          targetingRules: state.targetingRules || { rules: [] },
          remoteConfig: state.remoteConfig || {},
          rolloutRules: state.rolloutRules ? { rules: state.rolloutRules } : undefined,
          defaultVariation: state.defaultVariation,
        },
      });
      toast.success('Flag toggled');
    } catch (err: any) {
      toast.error('Failed to toggle flag');
    } finally {
      setTogglingStateId(null);
    }
  };

  const filteredFlags = flags.filter((f) => {
    const matchesSearch =
      f.key.toLowerCase().includes(searchTerm.toLowerCase()) ||
      f.description?.toLowerCase().includes(searchTerm.toLowerCase());
    // Assume flag has lifecycle status attached or default ACTIVE
    const matchesLifecycle =
      lifecycleFilter === 'ALL' || (f as any).lifecycleState === lifecycleFilter;
    const matchesTag = !tagFilter || (f.tags && f.tags.includes(tagFilter));
    return matchesSearch && matchesLifecycle && matchesTag;
  });

  const allTags = Array.from(new Set(flags.flatMap(f => f.tags || []))).sort();

  const handleDelete = async (id: string, key: string) => {
    if (confirm(`Are you sure you want to delete feature flag "${key}"?`)) {
      await deleteMutation.mutateAsync(id);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">Feature Flags</h2>
          <p className="text-sm text-slate-500">Manage flags defined for this project.</p>
        </div>
        <button
          onClick={() => setIsCreateOpen(true)}
          className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium px-3.5 py-2 rounded-lg shadow-sm transition-colors self-start sm:self-auto"
        >
          <Plus className="w-4 h-4" />
          <span>New Feature Flag</span>
        </button>
      </div>

      <div className="flex flex-col sm:flex-row items-center gap-3">
        <div className="flex items-center gap-3 bg-white px-3.5 py-2.5 rounded-lg border border-slate-200 shadow-sm max-w-md w-full focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-indigo-500 transition-all">
          <Search className="w-4 h-4 text-slate-400" />
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Filter flags by key..."
            className="w-full text-sm outline-none bg-transparent text-slate-900 placeholder:text-slate-400"
          />
        </div>

        <div className="flex items-center gap-2 bg-white px-3.5 py-2.5 rounded-lg border border-slate-200 shadow-sm text-sm focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-indigo-500 transition-all">
          <Filter className="w-4 h-4 text-slate-400" />
          <select
            value={lifecycleFilter}
            onChange={(e) => setLifecycleFilter(e.target.value as LifecycleState | 'ALL')}
            className="bg-transparent text-slate-700 outline-none text-sm font-medium"
          >
            <option value="ALL">All Statuses</option>
            <option value="ACTIVE">Active</option>
            <option value="STALE">Stale</option>
            <option value="DEPRECATED">Deprecated</option>
            <option value="ARCHIVED">Archived</option>
          </select>
        </div>

        <TagFilter
          allTags={allTags}
          selectedTag={tagFilter}
          onChange={setTagFilter}
        />
      </div>

      {isLoading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="w-6 h-6 animate-spin text-indigo-600" />
        </div>
      ) : isError ? (
        <div className="bg-red-50 text-red-600 p-4 rounded-xl border border-red-200 text-sm">
          Failed to load flags: {(error as Error).message}
        </div>
      ) : filteredFlags.length === 0 ? (
        <div className="bg-white border border-slate-200 rounded-xl p-10 text-center space-y-4">
          <div className="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-full flex items-center justify-center mx-auto">
            <Flag className="w-6 h-6" />
          </div>
          <div className="space-y-1">
            <h3 className="text-sm font-medium text-slate-900">
              {searchTerm || lifecycleFilter !== 'ALL' ? 'No results found' : 'No feature flags'}
            </h3>
            <p className="text-sm text-slate-500">
              {searchTerm || lifecycleFilter !== 'ALL' ? 'Try adjusting your search or filters.' : 'Get started by creating your first feature flag.'}
            </p>
          </div>
          {!searchTerm && lifecycleFilter === 'ALL' && (
            <div className="pt-2">
              <button
                onClick={() => setIsCreateOpen(true)}
                className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium px-4 py-2 rounded-lg shadow-sm transition-colors"
              >
                <Plus className="w-4 h-4" />
                <span>Create your first Feature Flag</span>
              </button>
            </div>
          )}
        </div>
      ) : (
        <div className="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm overflow-x-auto">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs uppercase tracking-wider font-semibold text-slate-500">
              <tr>
                <th className="px-6 py-3.5">Flag Key</th>
                <th className="px-6 py-3.5">Type</th>
                <th className="px-6 py-3.5">Lifecycle</th>
                <th className="px-6 py-3.5">Description</th>
                <th className="px-6 py-3.5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {filteredFlags.map((flag) => (
                <tr key={flag.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4 font-mono font-medium text-slate-900 flex items-center gap-2">
                    <Flag className="w-4 h-4 text-indigo-600 shrink-0" />
                    <Link to={`/projects/${projectId}/flags/${flag.id}`} className="text-indigo-600 hover:text-indigo-800 hover:underline">
                      {flag.key}
                    </Link>
                    {flag.parentFlagId && (
                      <span title="Depends on a parent flag" className="text-amber-500 flex items-center">
                        <Link2 className="w-3.5 h-3.5" />
                      </span>
                    )}
                    {flag.tags && flag.tags.length > 0 && (
                      <div className="flex gap-1 flex-wrap ml-2">
                        {flag.tags.map(t => (
                          <span key={t} className="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-medium bg-slate-100 text-slate-600 border border-slate-200">
                            {t}
                          </span>
                        ))}
                      </div>
                    )}
                  </td>
                  <td className="px-6 py-4">
                    <span className="inline-flex items-center gap-1 bg-slate-100 text-slate-700 text-xs font-semibold px-2.5 py-1 rounded-md">
                      <Tag className="w-3 h-3" />
                      {flag.type}
                    </span>
                  </td>
                  <td className="px-6 py-4">
                    <StaleBadge state={(flag as any).lifecycleState} />
                  </td>
                  <td className="px-6 py-4 text-slate-500 max-w-md truncate">
                    {flag.description || '—'}
                  </td>
                  <td className="px-6 py-4 text-right flex items-center justify-end gap-2">
                    {environmentId && (
                      (() => {
                        const state = flagStates.find(s => s.flagId === flag.id);
                        if (!state) return null;
                        const isEnabled = state.isEnabled;
                        return (
                          <div className="flex items-center ml-2">
                            <Switch
                              checked={isEnabled}
                              onChange={() => handleToggle(flag.id, isEnabled)}
                              loading={updateMutation.isPending && togglingStateId === flag.id}
                            />
                          </div>
                        );
                      })()
                    )}
                    {onNavigateToTargeting && (
                      <button
                        onClick={onNavigateToTargeting}
                        className="inline-flex items-center gap-1 text-xs font-semibold text-indigo-600 hover:text-indigo-700 bg-indigo-50 hover:bg-indigo-100 px-2.5 py-1.5 rounded transition-colors"
                        title="Configure Targeting & Toggle State"
                      >
                        <Sliders className="w-3.5 h-3.5" /> Targeting & Toggle
                      </button>
                    )}
                    <Link
                      to={`/projects/${projectId}/flags/${flag.id}`}
                      className="text-slate-400 hover:text-indigo-600 transition-colors p-1.5 rounded hover:bg-slate-100"
                      title="Edit Flag Details"
                    >
                      <Edit3 className="w-4 h-4" />
                    </Link>
                    <FlagActions
                      envId={environmentId}
                      flagId={flag.id}
                      currentLifecycle={(flag as any).lifecycleState}
                      onStateChanged={refetch}
                    />
                    <button
                      onClick={() => handleDelete(flag.id, flag.key)}
                      className="text-slate-400 hover:text-red-600 transition-colors p-1.5 rounded hover:bg-slate-100"
                      title="Delete Flag"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <CreateFlagDialog
        projectId={projectId}
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
      />
    </div>
  );
};
