import { useState } from 'react';
import { useCreateEnvironment } from '../../hooks/useEnvironments';
import { Plus, Loader2, X, Copy, Check, AlertTriangle, ShieldCheck } from 'lucide-react';

interface CreateEnvironmentDialogProps {
  projectId: string;
  isOpen: boolean;
  onClose: () => void;
}

export const CreateEnvironmentDialog = ({ projectId, isOpen, onClose }: CreateEnvironmentDialogProps) => {
  const [name, setName] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [createdKey, setCreatedKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const createMutation = useCreateEnvironment();

  if (!isOpen) return null;

  const handleClose = () => {
    setName('');
    setError(null);
    setCreatedKey(null);
    setCopied(false);
    onClose();
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      setError('Environment name is required (e.g. Production, Staging)');
      return;
    }

    try {
      setError(null);
      const res = await createMutation.mutateAsync({ projectId, name: name.trim() });
      if (res.apiKey) {
        setCreatedKey(res.apiKey);
      } else {
        handleClose();
      }
    } catch (err: unknown) {
      setError((err as Error).message || 'Failed to create environment');
    }
  };

  const handleCopyKey = () => {
    if (createdKey) {
      navigator.clipboard.writeText(createdKey);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="bg-white rounded-xl shadow-2xl border border-slate-200 max-w-md w-full p-6 space-y-5 animate-in fade-in zoom-in-95 duration-200">
        <div className="flex items-center justify-between border-b border-slate-100 pb-4">
          <div className="flex items-center gap-2 text-slate-900 font-semibold text-lg">
            <ShieldCheck className="w-5 h-5 text-indigo-600" />
            <span>{createdKey ? 'Environment Created!' : 'Add Environment'}</span>
          </div>
          <button
            onClick={handleClose}
            className="text-slate-400 hover:text-slate-600 transition-colors p-1 rounded-md hover:bg-slate-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {createdKey ? (
          <div className="space-y-4">
            <div className="bg-amber-50 border border-amber-200 rounded-lg p-3.5 flex items-start gap-3 text-amber-800 text-xs">
              <AlertTriangle className="w-4 h-4 text-amber-600 shrink-0 mt-0.5" />
              <div>
                <strong className="font-semibold block mb-0.5">Save your API key now!</strong>
                This key will <strong>NEVER</strong> be displayed again. If lost, you will need to generate a new environment key.
              </div>
            </div>

            <div>
              <label className="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-1.5">
                Environment API Key
              </label>
              <div className="flex items-center gap-2">
                <input
                  type="text"
                  readOnly
                  value={createdKey}
                  className="w-full bg-slate-900 text-emerald-400 font-mono text-xs p-3 rounded-lg border border-slate-800 outline-none select-all"
                />
                <button
                  type="button"
                  onClick={handleCopyKey}
                  className="bg-indigo-600 hover:bg-indigo-700 text-white p-3 rounded-lg transition-colors flex items-center justify-center shrink-0"
                  title="Copy API Key"
                >
                  {copied ? <Check className="w-4 h-4 text-emerald-300" /> : <Copy className="w-4 h-4" />}
                </button>
              </div>
            </div>

            <div className="pt-2 flex justify-end">
              <button
                type="button"
                onClick={handleClose}
                className="w-full bg-slate-900 hover:bg-slate-800 text-white font-medium text-sm py-2.5 rounded-lg transition-colors shadow-sm"
              >
                I have securely copied this key
              </button>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm">
                {error}
              </div>
            )}

            <div>
              <label className="block text-sm font-medium text-slate-700 mb-1">
                Environment Name <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. Production, Staging, QA"
                className="w-full px-3.5 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none text-slate-900 placeholder:text-slate-400 text-sm transition-all"
              />
            </div>

            <div className="flex items-center justify-end gap-3 pt-2">
              <button
                type="button"
                onClick={handleClose}
                className="px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 rounded-lg transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={createMutation.isPending}
                className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg transition-colors disabled:opacity-50 flex items-center gap-2 shadow-sm"
              >
                {createMutation.isPending ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <Plus className="w-4 h-4" />
                )}
                Create Environment
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};
