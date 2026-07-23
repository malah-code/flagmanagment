import { useState } from 'react';
import { useEnvironments, useDeleteEnvironment } from '../../hooks/useEnvironments';
import { CreateEnvironmentDialog } from './CreateEnvironmentDialog';
import { Plus, Server, Trash2, Loader2, Key } from 'lucide-react';

interface EnvironmentsListProps {
  projectId: string;
}

export const EnvironmentsList = ({ projectId }: EnvironmentsListProps) => {
  const { data: environments = [], isLoading, isError, error } = useEnvironments(projectId);
  const deleteMutation = useDeleteEnvironment(projectId);
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  const handleDelete = async (id: string, name: string) => {
    if (confirm(`Are you sure you want to delete environment "${name}"? Active flag configurations for this environment will be removed.`)) {
      await deleteMutation.mutateAsync(id);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold text-slate-900">Environments</h2>
          <p className="text-sm text-slate-500">Configure target deployment environments for feature flags.</p>
        </div>
        <button
          onClick={() => setIsCreateOpen(true)}
          className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium px-3.5 py-2 rounded-lg shadow-sm transition-colors"
        >
          <Plus className="w-4 h-4" />
          <span>Add Environment</span>
        </button>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="w-6 h-6 animate-spin text-indigo-600" />
        </div>
      ) : isError ? (
        <div className="bg-red-50 text-red-600 p-4 rounded-xl border border-red-200 text-sm">
          Failed to load environments: {(error as Error).message}
        </div>
      ) : environments.length === 0 ? (
        <div className="bg-white border border-slate-200 rounded-xl p-8 text-center space-y-3">
          <div className="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-full flex items-center justify-center mx-auto">
            <Server className="w-5 h-5" />
          </div>
          <p className="text-sm text-slate-500">No environments created yet for this project.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {environments.map((env) => (
            <div
              key={env.id}
              className="bg-white border border-slate-200 rounded-xl p-4 shadow-sm flex items-center justify-between group hover:border-slate-300 transition-all"
            >
              <div className="flex items-center gap-3">
                <div className="p-2 bg-slate-100 rounded-lg text-slate-700">
                  <Server className="w-5 h-5" />
                </div>
                <div>
                  <h4 className="font-semibold text-slate-900 text-sm">{env.name}</h4>
                  <div className="flex items-center gap-1.5 text-xs text-slate-400 mt-0.5">
                    <Key className="w-3 h-3" />
                    <span>API Key active</span>
                  </div>
                </div>
              </div>

              <button
                onClick={() => handleDelete(env.id, env.name)}
                className="text-slate-400 hover:text-red-600 transition-colors p-1.5 rounded hover:bg-slate-50"
                title="Delete Environment"
              >
                <Trash2 className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>
      )}

      <CreateEnvironmentDialog
        projectId={projectId}
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
      />
    </div>
  );
};
