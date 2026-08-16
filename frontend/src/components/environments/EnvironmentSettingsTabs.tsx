import { useState } from 'react';
import type { ReactNode } from 'react';
import { Key, Settings, Sliders } from 'lucide-react';
import { ClientSideKeyCard } from './ClientSideKeyCard';

interface EnvironmentSettingsTabsProps {
  envName: string;
  apiKey?: string;
  serverKeysPanel?: ReactNode;
  generalSettingsContent?: ReactNode;
  sdkSettingsContent?: ReactNode;
}

export const EnvironmentSettingsTabs = ({
  envName,
  apiKey,
  serverKeysPanel,
  generalSettingsContent,
  sdkSettingsContent,
}: EnvironmentSettingsTabsProps) => {
  const [activeTab, setActiveTab] = useState<'general' | 'keys' | 'sdk'>('keys');

  return (
    <div className="space-y-6">
      {/* Sub-Tabs Header */}
      <div className="border-b border-slate-200">
        <nav className="-mb-px flex gap-6" aria-label="Environment Settings Navigation">
          <button
            onClick={() => setActiveTab('general')}
            className={`flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-semibold transition-colors ${
              activeTab === 'general'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <Settings className="w-4 h-4" />
            <span>General</span>
          </button>

          <button
            onClick={() => setActiveTab('keys')}
            className={`flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-semibold transition-colors ${
              activeTab === 'keys'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <Key className="w-4 h-4" />
            <span>Keys</span>
          </button>

          <button
            onClick={() => setActiveTab('sdk')}
            className={`flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-semibold transition-colors ${
              activeTab === 'sdk'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <Sliders className="w-4 h-4" />
            <span>SDK Settings</span>
          </button>
        </nav>
      </div>

      {/* Tab Panels */}
      {activeTab === 'keys' && (
        <div className="space-y-6">
          <ClientSideKeyCard envName={envName} apiKey={apiKey} />
          {serverKeysPanel}
        </div>
      )}

      {activeTab === 'general' && (
        <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm">
          {generalSettingsContent || (
            <div className="text-sm text-slate-500">
              General settings for environment <strong className="text-slate-800">{envName}</strong>
              .
            </div>
          )}
        </div>
      )}

      {activeTab === 'sdk' && (
        <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm">
          {sdkSettingsContent || (
            <div className="text-sm text-slate-500">
              SDK configurations and stream settings for environment{' '}
              <strong className="text-slate-800">{envName}</strong>.
            </div>
          )}
        </div>
      )}
    </div>
  );
};
