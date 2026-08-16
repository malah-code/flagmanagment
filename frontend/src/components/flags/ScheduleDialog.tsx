import React, { useState, useEffect } from 'react';
import type { ScheduledChange, ScheduledChangeAction } from '../../types/scheduledChange';
import { scheduledChangesApi } from '../../services/scheduledChangesApi';

interface ScheduleDialogProps {
  isOpen: boolean;
  onClose: () => void;
  flagId: string;
  flagName: string;
  environmentId: string;
  existingSchedule: ScheduledChange | null | undefined;
  onSuccess: () => void;
}

export const ScheduleDialog: React.FC<ScheduleDialogProps> = ({
  isOpen,
  onClose,
  flagId,
  flagName,
  environmentId,
  existingSchedule,
  onSuccess,
}) => {
  const [scheduledTime, setScheduledTime] = useState('');
  const [action, setAction] = useState<ScheduledChangeAction>('ENABLE');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isModifying, setIsModifying] = useState(false);
  const [isConfirmingCancel, setIsConfirmingCancel] = useState(false);

  useEffect(() => {
    if (existingSchedule && existingSchedule.status === 'PENDING') {
      const localIso = new Date(existingSchedule.scheduled_for).toISOString().slice(0, 16);
      setScheduledTime(localIso);
      setAction(existingSchedule.action);
    } else {
      const nextHour = new Date(Date.now() + 3600 * 1000);
      const offsetMs = nextHour.getTimezoneOffset() * 60000;
      const localNextHour = new Date(nextHour.getTime() - offsetMs);
      setScheduledTime(localNextHour.toISOString().slice(0, 16));
      setAction('ENABLE');
    }
    setError(null);
    setIsModifying(false);
    setIsConfirmingCancel(false);
  }, [existingSchedule, isOpen]);

  if (!isOpen) return null;

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const utcString = new Date(scheduledTime).toISOString();
      await scheduledChangesApi.create(environmentId, {
        target_type: 'FLAG',
        target_id: flagId,
        action,
        scheduled_for: utcString,
      });
      onSuccess();
      onClose();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to create scheduled change');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!existingSchedule) return;
    setError(null);
    setLoading(true);
    try {
      const utcString = new Date(scheduledTime).toISOString();
      await scheduledChangesApi.update(existingSchedule.id, { scheduled_for: utcString });
      onSuccess();
      onClose();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to update scheduled change');
    } finally {
      setLoading(false);
    }
  };

  const handleCancelConfirmed = async () => {
    if (!existingSchedule) return;
    setError(null);
    setLoading(true);
    try {
      await scheduledChangesApi.cancel(existingSchedule.id);
      onSuccess();
      onClose();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to cancel scheduled change');
    } finally {
      setLoading(false);
      setIsConfirmingCancel(false);
    }
  };

  const hasPending = existingSchedule && existingSchedule.status === 'PENDING';

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
      <div className="w-full max-w-md bg-white dark:bg-gray-800 rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            Schedule Flag Change: {flagName}
          </h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
          >
            ✕
          </button>
        </div>

        <div className="p-6">
          {error && (
            <div className="mb-4 p-3 bg-red-100 border border-red-300 text-red-700 rounded text-sm dark:bg-red-900/30 dark:border-red-800 dark:text-red-300">
              {error}
            </div>
          )}

          {hasPending && !isModifying ? (
            <div className="space-y-4">
              <div className="p-4 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-md">
                <p className="text-sm font-medium text-amber-900 dark:text-amber-200">
                  Current Schedule:
                </p>
                <p className="text-sm text-amber-800 dark:text-amber-300 mt-1">
                  Action: <strong>{existingSchedule.action}</strong>
                </p>
                <p className="text-sm text-amber-800 dark:text-amber-300">
                  Target Time:{' '}
                  <strong>{new Date(existingSchedule.scheduled_for).toLocaleString()}</strong>
                </p>
              </div>

              {isConfirmingCancel ? (
                /* Inline confirmation — avoids window.confirm() which can be blocked */
                <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md">
                  <p className="text-sm text-red-700 dark:text-red-300 mb-3">
                    Are you sure you want to cancel this scheduled change? This cannot be undone.
                  </p>
                  <div className="flex justify-end gap-3">
                    <button
                      type="button"
                      onClick={() => setIsConfirmingCancel(false)}
                      disabled={loading}
                      className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600 rounded-md disabled:opacity-50"
                    >
                      No, go back
                    </button>
                    <button
                      type="button"
                      onClick={handleCancelConfirmed}
                      disabled={loading}
                      className="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-md disabled:opacity-50"
                    >
                      {loading ? 'Cancelling...' : 'Yes, cancel it'}
                    </button>
                  </div>
                </div>
              ) : (
                <div className="flex justify-end gap-3 pt-2">
                  <button
                    type="button"
                    onClick={() => setIsModifying(true)}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600 rounded-md"
                  >
                    Reschedule
                  </button>
                  <button
                    type="button"
                    onClick={() => setIsConfirmingCancel(true)}
                    className="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-md"
                  >
                    Cancel Schedule
                  </button>
                </div>
              )}
            </div>
          ) : (
            <form onSubmit={hasPending ? handleUpdate : handleCreate} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Action
                </label>
                <select
                  value={action}
                  onChange={(e) => setAction(e.target.value as ScheduledChangeAction)}
                  disabled={hasPending || loading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                >
                  <option value="ENABLE">Turn ON (Enable)</option>
                  <option value="DISABLE">Turn OFF (Disable)</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Scheduled Time (Local Time)
                </label>
                <input
                  type="datetime-local"
                  required
                  value={scheduledTime}
                  onChange={(e) => setScheduledTime(e.target.value)}
                  disabled={loading}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                />
              </div>

              <div className="flex justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
                <button
                  type="button"
                  onClick={onClose}
                  className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-200 rounded-md"
                >
                  Close
                </button>
                <button
                  type="submit"
                  disabled={loading}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md disabled:opacity-50"
                >
                  {loading ? 'Saving...' : hasPending ? 'Update Schedule' : 'Schedule Action'}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>
    </div>
  );
};
