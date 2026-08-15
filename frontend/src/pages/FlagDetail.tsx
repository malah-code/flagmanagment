import { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useFlags, useUpdateFlag } from '../hooks/useFlags';
import { useEnvironments } from '../hooks/useEnvironments';
import { useFlagStates, useUpdateFlagState, useInitFlagState } from '../hooks/useFlagStates';
import { ArrowLeft, Sliders, Settings, LayoutTemplate, Save, Loader2, CheckCircle2, XCircle } from 'lucide-react';
import toast from 'react-hot-toast';

type TabType = 'targeting' | 'variations' | 'settings';

export const FlagDetail = () => {
  const { projectId = '', flagId = '' } = useParams<{ projectId: string; flagId: string }>();
  const { data: flags = [], isLoading: isLoadingFlags } = useFlags(projectId);
  const { data: environments = [] } = useEnvironments(projectId);
  
  const [activeTab, setActiveTab] = useState<TabType>('targeting');
  const [selectedEnvId, setSelectedEnvId] = useState<string>('');

  const flag = flags.find(f => f.id === flagId);
  const currentEnvId = selectedEnvId || (environments[0]?.id ?? '');
  
  const { data: flagStates = [], isLoading: isLoadingStates } = useFlagStates(projectId, currentEnvId);
  const updateMutation = useUpdateFlagState(projectId, currentEnvId);
  const initMutation = useInitFlagState(projectId, currentEnvId);
  const updateFlagMutation = useUpdateFlag(projectId);

  const flagState = flagStates.find(s => s.flagId === flagId);

  const [settingsName, setSettingsName] = useState('');
  const [settingsDescription, setSettingsDescription] = useState('');
  const [settingsTags, setSettingsTags] = useState('');

  useEffect(() => {
    if (flag) {
      setSettingsName(flag.name || flag.key);
      setSettingsDescription(flag.description || '');
      setSettingsTags(flag.tags?.join(', ') || '');
    }
  }, [flag]);

  const handleUpdateSettings = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!flag) return;
    
    try {
      await updateFlagMutation.mutateAsync({
        flagId: flag.id,
        payload: {
          name: settingsName.trim(),
          description: settingsDescription.trim(),
          tags: settingsTags.split(',').map(t => t.trim()).filter(Boolean),
        }
      });
      toast.success('Flag settings updated');
    } catch (err) {
      toast.error('Failed to update flag settings');
    }
  };

  if (isLoadingFlags) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
      </div>
    );
  }

  if (!flag) {
    return (
      <div className="space-y-4">
        <Link to={`/projects/${projectId}`} className="inline-flex items-center gap-1 text-sm text-indigo-600 hover:underline">
          <ArrowLeft className="w-4 h-4" /> Back to Project
        </Link>
        <div className="bg-red-50 text-red-600 p-4 rounded-xl border border-red-200">
          Flag not found.
        </div>
      </div>
    );
  }

  const handleToggle = async () => {
    if (!flagState) return;
    try {
      await updateMutation.mutateAsync({
        flagId: flagState.flagId,
        payload: { 
          isEnabled: !flagState.isEnabled,
          targetingRules: flagState.targetingRules || { rules: [] },
          remoteConfig: flagState.remoteConfig || {},
          rolloutRules: flagState.rolloutRules ? { rules: flagState.rolloutRules } : undefined,
          defaultVariation: flagState.defaultVariation,
        },
      });
      toast.success('Flag state updated');
    } catch (err: any) {
      toast.error('Failed to update flag state');
    }
  };

  return (
    <div className="space-y-6">
      {/* Breadcrumb & Header */}
      <div>
        <Link
          to={`/projects/${projectId}`}
          className="inline-flex items-center gap-1 text-sm font-medium text-slate-500 hover:text-slate-900 transition-colors mb-3"
        >
          <ArrowLeft className="w-4 h-4" /> Back to Project
        </Link>

        <div className="flex items-center justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold text-slate-900 font-mono">{flag.key}</h1>
            <p className="text-sm text-slate-500 mt-1">{flag.description || 'No description provided.'}</p>
          </div>
          <div className="flex items-center gap-2 bg-slate-100 text-slate-600 text-xs font-semibold px-3 py-1.5 rounded-lg border border-slate-200">
            Type: {flag.type}
          </div>
        </div>
      </div>

      {/* Tabs & Env Selector */}
      <div className="border-b border-slate-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-0">
        <nav className="flex gap-6 -mb-px">
          <button
            onClick={() => setActiveTab('targeting')}
            className={`inline-flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-medium transition-colors ${
              activeTab === 'targeting'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <Sliders className="w-4 h-4" />
            <span>Targeting</span>
          </button>

          <button
            onClick={() => setActiveTab('variations')}
            className={`inline-flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-medium transition-colors ${
              activeTab === 'variations'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <LayoutTemplate className="w-4 h-4" />
            <span>Variations</span>
          </button>

          <button
            onClick={() => setActiveTab('settings')}
            className={`inline-flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-medium transition-colors ${
              activeTab === 'settings'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <Settings className="w-4 h-4" />
            <span>Settings</span>
          </button>
        </nav>

        {activeTab === 'targeting' && environments.length > 0 && (
          <div className="flex items-center gap-2 pb-2">
            <span className="text-xs font-semibold text-slate-500 uppercase">Environment:</span>
            <select
              value={currentEnvId}
              onChange={(e) => setSelectedEnvId(e.target.value)}
              className="text-sm border border-slate-300 rounded-lg px-3 py-1.5 bg-white text-slate-900 font-medium focus:ring-2 focus:ring-indigo-500 outline-none"
            >
              {environments.map((e) => (
                <option key={e.id} value={e.id}>
                  {e.name}
                </option>
              ))}
            </select>
          </div>
        )}
      </div>

      {/* Content */}
      <div className="pt-2">
        {activeTab === 'targeting' && (
          <div className="space-y-6">
            {!currentEnvId ? (
              <div className="bg-white border border-slate-200 rounded-xl p-8 text-center text-sm text-slate-500">
                Please create an environment first.
              </div>
            ) : isLoadingStates ? (
              <div className="flex justify-center p-8">
                <Loader2 className="w-6 h-6 animate-spin text-indigo-600" />
              </div>
            ) : !flagState ? (
              <div className="bg-amber-50 text-amber-700 p-4 rounded-xl border border-amber-200 text-sm flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                <span>This flag has not been initialized in the selected environment. Initialize it to start configuring targeting rules and serving variations.</span>
                <button
                  onClick={() => {
                    initMutation.mutate({
                      flagId, 
                      payload: { isEnabled: false, targetingRules: { rules: [] }, remoteConfig: {} }
                    });
                  }}
                  disabled={initMutation.isPending}
                  className="shrink-0 bg-white hover:bg-amber-100 text-amber-700 font-medium px-4 py-2 rounded-lg border border-amber-300 transition-colors shadow-sm disabled:opacity-50 flex items-center gap-2"
                >
                  {initMutation.isPending && <Loader2 className="w-4 h-4 animate-spin" />}
                  Initialize Flag
                </button>
              </div>
            ) : (
              <div className="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm">
                <div className="p-6 border-b border-slate-200 flex items-center justify-between bg-slate-50">
                  <div>
                    <h3 className="text-base font-semibold text-slate-900">Environment Status</h3>
                    <p className="text-sm text-slate-500">Enable or disable this flag in {environments.find(e => e.id === currentEnvId)?.name}.</p>
                  </div>
                  <button
                    onClick={handleToggle}
                    disabled={updateMutation.isPending}
                    className={`inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors disabled:opacity-50 ${
                      flagState.isEnabled
                        ? 'bg-emerald-50 text-emerald-700 border border-emerald-200 hover:bg-emerald-100'
                        : 'bg-slate-100 text-slate-600 border border-slate-200 hover:bg-slate-200'
                    }`}
                  >
                    {updateMutation.isPending ? (
                      <Loader2 className="w-4 h-4 animate-spin" />
                    ) : flagState.isEnabled ? (
                      <>
                        <CheckCircle2 className="w-4 h-4" /> Enabled
                      </>
                    ) : (
                      <>
                        <XCircle className="w-4 h-4" /> Disabled
                      </>
                    )}
                  </button>
                </div>
                <div className="p-6">
                  <h3 className="text-sm font-semibold text-slate-900 mb-4">Targeting Rules</h3>
                  <div className="bg-slate-50 rounded-lg border border-slate-200 p-8 text-center">
                    <p className="text-sm text-slate-500 mb-4">
                      Targeting rules allow you to serve different variations to specific segments of users.
                    </p>
                    <button className="inline-flex items-center gap-2 bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 text-sm font-medium px-4 py-2 rounded-lg shadow-sm transition-colors">
                      <Sliders className="w-4 h-4" /> Manage Rules in Flags List
                    </button>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'variations' && (
          <div className="bg-white border border-slate-200 rounded-xl shadow-sm p-6">
            <h3 className="text-base font-semibold text-slate-900 mb-4">Flag Variations</h3>
            {flag.variations && flag.variations.length > 0 ? (
              <div className="space-y-4">
                {flag.variations.map((v) => (
                  <div key={v.id} className="flex items-center justify-between p-4 bg-slate-50 rounded-lg border border-slate-200">
                    <div>
                      <div className="font-semibold text-slate-900 text-sm">{v.name}</div>
                      <div className="text-sm text-slate-500 font-mono mt-1">{String(v.value)}</div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-8">
                <p className="text-sm text-slate-500">This flag is a simple boolean flag and does not have additional variations.</p>
              </div>
            )}
          </div>
        )}

        {activeTab === 'settings' && (
          <form onSubmit={handleUpdateSettings} className="bg-white border border-slate-200 rounded-xl shadow-sm p-6 max-w-2xl">
            <h3 className="text-base font-semibold text-slate-900 mb-6">Flag Settings</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Flag Key</label>
                <input
                  type="text"
                  disabled
                  value={flag.key}
                  className="w-full text-sm border border-slate-300 rounded-lg px-3 py-2 bg-slate-50 text-slate-500 font-mono"
                />
                <p className="text-xs text-slate-500 mt-1">The key cannot be changed after creation.</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Name</label>
                <input
                  type="text"
                  value={settingsName}
                  onChange={(e) => setSettingsName(e.target.value)}
                  className="w-full text-sm border border-slate-300 rounded-lg px-3 py-2 bg-white text-slate-900 focus:ring-2 focus:ring-indigo-500 outline-none"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Description</label>
                <textarea
                  value={settingsDescription}
                  onChange={(e) => setSettingsDescription(e.target.value)}
                  rows={3}
                  className="w-full text-sm border border-slate-300 rounded-lg px-3 py-2 bg-white text-slate-900 focus:ring-2 focus:ring-indigo-500 outline-none"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Tags (comma separated)</label>
                <input
                  type="text"
                  value={settingsTags}
                  onChange={(e) => setSettingsTags(e.target.value)}
                  className="w-full text-sm border border-slate-300 rounded-lg px-3 py-2 bg-white text-slate-900 focus:ring-2 focus:ring-indigo-500 outline-none"
                />
              </div>

              <div className="pt-4 border-t border-slate-100 flex justify-end">
                <button
                  type="submit"
                  disabled={updateFlagMutation.isPending}
                  className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white text-sm font-medium px-4 py-2 rounded-lg shadow-sm transition-colors"
                >
                  {updateFlagMutation.isPending ? <Loader2 className="w-4 h-4 animate-spin" /> : <Save className="w-4 h-4" />}
                  Save Settings
                </button>
              </div>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};
