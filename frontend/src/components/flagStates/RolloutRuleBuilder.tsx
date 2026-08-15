import React, { useState } from 'react';
import { X, Save, Percent } from 'lucide-react';
import { useUpdateFlagState } from '../../hooks/useFlagStates';
import type { Variation, RolloutRule } from '../../types';

interface RolloutRuleBuilderProps {
  isOpen: boolean;
  onClose: () => void;
  envId: string;
  projectId: string;
  flagId: string;
  flagKey: string;
  variations: Variation[];
  initialRollouts?: RolloutRule[];
  targetingRules?: any;
  remoteConfig?: any;
}

export const RolloutRuleBuilder: React.FC<RolloutRuleBuilderProps> = ({
  isOpen,
  onClose,
  envId,
  projectId,
  flagId,
  flagKey,
  variations = [],
  initialRollouts = [],
  targetingRules = { rules: [] },
  remoteConfig = {},
}) => {
  // Map percentages (in basis points 0-10000)
  const [rollouts, setRollouts] = useState<{ [variationId: string]: number }>(() => {
    const map: { [key: string]: number } = {};
    if (initialRollouts && initialRollouts.length > 0) {
      initialRollouts.forEach((r) => {
        map[r.variationId] = r.percentage / 100; // Convert to %
      });
    } else {
      const equalShare = Math.floor(100 / (variations.length || 1));
      variations.forEach((v) => {
        map[v.id] = equalShare;
      });
    }
    return map;
  });

  const [error, setError] = useState<string | null>(null);
  const updateMutation = useUpdateFlagState(projectId, envId);

  if (!isOpen) return null;

  const totalPercentage = Object.values(rollouts).reduce((acc, val) => acc + (Number(val) || 0), 0);

  const handlePercentageChange = (varId: string, val: number) => {
    setRollouts({
      ...rollouts,
      [varId]: val,
    });
  };

  const handleSave = async () => {
    if (Math.round(totalPercentage * 100) !== 10000) {
      setError(`Total percentage must equal exactly 100% (currently ${totalPercentage}%).`);
      return;
    }

    const rules: RolloutRule[] = Object.entries(rollouts).map(([variationId, pct]) => ({
      variationId,
      percentage: Math.round(pct * 100), // convert to basis points
    }));

    try {
      setError(null);
      await updateMutation.mutateAsync({
        flagId,
        payload: {
          targetingRules,
          remoteConfig,
          rolloutRules: { rules },
        },
      });
      onClose();
    } catch (err: unknown) {
      setError((err as Error).message || 'Failed to save rollout rules');
    }
  };

  return (
    <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-xl flex flex-col">
        <div className="flex items-center justify-between p-6 border-b border-slate-200">
          <div>
            <h2 className="text-xl font-bold text-slate-900 flex items-center gap-2">
              <Percent className="w-5 h-5 text-indigo-600" />
              Percentage Rollout (A/B/n Split)
            </h2>
            <p className="text-sm text-slate-500 mt-1">
              Configure traffic distribution for <span className="font-mono text-slate-700">{flagKey}</span>
            </p>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-slate-500 transition-colors rounded-full p-1 hover:bg-slate-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6 space-y-6">
          {error && (
            <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm">
              {error}
            </div>
          )}

          <div className="space-y-4">
            {variations.map((v) => (
              <div key={v.id} className="bg-slate-50 border border-slate-200 p-4 rounded-xl space-y-2">
                <div className="flex justify-between items-center text-sm">
                  <span className="font-semibold text-slate-900">{v.name}</span>
                  <div className="flex items-center gap-1">
                    <input
                      type="number"
                      min="0"
                      max="100"
                      step="0.01"
                      value={rollouts[v.id] ?? 0}
                      onChange={(e) => handlePercentageChange(v.id, parseFloat(e.target.value) || 0)}
                      className="w-20 px-2 py-1 border border-slate-300 rounded text-right text-sm font-mono font-bold focus:ring-1 focus:ring-indigo-500 outline-none"
                    />
                    <span className="text-slate-500 font-bold">%</span>
                  </div>
                </div>
                <input
                  type="range"
                  min="0"
                  max="100"
                  step="1"
                  value={rollouts[v.id] ?? 0}
                  onChange={(e) => handlePercentageChange(v.id, parseFloat(e.target.value) || 0)}
                  className="w-full accent-indigo-600 h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer"
                />
              </div>
            ))}
          </div>

          <div className="flex justify-between items-center p-3 bg-indigo-50 rounded-lg text-sm border border-indigo-100">
            <span className="font-medium text-indigo-900">Total Traffic Allocation:</span>
            <span className={`font-bold font-mono text-base ${totalPercentage === 100 ? 'text-emerald-600' : 'text-amber-600'}`}>
              {totalPercentage}%
            </span>
          </div>
        </div>

        <div className="p-6 border-t border-slate-200 bg-white flex justify-end gap-3 rounded-b-xl">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-lg hover:bg-slate-50 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={updateMutation.isPending}
            className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 flex items-center gap-2"
          >
            {updateMutation.isPending ? (
              <span className="animate-spin text-white">⟳</span>
            ) : (
              <Save className="w-4 h-4" />
            )}
            Save Rollout Split
          </button>
        </div>
      </div>
    </div>
  );
};
