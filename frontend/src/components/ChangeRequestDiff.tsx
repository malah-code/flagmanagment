import React from 'react';
import ReactDiffViewer from 'react-diff-viewer-continued';

interface ChangeRequestDiffProps {
  proposedChanges: Record<string, any>;
  currentState?: Record<string, any>;
}

export const ChangeRequestDiff: React.FC<ChangeRequestDiffProps> = ({
  proposedChanges,
  currentState = {},
}) => {
  const oldCode = JSON.stringify(currentState, null, 2);
  const newCode = JSON.stringify(proposedChanges, null, 2);

  return (
    <div className="my-4 border rounded shadow-sm overflow-hidden text-sm font-mono">
      <ReactDiffViewer 
        oldValue={oldCode} 
        newValue={newCode} 
        splitView={true}
        useDarkTheme={false}
        leftTitle="Current Configuration"
        rightTitle="Proposed Configuration"
      />
    </div>
  );
};
