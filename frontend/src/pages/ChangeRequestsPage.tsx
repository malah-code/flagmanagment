import React, { useEffect, useState } from 'react';
import { GitPullRequest, Check, X, Clock, AlertCircle, Loader2, Filter } from 'lucide-react';
import type { ChangeRequest, Environment } from '../types';
import { changeRequestApi } from '../services/changeRequestApi';
import { ChangeRequestDiff } from '../components/ChangeRequestDiff';

interface ChangeRequestsPageProps {
  environments?: Environment[];
  environmentId?: string;
}

export const ChangeRequestsPage: React.FC<ChangeRequestsPageProps> = ({
  environments = [],
  environmentId: initialEnvId,
}) => {
  const [selectedEnvId, setSelectedEnvId] = useState<string>(
    initialEnvId || environments[0]?.id || '',
  );
  const [requests, setRequests] = useState<ChangeRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [rejectReason, setRejectReason] = useState<string>('');
  const [selectedRejectId, setSelectedRejectId] = useState<string | null>(null);
  const [statusFilter, setStatusFilter] = useState<string>('ALL');

  useEffect(() => {
    if (!selectedEnvId && environments.length > 0) {
      setSelectedEnvId(environments[0].id);
    }
  }, [environments, selectedEnvId]);

  const fetchRequests = async () => {
    if (!selectedEnvId) {
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      const data = await changeRequestApi.listByEnvironment(selectedEnvId);
      setRequests(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to load change requests');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRequests();
  }, [selectedEnvId]);

  const handleApprove = async (id: string) => {
    setActionError(null);
    try {
      await changeRequestApi.approve(id);
      await fetchRequests();
    } catch (err: any) {
      setActionError(err.message || 'Failed to approve change request');
    }
  };

  const handleReject = async (id: string) => {
    setActionError(null);
    try {
      await changeRequestApi.reject(id, rejectReason);
      setSelectedRejectId(null);
      setRejectReason('');
      await fetchRequests();
    } catch (err: any) {
      setActionError(err.message || 'Failed to reject change request');
    }
  };

  const filteredRequests = requests.filter((req) => {
    if (statusFilter === 'ALL') return true;
    return req.status === statusFilter;
  });

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'APPLIED':
      case 'APPROVED':
        return 'bg-emerald-50 text-emerald-700 border-emerald-200';
      case 'REJECTED':
        return 'bg-rose-50 text-rose-700 border-rose-200';
      default:
        return 'bg-amber-50 text-amber-700 border-amber-200';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900 flex items-center gap-2">
            <GitPullRequest className="w-6 h-6 text-indigo-600" />
            Change Requests
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Review, approve, and promote proposed configuration changes
          </p>
        </div>

        {/* Environment Selector Dropdown */}
        {environments.length > 0 && (
          <div className="flex items-center gap-2">
            <label className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Environment:
            </label>
            <select
              value={selectedEnvId}
              onChange={(e) => setSelectedEnvId(e.target.value)}
              className="px-3 py-1.5 text-sm font-medium bg-white border border-slate-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              {environments.map((env) => (
                <option key={env.id} value={env.id}>
                  {env.name}
                </option>
              ))}
            </select>
          </div>
        )}
      </div>

      {actionError && (
        <div className="p-4 text-sm text-rose-800 bg-rose-50 border border-rose-200 rounded-xl flex items-center gap-2">
          <AlertCircle className="w-4 h-4 text-rose-600 flex-shrink-0" />
          <span>{actionError}</span>
        </div>
      )}

      {/* Filter Bar */}
      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-4 flex items-center justify-between gap-4">
        <div className="flex items-center gap-2">
          <Filter className="w-4 h-4 text-slate-400" />
          <span className="text-xs font-semibold text-slate-500 uppercase">Filter by status:</span>
          <div className="flex gap-1">
            {['ALL', 'PENDING', 'APPROVED', 'REJECTED', 'APPLIED'].map((st) => (
              <button
                key={st}
                onClick={() => setStatusFilter(st)}
                className={`px-2.5 py-1 text-xs font-medium rounded-lg transition-colors ${
                  statusFilter === st
                    ? 'bg-indigo-600 text-white'
                    : 'text-slate-600 hover:bg-slate-100'
                }`}
              >
                {st}
              </button>
            ))}
          </div>
        </div>
        <div className="text-xs text-slate-500 font-medium">
          Showing {filteredRequests.length} of {requests.length} requests
        </div>
      </div>

      {/* Content Area */}
      {loading ? (
        <div className="flex flex-col items-center justify-center p-12 bg-white rounded-xl border border-slate-200 text-slate-400">
          <Loader2 className="w-8 h-8 animate-spin text-indigo-600 mb-2" />
          <p className="text-sm font-medium">Loading change requests...</p>
        </div>
      ) : error ? (
        <div className="p-8 text-center bg-rose-50 rounded-xl border border-rose-200 text-rose-700">
          <AlertCircle className="w-8 h-8 text-rose-500 mx-auto mb-2" />
          <p className="font-semibold text-sm">Failed to load change requests</p>
          <p className="text-xs mt-1 text-rose-600">{error}</p>
        </div>
      ) : filteredRequests.length === 0 ? (
        <div className="bg-white rounded-xl border border-slate-200 p-12 text-center text-slate-500">
          <Clock className="w-12 h-12 text-slate-300 mx-auto mb-3" />
          <p className="font-medium text-slate-700">No change requests found</p>
          <p className="text-xs text-slate-400 mt-1">
            Pending changes submitted for review in this environment will appear here.
          </p>
        </div>
      ) : (
        <div className="space-y-4">
          {filteredRequests.map((req) => (
            <div
              key={req.id}
              className="bg-white border border-slate-200 rounded-xl p-6 shadow-sm space-y-4 hover:border-slate-300 transition-colors"
            >
              <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-3">
                <div>
                  <h3 className="font-semibold text-lg text-slate-900">
                    {req.title || 'Proposed Flag Configuration Change'}
                  </h3>
                  <p className="text-xs text-slate-500 mt-0.5">
                    Submitted on {new Date(req.createdAt).toLocaleString()} by{' '}
                    <span className="font-mono text-slate-700">{req.createdBy || 'Unknown'}</span>
                  </p>
                </div>
                <span
                  className={`inline-flex items-center text-xs font-semibold px-2.5 py-1 rounded-full border ${getStatusBadge(
                    req.status,
                  )}`}
                >
                  {req.status}
                </span>
              </div>

              {req.description && <p className="text-sm text-slate-600">{req.description}</p>}

              {/* Visual Diff */}
              <div className="pt-2 border-t border-slate-100">
                <ChangeRequestDiff
                  proposedChanges={req.proposedChanges}
                  currentState={req.currentState}
                />
              </div>

              {/* Action Buttons for Pending */}
              {req.status === 'PENDING' && (
                <div className="flex gap-3 pt-3 border-t border-slate-100">
                  <button
                    onClick={() => handleApprove(req.id)}
                    className="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-white bg-emerald-600 hover:bg-emerald-700 rounded-lg shadow-sm transition-colors"
                  >
                    <Check className="w-4 h-4" />
                    Approve & Apply
                  </button>
                  <button
                    onClick={() => setSelectedRejectId(req.id)}
                    className="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-rose-700 bg-rose-50 hover:bg-rose-100 border border-rose-200 rounded-lg transition-colors"
                  >
                    <X className="w-4 h-4" />
                    Reject
                  </button>
                </div>
              )}

              {/* Reject Reason Form */}
              {selectedRejectId === req.id && (
                <div className="mt-4 p-4 border border-rose-200 rounded-xl bg-rose-50/50 space-y-3">
                  <label className="block text-xs font-semibold text-rose-900 uppercase">
                    Reason for Rejection:
                  </label>
                  <input
                    type="text"
                    value={rejectReason}
                    onChange={(e) => setRejectReason(e.target.value)}
                    placeholder="Provide context on why this change is being rejected..."
                    className="w-full px-3 py-2 text-sm bg-white border border-rose-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-rose-500 focus:border-rose-500"
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleReject(req.id)}
                      className="px-3 py-1.5 text-xs font-medium text-white bg-rose-600 hover:bg-rose-700 rounded-lg"
                    >
                      Confirm Rejection
                    </button>
                    <button
                      onClick={() => {
                        setSelectedRejectId(null);
                        setRejectReason('');
                      }}
                      className="px-3 py-1.5 text-xs font-medium text-slate-700 bg-white border border-slate-300 hover:bg-slate-50 rounded-lg"
                    >
                      Cancel
                    </button>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
