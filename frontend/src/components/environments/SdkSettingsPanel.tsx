import { useState, useEffect } from 'react';
import type { FormEvent } from 'react';
import { useUpdateEnvironment } from '../../hooks/useEnvironments';
import { Save, Radio, Activity, BarChart3 } from 'lucide-react';
import toast from 'react-hot-toast';
import type { Environment } from '../../types';

interface SdkSettingsPanelProps {
  projectId: string;
  environment: Environment;
}

export const SdkSettingsPanel = ({ projectId, environment }: SdkSettingsPanelProps) => {
  const updateMutation = useUpdateEnvironment(projectId);

  const [pollingInterval, setPollingInterval] = useState<number>(30);
  const [enableStreaming, setEnableStreaming] = useState<boolean>(true);
  const [enableAnalytics, setEnableAnalytics] = useState<boolean>(false);

  useEffect(() => {
    if (environment.sdkSettings) {
      setPollingInterval(environment.sdkSettings.pollingInterval ?? 30);
      setEnableStreaming(environment.sdkSettings.enableStreaming ?? true);
      setEnableAnalytics(environment.sdkSettings.enableAnalytics ?? false);
    }
  }, [environment]);

  const handleUpdate = async (e: FormEvent) => {
    e.preventDefault();

    try {
      const updatedSdkSettings = {
        ...environment.sdkSettings,
        pollingInterval,
        enableStreaming,
        enableAnalytics,
      };

      await updateMutation.mutateAsync({
        envId: environment.id,
        name: environment.name,
        isProtected: environment.isProtected || false,
        sdkSettings: updatedSdkSettings,
      });
      toast.success('SDK settings updated successfully');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to update SDK settings';
      toast.error(message);
    }
  };

  const currentSettings = environment.sdkSettings || {};
  const isDirty = 
    pollingInterval !== (currentSettings.pollingInterval ?? 30) ||
    enableStreaming !== (currentSettings.enableStreaming ?? true) ||
    enableAnalytics !== (currentSettings.enableAnalytics ?? false);

  return (
    <form onSubmit={handleUpdate} className="bg-white border border-slate-200 rounded-2xl shadow-sm overflow-hidden max-w-3xl">
      <div className="p-6 border-b border-slate-100">
        <h3 className="text-base font-semibold text-slate-900">SDK Configurations</h3>
        <p className="text-sm text-slate-500 mt-1">
          Control how client and server SDKs connect, fetch updates, and report analytics for this environment.
        </p>
      </div>

      <div className="p-6 space-y-8">
        {/* Polling Interval */}
        <div className="space-y-3">
          <div className="flex items-center gap-2">
            <Radio className="w-5 h-5 text-indigo-600" />
            <h4 className="text-sm font-semibold text-slate-900">Event Polling Interval</h4>
          </div>
          <p className="text-xs text-slate-500">
            How often older SDKs should poll the server for feature flag updates (when streaming is disabled or unavailable).
          </p>
          <select
            value={pollingInterval}
            onChange={(e) => setPollingInterval(Number(e.target.value))}
            className="w-full max-w-xs px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-900 focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-colors"
          >
            <option value={15}>15 seconds (High Load)</option>
            <option value={30}>30 seconds (Standard)</option>
            <option value={60}>60 seconds (1 minute)</option>
            <option value={300}>300 seconds (5 minutes)</option>
            <option value={600}>600 seconds (10 minutes)</option>
          </select>
        </div>

        {/* Streaming */}
        <div className="flex items-start gap-4 p-4 rounded-xl border border-slate-200 bg-slate-50">
          <div className="pt-0.5">
            <button
              type="button"
              role="switch"
              aria-checked={enableStreaming}
              onClick={() => setEnableStreaming(!enableStreaming)}
              className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2 ${
                enableStreaming ? 'bg-indigo-600' : 'bg-slate-200'
              }`}
            >
              <span
                aria-hidden="true"
                className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                  enableStreaming ? 'translate-x-5' : 'translate-x-0'
                }`}
              />
            </button>
          </div>
          <div>
            <div className="flex items-center gap-2">
              <label className="text-sm font-medium text-slate-900 cursor-pointer" onClick={() => setEnableStreaming(!enableStreaming)}>
                Real-time Updates (SSE / Streaming)
              </label>
              <Activity className="w-4 h-4 text-emerald-600" />
            </div>
            <p className="text-sm text-slate-500 mt-1">
              Enable Server-Sent Events to push instantaneous feature flag updates to connected clients. Disabling this forces all clients to fall back to standard HTTP polling.
            </p>
          </div>
        </div>

        {/* Analytics */}
        <div className="flex items-start gap-4 p-4 rounded-xl border border-slate-200 bg-slate-50">
          <div className="pt-0.5">
            <button
              type="button"
              role="switch"
              aria-checked={enableAnalytics}
              onClick={() => setEnableAnalytics(!enableAnalytics)}
              className={`relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2 ${
                enableAnalytics ? 'bg-indigo-600' : 'bg-slate-200'
              }`}
            >
              <span
                aria-hidden="true"
                className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                  enableAnalytics ? 'translate-x-5' : 'translate-x-0'
                }`}
              />
            </button>
          </div>
          <div>
            <div className="flex items-center gap-2">
              <label className="text-sm font-medium text-slate-900 cursor-pointer" onClick={() => setEnableAnalytics(!enableAnalytics)}>
                SDK Analytics Collection
              </label>
              <BarChart3 className="w-4 h-4 text-indigo-500" />
            </div>
            <p className="text-sm text-slate-500 mt-1">
              Allow SDKs to periodically report diagnostic evaluation data back to the server. Disabling this reduces server load but prevents usage metrics from being captured.
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
          {updateMutation.isPending ? 'Saving...' : 'Save Settings'}
        </button>
      </div>
    </form>
  );
};
