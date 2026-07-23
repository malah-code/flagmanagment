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
