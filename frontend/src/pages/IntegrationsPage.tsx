import React, { useState } from 'react';
import { Settings, Plus, Activity, ExternalLink, Globe, Key, AlertCircle } from 'lucide-react';
import { useWebhooks, useCreateWebhook } from '../hooks/useWebhooks';

const INTEGRATION_TYPES = [
  {
    id: 'posthog',
    name: 'PostHog',
    description: 'Send flag evaluation data to PostHog for product analytics and A/B testing.',
    icon: <Activity className="w-8 h-8 text-orange-500" />,
    fields: [
      { id: 'url', label: 'PostHog Host', placeholder: 'https://app.posthog.com' },
      { id: 'secret_key', label: 'Project API Key', placeholder: 'phc_...' },
    ],
  },
  {
    id: 'datadog',
    name: 'Datadog',
    description: 'Track feature flag toggles and operational metrics in Datadog.',
    icon: <Globe className="w-8 h-8 text-purple-500" />,
    fields: [
      {
        id: 'url',
        label: 'Datadog Site (Intake URL)',
        placeholder: 'https://http-intake.logs.datadoghq.com',
      },
      { id: 'secret_key', label: 'API Key', placeholder: 'dd_...' },
    ],
  },
  {
    id: 'custom',
    name: 'Custom Webhook',
    description: 'Receive HTTP POST requests for any flag state change or audit event.',
    icon: <Settings className="w-8 h-8 text-slate-500" />,
    fields: [
      { id: 'url', label: 'Webhook URL', placeholder: 'https://api.example.com/webhooks' },
      { id: 'secret_key', label: 'Secret (Optional HMAC)', placeholder: 'Optional signing secret' },
    ],
  },
];

export const IntegrationsPage = ({ projectId }: { projectId: string }) => {
  const { data: integrations = [], isLoading } = useWebhooks(projectId);
  const createWebhook = useCreateWebhook(projectId);

  const [selectedType, setSelectedType] = useState<string | null>(null);
  const [formData, setFormData] = useState({ name: '', url: '', secret_key: '' });

  const activeIntegrationDef = INTEGRATION_TYPES.find((t) => t.id === selectedType);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeIntegrationDef) return;

    try {
      await createWebhook.mutateAsync({
        name: formData.name || activeIntegrationDef.name,
        integration_type: activeIntegrationDef.id,
        url: formData.url,
        secret_key: formData.secret_key || undefined,
        events: ['flag.updated', 'flag.created'], // default events
        is_active: true,
      });
      setSelectedType(null);
      setFormData({ name: '', url: '', secret_key: '' });
    } catch (err) {
      console.error(err);
    }
  };

  if (isLoading) {
    return <div className="p-8">Loading integrations...</div>;
  }

  return (
    <div className="max-w-6xl mx-auto p-8 space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Integrations</h1>
        <p className="text-slate-500 mt-2">
          Connect FlagManagment to your existing observability and product analytics tools.
        </p>
      </div>

      {/* Catalog */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {INTEGRATION_TYPES.map((type) => (
          <div
            key={type.id}
            className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm hover:shadow-md transition-shadow"
          >
            <div className="flex items-start justify-between mb-4">
              <div className="p-3 bg-slate-50 rounded-xl border border-slate-100">{type.icon}</div>
              <button
                onClick={() => setSelectedType(type.id)}
                className="flex items-center gap-1.5 text-sm font-semibold text-indigo-600 bg-indigo-50 px-3 py-1.5 rounded-lg hover:bg-indigo-100 transition-colors"
              >
                <Plus className="w-4 h-4" /> Connect
              </button>
            </div>
            <h3 className="font-bold text-slate-900 mb-2">{type.name}</h3>
            <p className="text-sm text-slate-500 leading-relaxed">{type.description}</p>
          </div>
        ))}
      </div>

      {/* Configured Integrations */}
      <div className="pt-8">
        <h2 className="text-xl font-bold text-slate-900 mb-6">Active Integrations</h2>
        {integrations.length === 0 ? (
          <div className="text-center py-12 bg-slate-50 rounded-2xl border border-slate-200 border-dashed">
            <Globe className="w-12 h-12 text-slate-300 mx-auto mb-3" />
            <p className="text-slate-500">No active integrations found.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {integrations.map((inv) => {
              const def =
                INTEGRATION_TYPES.find((t) => t.id === inv.integration_type) ||
                INTEGRATION_TYPES[2];
              return (
                <div
                  key={inv.id}
                  className="bg-white border border-slate-200 rounded-xl p-5 flex items-center justify-between shadow-sm"
                >
                  <div className="flex items-center gap-4">
                    <div className="p-2 bg-slate-50 rounded-lg border border-slate-100">
                      {def.icon}
                    </div>
                    <div>
                      <h4 className="font-bold text-slate-900 flex items-center gap-2">
                        {inv.name}
                        <span
                          className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${inv.is_active ? 'bg-emerald-100 text-emerald-700' : 'bg-slate-100 text-slate-600'}`}
                        >
                          {inv.is_active ? 'Active' : 'Disabled'}
                        </span>
                      </h4>
                      <div className="flex items-center gap-3 mt-1 text-sm text-slate-500">
                        <span className="flex items-center gap-1">
                          <ExternalLink className="w-3.5 h-3.5" />
                          {inv.url}
                        </span>
                      </div>
                    </div>
                  </div>
                  <button className="text-sm font-semibold text-slate-500 hover:text-slate-900">
                    Configure
                  </button>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Configure Modal */}
      {selectedType && activeIntegrationDef && (
        <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-lg overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between bg-slate-50">
              <div className="flex items-center gap-3">
                {activeIntegrationDef.icon}
                <h2 className="font-bold text-slate-900">Configure {activeIntegrationDef.name}</h2>
              </div>
              <button
                onClick={() => setSelectedType(null)}
                className="text-slate-400 hover:text-slate-600"
              >
                ×
              </button>
            </div>

            <form onSubmit={handleSubmit} className="p-6 space-y-5">
              <div>
                <label className="block text-sm font-semibold text-slate-700 mb-1.5">
                  Connection Name
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder={`My ${activeIntegrationDef.name} Integration`}
                  className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
                />
              </div>

              {activeIntegrationDef.fields.map((field) => (
                <div key={field.id}>
                  <label className="block text-sm font-semibold text-slate-700 mb-1.5 flex items-center gap-2">
                    {field.id === 'secret_key' && <Key className="w-3.5 h-3.5 text-slate-400" />}
                    {field.label}
                  </label>
                  <input
                    type={field.id === 'secret_key' ? 'password' : 'text'}
                    required={field.id === 'url'}
                    value={formData[field.id as keyof typeof formData] || ''}
                    onChange={(e) => setFormData({ ...formData, [field.id]: e.target.value })}
                    placeholder={field.placeholder}
                    className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 text-sm outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
                  />
                </div>
              ))}

              <div className="bg-amber-50 text-amber-800 p-4 rounded-xl border border-amber-200/50 flex gap-3 text-sm">
                <AlertCircle className="w-5 h-5 shrink-0 text-amber-500" />
                <p>
                  Events will be synced securely. Make sure your API keys have the necessary
                  permissions in {activeIntegrationDef.name}.
                </p>
              </div>

              <div className="flex items-center gap-3 pt-4 border-t border-slate-100">
                <button
                  type="button"
                  onClick={() => setSelectedType(null)}
                  className="flex-1 px-4 py-2.5 text-sm font-semibold text-slate-600 bg-slate-100 hover:bg-slate-200 rounded-xl transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={createWebhook.isPending}
                  className="flex-1 px-4 py-2.5 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 rounded-xl transition-colors disabled:opacity-50"
                >
                  {createWebhook.isPending ? 'Saving...' : 'Connect Integration'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
