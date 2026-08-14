import React, { useState, useEffect } from 'react';
import { flagService } from '../../services/flags';
import { environmentService } from '../../services/environments';
import type { Environment } from '../../types';

interface PromoteFlagModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  flagId: string;
  sourceEnvId: string;
  sourceEnvName: string;
  onSuccess: () => void;
}

export const PromoteFlagModal: React.FC<PromoteFlagModalProps> = ({
  isOpen,
  onClose,
  projectId,
  flagId,
  sourceEnvId,
  sourceEnvName,
  onSuccess,
}) => {
  const [environments, setEnvironments] = useState<Environment[]>([]);
  const [targetEnvId, setTargetEnvId] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      environmentService.getByProject(projectId)
        .then((envs) => {
          setEnvironments(envs.filter(env => env.id !== sourceEnvId));
          if (envs.length > 0) {
            const defaultTarget = envs.find(env => env.id !== sourceEnvId);
            if (defaultTarget) setTargetEnvId(defaultTarget.id);
          }
        })
        .catch(() => setError('Failed to load environments'));
    }
  }, [isOpen, projectId, sourceEnvId]);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!targetEnvId) return;

    setIsSubmitting(true);
    setError(null);

    try {
      await flagService.promote(projectId, flagId, sourceEnvId, targetEnvId);
      onSuccess();
      onClose();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to promote flag');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center">
      <div className="relative bg-white rounded-lg shadow-xl max-w-md w-full p-6">
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-medium text-gray-900">Promote Flag</h3>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-500 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 rounded"
          >
            <span className="sr-only">Close</span>
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Source Environment
            </label>
            <input
              type="text"
              disabled
              value={sourceEnvName}
              className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 bg-gray-50 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
            />
          </div>

          <div className="mb-4">
            <label htmlFor="targetEnv" className="block text-sm font-medium text-gray-700 mb-1">
              Target Environment
            </label>
            <select
              id="targetEnv"
              value={targetEnvId}
              onChange={(e) => setTargetEnvId(e.target.value)}
              className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
              required
            >
              <option value="" disabled>Select an environment</option>
              {environments.map((env) => (
                <option key={env.id} value={env.id}>
                  {env.name} {env.isProtected ? '(Protected)' : ''}
                </option>
              ))}
            </select>
          </div>
          
          <div className="mt-2 mb-4 text-sm text-gray-500">
            {environments.find(e => e.id === targetEnvId)?.isProtected 
              ? "This environment is protected. A Change Request will be generated for approval."
              : "This will instantly update the flag configuration in the target environment."}
          </div>

          <div className="flex justify-end gap-3 mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !targetEnvId}
              className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              {isSubmitting ? 'Promoting...' : 'Promote'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
