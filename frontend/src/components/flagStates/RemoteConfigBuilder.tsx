import React, { useState, useEffect } from 'react';
import { X, Plus, Save, Settings, Code, FileText } from 'lucide-react';
import { useUpdateFlagState } from '../../hooks/useFlagStates';

interface RemoteConfigBuilderProps {
  isOpen: boolean;
  onClose: () => void;
  envId: string;
  projectId: string;
  flagId: string;
  flagKey: string;
  initialConfig?: Record<string, unknown>;
}

type EditorMode = 'VISUAL' | 'JSON';

interface VisualField {
  id: string;
  key: string;
  type: 'STRING' | 'NUMBER' | 'BOOLEAN';
  value: any;
}

export const RemoteConfigBuilder: React.FC<RemoteConfigBuilderProps> = ({
  isOpen,
  onClose,
  envId,
  projectId,
  flagId,
  flagKey,
  initialConfig = {},
}) => {
  const [mode, setMode] = useState<EditorMode>('VISUAL');
  const [jsonText, setJsonText] = useState('');
  const [fields, setFields] = useState<VisualField[]>([]);
  const [jsonError, setJsonError] = useState<string | null>(null);

  const updateMutation = useUpdateFlagState(projectId, envId);

  useEffect(() => {
    if (isOpen) {
      setJsonText(JSON.stringify(initialConfig || {}, null, 2));
      setFields(parseConfigToFields(initialConfig || {}));
      setJsonError(null);
    }
  }, [isOpen, initialConfig]);

  if (!isOpen) return null;

  // Sync Visual -> JSON
  const handleModeSwitch = (newMode: EditorMode) => {
    if (newMode === 'JSON') {
      const config = buildConfigFromFields(fields);
      setJsonText(JSON.stringify(config, null, 2));
    } else {
      try {
        const parsed = JSON.parse(jsonText);
        if (typeof parsed !== 'object' || Array.isArray(parsed) || parsed === null) {
          setJsonError('Payload must be a JSON object.');
          return;
        }
        setFields(parseConfigToFields(parsed));
        setJsonError(null);
      } catch (err) {
        setJsonError('Invalid JSON format.');
        return; // Prevent switching mode if invalid
      }
    }
    setMode(newMode);
  };

  const handleAddField = () => {
    setFields([...fields, { id: `f-${Date.now()}`, key: '', type: 'STRING', value: '' }]);
  };

  const handleRemoveField = (id: string) => {
    setFields(fields.filter((f) => f.id !== id));
  };

  const handleFieldChange = (id: string, keyName: keyof VisualField, val: any) => {
    setFields((prev) =>
      prev.map((f) => {
        if (f.id !== id) return f;
        const updated = { ...f, [keyName]: val };
        // If type changed, reset value safely
        if (keyName === 'type') {
          if (val === 'BOOLEAN') updated.value = false;
          else if (val === 'NUMBER') updated.value = 0;
          else updated.value = '';
        }
        return updated;
      }),
    );
  };

  const handleSave = async () => {
    let payloadConfig: Record<string, unknown>;

    if (mode === 'JSON') {
      try {
        const parsed = JSON.parse(jsonText);
        if (typeof parsed !== 'object' || Array.isArray(parsed) || parsed === null) {
          setJsonError('Payload must be a JSON object.');
          return;
        }
        payloadConfig = parsed;
      } catch (err) {
        setJsonError('Invalid JSON format.');
        return;
      }
    } else {
      payloadConfig = buildConfigFromFields(fields);
    }

    await updateMutation.mutateAsync({
      flagId,
      payload: {
        remoteConfig: payloadConfig,
      },
    });
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-3xl max-h-[90vh] flex flex-col">
        <div className="flex items-center justify-between p-6 border-b border-slate-200">
          <div>
            <h2 className="text-xl font-bold text-slate-900 flex items-center gap-2">
              <Settings className="w-5 h-5 text-indigo-600" />
              Remote Config Payload
            </h2>
            <p className="text-sm text-slate-500 mt-1">
              Attach dynamic payload configuration to{' '}
              <span className="font-mono text-slate-700">{flagKey}</span>
            </p>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-slate-500 transition-colors rounded-full p-1 hover:bg-slate-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="bg-slate-50 px-6 py-3 border-b border-slate-200 flex items-center justify-between">
          <div className="flex bg-slate-200/50 p-1 rounded-lg">
            <button
              onClick={() => handleModeSwitch('VISUAL')}
              className={`flex items-center gap-2 px-3 py-1.5 text-sm font-semibold rounded-md transition-colors ${
                mode === 'VISUAL'
                  ? 'bg-white text-indigo-700 shadow-sm'
                  : 'text-slate-600 hover:text-slate-900'
              }`}
            >
              <FileText className="w-4 h-4" /> Visual Editor
            </button>
            <button
              onClick={() => handleModeSwitch('JSON')}
              className={`flex items-center gap-2 px-3 py-1.5 text-sm font-semibold rounded-md transition-colors ${
                mode === 'JSON'
                  ? 'bg-white text-indigo-700 shadow-sm'
                  : 'text-slate-600 hover:text-slate-900'
              }`}
            >
              <Code className="w-4 h-4" /> Raw JSON
            </button>
          </div>
          {jsonError && <span className="text-sm text-red-600 font-medium">{jsonError}</span>}
        </div>

        <div className="p-6 overflow-y-auto flex-1 bg-slate-50/50">
          {mode === 'VISUAL' ? (
            <div className="space-y-4">
              {fields.length === 0 ? (
                <div className="text-center py-10 border-2 border-dashed border-slate-200 rounded-xl bg-white">
                  <p className="text-sm text-slate-500 mb-3">No config properties defined yet.</p>
                  <button
                    onClick={handleAddField}
                    className="inline-flex items-center gap-1.5 text-sm font-medium text-indigo-600 hover:text-indigo-700"
                  >
                    <Plus className="w-4 h-4" /> Add Property
                  </button>
                </div>
              ) : (
                <div className="space-y-3">
                  <div className="grid grid-cols-[2fr,1fr,2fr,auto] gap-3 px-1">
                    <label className="text-xs font-semibold text-slate-500 uppercase">
                      Key Name
                    </label>
                    <label className="text-xs font-semibold text-slate-500 uppercase">Type</label>
                    <label className="text-xs font-semibold text-slate-500 uppercase">Value</label>
                  </div>
                  {fields.map((field) => (
                    <div
                      key={field.id}
                      className="grid grid-cols-[2fr,1fr,2fr,auto] gap-3 items-center"
                    >
                      <input
                        type="text"
                        placeholder="e.g. primaryColor"
                        value={field.key}
                        onChange={(e) => handleFieldChange(field.id, 'key', e.target.value)}
                        className="rounded-lg border-slate-300 border px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none"
                      />
                      <select
                        value={field.type}
                        onChange={(e) => handleFieldChange(field.id, 'type', e.target.value)}
                        className="rounded-lg border-slate-300 border px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none bg-white"
                      >
                        <option value="STRING">String</option>
                        <option value="NUMBER">Number</option>
                        <option value="BOOLEAN">Boolean</option>
                      </select>

                      {field.type === 'BOOLEAN' ? (
                        <select
                          value={field.value.toString()}
                          onChange={(e) =>
                            handleFieldChange(field.id, 'value', e.target.value === 'true')
                          }
                          className="rounded-lg border-slate-300 border px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none bg-white"
                        >
                          <option value="true">True</option>
                          <option value="false">False</option>
                        </select>
                      ) : field.type === 'NUMBER' ? (
                        <input
                          type="number"
                          placeholder="0"
                          value={field.value}
                          onChange={(e) =>
                            handleFieldChange(field.id, 'value', Number(e.target.value))
                          }
                          className="rounded-lg border-slate-300 border px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none"
                        />
                      ) : (
                        <input
                          type="text"
                          placeholder="Value"
                          value={field.value}
                          onChange={(e) => handleFieldChange(field.id, 'value', e.target.value)}
                          className="rounded-lg border-slate-300 border px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-600 focus:border-indigo-600 outline-none"
                        />
                      )}

                      <button
                        onClick={() => handleRemoveField(field.id)}
                        className="p-2 text-slate-400 hover:text-red-500 transition-colors"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  ))}

                  <button
                    onClick={handleAddField}
                    className="inline-flex items-center gap-1.5 text-sm font-medium text-indigo-600 hover:text-indigo-700 mt-2"
                  >
                    <Plus className="w-4 h-4" /> Add Property
                  </button>
                </div>
              )}
            </div>
          ) : (
            <textarea
              value={jsonText}
              onChange={(e) => {
                setJsonText(e.target.value);
                setJsonError(null);
              }}
              className="w-full h-64 p-4 font-mono text-sm bg-slate-900 text-slate-100 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="{}"
            />
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
            Save Payload
          </button>
        </div>
      </div>
    </div>
  );
};

// Utilities for parsing visual fields from/to record
function parseConfigToFields(config: Record<string, unknown>): VisualField[] {
  return Object.keys(config).map((key, i) => {
    const val = config[key];
    let type: 'STRING' | 'NUMBER' | 'BOOLEAN' = 'STRING';
    let safeVal = val;

    if (typeof val === 'boolean') {
      type = 'BOOLEAN';
    } else if (typeof val === 'number') {
      type = 'NUMBER';
    } else if (typeof val === 'object') {
      type = 'STRING';
      safeVal = JSON.stringify(val);
    }

    return {
      id: `f-${Date.now()}-${i}`,
      key,
      type,
      value: safeVal,
    };
  });
}

function buildConfigFromFields(fields: VisualField[]): Record<string, unknown> {
  const result: Record<string, unknown> = {};
  for (const f of fields) {
    if (!f.key.trim()) continue;
    let val = f.value;
    if (f.type === 'NUMBER') val = Number(val);
    else if (f.type === 'BOOLEAN') val = Boolean(val);
    result[f.key.trim()] = val;
  }
  return result;
}
