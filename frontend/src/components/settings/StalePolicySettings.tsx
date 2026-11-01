import React, { useState, useEffect } from 'react';
import { getStalePolicy, setStalePolicy, runStaleScan } from '../../services/api';
import { Clock, Save, Loader2, Bot, Play, CheckCircle2, Bell, GitPullRequest } from 'lucide-react';
import toast from 'react-hot-toast';

interface StalePolicySettingsProps {
  projectId: string;
}

interface ScanResult {
  flags_scanned: number;
  flags_marked_stale: number;
  change_requests_created: number;
  notifications_sent: number;
}

export const StalePolicySettings: React.FC<StalePolicySettingsProps> = ({ projectId }) => {
  const [staleAfterDays, setStaleDays] = useState<number>(30);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [scanning, setScanning] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [scanResult, setScanResult] = useState<ScanResult | null>(null);

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
      toast.success('Stale flag policy saved');
    } catch {
      setMessage('Failed to save policy.');
      toast.error('Failed to save policy');
    } finally {
      setSaving(false);
    }
  };

  const handleTriggerScan = async () => {
    setScanning(true);
    try {
      const res = await runStaleScan(projectId);
      setScanResult(res);
      toast.success('Governance bot scan completed');
    } catch {
      toast.error('Failed to run governance scan');
    } finally {
      setScanning(false);
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
    <div className="space-y-6 max-w-2xl">
      {/* Policy Configuration Card */}
      <div className="bg-white border border-slate-200 rounded-xl p-6 shadow-sm space-y-4">
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
              Default is 30 days. Flags with no evaluations or state changes exceeding this window
              will be identified by the governance bot.
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
            className="inline-flex items-center gap-2 bg-amber-600 hover:bg-amber-700 disabled:opacity-50 text-white text-xs font-medium px-4 py-2 rounded-lg transition-colors"
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

      {/* Automated Governance Bot Card */}
      <div className="bg-white border border-slate-200 rounded-xl p-6 shadow-sm space-y-5">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 bg-indigo-50 text-indigo-600 rounded-lg flex items-center justify-center">
              <Bot className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base font-semibold text-slate-900 flex items-center gap-2">
                Automated Governance Bot
                <span className="inline-flex items-center px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider bg-emerald-50 text-emerald-700 border border-emerald-200">
                  Active
                </span>
              </h3>
              <p className="text-xs text-slate-500">
                Runs hourly policy checks to automate dead flag cleanup and team alerting.
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={handleTriggerScan}
            disabled={scanning}
            className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white text-xs font-medium px-3.5 py-2 rounded-lg shadow-xs transition-colors shrink-0"
          >
            {scanning ? (
              <Loader2 className="w-3.5 h-3.5 animate-spin" />
            ) : (
              <Play className="w-3.5 h-3.5 fill-white" />
            )}
            <span>Run Governance Bot Scan</span>
          </button>
        </div>

        {/* Feature Badges */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
          <div className="p-3 bg-slate-50 rounded-lg border border-slate-100 flex items-start gap-2.5">
            <GitPullRequest className="w-4 h-4 text-indigo-600 shrink-0 mt-0.5" />
            <div>
              <span className="font-semibold text-slate-800">Auto Deprecation Proposals</span>
              <p className="text-[11px] text-slate-500 mt-0.5">
                Automatically drafts Change Requests to safely disable stale flags.
              </p>
            </div>
          </div>

          <div className="p-3 bg-slate-50 rounded-lg border border-slate-100 flex items-start gap-2.5">
            <Bell className="w-4 h-4 text-amber-500 shrink-0 mt-0.5" />
            <div>
              <span className="font-semibold text-slate-800">Slack & Webhook Alerts</span>
              <p className="text-[11px] text-slate-500 mt-0.5">
                Dispatches notification cards to team channels when flags become stale.
              </p>
            </div>
          </div>
        </div>

        {/* Last Scan Result Output */}
        {scanResult && (
          <div className="bg-slate-900 text-slate-100 rounded-lg p-4 text-xs space-y-2 border border-slate-800">
            <div className="flex items-center justify-between font-semibold text-emerald-400 border-b border-slate-800 pb-2">
              <span className="flex items-center gap-1.5">
                <CheckCircle2 className="w-3.5 h-3.5" /> Latest Governance Scan Result
              </span>
              <span className="text-[11px] font-mono text-slate-400">Just now</span>
            </div>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 pt-1 font-mono">
              <div>
                <span className="text-slate-400 block text-[10px]">SCANNED</span>
                <span className="font-bold text-slate-200">{scanResult.flags_scanned}</span>
              </div>
              <div>
                <span className="text-slate-400 block text-[10px]">MARKED STALE</span>
                <span className="font-bold text-amber-400">{scanResult.flags_marked_stale}</span>
              </div>
              <div>
                <span className="text-slate-400 block text-[10px]">CRS DRAFTED</span>
                <span className="font-bold text-indigo-400">
                  {scanResult.change_requests_created}
                </span>
              </div>
              <div>
                <span className="text-slate-400 block text-[10px]">NOTIFIED</span>
                <span className="font-bold text-emerald-400">{scanResult.notifications_sent}</span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
