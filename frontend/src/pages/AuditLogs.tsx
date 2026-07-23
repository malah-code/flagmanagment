import React, { useEffect, useState } from 'react';

interface AuditLog {
  id: string;
  action: string;
  target_type: string;
  target_id: string;
  created_at: string;
}

export const AuditLogs: React.FC = () => {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Mock audit logs fetch for UI
    setLogs([
      {
        id: '1',
        action: 'FLAG_TOGGLED',
        target_type: 'FLAG',
        target_id: 'new-checkout',
        created_at: new Date().toISOString(),
      },
      {
        id: '2',
        action: 'ENVIRONMENT_CREATED',
        target_type: 'ENVIRONMENT',
        target_id: 'Production',
        created_at: new Date(Date.now() - 3600000).toISOString(),
      },
    ]);
    setLoading(false);
  }, []);

  return (
    <div className="p-6 space-y-4">
      <h1 className="text-2xl font-bold text-slate-100">Audit Logs</h1>
      {loading ? (
        <div className="text-slate-400">Loading audit logs...</div>
      ) : (
        <div className="rounded-lg border border-slate-800 bg-slate-900 overflow-hidden">
          <table className="w-full text-left text-sm text-slate-300">
            <thead className="bg-slate-950 text-xs text-slate-400 uppercase">
              <tr>
                <th className="px-4 py-3">Action</th>
                <th className="px-4 py-3">Resource Type</th>
                <th className="px-4 py-3">Target ID</th>
                <th className="px-4 py-3">Timestamp</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800">
              {logs.map((log) => (
                <tr key={log.id} className="hover:bg-slate-800/50">
                  <td className="px-4 py-3 font-mono font-medium text-indigo-400">{log.action}</td>
                  <td className="px-4 py-3">{log.target_type}</td>
                  <td className="px-4 py-3 font-mono">{log.target_id}</td>
                  <td className="px-4 py-3 text-slate-400">{new Date(log.created_at).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};
