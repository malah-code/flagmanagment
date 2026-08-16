import React from 'react';
import type { ScheduledChange } from '../../types/scheduledChange';

interface ScheduledChangeBadgeProps {
  scheduledChange: ScheduledChange | null | undefined;
}

export const ScheduledChangeBadge: React.FC<ScheduledChangeBadgeProps> = ({ scheduledChange }) => {
  if (!scheduledChange || scheduledChange.status !== 'PENDING') {
    return null;
  }

  const localTime = new Date(scheduledChange.scheduled_for).toLocaleString();

  return (
    <span
      title={`UTC: ${scheduledChange.scheduled_for}`}
      className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium bg-amber-100 text-amber-800 border border-amber-300 dark:bg-amber-900/40 dark:text-amber-300 dark:border-amber-700"
    >
      <span role="img" aria-label="clock">
        ⏰
      </span>
      <span>
        Scheduled: <strong>{scheduledChange.action}</strong> @ {localTime}
      </span>
    </span>
  );
};
