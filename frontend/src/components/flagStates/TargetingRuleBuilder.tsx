import React, { useState } from 'react';
import { X, Plus, Save, Target } from 'lucide-react';
import { useUpdateFlagState } from '../../hooks/useFlagStates';

export interface TargetingCondition {
  attribute: string;
  operator: 'EQUALS' | 'CONTAINS' | 'REGEX';
  value: string;
}

export interface TargetingRule {
  id: string;
  conditions: TargetingCondition[];
  variation: boolean;
}

interface TargetingRuleBuilderProps {
  isOpen: boolean;
  onClose: () => void;
  envId: string;
  projectId: string;
  flagId: string;
  initialRules: TargetingRule[];
  flagKey: string;
}

export const TargetingRuleBuilder: React.FC<TargetingRuleBuilderProps> = ({
  isOpen,
  onClose,
  envId,
  projectId,
  flagId,
  initialRules,
  flagKey,
}) => {
  const [rules, setRules] = useState<TargetingRule[]>(initialRules || []);
  const updateMutation = useUpdateFlagState(projectId, envId);

  if (!isOpen) return null;

  const handleAddRule = () => {
    setRules([
      ...rules,
      {
        id: `rule-${Date.now()}`,
        conditions: [{ attribute: '', operator: 'EQUALS', value: '' }],
        variation: true,
      },
    ]);
  };

  const handleRemoveRule = (ruleId: string) => {
    setRules(rules.filter((r) => r.id !== ruleId));
  };

  const handleAddCondition = (ruleId: string) => {
    setRules(
      rules.map((r) =>
        r.id === ruleId
          ? {
              ...r,
              conditions: [...r.conditions, { attribute: '', operator: 'EQUALS', value: '' }],
            }
          : r
      )
    );
  };

  const handleRemoveCondition = (ruleId: string, condIndex: number) => {
    setRules(
      rules.map((r) =>
        r.id === ruleId
          ? {
              ...r,
              conditions: r.conditions.filter((_, i) => i !== condIndex),
            }
          : r
      )
    );
  };

  const handleConditionChange = (
    ruleId: string,
    condIndex: number,
    field: keyof TargetingCondition,
    value: string
  ) => {
    setRules(
      rules.map((r) =>
        r.id === ruleId
          ? {
              ...r,
              conditions: r.conditions.map((c, i) => (i === condIndex ? { ...c, [field]: value } : c)),
            }
          : r
      )
    );
  };

  const handleVariationChange = (ruleId: string, value: boolean) => {
    setRules(
      rules.map((r) =>
        r.id === ruleId ? { ...r, variation: value } : r
      )
    );
  };

  const handleSave = async () => {
    await updateMutation.mutateAsync({
      flagId,
      payload: {
        targetingRules: { rules },
      },
    });
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-4xl max-h-[90vh] flex flex-col">
        <div className="flex items-center justify-between p-6 border-b border-slate-200">
          <div>
            <h2 className="text-xl font-bold text-slate-900 flex items-center gap-2">
              <Target className="w-5 h-5 text-indigo-600" />
              Targeting Rules
            </h2>
            <p className="text-sm text-slate-500 mt-1">
              Configure contextual targeting for <span className="font-mono text-slate-700">{flagKey}</span>
            </p>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-slate-500 transition-colors rounded-full p-1 hover:bg-slate-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6 overflow-y-auto flex-1 space-y-6 bg-slate-50/50">
          {rules.length === 0 ? (
            <div className="text-center py-12 bg-white rounded-xl border border-slate-200 border-dashed">
              <Target className="w-12 h-12 text-slate-300 mx-auto mb-3" />
              <h3 className="text-sm font-medium text-slate-900">No targeting rules</h3>
              <p className="text-sm text-slate-500 mt-1 mb-4">
                Add rules to target specific users based on their context.
              </p>
              <button
                onClick={handleAddRule}
                className="inline-flex items-center gap-2 px-4 py-2 bg-white border border-slate-300 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-50 transition-colors"
              >
                <Plus className="w-4 h-4" /> Add Rule
              </button>
            </div>
          ) : (
            rules.map((rule, index) => (
              <div key={rule.id} className="bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden">
                <div className="bg-slate-50 px-4 py-3 border-b border-slate-200 flex justify-between items-center">
                  <div className="flex items-center gap-2 text-sm font-medium text-slate-700">
                    <span className="bg-indigo-100 text-indigo-700 px-2 py-0.5 rounded text-xs font-bold">
                      Rule {index + 1}
                    </span>
                    If ALL of the following conditions match:
                  </div>
                  <button
                    onClick={() => handleRemoveRule(rule.id)}
                    className="text-red-500 hover:text-red-700 text-sm font-medium"
                  >
                    Delete Rule
                  </button>
                </div>
                
                <div className="p-4 space-y-3">
                  {rule.conditions.map((cond, condIndex) => (
                    <div key={condIndex} className="flex items-center gap-3">
                      <input
                        type="text"
                        placeholder="Attribute (e.g. email)"
                        value={cond.attribute}
                        onChange={(e) => handleConditionChange(rule.id, condIndex, 'attribute', e.target.value)}
                        className="flex-1 rounded-lg border-slate-300 border px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none"
                      />
                      <select
                        value={cond.operator}
                        onChange={(e) => handleConditionChange(rule.id, condIndex, 'operator', e.target.value as any)}
                        className="rounded-lg border-slate-300 border px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none bg-white"
                      >
                        <option value="EQUALS">Equals</option>
                        <option value="CONTAINS">Contains</option>
                        <option value="REGEX">Matches Regex</option>
                      </select>
                      <input
                        type="text"
                        placeholder="Value"
                        value={cond.value}
                        onChange={(e) => handleConditionChange(rule.id, condIndex, 'value', e.target.value)}
                        className="flex-1 rounded-lg border-slate-300 border px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none"
                      />
                      <button
                        onClick={() => handleRemoveCondition(rule.id, condIndex)}
                        className="p-2 text-slate-400 hover:text-red-500 transition-colors"
                        disabled={rule.conditions.length === 1}
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  ))}

                  <button
                    onClick={() => handleAddCondition(rule.id)}
                    className="inline-flex items-center gap-1.5 text-sm font-medium text-indigo-600 hover:text-indigo-700 mt-2"
                  >
                    <Plus className="w-4 h-4" /> Add Condition (AND)
                  </button>
                </div>

                <div className="bg-slate-50/80 px-4 py-3 border-t border-slate-200 flex items-center justify-between">
                  <span className="text-sm font-medium text-slate-700">Then serve variation:</span>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleVariationChange(rule.id, true)}
                      className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
                        rule.variation ? 'bg-emerald-100 text-emerald-700 ring-1 ring-emerald-500' : 'bg-white border border-slate-300 text-slate-600 hover:bg-slate-50'
                      }`}
                    >
                      TRUE
                    </button>
                    <button
                      onClick={() => handleVariationChange(rule.id, false)}
                      className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${
                        !rule.variation ? 'bg-red-100 text-red-700 ring-1 ring-red-500' : 'bg-white border border-slate-300 text-slate-600 hover:bg-slate-50'
                      }`}
                    >
                      FALSE
                    </button>
                  </div>
                </div>
              </div>
            ))
          )}

          {rules.length > 0 && (
            <button
              onClick={handleAddRule}
              className="w-full py-3 border-2 border-dashed border-slate-300 rounded-xl text-sm font-medium text-slate-500 hover:text-slate-700 hover:border-slate-400 hover:bg-slate-50 transition-all flex items-center justify-center gap-2"
            >
              <Plus className="w-5 h-5" /> Add Another Rule (OR)
            </button>
          )}
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
            Save Rules
          </button>
        </div>
      </div>
    </div>
  );
};
