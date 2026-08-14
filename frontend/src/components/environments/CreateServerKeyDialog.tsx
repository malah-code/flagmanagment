import { useState } from 'react';
import type { FormEvent } from 'react';
import { X, Key, Copy, Check, AlertCircle } from 'lucide-react';
import toast from 'react-hot-toast';
import { environmentService } from '../../services/environments';

interface CreateServerKeyDialogProps {
  projectId: string;
  envId: string;
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export const CreateServerKeyDialog = ({
  projectId,
  envId,
  isOpen,
  onClose,
  onSuccess,
}: CreateServerKeyDialogProps) => {
  const [name, setName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [createdKey, setCreatedKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setIsLoading(true);
    try {
      const res = await environmentService.createServerKey(projectId, envId, name.trim());
      setCreatedKey(res.key);
      toast.success('Server-side key created successfully!');
      onSuccess();
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to create server key';
      toast.error(message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopy = () => {
    if (!createdKey) return;
    navigator.clipboard.writeText(createdKey);
    setCopied(true);
    toast.success('Server key copied!');
    setTimeout(() => setCopied(false), 2000);
  };

  const handleClose = () => {
    setName('');
    setCreatedKey(null);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/40 backdrop-blur-sm">
      <div className="bg-white border border-slate-200 rounded-2xl max-w-md w-full p-6 shadow-xl space-y-5 animate-in fade-in zoom-in-95 duration-150">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <div className="p-2 bg-indigo-50 text-indigo-600 rounded-xl">
              <Key className="w-5 h-5" />
            </div>
            <h3 className="font-bold text-slate-900 text-lg">
              {createdKey ? 'Server Key Created' : 'Create Server-side Key'}
            </h3>
          </div>
          <button
            onClick={handleClose}
            className="text-slate-400 hover:text-slate-600 transition-colors p-1 rounded-lg hover:bg-slate-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {createdKey ? (
          <div className="space-y-4">
            <div className="p-3.5 bg-amber-50 border border-amber-200/80 rounded-xl text-amber-800 text-xs flex items-start gap-2.5">
              <AlertCircle className="w-4 h-4 text-amber-600 shrink-0 mt-0.5" />
              <div>
                <strong className="font-semibold">Save this secret key!</strong>
                <p className="mt-0.5 text-amber-700">
                  This key will only be displayed once. Store it securely in your server configuration or environment variables.
                </p>
              </div>
            </div>

            <div className="space-y-1.5">
              <label className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
                Generated Key
              </label>
              <div className="flex items-center gap-2">
                <div className="flex-1 font-mono text-xs text-slate-800 bg-slate-50 border border-slate-200 rounded-xl px-3 py-2 truncate select-all">
                  {createdKey}
                </div>
                <button
                  onClick={handleCopy}
                  className="inline-flex items-center gap-1.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-semibold px-3.5 py-2 rounded-xl transition-colors shrink-0"
                >
                  {copied ? (
                    <>
                      <Check className="w-3.5 h-3.5" />
                      <span>Copied!</span>
                    </>
                  ) : (
                    <>
                      <Copy className="w-3.5 h-3.5" />
                      <span>Copy</span>
                    </>
                  )}
                </button>
              </div>
            </div>

            <div className="pt-2">
              <button
                onClick={handleClose}
                className="w-full bg-slate-900 hover:bg-slate-800 text-white text-sm font-semibold py-2.5 rounded-xl transition-colors"
              >
                Done
              </button>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-1.5">
              <label htmlFor="keyName" className="text-sm font-semibold text-slate-700">
                Key Name
              </label>
              <input
                id="keyName"
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="e.g. billing-service, staging-k8s-pod"
                className="w-full px-3.5 py-2.5 border border-slate-300 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                required
                autoFocus
              />
              <p className="text-xs text-slate-500">
                A descriptive name to identify which microservice or backend environment uses this token.
              </p>
            </div>

            <div className="flex justify-end gap-3 pt-2">
              <button
                type="button"
                onClick={handleClose}
                className="px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-100 rounded-xl transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isLoading || !name.trim()}
                className="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 rounded-xl transition-colors shadow-sm"
              >
                {isLoading ? 'Creating...' : 'Create Server Key'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};
