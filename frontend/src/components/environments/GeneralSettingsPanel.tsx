import { useState, useEffect } from 'react';
import type { FormEvent } from 'react';
import { useUpdateEnvironment, useDeleteEnvironment } from '../../hooks/useEnvironments';
import { Shield, ShieldAlert, Trash2, Save, AlertTriangle } from 'lucide-react';
import toast from 'react-hot-toast';
import type { Environment } from '../../types';

interface GeneralSettingsPanelProps {
  projectId: string;
  environment: Environment;
  onEnvironmentDeleted?: () => void;
}

export const GeneralSettingsPanel = ({
  projectId,
  environment,
  onEnvironmentDeleted,
}: GeneralSettingsPanelProps) => {
  const [name, setName] = useState(environment.name);
  const [isProtected, setIsProtected] = useState(environment.isProtected || false);

  const updateMutation = useUpdateEnvironment(projectId);
  const deleteMutation = useDeleteEnvironment(projectId);

  // Sync state if environment prop changes
  useEffect(() => {
    setName(environment.name);
    setIsProtected(environment.isProtected || false);
  }, [environment]);

  const handleUpdate = async (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    try {
      await updateMutation.mutateAsync({
        envId: environment.id,
        name: name.trim(),
        isProtected,
        sdkSettings: environment.sdkSettings,
      });
      toast.success('Environment general settings updated successfully');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to update environment';
      toast.error(message);
    }
  };

  const handleDelete = async () => {
    if (environment.isProtected) {
      toast.error('Cannot delete a protected environment. Disable protection first.');
      return;
    }

    if (
      confirm(
        `Are you absolutely sure you want to delete "${environment.name}"? This action cannot be undone.`,
      )
    ) {
      try {
        await deleteMutation.mutateAsync(environment.id);
        toast.success(`Environment "${environment.name}" deleted`);
        if (onEnvironmentDeleted) {
          onEnvironmentDeleted();
        }
      } catch (err: unknown) {
        const message = err instanceof Error ? err.message : 'Failed to delete environment';
        toast.error(message);
      }
    }
  };

  const isDirty =
    name.trim() !== environment.name || isProtected !== (environment.isProtected || false);

  return (
    <div className="space-y-8 max-w-3xl">
      <form
        onSubmit={handleUpdate}
        className="bg-white border border-slate-200 rounded-2xl shadow-sm overflow-hidden"
      >
        <div className="p-6 border-b border-slate-100">
          <h3 className="text-base font-semibold text-slate-900">Environment Details</h3>
          <p className="text-sm text-slate-500 mt-1">
            Update basic information and protection status for this environment.
          </p>
        </div>

        <div className="p-6 space-y-6">
          <div className="space-y-2">
            <label htmlFor="envName" className="block text-sm font-medium text-slate-700">
              Environment Name
            </label>
            <input
              id="envName"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-900 focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              placeholder="e.g. Production"
              required
            />
          </div>

          <div className="flex items-start gap-4 p-4 rounded-xl border border-slate-200 bg-slate-50">
            <div className="pt-0.5">
              <button
                type="button"
                role="switch"
                aria-checked={isProtected}
                onClick={() => setIsProtected(!isProtected)}
                className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2 ${
                  isProtected ? 'bg-indigo-600' : 'bg-slate-200'
                }`}
              >
                <span
                  aria-hidden="true"
                  className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                    isProtected ? 'translate-x-5' : 'translate-x-0'
                  }`}
                />
              </button>
            </div>
            <div>
              <div className="flex items-center gap-2">
                <label
                  className="text-sm font-medium text-slate-900 cursor-pointer"
                  onClick={() => setIsProtected(!isProtected)}
                >
                  Protected Environment
                </label>
                {isProtected ? (
                  <Shield className="w-4 h-4 text-indigo-600" />
                ) : (
                  <ShieldAlert className="w-4 h-4 text-slate-400" />
                )}
              </div>
              <p className="text-sm text-slate-500 mt-1">
                When enabled, this environment cannot be deleted. Recommended for critical
                environments like Production.
              </p>
            </div>
          </div>
        </div>

        <div className="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3">
          <button
            type="submit"
            disabled={!isDirty || updateMutation.isPending}
            className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm font-medium px-4 py-2 rounded-xl shadow-sm transition-colors"
          >
            <Save className="w-4 h-4" />
            {updateMutation.isPending ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </form>

      {/* Danger Zone */}
      <div className="border border-red-200 rounded-2xl overflow-hidden shadow-sm">
        <div className="bg-red-50 p-6 border-b border-red-200">
          <div className="flex items-center gap-2 text-red-800">
            <AlertTriangle className="w-5 h-5" />
            <h3 className="text-base font-semibold">Danger Zone</h3>
          </div>
        </div>
        <div className="bg-white p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-6">
          <div className="space-y-1 max-w-lg">
            <h4 className="text-sm font-semibold text-slate-900">Delete Environment</h4>
            <p className="text-sm text-slate-500">
              Permanently delete this environment and all its active flag states, rules, and server
              keys. This action cannot be undone.
            </p>
          </div>
          <button
            type="button"
            onClick={handleDelete}
            disabled={environment.isProtected || deleteMutation.isPending}
            className="shrink-0 inline-flex items-center justify-center gap-2 px-4 py-2 bg-red-50 hover:bg-red-100 text-red-700 border border-red-200 hover:border-red-300 disabled:opacity-50 disabled:cursor-not-allowed rounded-xl text-sm font-semibold transition-colors"
          >
            <Trash2 className="w-4 h-4" />
            Delete Environment
          </button>
        </div>
      </div>
    </div>
  );
};
