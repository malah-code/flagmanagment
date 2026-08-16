import { useState, useEffect, useCallback } from 'react';
import { Key, Plus, Trash2, Search, Eye, EyeOff, ShieldAlert, Loader2 } from 'lucide-react';
import toast from 'react-hot-toast';
import type { ServerKey } from '../../types';
import { environmentService } from '../../services/environments';
import { CreateServerKeyDialog } from './CreateServerKeyDialog';

interface ServerSideKeysPanelProps {
  projectId: string;
  envId: string;
}

export const ServerSideKeysPanel = ({ projectId, envId }: ServerSideKeysPanelProps) => {
  const [keys, setKeys] = useState<ServerKey[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [visibleKeyIds, setVisibleKeyIds] = useState<Record<string, boolean>>({});
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  const fetchKeys = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await environmentService.listServerKeys(projectId, envId);
      setKeys(data);
    } catch {
      toast.error('Failed to load server-side keys');
    } finally {
      setIsLoading(false);
    }
  }, [projectId, envId]);

  useEffect(() => {
    fetchKeys();
  }, [fetchKeys]);

  const handleDelete = async (keyId: string, name: string) => {
    if (
      !confirm(
        `Are you sure you want to revoke server key "${name}"? Backend SDKs using this key will lose access.`,
      )
    ) {
      return;
    }

    try {
      await environmentService.deleteServerKey(projectId, envId, keyId);
      toast.success(`Server key "${name}" revoked`);
      fetchKeys();
    } catch {
      toast.error('Failed to revoke server key');
    }
  };

  const toggleVisibility = (keyId: string) => {
    setVisibleKeyIds((prev) => ({
      ...prev,
      [keyId]: !prev[keyId],
    }));
  };

  const filteredKeys = keys.filter((k) => k.name.toLowerCase().includes(searchQuery.toLowerCase()));

  return (
    <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm space-y-5">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h3 className="text-base font-bold text-slate-900">Server-side Environment Keys</h3>
            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-amber-50 text-amber-700 border border-amber-200">
              <ShieldAlert className="w-3.5 h-3.5" /> Private / Secret Keys
            </span>
          </div>
          <p className="text-sm text-slate-500 mt-1">
            Privileged SDK keys for backend local evaluation (Go, Node.js, Python, Java). Kept
            secret.
          </p>
        </div>

        <button
          onClick={() => setIsCreateOpen(true)}
          className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs px-3.5 py-2.5 rounded-xl shadow-sm transition-colors shrink-0"
        >
          <Plus className="w-4 h-4" />
          <span>Create Server-side Key</span>
        </button>
      </div>

      {/* Search and Filters */}
      <div className="flex items-center gap-3">
        <div className="relative flex-1">
          <Search className="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search server keys by name..."
            className="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-xl text-xs focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          />
        </div>
      </div>

      {/* Table / List */}
      {isLoading ? (
        <div className="flex justify-center p-8">
          <Loader2 className="w-6 h-6 animate-spin text-indigo-600" />
        </div>
      ) : filteredKeys.length === 0 ? (
        <div className="text-center py-8 bg-slate-50 border border-slate-200/80 rounded-xl space-y-2">
          <Key className="w-8 h-8 text-slate-300 mx-auto" />
          <p className="text-xs font-medium text-slate-500">
            {searchQuery ? 'No server keys matching query.' : 'No server-side keys created yet.'}
          </p>
        </div>
      ) : (
        <div className="overflow-x-auto border border-slate-200 rounded-xl">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 border-b border-slate-200 text-[11px] font-semibold text-slate-500 uppercase tracking-wider">
                <th className="py-3 px-4">Key Name</th>
                <th className="py-3 px-4">Secret Value</th>
                <th className="py-3 px-4">Created</th>
                <th className="py-3 px-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 text-xs">
              {filteredKeys.map((keyItem) => {
                const isRevealed = !!visibleKeyIds[keyItem.id];
                return (
                  <tr key={keyItem.id} className="hover:bg-slate-50/80 transition-colors">
                    <td className="py-3 px-4 font-semibold text-slate-900">{keyItem.name}</td>
                    <td className="py-3 px-4">
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-slate-600 bg-slate-100 px-2 py-1 rounded border border-slate-200 text-[11px]">
                          {isRevealed
                            ? `env_server_key_${keyItem.id.slice(0, 8)}...`
                            : '••••••••••••••••••••••••'}
                        </span>
                        <button
                          onClick={() => toggleVisibility(keyItem.id)}
                          className="text-slate-400 hover:text-slate-600 p-1 rounded transition-colors"
                          title={isRevealed ? 'Mask key' : 'Show key hint'}
                        >
                          {isRevealed ? (
                            <EyeOff className="w-3.5 h-3.5" />
                          ) : (
                            <Eye className="w-3.5 h-3.5" />
                          )}
                        </button>
                      </div>
                    </td>
                    <td className="py-3 px-4 text-slate-500">
                      {new Date(keyItem.createdAt).toLocaleDateString()}
                    </td>
                    <td className="py-3 px-4 text-right">
                      <button
                        onClick={() => handleDelete(keyItem.id, keyItem.name)}
                        className="text-slate-400 hover:text-red-600 p-1.5 rounded-lg hover:bg-slate-100 transition-colors"
                        title="Revoke / Delete Key"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      <CreateServerKeyDialog
        projectId={projectId}
        envId={envId}
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSuccess={fetchKeys}
      />
    </div>
  );
};
