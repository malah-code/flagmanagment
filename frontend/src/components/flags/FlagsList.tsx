import { useState } from 'react';
import { useFlags, useDeleteFlag } from '../../hooks/useFlags';
import { CreateFlagDialog } from './CreateFlagDialog';
import { Plus, Flag, Trash2, Loader2, Search, Tag } from 'lucide-react';

interface FlagsListProps {
  projectId: string;
}

export const FlagsList = ({ projectId }: FlagsListProps) => {
  const { data: flags = [], isLoading, isError, error } = useFlags(projectId);
  const deleteMutation = useDeleteFlag(projectId);

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const filteredFlags = flags.filter(
    (f) =>
      f.key.toLowerCase().includes(searchTerm.toLowerCase()) ||
      f.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

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

      <div className="flex items-center gap-3 bg-white px-3.5 py-2.5 rounded-lg border border-slate-200 shadow-sm max-w-md">
        <Search className="w-4 h-4 text-slate-400" />
        <input
          type="text"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          placeholder="Filter flags by key..."
          className="w-full text-sm outline-none bg-transparent text-slate-900 placeholder:text-slate-400"
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
        <div className="bg-white border border-slate-200 rounded-xl p-8 text-center space-y-3">
          <div className="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-full flex items-center justify-center mx-auto">
            <Flag className="w-5 h-5" />
          </div>
          <p className="text-sm text-slate-500">
            {searchTerm ? 'No flags match your filter.' : 'No feature flags created yet for this project.'}
          </p>
        </div>
      ) : (
        <div className="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs uppercase tracking-wider font-semibold text-slate-500">
              <tr>
                <th className="px-6 py-3.5">Flag Key</th>
                <th className="px-6 py-3.5">Type</th>
                <th className="px-6 py-3.5">Description</th>
                <th className="px-6 py-3.5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {filteredFlags.map((flag) => (
                <tr key={flag.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4 font-mono font-medium text-slate-900 flex items-center gap-2">
                    <Flag className="w-4 h-4 text-indigo-600 shrink-0" />
                    <span>{flag.key}</span>
                  </td>
                  <td className="px-6 py-4">
                    <span className="inline-flex items-center gap-1 bg-slate-100 text-slate-700 text-xs font-semibold px-2.5 py-1 rounded-md">
                      <Tag className="w-3 h-3" />
                      {flag.type}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-500 max-w-md truncate">
                    {flag.description || '—'}
                  </td>
                  <td className="px-6 py-4 text-right">
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
