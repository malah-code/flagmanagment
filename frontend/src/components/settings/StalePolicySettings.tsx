import React, { useState, useEffect } from 'react';
import { getStalePolicy, setStalePolicy } from '../../services/api';
import { Clock, Save, Loader2 } from 'lucide-react';

interface StalePolicySettingsProps {
  projectId: string;
}

export const StalePolicySettings: React.FC<StalePolicySettingsProps> = ({ projectId }) => {
  const [staleAfterDays, setStaleDays] = useState<number>(30);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    const loadPolicy = async () => {
      setLoading(true);
      try {
        const policy = await getStalePolicy(projectId);
        if (policy?.stale_after_days) {
          setStaleDays(policy.stale_after_days);
        }
      } catch (err) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    };
    loadPolicy();
  }, [projectId]);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setMessage(null);
    try {
      await setStalePolicy(projectId, staleAfterDays);
      setMessage('Stale flag threshold saved successfully.');
    } catch (err) {
      setMessage('Failed to save policy.');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center gap-2 text-slate-500 text-sm p-4">
        <Loader2 className="w-4 h-4 animate-spin" />
        <span>Loading staleness policy...</span>
      </div>
    );
  }

  return (
    <div className="bg-white border border-slate-200 rounded-xl p-6 shadow-sm space-y-4 max-w-xl">
      <div className="flex items-center gap-3">
        <div className="w-9 h-9 bg-amber-50 text-amber-600 rounded-lg flex items-center justify-center">
          <Clock className="w-5 h-5" />
        </div>
        <div>
          <h3 className="text-base font-semibold text-slate-900">Stale Flag Detection Policy</h3>
          <p className="text-xs text-slate-500">
            Define after how many inactive days feature flags are marked as Stale.
          </p>
        </div>
      </div>

      <form onSubmit={handleSave} className="space-y-4 pt-2">
        <div>
          <label className="block text-xs font-semibold uppercase text-slate-500 mb-1">
            Staleness Threshold (Days)
          </label>
          <input
            type="number"
            min={1}
            max={365}
            value={staleAfterDays}
            onChange={(e) => setStaleDays(parseInt(e.target.value, 10) || 30)}
            className="w-full text-sm px-3 py-2 border border-slate-200 rounded-lg outline-none focus:ring-2 focus:ring-amber-500/20 focus:border-amber-500"
          />
          <p className="text-xs text-slate-400 mt-1">
            Default is 30 days. Flags at 100% rollout or inactive for longer than this will be
            marked STALE.
          </p>
        </div>

        {message && (
          <div className="text-xs font-medium text-amber-700 bg-amber-50 px-3 py-2 rounded-md">
            {message}
          </div>
        )}

        <button
          type="submit"
          disabled={saving}
          className="inline-flex items-center gap-2 bg-amber-600 hover:bg-amber-700 text-white text-xs font-medium px-4 py-2 rounded-lg transition-colors"
        >
          {saving ? (
            <Loader2 className="w-3.5 h-3.5 animate-spin" />
          ) : (
            <Save className="w-3.5 h-3.5" />
          )}
          <span>Save Threshold</span>
        </button>
      </form>
    </div>
  );
};
