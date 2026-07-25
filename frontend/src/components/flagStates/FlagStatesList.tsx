import { useFlagStates, useUpdateFlagState } from '../../hooks/useFlagStates';
import { useFlags } from '../../hooks/useFlags';
import { ToggleLeft, ToggleRight, Loader2, CheckCircle2, XCircle } from 'lucide-react';
import { useState } from 'react';
import { KillSwitchForm } from '../KillSwitchForm';

interface FlagStatesListProps {
  projectId: string;
  environmentId: string;
}

export const FlagStatesList = ({ projectId, environmentId }: FlagStatesListProps) => {
  const { data: flags = [], isLoading: isLoadingFlags } = useFlags(projectId);
  const { data: flagStates = [], isLoading: isLoadingStates, isError, error } = useFlagStates(environmentId);
  const updateMutation = useUpdateFlagState(environmentId);
  const [selectedFlagForKS, setSelectedFlagForKS] = useState<string | null>(null);

  const isLoading = isLoadingFlags || isLoadingStates;

  const handleToggle = async (stateId: string, currentEnabled: boolean) => {
    await updateMutation.mutateAsync({
      id: stateId,
      payload: { isEnabled: !currentEnabled },
    });
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
        <p className="text-sm text-slate-500">Toggle flag states directly for the selected environment.</p>
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
          No feature flags exist for this project yet. Add flags under the "Feature Flags" tab first.
        </div>
      ) : (
        <div className="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm">
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
                      {flag.key}
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
                        <div className="flex items-center justify-end gap-4">
                          <button
                            onClick={() => setSelectedFlagForKS(flag.id)}
                            className="text-xs font-semibold text-red-600 hover:text-red-700 bg-red-50 hover:bg-red-100 px-2.5 py-1.5 rounded transition-colors"
                          >
                            Kill Switches
                          </button>
                          <button
                            onClick={() => handleToggle(state.id, isEnabled)}
                            disabled={updateMutation.isPending}
                            className="focus:outline-none transition-transform active:scale-95 disabled:opacity-50 flex-shrink-0"
                            title={isEnabled ? 'Turn OFF' : 'Turn ON'}
                          >
                            {isEnabled ? (
                              <ToggleRight className="w-8 h-8 text-indigo-600 hover:text-indigo-700 transition-colors" />
                            ) : (
                              <ToggleLeft className="w-8 h-8 text-slate-300 hover:text-slate-400 transition-colors" />
                            )}
                          </button>
                        </div>
                      ) : (
                        <span className="text-xs text-slate-400 italic">No state record</span>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      {selectedFlagForKS && (
        <KillSwitchForm envId={environmentId} flagId={selectedFlagForKS} />
      )}
    </div>
  );
};
