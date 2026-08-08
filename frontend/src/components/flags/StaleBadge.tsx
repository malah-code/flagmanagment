import React from 'react';
import type { LifecycleState } from '../../types';
import { AlertTriangle, Archive, AlertCircle, CheckCircle } from 'lucide-react';

interface StaleBadgeProps {
  state?: LifecycleState;
}

export const StaleBadge: React.FC<StaleBadgeProps> = ({ state = 'ACTIVE' }) => {
  switch (state) {
    case 'STALE':
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-800 border border-amber-200">
          <AlertTriangle className="w-3 h-3 text-amber-600" />
          Stale
        </span>
      );
    case 'DEPRECATED':
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-rose-100 text-rose-800 border border-rose-200">
          <AlertCircle className="w-3 h-3 text-rose-600" />
          Deprecated
        </span>
      );
    case 'ARCHIVED':
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-600 border border-slate-200">
          <Archive className="w-3 h-3 text-slate-500" />
          Archived
        </span>
      );
    case 'ACTIVE':
    default:
      return (
        <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700 border border-emerald-200">
          <CheckCircle className="w-3 h-3 text-emerald-500" />
          Active
        </span>
      );
  }
};
