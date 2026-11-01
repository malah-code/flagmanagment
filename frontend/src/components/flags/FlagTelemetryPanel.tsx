import React, { useState, useEffect } from 'react';
import {
  Activity,
  Zap,
  Clock,
  CheckCircle2,
  Play,
  Pause,
  RefreshCw,
  BarChart3,
  Server,
  Layers,
} from 'lucide-react';

export interface FlagTelemetryPanelProps {
  flagKey: string;
  flagType: string;
  currentEnvName?: string;
  variations?: Array<{ id: string; name: string; value: any }>;
}

interface EvaluationEvent {
  id: string;
  timestamp: string;
  contextKey: string;
  variation: string;
  latencyUs: number;
  reason: 'TARGETING_MATCH' | 'PERCENTAGE_ROLLOUT' | 'DEFAULT_FALLBACK' | 'OVERRIDE';
}

export const FlagTelemetryPanel: React.FC<FlagTelemetryPanelProps> = ({
  flagKey,
  flagType,
  currentEnvName = 'Production',
  variations = [],
}) => {
  const [timeRange, setTimeRange] = useState<'1h' | '24h' | '7d'>('24h');
  const [isLiveStreaming, setIsLiveStreaming] = useState<boolean>(true);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);

  // Simulated metrics based on time range
  const metrics = {
    '1h': { total: '48,290', rps: '13.4/s', p99: '8.2 µs', cacheHit: '99.99%', trend: '+4.2%' },
    '24h': {
      total: '1,158,400',
      rps: '14.8/s',
      p99: '8.6 µs',
      cacheHit: '99.98%',
      trend: '+12.5%',
    },
    '7d': { total: '8,420,150', rps: '13.9/s', p99: '8.9 µs', cacheHit: '99.96%', trend: '+18.1%' },
  }[timeRange];

  // Default variation breakdown if no specific variations configured
  const distributionData =
    variations.length >= 2
      ? [
          { name: variations[0].name, percentage: 56.4, count: '653,337', color: 'bg-indigo-600' },
          { name: variations[1].name, percentage: 43.6, count: '505,063', color: 'bg-emerald-500' },
        ]
      : [
          { name: 'Treatment (true)', percentage: 52.8, count: '611,635', color: 'bg-indigo-600' },
          { name: 'Control (false)', percentage: 47.2, count: '546,765', color: 'bg-slate-400' },
        ];

  // Live evaluation stream data
  const [events, setEvents] = useState<EvaluationEvent[]>([
    {
      id: 'evt-1',
      timestamp: 'Just now',
      contextKey: 'usr_89f01ab2',
      variation: 'Treatment (true)',
      latencyUs: 7.8,
      reason: 'PERCENTAGE_ROLLOUT',
    },
    {
      id: 'evt-2',
      timestamp: '2s ago',
      contextKey: 'usr_33c84be9',
      variation: 'Control (false)',
      latencyUs: 6.4,
      reason: 'DEFAULT_FALLBACK',
    },
    {
      id: 'evt-3',
      timestamp: '5s ago',
      contextKey: 'usr_premium_9921',
      variation: 'Treatment (true)',
      latencyUs: 8.9,
      reason: 'TARGETING_MATCH',
    },
    {
      id: 'evt-4',
      timestamp: '9s ago',
      contextKey: 'usr_31fa9401',
      variation: 'Treatment (true)',
      latencyUs: 7.2,
      reason: 'PERCENTAGE_ROLLOUT',
    },
    {
      id: 'evt-5',
      timestamp: '14s ago',
      contextKey: 'usr_02b9e844',
      variation: 'Control (false)',
      latencyUs: 6.9,
      reason: 'DEFAULT_FALLBACK',
    },
  ]);

  // Periodic heartbeat when live streaming is enabled
  useEffect(() => {
    if (!isLiveStreaming) return;
    const interval = setInterval(() => {
      const now = new Date();
      const timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`;
      const randomId = Math.random().toString(36).substring(2, 10);
      const isTreatment = Math.random() > 0.45;
      const newEvent: EvaluationEvent = {
        id: `evt-${Date.now()}`,
        timestamp: timeStr,
        contextKey: `usr_${randomId}`,
        variation: isTreatment ? 'Treatment (true)' : 'Control (false)',
        latencyUs: parseFloat((6.0 + Math.random() * 4.5).toFixed(1)),
        reason: isTreatment ? 'PERCENTAGE_ROLLOUT' : 'DEFAULT_FALLBACK',
      };
      setEvents((prev) => [newEvent, ...prev.slice(0, 7)]);
    }, 3500);
    return () => clearInterval(interval);
  }, [isLiveStreaming]);

  const handleManualRefresh = () => {
    setIsRefreshing(true);
    setTimeout(() => {
      setIsRefreshing(false);
    }, 600);
  };

  return (
    <div className="space-y-6">
      {/* Top Controls Bar */}
      <div className="bg-white border border-slate-200 rounded-xl p-4 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2">
            <span
              className={`inline-block w-2.5 h-2.5 rounded-full ${
                isLiveStreaming ? 'bg-emerald-500 animate-pulse' : 'bg-slate-400'
              }`}
            />
            <span className="text-sm font-semibold text-slate-900">
              {isLiveStreaming ? 'Live Stream Active' : 'Live Stream Paused'}
            </span>
          </div>
          <span className="text-slate-300">|</span>
          <span className="text-xs text-slate-500">
            Environment: <strong className="text-slate-800 font-semibold">{currentEnvName}</strong>
          </span>
        </div>

        <div className="flex items-center gap-2">
          {/* Time Range Selector */}
          <div className="inline-flex rounded-lg border border-slate-200 p-0.5 bg-slate-50 text-xs font-medium text-slate-600">
            {(['1h', '24h', '7d'] as const).map((r) => (
              <button
                key={r}
                onClick={() => setTimeRange(r)}
                className={`px-3 py-1 rounded-md transition-all ${
                  timeRange === r
                    ? 'bg-white text-indigo-600 shadow-sm font-semibold'
                    : 'hover:text-slate-900'
                }`}
              >
                {r.toUpperCase()}
              </button>
            ))}
          </div>

          {/* Pause / Play Live Stream */}
          <button
            onClick={() => setIsLiveStreaming(!isLiveStreaming)}
            className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-slate-200 text-xs font-medium text-slate-700 hover:bg-slate-50 transition-colors"
          >
            {isLiveStreaming ? (
              <>
                <Pause className="w-3.5 h-3.5 text-slate-500" /> Pause
              </>
            ) : (
              <>
                <Play className="w-3.5 h-3.5 text-emerald-600" /> Resume
              </>
            )}
          </button>

          {/* Manual Refresh */}
          <button
            onClick={handleManualRefresh}
            aria-label="Refresh telemetry"
            className="p-1.5 rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 transition-colors"
          >
            <RefreshCw
              className={`w-4 h-4 ${isRefreshing ? 'animate-spin text-indigo-600' : ''}`}
            />
          </button>
        </div>
      </div>

      {/* KPI Stat Cards Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Total Evaluations */}
        <div className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
          <div className="flex items-center justify-between text-slate-500 text-xs font-semibold uppercase tracking-wider mb-2">
            <span>Total Evaluations</span>
            <Activity className="w-4 h-4 text-indigo-500" />
          </div>
          <div className="text-2xl font-bold text-slate-900 font-mono">{metrics.total}</div>
          <div className="mt-2 text-xs text-emerald-600 font-medium flex items-center gap-1">
            <span>{metrics.trend}</span> vs prior period
          </div>
        </div>

        {/* Real-time Throughput */}
        <div className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
          <div className="flex items-center justify-between text-slate-500 text-xs font-semibold uppercase tracking-wider mb-2">
            <span>Evaluation Rate</span>
            <Zap className="w-4 h-4 text-amber-500" />
          </div>
          <div className="text-2xl font-bold text-slate-900 font-mono">{metrics.rps}</div>
          <div className="mt-2 text-xs text-slate-500 flex items-center gap-1">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" /> Steady load profile
          </div>
        </div>

        {/* Evaluation Latency (Sub-millisecond) */}
        <div className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
          <div className="flex items-center justify-between text-slate-500 text-xs font-semibold uppercase tracking-wider mb-2">
            <span>P99 Eval Latency</span>
            <Clock className="w-4 h-4 text-blue-500" />
          </div>
          <div className="text-2xl font-bold text-emerald-600 font-mono">{metrics.p99}</div>
          <div className="mt-2 text-xs text-slate-500">Target &lt; 1,000 µs (1.0 ms)</div>
        </div>

        {/* Edge Cache Hit Ratio */}
        <div className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
          <div className="flex items-center justify-between text-slate-500 text-xs font-semibold uppercase tracking-wider mb-2">
            <span>Cache Hit Ratio</span>
            <CheckCircle2 className="w-4 h-4 text-emerald-500" />
          </div>
          <div className="text-2xl font-bold text-indigo-600 font-mono">{metrics.cacheHit}</div>
          <div className="mt-2 text-xs text-slate-500">In-memory evaluation cache</div>
        </div>
      </div>

      {/* Variation Split & Distribution Chart */}
      <div className="bg-white border border-slate-200 rounded-xl p-6 shadow-sm">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-2">
            <BarChart3 className="w-5 h-5 text-indigo-600" />
            <h3 className="text-base font-semibold text-slate-900">
              Variation Traffic Distribution
            </h3>
          </div>
          <span className="text-xs text-slate-500">
            Aggregated across active evaluation contexts
          </span>
        </div>

        {/* Stacked Progress Bar */}
        <div className="w-full h-4 bg-slate-100 rounded-full overflow-hidden flex mb-4">
          {distributionData.map((d, idx) => (
            <div
              key={idx}
              style={{ width: `${d.percentage}%` }}
              className={`${d.color} transition-all duration-500`}
              title={`${d.name}: ${d.percentage}%`}
            />
          ))}
        </div>

        {/* Legend & Exact Percentages */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {distributionData.map((d, idx) => (
            <div
              key={idx}
              className="flex items-center justify-between p-3.5 bg-slate-50 rounded-lg border border-slate-200"
            >
              <div className="flex items-center gap-2.5">
                <span className={`w-3 h-3 rounded-full ${d.color}`} />
                <span className="text-sm font-medium text-slate-800">{d.name}</span>
              </div>
              <div className="text-right">
                <div className="text-sm font-bold text-slate-900">{d.percentage}%</div>
                <div className="text-xs text-slate-500 font-mono">{d.count} calls</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Live Stream Table */}
      <div className="bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden">
        <div className="p-5 border-b border-slate-200 flex items-center justify-between bg-slate-50">
          <div>
            <h3 className="text-sm font-semibold text-slate-900 flex items-center gap-2">
              <Server className="w-4 h-4 text-slate-600" />
              Recent Evaluation Stream
            </h3>
            <p className="text-xs text-slate-500 mt-0.5">
              Live audit events captured across edge and microservice SDK clients
            </p>
          </div>
          <span className="text-xs font-mono px-2 py-1 bg-white border border-slate-200 rounded text-slate-600">
            Flag: {flagKey} ({flagType})
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-100/75 text-slate-600 border-b border-slate-200">
              <tr>
                <th className="py-2.5 px-4 font-semibold">Timestamp</th>
                <th className="py-2.5 px-4 font-semibold">Target Context</th>
                <th className="py-2.5 px-4 font-semibold">Variation Served</th>
                <th className="py-2.5 px-4 font-semibold">Reason</th>
                <th className="py-2.5 px-4 font-semibold text-right">Latency</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {events.map((evt) => (
                <tr key={evt.id} className="hover:bg-slate-50/75 transition-colors">
                  <td className="py-2.5 px-4 font-mono text-slate-500 whitespace-nowrap">
                    {evt.timestamp}
                  </td>
                  <td className="py-2.5 px-4 font-mono text-slate-800 font-medium whitespace-nowrap">
                    {evt.contextKey}
                  </td>
                  <td className="py-2.5 px-4 whitespace-nowrap">
                    <span
                      className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold ${
                        evt.variation.includes('Treatment') || evt.variation.includes('true')
                          ? 'bg-indigo-50 text-indigo-700 border border-indigo-200'
                          : 'bg-slate-100 text-slate-700 border border-slate-200'
                      }`}
                    >
                      {evt.variation}
                    </span>
                  </td>
                  <td className="py-2.5 px-4 text-slate-600 whitespace-nowrap">
                    <span className="inline-flex items-center gap-1 font-mono text-[11px] bg-slate-100 px-1.5 py-0.5 rounded">
                      <Layers className="w-3 h-3 text-slate-400" />
                      {evt.reason}
                    </span>
                  </td>
                  <td className="py-2.5 px-4 text-right font-mono font-medium text-emerald-600 whitespace-nowrap">
                    {evt.latencyUs} µs
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
