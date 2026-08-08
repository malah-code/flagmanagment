import React, { useState } from 'react';
import { transitionLifecycle } from '../../services/api';
import type { LifecycleState } from '../../types';
import { Archive, RotateCcw, AlertOctagon, MoreVertical } from 'lucide-react';

interface FlagActionsProps {
  envId?: string;
  flagId: string;
  currentLifecycle?: LifecycleState;
  onStateChanged?: () => void;
}

export const FlagActions: React.FC<FlagActionsProps> = ({
  envId,
  flagId,
  currentLifecycle = 'ACTIVE',
  onStateChanged,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  if (!envId) return null;

  const handleAction = async (action: 'ARCHIVE' | 'DEPRECATE' | 'RESTORE' | 'MARK_STALE') => {
    setLoading(true);
    try {
      await transitionLifecycle(envId, flagId, action);
      setIsOpen(false);
      if (onStateChanged) onStateChanged();
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative inline-block text-left">
      <button
        onClick={() => setIsOpen(!isOpen)}
        disabled={loading}
        className="p-1 rounded hover:bg-slate-100 text-slate-500 hover:text-slate-700 transition-colors"
        title="Lifecycle Actions"
      >
        <MoreVertical className="w-4 h-4" />
      </button>

      {isOpen && (
        <div className="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-20 py-1 divide-y divide-slate-100">
          {currentLifecycle === 'ARCHIVED' ? (
            <button
              onClick={() => handleAction('RESTORE')}
              className="w-full text-left px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 flex items-center gap-2"
            >
              <RotateCcw className="w-4 h-4 text-emerald-600" />
              Restore Flag
            </button>
          ) : (
            <>
              {currentLifecycle !== 'DEPRECATED' && (
                <button
                  onClick={() => handleAction('DEPRECATE')}
                  className="w-full text-left px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 flex items-center gap-2"
                >
                  <AlertOctagon className="w-4 h-4 text-amber-600" />
                  Mark Deprecated
                </button>
              )}
              <button
                onClick={() => handleAction('ARCHIVE')}
                className="w-full text-left px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 flex items-center gap-2"
              >
                <Archive className="w-4 h-4 text-slate-600" />
                Archive Flag
              </button>
            </>
          )}
        </div>
      )}
    </div>
  );
};
