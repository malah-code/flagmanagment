import toast from 'react-hot-toast';
import { useFlagStates, useUpdateFlagState, useInitFlagState } from '../../hooks/useFlagStates';
import { useFlags } from '../../hooks/useFlags';
import { useEnvironments } from '../../hooks/useEnvironments';
import { Loader2, CheckCircle2, XCircle, ArrowUpRight, Target, Clock } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Switch } from '../ui/Switch';
import { useState, useEffect, useCallback } from 'react';
import { KillSwitchForm } from '../KillSwitchForm';
import { SlackConfigForm } from '../SlackConfigForm';
import { PromoteFlagModal } from './PromoteFlagModal';
import { TargetingRuleBuilder } from './TargetingRuleBuilder';
import type { TargetingRule } from './TargetingRuleBuilder';
import { ScheduledChangeBadge } from '../flags/ScheduledChangeBadge';
import { ScheduleDialog } from '../flags/ScheduleDialog';
import { scheduledChangesApi } from '../../services/scheduledChangesApi';
import type { ScheduledChange } from '../../types/scheduledChange';

interface FlagStatesListProps {
  projectId: string;
  environmentId: string;
}

export const FlagStatesList = ({ projectId, environmentId }: FlagStatesListProps) => {
  const { data: flags = [], isLoading: isLoadingFlags } = useFlags(projectId);
  const {
    data: flagStates = [],
    isLoading: isLoadingStates,
    isError,
    error,
  } = useFlagStates(projectId, environmentId);
  const { data: environments = [] } = useEnvironments(projectId);
  const updateMutation = useUpdateFlagState(projectId, environmentId);
  const initMutation = useInitFlagState(projectId, environmentId);
  const [selectedFlagForKS, setSelectedFlagForKS] = useState<string | null>(null);
  const [selectedFlagForPromote, setSelectedFlagForPromote] = useState<string | null>(null);
  const [editingRulesState, setEditingRulesState] = useState<{
    flagId: string;
    key: string;
    rules: TargetingRule[];
  } | null>(null);
  const [scheduledChanges, setScheduledChanges] = useState<Record<string, ScheduledChange>>({});
  const [selectedFlagForSchedule, setSelectedFlagForSchedule] = useState<{
    id: string;
    name: string;
  } | null>(null);
  const [togglingStateId, setTogglingStateId] = useState<string | null>(null);
  const [isTransitioning, setIsTransitioning] = useState(false);

  useEffect(() => {
    setIsTransitioning(true);
    const timer = setTimeout(() => {
      setIsTransitioning(false);
    }, 200);
    return () => clearTimeout(timer);
  }, [environmentId]);

  const loadSchedules = useCallback(async () => {
    if (!environmentId) return;
    try {
      const res = await scheduledChangesApi.list(environmentId, 'PENDING');
      const map: Record<string, ScheduledChange> = {};
      if (res.data) {
        res.data.forEach((sc) => {
          if (sc.target_type === 'FLAG') {
            map[sc.target_id] = sc;
          }
        });
      }
      setScheduledChanges(map);
    } catch {
      // Ignore list error
    }
  }, [environmentId]);

  useEffect(() => {
    loadSchedules();
  }, [loadSchedules]);

  const isLoading = isLoadingFlags || isLoadingStates;

  const currentEnv = environments.find((e) => e.id === environmentId);

  const handleToggle = async (flagId: string, currentEnabled: boolean) => {
    const state = flagStates.find((s) => s.flagId === flagId);
    if (!state) return;

    setTogglingStateId(flagId);
    try {
      await updateMutation.mutateAsync({
        flagId,
        payload: {
          isEnabled: !currentEnabled,
          targetingRules: state.targetingRules || { rules: [] },
          remoteConfig: state.remoteConfig || {},
          rolloutRules: state.rolloutRules ? { rules: state.rolloutRules } : undefined,
          defaultVariation: state.defaultVariation,
        },
      });
      toast.success('Flag updated successfully');
    } catch (err: any) {
      toast.error('Failed to update flag state');
    } finally {
      setTogglingStateId(null);
    }
  };

  const handleInit = async (flagId: string) => {
    try {
      await initMutation.mutateAsync({
        flagId,
        payload: {
          isEnabled: false,
          targetingRules: { rules: [] },
          remoteConfig: {},
        },
      });
      toast.success('Flag state initialized successfully');
    } catch (err: any) {
      toast.error('Failed to initialize flag state');
    }
  };

  // Map flags with their environment state
  const flagItems = flags.map((flag) => {
    const state = flagStates.find((s) => s.flagId === flag.id);
    return {
      flag,
      state,
    };
  });

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold text-slate-900">Environment Targeting & Status</h2>
        <p className="text-sm text-slate-500">
          Toggle flag states directly for the selected environment.
        </p>
      </div>

      {isLoading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="w-6 h-6 animate-spin text-indigo-600" />
        </div>
      ) : isError ? (
        <div className="bg-red-50 text-red-600 p-4 rounded-xl border border-red-200 text-sm">
          Failed to load flag states: {(error as Error).message}
        </div>
      ) : flagItems.length === 0 ? (
        <div className="bg-white border border-slate-200 rounded-xl p-8 text-center text-sm text-slate-500">
          No feature flags exist for this project yet. Add flags under the "Feature Flags" tab
          first.
        </div>
      ) : (
        <div
          className={`bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm overflow-x-auto transition-opacity duration-200 ${isTransitioning ? 'opacity-50' : 'opacity-100'}`}
        >
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 border-b border-slate-200 text-xs uppercase tracking-wider font-semibold text-slate-500">
              <tr>
                <th className="px-6 py-3.5">Flag Key</th>
                <th className="px-6 py-3.5">Type</th>
                <th className="px-6 py-3.5">Status</th>
                <th className="px-6 py-3.5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {flagItems.map(({ flag, state }) => {
                const isEnabled = state?.isEnabled ?? false;
                return (
                  <tr key={flag.id} className="hover:bg-slate-50/50 transition-colors">
                    <td className="px-6 py-4 font-mono font-medium text-slate-900">
                      <div className="flex items-center gap-2">
                        <Link
                          to={`/projects/${projectId}/flags/${flag.id}`}
                          className="text-indigo-600 hover:text-indigo-800 hover:underline"
                        >
                          {flag.key}
                        </Link>
                        <ScheduledChangeBadge scheduledChange={scheduledChanges[flag.id]} />
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className="bg-slate-100 text-slate-700 text-xs font-semibold px-2.5 py-1 rounded-md">
                        {flag.type}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      {isEnabled ? (
                        <span className="inline-flex items-center gap-1.5 text-emerald-600 bg-emerald-50 text-xs font-semibold px-2.5 py-1 rounded-full border border-emerald-200">
                          <CheckCircle2 className="w-3.5 h-3.5" />
                          ACTIVE (ON)
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1.5 text-slate-500 bg-slate-100 text-xs font-semibold px-2.5 py-1 rounded-full border border-slate-200">
                          <XCircle className="w-3.5 h-3.5" />
                          INACTIVE (OFF)
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-right">
                      {state ? (
                        <div className="flex items-center justify-end gap-3">
                          <button
                            onClick={() => {
                              const existingRules =
                                (state.targetingRules?.rules as TargetingRule[]) || [];
                              setEditingRulesState({
                                flagId: flag.id,
                                key: flag.key,
                                rules: existingRules,
                              });
                            }}
                            className="inline-flex items-center gap-1 text-xs font-semibold text-slate-600 hover:text-slate-800 bg-slate-100 hover:bg-slate-200 px-2.5 py-1.5 rounded transition-colors"
                          >
                            <Target className="w-3.5 h-3.5" /> Targeting
                          </button>
                          <button
                            onClick={() =>
                              setSelectedFlagForSchedule({ id: flag.id, name: flag.key })
                            }
                            className="inline-flex items-center gap-1 text-xs font-semibold text-amber-700 hover:text-amber-800 bg-amber-50 hover:bg-amber-100 border border-amber-200 px-2.5 py-1.5 rounded transition-colors"
                          >
                            <Clock className="w-3.5 h-3.5" /> Schedule
                          </button>
                          <button
                            onClick={() => setSelectedFlagForPromote(flag.id)}
                            className="inline-flex items-center gap-1 text-xs font-semibold text-blue-600 hover:text-blue-700 bg-blue-50 hover:bg-blue-100 px-2.5 py-1.5 rounded transition-colors"
                          >
                            <ArrowUpRight className="w-3.5 h-3.5" /> Promote
                          </button>
                          <button
                            onClick={() => setSelectedFlagForKS(flag.id)}
                            className="text-xs font-semibold text-red-600 hover:text-red-700 bg-red-50 hover:bg-red-100 px-2.5 py-1.5 rounded transition-colors"
                          >
                            Kill Switches
                          </button>
                          <div className="flex-shrink-0 ml-1">
                            <Switch
                              checked={isEnabled}
                              onChange={() => handleToggle(flag.id, isEnabled)}
                              loading={updateMutation.isPending && togglingStateId === flag.id}
                            />
                          </div>
                        </div>
                      ) : (
                        <button
                          onClick={() => handleInit(flag.id)}
                          disabled={initMutation.isPending}
                          className="inline-flex items-center gap-1 text-xs font-semibold text-indigo-600 hover:text-indigo-700 bg-indigo-50 hover:bg-indigo-100 border border-indigo-200 px-3 py-1.5 rounded-lg transition-colors focus:ring-2 focus:ring-indigo-500 outline-none disabled:opacity-50"
                        >
                          {initMutation.isPending ? (
                            <Loader2 className="w-3.5 h-3.5 animate-spin" />
                          ) : (
                            '+ Initialize State'
                          )}
                        </button>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {selectedFlagForKS && <KillSwitchForm envId={environmentId} flagId={selectedFlagForKS} />}

      <PromoteFlagModal
        isOpen={selectedFlagForPromote !== null}
        onClose={() => setSelectedFlagForPromote(null)}
        projectId={projectId}
        flagId={selectedFlagForPromote!}
        sourceEnvId={environmentId}
        sourceEnvName={currentEnv?.name || 'Current Environment'}
        onSuccess={() => {
          // Re-fetch states if needed
        }}
      />

      <ScheduleDialog
        isOpen={selectedFlagForSchedule !== null}
        onClose={() => setSelectedFlagForSchedule(null)}
        flagId={selectedFlagForSchedule?.id || ''}
        flagName={selectedFlagForSchedule?.name || ''}
        environmentId={environmentId}
        existingSchedule={
          selectedFlagForSchedule ? scheduledChanges[selectedFlagForSchedule.id] : null
        }
        onSuccess={loadSchedules}
      />

      <SlackConfigForm environmentId={environmentId} />

      {editingRulesState && (
        <TargetingRuleBuilder
          isOpen={true}
          onClose={() => setEditingRulesState(null)}
          envId={environmentId}
          projectId={projectId}
          flagId={editingRulesState.flagId}
          flagKey={editingRulesState.key}
          initialRules={editingRulesState.rules}
        />
      )}
    </div>
  );
};
