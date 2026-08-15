import { useState } from 'react';
import { useEnvironments, useDeleteEnvironment, useCloneEnvironment } from '../../hooks/useEnvironments';
import { CreateEnvironmentDialog } from './CreateEnvironmentDialog';
import { SDKIntegrationModal } from './SDKIntegrationModal';
import { EnvironmentSettingsTabs } from './EnvironmentSettingsTabs';
import { ServerSideKeysPanel } from './ServerSideKeysPanel';
import { GeneralSettingsPanel } from './GeneralSettingsPanel';
import { SdkSettingsPanel } from './SdkSettingsPanel';
import { Plus, Server, Trash2, Loader2, Copy, Check, Code2, ShieldCheck, Key, ArrowLeft } from 'lucide-react';
import toast from 'react-hot-toast';
import type { Environment } from '../../types';

interface EnvironmentsListProps {
  projectId: string;
}

export const EnvironmentsList = ({ projectId }: EnvironmentsListProps) => {
  const { data: environments = [], isLoading, isError, error } = useEnvironments(projectId);
  const deleteMutation = useDeleteEnvironment(projectId);
  const cloneMutation = useCloneEnvironment(projectId);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [copiedKeyId, setCopiedKeyId] = useState<string | null>(null);
  const [selectedEnvForGuide, setSelectedEnvForGuide] = useState<{ name: string; key: string } | null>(null);
  const [activeSettingsEnv, setActiveSettingsEnv] = useState<Environment | null>(null);

  const handleDelete = async (id: string, name: string) => {
    if (confirm(`Are you sure you want to delete environment "${name}"? Active flag configurations for this environment will be removed.`)) {
      await deleteMutation.mutateAsync(id);
    }
  };

  const handleClone = async (envId: string, envName: string) => {
    const cloneName = prompt(`Enter a name for the clone of "${envName}":`, `${envName} (Clone)`);
    if (cloneName) {
      toast.promise(cloneMutation.mutateAsync({ envId, name: cloneName }), {
        loading: 'Cloning environment...',
        success: 'Environment cloned successfully!',
        error: 'Failed to clone environment',
      });
    }
  };

  const handleCopyKey = (envId: string, apiKey?: string, name?: string) => {
    const keyToCopy = apiKey || `env_${name?.toLowerCase().replace(/[^a-z0-9]/g, '_')}_token`;
    navigator.clipboard.writeText(keyToCopy);
    setCopiedKeyId(envId);
    toast.success('Environment Client SDK Key copied');
    setTimeout(() => setCopiedKeyId(null), 2000);
  };

  if (activeSettingsEnv) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button
              onClick={() => setActiveSettingsEnv(null)}
              className="p-2 rounded-xl text-slate-400 hover:text-slate-700 hover:bg-slate-100 transition-colors"
              title="Back to Environments"
            >
              <ArrowLeft className="w-5 h-5" />
            </button>
            <div>
              <h2 className="text-xl font-bold text-slate-900">{activeSettingsEnv.name} Settings</h2>
              <p className="text-xs text-slate-500">Manage client-side & server-side SDK credentials and configurations.</p>
            </div>
          </div>
        </div>

        <EnvironmentSettingsTabs
          envName={activeSettingsEnv.name}
          apiKey={activeSettingsEnv.apiKey}
          serverKeysPanel={
            <ServerSideKeysPanel projectId={projectId} envId={activeSettingsEnv.id} />
          }
          generalSettingsContent={
            <GeneralSettingsPanel 
              projectId={projectId} 
              environment={activeSettingsEnv} 
              onEnvironmentDeleted={() => setActiveSettingsEnv(null)} 
            />
          }
          sdkSettingsContent={
            <SdkSettingsPanel 
              projectId={projectId} 
              environment={activeSettingsEnv} 
            />
          }
        />
      </div>
    );
  }

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
        <div className="bg-white border border-slate-200 rounded-xl p-10 text-center space-y-4">
          <div className="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-full flex items-center justify-center mx-auto">
            <Server className="w-6 h-6" />
          </div>
          <div className="space-y-1">
            <h3 className="text-sm font-medium text-slate-900">No environments</h3>
            <p className="text-sm text-slate-500">Get started by creating your first environment.</p>
          </div>
          <div className="pt-2">
            <button
              onClick={() => setIsCreateOpen(true)}
              className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium px-4 py-2 rounded-lg shadow-sm transition-colors"
            >
              <Plus className="w-4 h-4" />
              <span>Create your first Environment</span>
            </button>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          {environments.map((env) => {
            const displayKey = env.apiKey || `env_${env.name.toLowerCase().replace(/[^a-z0-9]/g, '_')}_token`;
            return (
              <div
                key={env.id}
                className="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm hover:shadow-md hover:border-slate-300 transition-all flex flex-col justify-between space-y-4 group"
              >
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div className="p-2.5 bg-indigo-50 text-indigo-600 rounded-xl border border-indigo-100">
                      <Server className="w-5 h-5" />
                    </div>
                    <div>
                      <h4 className="font-bold text-slate-900 text-base leading-tight">{env.name}</h4>
                      <span className="inline-flex items-center gap-1 mt-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                        <ShieldCheck className="w-3 h-3" /> Public Client SDK Key
                      </span>
                    </div>
                  </div>

                  <div className="flex items-center gap-1">
                    <button
                      onClick={() => handleClone(env.id, env.name)}
                      className="text-slate-400 hover:text-indigo-600 transition-colors p-1.5 rounded-lg hover:bg-slate-100"
                      title="Clone Environment"
                    >
                      <Copy className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(env.id, env.name)}
                      className="text-slate-400 hover:text-red-600 transition-colors p-1.5 rounded-lg hover:bg-slate-100"
                      title="Delete Environment"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>

                {/* Permanent Client Key Box */}
                <div className="bg-slate-50 border border-slate-200 rounded-xl p-3 space-y-1.5">
                  <div className="flex items-center justify-between text-[11px] font-semibold uppercase tracking-wider text-slate-400">
                    <span>SDK Client Key</span>
                    <button
                      onClick={() => handleCopyKey(env.id, env.apiKey, env.name)}
                      className="text-indigo-600 hover:text-indigo-800 flex items-center gap-1 text-[11px] font-semibold transition-colors"
                      title="Copy Key"
                    >
                      {copiedKeyId === env.id ? (
                        <>
                          <Check className="w-3 h-3 text-emerald-600" />
                          <span className="text-emerald-600">Copied!</span>
                        </>
                      ) : (
                        <>
                          <Copy className="w-3 h-3" />
                          <span>Copy Key</span>
                        </>
                      )}
                    </button>
                  </div>
                  <div className="font-mono text-xs text-slate-700 bg-white border border-slate-200 rounded-lg px-2.5 py-1.5 truncate select-all">
                    {displayKey}
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="grid grid-cols-2 gap-2">
                  <button
                    onClick={() => setActiveSettingsEnv(env)}
                    className="inline-flex items-center justify-center gap-1.5 bg-slate-900 hover:bg-slate-800 text-white text-xs font-semibold py-2 px-3 rounded-xl transition-colors"
                  >
                    <Key className="w-3.5 h-3.5" />
                    <span>Keys & Settings</span>
                  </button>
                  <button
                    onClick={() => setSelectedEnvForGuide({ name: env.name, key: displayKey })}
                    className="inline-flex items-center justify-center gap-1.5 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 text-xs font-semibold py-2 px-3 rounded-xl border border-indigo-200/60 transition-colors"
                  >
                    <Code2 className="w-3.5 h-3.5 text-indigo-600" />
                    <span>SDK Guide</span>
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      <CreateEnvironmentDialog
        projectId={projectId}
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
      />

      {selectedEnvForGuide && (
        <SDKIntegrationModal
          isOpen={true}
          onClose={() => setSelectedEnvForGuide(null)}
          envName={selectedEnvForGuide.name}
          envKey={selectedEnvForGuide.key}
        />
      )}
    </div>
  );
};
