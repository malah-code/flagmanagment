import React, { useEffect, useState } from 'react';
import type { ChangeRequest } from '../types';
import { changeRequestApi } from '../services/changeRequestApi';
import { ChangeRequestDiff } from '../components/ChangeRequestDiff';

interface ChangeRequestsPageProps {
  environmentId: string;
}

export const ChangeRequestsPage: React.FC<ChangeRequestsPageProps> = ({ environmentId }) => {
  const [requests, setRequests] = useState<ChangeRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [actionError, setActionError] = useState<string | null>(null);
  const [rejectReason, setRejectReason] = useState<string>('');
  const [selectedRejectId, setSelectedRejectId] = useState<string | null>(null);

  const fetchRequests = async () => {
    try {
      setLoading(true);
      const data = await changeRequestApi.listByEnvironment(environmentId);
      setRequests(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to load change requests');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (environmentId) {
      fetchRequests();
    }
  }, [environmentId]);

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

  if (loading) return <div className="p-4">Loading change requests...</div>;
  if (error) return <div className="p-4 text-red-600">{error}</div>;

  return (
    <div className="p-6 max-w-5xl mx-auto">
      <h2 className="text-2xl font-bold mb-4">Change Requests</h2>

      {actionError && (
        <div className="p-3 mb-4 text-sm text-red-800 bg-red-100 rounded-lg">
          {actionError}
        </div>
      )}

      {requests.length === 0 ? (
        <p className="text-gray-500">No change requests found for this environment.</p>
      ) : (
        <div className="space-y-6">
          {requests.map((req) => (
            <div key={req.id} className="border rounded-lg p-5 shadow-sm bg-white dark:bg-gray-900">
              <div className="flex justify-between items-start mb-2">
                <div>
                  <h3 className="font-semibold text-lg">{req.title || 'Proposed Change'}</h3>
                  <p className="text-sm text-gray-500">
                    Created by {req.createdBy} on {new Date(req.createdAt).toLocaleString()}
                  </p>
                </div>
                <span
                  className={`px-2.5 py-1 text-xs font-semibold rounded-full ${
                    req.status === 'APPLIED' || req.status === 'APPROVED'
                      ? 'bg-green-100 text-green-800'
                      : req.status === 'REJECTED'
                      ? 'bg-red-100 text-red-800'
                      : 'bg-yellow-100 text-yellow-800'
                  }`}
                >
                  {req.status}
                </span>
              </div>

              <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">{req.description}</p>

              <ChangeRequestDiff proposedChanges={req.proposedChanges} />

              {req.status === 'PENDING' && (
                <div className="flex gap-3 mt-4">
                  <button
                    onClick={() => handleApprove(req.id)}
                    className="px-4 py-2 text-sm font-medium text-white bg-green-600 hover:bg-green-700 rounded-md shadow-sm"
                  >
                    Approve & Apply
                  </button>
                  <button
                    onClick={() => setSelectedRejectId(req.id)}
                    className="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-md shadow-sm"
                  >
                    Reject
                  </button>
                </div>
              )}

              {selectedRejectId === req.id && (
                <div className="mt-4 p-4 border rounded bg-gray-50 dark:bg-gray-800">
                  <label className="block text-sm font-medium mb-1">Reason for Rejection</label>
                  <input
                    type="text"
                    value={rejectReason}
                    onChange={(e) => setRejectReason(e.target.value)}
                    placeholder="Enter reason..."
                    className="w-full p-2 border rounded mb-3 text-sm"
                  />
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleReject(req.id)}
                      className="px-3 py-1.5 text-xs font-medium text-white bg-red-600 hover:bg-red-700 rounded"
                    >
                      Confirm Rejection
                    </button>
                    <button
                      onClick={() => setSelectedRejectId(null)}
                      className="px-3 py-1.5 text-xs font-medium text-gray-700 bg-gray-200 hover:bg-gray-300 rounded"
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
