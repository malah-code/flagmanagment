import React, { useState, useEffect } from 'react';
import { killSwitchApi, type KillSwitchRule } from '../services/killSwitchApi';

interface Props {
  envId: string;
  flagId: string;
}

export const KillSwitchForm: React.FC<Props> = ({ envId, flagId }) => {
  const [rules, setRules] = useState<KillSwitchRule[]>([]);
  const [alertIdentifier, setAlertIdentifier] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchRules();
  }, [envId, flagId]);

  const fetchRules = async () => {
    try {
      const data = await killSwitchApi.list(envId, flagId);
      setRules(data || []);
    } catch (err) {
      console.error(err);
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!alertIdentifier) return;
    setLoading(true);
    setError(null);
    try {
      await killSwitchApi.create(envId, flagId, alertIdentifier);
      setAlertIdentifier('');
      fetchRules();
    } catch (err: any) {
      setError(err.message || 'Failed to create kill switch');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to remove this kill switch?')) return;
    try {
      await killSwitchApi.delete(envId, flagId, id);
      fetchRules();
    } catch (err) {
      console.error('Failed to delete kill switch');
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-200 mt-6">
      <h3 className="text-lg font-medium mb-4">Automated Kill Switches</h3>
      <p className="text-sm text-gray-500 mb-4">
        Automatically disable this flag when a specific APM alert webhook is received.
      </p>

      {error && <div className="text-red-500 text-sm mb-4">{error}</div>}

      <form onSubmit={handleCreate} className="flex gap-2 mb-6">
        <input
          type="text"
          value={alertIdentifier}
          onChange={(e) => setAlertIdentifier(e.target.value)}
          placeholder="e.g. high_error_rate_payment_service"
          className="flex-1 border rounded px-3 py-2 text-sm"
        />
        <button
          type="submit"
          disabled={loading || !alertIdentifier}
          className="bg-red-600 text-white px-4 py-2 rounded text-sm disabled:opacity-50 hover:bg-red-700"
        >
          {loading ? 'Adding...' : 'Add Kill Switch'}
        </button>
      </form>

      {rules.length > 0 ? (
        <div className="border rounded divide-y">
          {rules.map((rule) => (
            <div key={rule.id} className="p-3 flex justify-between items-center text-sm">
              <div>
                <span className="font-mono bg-gray-100 px-2 py-1 rounded text-red-600">
                  {rule.alert_identifier}
                </span>
                <span className="ml-3 text-gray-500">Action: DISABLE</span>
              </div>
              <button
                onClick={() => handleDelete(rule.id)}
                className="text-gray-400 hover:text-red-600"
              >
                Remove
              </button>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-sm text-gray-500 text-center py-4 border rounded bg-gray-50">
          No kill switches configured.
        </div>
      )}
    </div>
  );
};
