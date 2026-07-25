import React, { useEffect, useState } from 'react';
import { slackApi, type SlackWebhookConfig } from '../services/slackApi';

interface SlackConfigFormProps {
  environmentId: string;
}

export const SlackConfigForm: React.FC<SlackConfigFormProps> = ({ environmentId }) => {
  const [config, setConfig] = useState<SlackWebhookConfig>({
    environment_id: environmentId,
    webhook_url: '',
    enabled: false,
  });
  const [loading, setLoading] = useState<boolean>(true);
  const [saving, setSaving] = useState<boolean>(false);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    fetchConfig();
  }, [environmentId]);

  const fetchConfig = async () => {
    setLoading(true);
    try {
      const data = await slackApi.getSlackConfig(environmentId);
      setConfig(data);
    } catch (err) {
      console.error('Failed to load Slack config', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setMessage(null);
    try {
      const updated = await slackApi.saveSlackConfig(environmentId, config.webhook_url, config.enabled);
      setConfig(updated);
      setMessage('Slack notification settings saved successfully!');
    } catch (err: any) {
      setMessage('Failed to save settings: ' + (err.response?.data?.error || err.message));
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="text-gray-500 text-sm">Loading Slack Settings...</div>;
  }

  return (
    <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 mt-4">
      <h3 className="text-md font-semibold text-gray-800 mb-2 flex items-[#4A154B] items-center gap-2">
        <span>💬</span> Slack Notification Webhook
      </h3>
      <form onSubmit={handleSave} className="space-y-3">
        <div>
          <label className="block text-xs font-medium text-gray-700 mb-1">
            Slack Incoming Webhook URL
          </label>
          <input
            type="url"
            required
            placeholder="https://hooks.slack.com/services/..."
            value={config.webhook_url}
            onChange={(e) => setConfig({ ...config, webhook_url: e.target.value })}
            className="w-full px-3 py-1.5 border border-gray-300 rounded text-sm focus:ring-2 focus:ring-purple-500 focus:outline-none"
          />
        </div>
        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            id="slack-enabled"
            checked={config.enabled}
            onChange={(e) => setConfig({ ...config, enabled: e.target.checked })}
            className="rounded text-purple-600 focus:ring-purple-500 h-4 w-4"
          />
          <label htmlFor="slack-enabled" className="text-xs font-medium text-gray-700">
            Enable Slack notifications for flag updates in this environment
          </label>
        </div>

        {message && (
          <div className={`text-xs p-2 rounded ${message.includes('Failed') ? 'bg-red-50 text-red-700' : 'bg-green-50 text-green-700'}`}>
            {message}
          </div>
        )}

        <button
          type="submit"
          disabled={saving}
          className="px-3 py-1.5 bg-[#4A154B] hover:bg-[#3B113C] text-white text-xs font-medium rounded shadow-sm transition-colors"
        >
          {saving ? 'Saving...' : 'Save Slack Settings'}
        </button>
      </form>
    </div>
  );
};
