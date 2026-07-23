import { useState } from 'react';
import { useCreateFlag } from '../../hooks/useFlags';
import type { FlagType } from '../../types';
import { Flag, Loader2, X } from 'lucide-react';

interface CreateFlagDialogProps {
  projectId: string;
  isOpen: boolean;
  onClose: () => void;
}

export const CreateFlagDialog = ({ projectId, isOpen, onClose }: CreateFlagDialogProps) => {
  const [key, setKey] = useState('');
  const [description, setDescription] = useState('');
  const [type, setType] = useState<FlagType>('BOOLEAN');
  const [error, setError] = useState<string | null>(null);

  const createMutation = useCreateFlag();

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!key.trim()) {
      setError('Flag key is required (e.g. enable-new-checkout)');
      return;
    }

    try {
      setError(null);
      await createMutation.mutateAsync({
        projectId,
        key: key.trim(),
        description: description.trim(),
        type,
      });
      setKey('');
      setDescription('');
      setType('BOOLEAN');
      onClose();
    } catch (err: unknown) {
      setError((err as Error).message || 'Failed to create feature flag');
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="bg-white rounded-xl shadow-2xl border border-slate-200 max-w-md w-full p-6 space-y-5 animate-in fade-in zoom-in-95 duration-200">
        <div className="flex items-center justify-between border-b border-slate-100 pb-4">
          <div className="flex items-center gap-2 text-slate-900 font-semibold text-lg">
            <Flag className="w-5 h-5 text-indigo-600" />
            <span>Create Feature Flag</span>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-slate-600 transition-colors p-1 rounded-md hover:bg-slate-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="bg-red-50 text-red-600 border border-red-200 p-3 rounded-lg text-sm">
              {error}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">
              Flag Key <span className="text-red-500">*</span>
            </label>
            <input
              type="text"
              value={key}
              onChange={(e) => setKey(e.target.value)}
              placeholder="e.g. new-checkout-v2"
              className="w-full px-3.5 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none text-slate-900 placeholder:text-slate-400 text-sm font-mono transition-all"
            />
            <p className="text-xs text-slate-400 mt-1">Unique key used by SDKs to evaluate this flag.</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Flag Value Type</label>
            <select
              value={type}
              onChange={(e) => setType(e.target.value as FlagType)}
              className="w-full px-3.5 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none text-slate-900 text-sm bg-white transition-all"
            >
              <option value="BOOLEAN">Boolean (true / false)</option>
              <option value="STRING">String (Text)</option>
              <option value="NUMBER">Number (Numeric)</option>
              <option value="JSON">JSON Object</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700 mb-1">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Explain the purpose of this flag..."
              rows={3}
              className="w-full px-3.5 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none text-slate-900 placeholder:text-slate-400 text-sm transition-all resize-none"
            />
          </div>

          <div className="flex items-center justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100 rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={createMutation.isPending}
              className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg transition-colors disabled:opacity-50 flex items-center gap-2 shadow-sm"
            >
              {createMutation.isPending && <Loader2 className="w-4 h-4 animate-spin" />}
              Create Flag
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
