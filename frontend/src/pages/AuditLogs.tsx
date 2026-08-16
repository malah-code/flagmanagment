import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Shield, Search, Loader2, Calendar, User, Activity } from 'lucide-react';
import { useAuditLogs } from '../hooks/useAuditLogs';
import type { AuditLog } from '../services/audit';

interface AuditLogsProps {
  projectId?: string;
}

export const AuditLogs: React.FC<AuditLogsProps> = ({ projectId: propProjectId }) => {
  const { projectId: routeProjectId } = useParams<{ projectId: string }>();
  const effectiveProjectId = propProjectId || routeProjectId || '';

  const { data, isLoading, error } = useAuditLogs(effectiveProjectId);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null);

  const logs = data?.data || [];

  const filteredLogs = logs.filter((log) => {
    const query = searchQuery.toLowerCase();
    return (
      log.action.toLowerCase().includes(query) ||
      log.target_type.toLowerCase().includes(query) ||
      log.target_id.toLowerCase().includes(query) ||
      log.actor_id.toLowerCase().includes(query)
    );
  });

  const getActionBadgeClass = (action: string) => {
    const act = action.toUpperCase();
    if (act.includes('CREATE') || act.includes('INSERT')) {
      return 'bg-emerald-50 text-emerald-700 border-emerald-200';
    }
    if (act.includes('UPDATE') || act.includes('TOGGLE')) {
      return 'bg-indigo-50 text-indigo-700 border-indigo-200';
    }
    if (act.includes('DELETE') || act.includes('REMOVE') || act.includes('KILL')) {
      return 'bg-rose-50 text-rose-700 border-rose-200';
    }
    return 'bg-slate-100 text-slate-700 border-slate-200';
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900 flex items-center gap-2">
            <Shield className="w-6 h-6 text-indigo-600" />
            Audit Logs
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Tamper-evident chronological activity trail for this workspace
          </p>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <div className="p-4 border-b border-slate-200 bg-slate-50/50 flex flex-col sm:flex-row items-center justify-between gap-3">
          <div className="relative w-full sm:w-80">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search by action, resource, target..."
              className="w-full pl-9 pr-4 py-2 bg-white text-sm border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div className="text-xs text-slate-500 font-medium">
            Showing {filteredLogs.length} of {logs.length} events
          </div>
        </div>

        {isLoading ? (
          <div className="flex flex-col items-center justify-center p-12 text-slate-400">
            <Loader2 className="w-8 h-8 animate-spin text-indigo-600 mb-2" />
            <p className="text-sm font-medium">Loading audit events...</p>
          </div>
        ) : error ? (
          <div className="p-8 text-center text-rose-600 bg-rose-50/50">
            <p className="font-semibold text-sm">Failed to load audit logs.</p>
            <p className="text-xs mt-1 text-rose-500">
              {(error as any)?.message || 'Please try again.'}
            </p>
          </div>
        ) : filteredLogs.length === 0 ? (
          <div className="p-12 text-center text-slate-500">
            <Activity className="w-12 h-12 text-slate-300 mx-auto mb-3" />
            <p className="font-medium text-slate-700">No audit logs found</p>
            <p className="text-xs text-slate-400 mt-1">
              Actions performed on flags and environments will appear here.
            </p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="bg-slate-50 text-xs font-semibold text-slate-500 uppercase tracking-wider border-b border-slate-200">
                <tr>
                  <th className="px-6 py-3">Action</th>
                  <th className="px-6 py-3">Resource Type</th>
                  <th className="px-6 py-3">Target ID</th>
                  <th className="px-6 py-3">Actor</th>
                  <th className="px-6 py-3">Timestamp</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 font-sans">
                {filteredLogs.map((log) => (
                  <tr
                    key={log.id}
                    onClick={() => setSelectedLog(log)}
                    className="hover:bg-slate-50/75 transition-colors cursor-pointer"
                  >
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span
                        className={`inline-flex items-center text-xs font-semibold px-2.5 py-1 rounded-full border ${getActionBadgeClass(log.action)}`}
                      >
                        {log.action}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap font-medium text-slate-700">
                      {log.target_type}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap font-mono text-xs text-indigo-600">
                      {log.target_id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-xs text-slate-500">
                      <span className="inline-flex items-center gap-1">
                        <User className="w-3.5 h-3.5 text-slate-400" />
                        {log.actor_id ? log.actor_id.slice(0, 8) + '...' : 'System'}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-xs text-slate-500">
                      <span className="inline-flex items-center gap-1.5">
                        <Calendar className="w-3.5 h-3.5 text-slate-400" />
                        {new Date(log.created_at).toLocaleString()}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Log Details Modal */}
      {selectedLog && (
        <div className="fixed inset-0 bg-slate-900/50 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[85vh] flex flex-col overflow-hidden border border-slate-200">
            <div className="px-6 py-4 border-b border-slate-200 flex items-center justify-between bg-slate-50">
              <h3 className="font-semibold text-slate-900">Audit Event Details</h3>
              <button
                onClick={() => setSelectedLog(null)}
                className="text-slate-400 hover:text-slate-600 text-lg leading-none"
              >
                &times;
              </button>
            </div>
            <div className="p-6 overflow-y-auto space-y-4 text-sm font-sans">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-xs font-semibold text-slate-400 uppercase">Event ID</label>
                  <p className="font-mono text-xs text-slate-800 break-all">{selectedLog.id}</p>
                </div>
                <div>
                  <label className="text-xs font-semibold text-slate-400 uppercase">Action</label>
                  <p className="font-semibold text-slate-900">{selectedLog.action}</p>
                </div>
                <div>
                  <label className="text-xs font-semibold text-slate-400 uppercase">
                    Resource Type
                  </label>
                  <p className="text-slate-800">{selectedLog.target_type}</p>
                </div>
                <div>
                  <label className="text-xs font-semibold text-slate-400 uppercase">
                    Target ID
                  </label>
                  <p className="font-mono text-xs text-indigo-600 break-all">
                    {selectedLog.target_id}
                  </p>
                </div>
                <div>
                  <label className="text-xs font-semibold text-slate-400 uppercase">Actor ID</label>
                  <p className="font-mono text-xs text-slate-800 break-all">
                    {selectedLog.actor_id || 'System'}
                  </p>
                </div>
                <div>
                  <label className="text-xs font-semibold text-slate-400 uppercase">
                    Timestamp
                  </label>
                  <p className="text-slate-800">
                    {new Date(selectedLog.created_at).toLocaleString()}
                  </p>
                </div>
              </div>

              {(selectedLog.previous_state || selectedLog.new_state) && (
                <div className="pt-2 border-t border-slate-100 space-y-2">
                  <label className="text-xs font-semibold text-slate-400 uppercase">
                    State Payloads
                  </label>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3 font-mono text-xs">
                    {selectedLog.previous_state && (
                      <div className="bg-slate-50 p-3 rounded-lg border border-slate-200 overflow-x-auto">
                        <span className="text-xs font-bold text-slate-500 block mb-1">
                          Previous State
                        </span>
                        <pre className="text-slate-700">
                          {JSON.stringify(selectedLog.previous_state, null, 2)}
                        </pre>
                      </div>
                    )}
                    {selectedLog.new_state && (
                      <div className="bg-slate-50 p-3 rounded-lg border border-slate-200 overflow-x-auto">
                        <span className="text-xs font-bold text-slate-500 block mb-1">
                          New State
                        </span>
                        <pre className="text-slate-700">
                          {JSON.stringify(selectedLog.new_state, null, 2)}
                        </pre>
                      </div>
                    )}
                  </div>
                </div>
              )}
            </div>
            <div className="px-6 py-3 border-t border-slate-100 bg-slate-50 flex justify-end">
              <button
                onClick={() => setSelectedLog(null)}
                className="px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-lg hover:bg-slate-50"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
