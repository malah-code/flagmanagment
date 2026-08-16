export interface CheckResult {
  status: string;
  latency_ms?: number;
  error?: string;
}

export interface HealthResponse {
  status: string;
  version: string;
  uptime_seconds: number;
  checks: Record<string, CheckResult>;
}

const API_BASE = '/api';

export const fetchHealth = async (): Promise<HealthResponse> => {
  const response = await fetch(`${API_BASE}/healthz`);

  if (!response.ok && response.status !== 503) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  return response.json();
};

export const transitionLifecycle = async (
  envId: string,
  flagId: string,
  action: 'ARCHIVE' | 'DEPRECATE' | 'RESTORE' | 'MARK_STALE',
) => {
  const response = await fetch(`${API_BASE}/v1/environments/${envId}/flags/${flagId}/lifecycle`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ action }),
  });
  if (!response.ok) {
    throw new Error(`Failed to perform lifecycle transition: ${response.status}`);
  }
  return response.json();
};

export const getStalePolicy = async (projectId: string) => {
  const response = await fetch(`${API_BASE}/v1/projects/${projectId}/stale-policy`);
  if (!response.ok) {
    throw new Error(`Failed to fetch stale policy: ${response.status}`);
  }
  return response.json();
};

export const setStalePolicy = async (projectId: string, staleAfterDays: number) => {
  const response = await fetch(`${API_BASE}/v1/projects/${projectId}/stale-policy`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ stale_after_days: staleAfterDays }),
  });
  if (!response.ok) {
    throw new Error(`Failed to update stale policy: ${response.status}`);
  }
  return response.json();
};
