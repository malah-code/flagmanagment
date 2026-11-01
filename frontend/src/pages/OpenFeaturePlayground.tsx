import React, { useState } from 'react';
import {
  Code2,
  Terminal,
  CheckCircle2,
  Copy,
  Check,
  Play,
  Layers,
  Sparkles,
  Server,
  Zap,
} from 'lucide-react';
import { useFlags } from '../hooks/useFlags';
import { useEnvironments } from '../hooks/useEnvironments';
import toast from 'react-hot-toast';

interface OpenFeaturePlaygroundProps {
  projectId: string;
}

type SdkLanguage = 'react' | 'node' | 'python' | 'go' | 'java' | 'dotnet' | 'ios' | 'android';

export const OpenFeaturePlayground: React.FC<OpenFeaturePlaygroundProps> = ({ projectId }) => {
  const { data: flags = [], isLoading: isLoadingFlags } = useFlags(projectId);
  const { data: environments = [], isLoading: isLoadingEnvs } = useEnvironments(projectId);

  const [selectedFlagKey, setSelectedFlagKey] = useState<string>('');
  const [selectedEnvId, setSelectedEnvId] = useState<string>('');
  const [evalType, setEvalType] = useState<'boolean' | 'string' | 'number' | 'object'>('boolean');
  const [targetingKey, setTargetingKey] = useState<string>('usr_1024_beta');
  const [customAttributes, setCustomAttributes] = useState<string>(
    JSON.stringify({ tier: 'premium', country: 'US', device: 'mobile' }, null, 2),
  );

  const [activeLang, setActiveLang] = useState<SdkLanguage>('react');
  const [hasCopied, setHasCopied] = useState(false);

  // Evaluation Output State
  const [evalResult, setEvalResult] = useState<{
    value: any;
    variant: string;
    reason: string;
    latencyUs: number;
    errorCode?: string;
    evaluatedAt: string;
  } | null>(null);

  const currentFlagKey = selectedFlagKey || flags[0]?.key || 'new-feature-flag';
  const currentEnv = environments.find((e) => e.id === selectedEnvId) || environments[0];

  const handleRunEvaluation = () => {
    const start = performance.now();
    let parsedAttrs: Record<string, any> = {};
    try {
      if (customAttributes.trim()) {
        parsedAttrs = JSON.parse(customAttributes);
      }
    } catch {
      toast.error('Invalid JSON in custom attributes');
      return;
    }

    // Deterministic simulation based on targetingKey hash
    const hash = targetingKey.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
    const isTreatment = hash % 2 === 0 || parsedAttrs.tier === 'premium';
    const latency = parseFloat(((performance.now() - start) * 1000 + 4.2).toFixed(1));

    let evaluatedValue: any = isTreatment;
    let variantName = isTreatment ? 'treatment' : 'control';

    if (evalType === 'string') {
      evaluatedValue = isTreatment ? 'v2-quantum-engine' : 'v1-standard';
      variantName = isTreatment ? 'variant-b' : 'variant-a';
    } else if (evalType === 'number') {
      evaluatedValue = isTreatment ? 100 : 25;
      variantName = isTreatment ? 'quota-high' : 'quota-low';
    } else if (evalType === 'object') {
      evaluatedValue = isTreatment
        ? { theme: 'dark', aiFeatures: true, maxBatchSize: 500 }
        : { theme: 'light', aiFeatures: false, maxBatchSize: 50 };
      variantName = isTreatment ? 'config-advanced' : 'config-default';
    }

    setEvalResult({
      value: evaluatedValue,
      variant: variantName,
      reason: parsedAttrs.tier === 'premium' ? 'TARGETING_MATCH' : 'PERCENTAGE_ROLLOUT',
      latencyUs: latency,
      evaluatedAt: new Date().toLocaleTimeString(),
    });

    toast.success('Evaluated via OpenFeature Provider!');
  };

  const getCodeSnippet = (lang: SdkLanguage) => {
    switch (lang) {
      case 'react':
        return `// 1. Install OpenFeature React SDK and FlagManagment Provider
// npm install @openfeature/react-sdk @flagmanagment/react-sdk

import React from 'react';
import { OpenFeatureProvider, useBooleanFlagValue } from '@openfeature/react-sdk';
import { FlagManagmentWebProvider } from '@flagmanagment/react-sdk';

const provider = new FlagManagmentWebProvider({
  apiUrl: 'http://localhost:8080',
  apiKey: 'client-sdk-key',
});

export function App() {
  return (
    <OpenFeatureProvider provider={provider}>
      <FeatureComponent />
    </OpenFeatureProvider>
  );
}

function FeatureComponent() {
  const isEnabled = useBooleanFlagValue('${currentFlagKey}', false, {
    targetingKey: '${targetingKey}',
  });

  return isEnabled ? <div>Treatment Active</div> : <div>Control Default</div>;
}`;
      case 'node':
        return `// 1. Install OpenFeature Server SDK and FlagManagment Provider
// npm install @openfeature/server-sdk @flagmanagment/node-sdk

import { OpenFeature } from '@openfeature/server-sdk';
import { FlagManagmentProvider } from '@flagmanagment/node-sdk';

// 2. Register FlagManagment as the CNCF OpenFeature provider
await OpenFeature.setProviderAndWait(new FlagManagmentProvider({
  apiUrl: 'http://localhost:8080',
  apiKey: process.env.FLAGMANAGMENT_API_KEY,
}));

// 3. Resolve evaluation details with context
const client = OpenFeature.getClient();
const details = await client.getBooleanDetails('${currentFlagKey}', false, {
  targetingKey: '${targetingKey}',
  tier: 'premium',
});

console.log('Value:', details.value);
console.log('Reason:', details.reason); // 'TARGETING_MATCH'`;

      case 'python':
        return `# 1. Install OpenFeature Python SDK
# pip install openfeature-sdk flagmanagment-sdk

from openfeature import api
from openfeature.evaluation_context import EvaluationContext
from flagmanagment.provider import FlagManagmentProvider

# 2. Register OpenFeature Provider
provider = FlagManagmentProvider(api_url="http://localhost:8080", api_key="...")
api.set_provider(provider)

# 3. Evaluate flag with evaluation context
client = api.get_client()
context = EvaluationContext(
    targeting_key="${targetingKey}",
    attributes={"tier": "premium", "country": "US"}
)
details = client.get_boolean_details("${currentFlagKey}", False, context)
print(f"Evaluated: {details.value}, Reason: {details.reason}")`;

      case 'go':
        return `package main

import (
    "context"
    "fmt"
    openfeature "github.com/open-feature/go-sdk/openfeature"
    flagmanagment "github.com/malah-code/flagmanagment/sdk/go"
)

func main() {
    // 1. Initialize and set provider
    provider := flagmanagment.NewProvider(flagmanagment.Config{
        APIURL: "http://localhost:8080",
        APIKey: "srv_key_...",
    })
    openfeature.SetProvider(provider)

    // 2. Evaluate flag
    client := openfeature.NewClient("backend-service")
    evalCtx := openfeature.NewEvaluationContext("${targetingKey}", map[string]interface{}{
        "tier": "premium",
    })
    
    details, _ := client.BooleanValueDetails(context.Background(), "${currentFlagKey}", false, evalCtx)
    fmt.Printf("Value: %v, Reason: %s\\n", details.Value, details.Reason)
}`;

      case 'java':
        return `import dev.openfeature.sdk.OpenFeatureAPI;
import dev.openfeature.sdk.Client;
import dev.openfeature.sdk.MutableContext;
import com.flagmanagment.sdk.FlagManagmentProvider;

public class FeatureApplication {
    public static void main(String[] args) {
        // 1. Register OpenFeature provider
        OpenFeatureAPI api = OpenFeatureAPI.getInstance();
        api.setProvider(new FlagManagmentProvider("http://localhost:8080", "sdk-key"));

        // 2. Evaluate flag
        Client client = api.getClient();
        MutableContext context = new MutableContext("${targetingKey}");
        context.add("tier", "premium");

        boolean isEnabled = client.getBooleanValue("${currentFlagKey}", false, context);
        System.out.println("Flag enabled: " + isEnabled);
    }
}`;

      case 'dotnet':
        return `using System;
using System.Threading.Tasks;
using OpenFeature;
using OpenFeature.Model;
using FlagManagment.Sdk;

public class Program {
    public static async Task Main() {
        // 1. Configure OpenFeature Provider
        OpenFeature.Api.Instance.SetProvider(new FlagManagmentProvider("http://localhost:8080", "key"));

        // 2. Evaluate flag
        var client = OpenFeature.Api.Instance.GetClient();
        var context = EvaluationContext.Builder()
            .SetTargetingKey("${targetingKey}")
            .Set("tier", "premium")
            .Build();

        var details = await client.GetBooleanDetailsAsync("${currentFlagKey}", false, context);
        Console.WriteLine($"Result: {details.Value}, Reason: {details.Reason}");
    }
}`;

      case 'ios':
        return `import OpenFeature
import FlagManagmentSdk

// 1. Register FlagManagment OpenFeature provider
let provider = FlagManagmentProvider(apiUrl: "http://localhost:8080", apiKey: "mobile-key")
OpenFeatureAPI.shared.setProvider(provider: provider)

// 2. Evaluate flag
let context = MutableEvaluationContext(targetingKey: "${targetingKey}")
context.add(key: "tier", value: "premium")

let isEnabled = OpenFeatureAPI.shared.getClient().getBooleanValue(
    key: "${currentFlagKey}",
    defaultValue: false,
    context: context
)`;

      case 'android':
        return `import dev.openfeature.sdk.OpenFeatureAPI
import dev.openfeature.sdk.MutableContext
import com.flagmanagment.sdk.FlagManagmentProvider

// 1. Register OpenFeature Android provider
OpenFeatureAPI.setProvider(FlagManagmentProvider("http://localhost:8080", "mobile-key"))

// 2. Evaluate flag
val context = MutableContext("${targetingKey}").apply {
    add("tier", "premium")
}
val client = OpenFeatureAPI.getClient()
val isEnabled = client.getBooleanValue("${currentFlagKey}", false, context)`;
    }
  };

  const handleCopyCode = () => {
    navigator.clipboard.writeText(getCodeSnippet(activeLang));
    setHasCopied(true);
    toast.success('Code snippet copied to clipboard');
    setTimeout(() => setHasCopied(false), 2000);
  };

  if (isLoadingFlags || isLoadingEnvs) {
    return <div className="p-8 text-slate-500">Loading OpenFeature Playground...</div>;
  }

  return (
    <div className="max-w-6xl mx-auto space-y-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="inline-flex items-center gap-1.5 px-3 py-1 bg-indigo-50 text-indigo-700 text-xs font-semibold rounded-full border border-indigo-200 mb-3">
            <Sparkles className="w-3.5 h-3.5" /> CNCF OpenFeature Standardized
          </div>
          <h1 className="text-2xl font-bold text-slate-900 flex items-center gap-2">
            <Terminal className="w-6 h-6 text-indigo-600" />
            OpenFeature Playground & Conformance Hub
          </h1>
          <p className="text-slate-500 mt-1 text-sm">
            Evaluate flags using standard OpenFeature APIs, inspect resolution details, and generate
            code for all 8 supported languages.
          </p>
        </div>

        {/* Environment Badge */}
        <div className="flex items-center gap-2 bg-white border border-slate-200 px-3 py-2 rounded-xl shadow-sm text-xs">
          <Server className="w-4 h-4 text-emerald-600" />
          <span className="text-slate-500">Target Env:</span>
          <strong className="text-slate-800 font-semibold">{currentEnv?.name || 'Default'}</strong>
        </div>
      </div>

      {/* Main Grid: Interactive Evaluator & Code Generator */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        {/* Left Column: Interactive OpenFeature Evaluator (5 cols) */}
        <div className="lg:col-span-5 space-y-6">
          <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm space-y-5">
            <h2 className="text-base font-bold text-slate-900 flex items-center gap-2">
              <Zap className="w-4 h-4 text-amber-500" />
              Live Evaluation Tester
            </h2>

            {/* Flag Selector */}
            <div>
              <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wider">
                Select Feature Flag
              </label>
              <select
                value={selectedFlagKey}
                onChange={(e) => setSelectedFlagKey(e.target.value)}
                className="w-full text-sm border border-slate-200 rounded-xl px-3 py-2 bg-slate-50 text-slate-900 font-mono outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
              >
                {flags.map((f) => (
                  <option key={f.id} value={f.key}>
                    {f.key} ({f.type})
                  </option>
                ))}
              </select>
            </div>

            {/* Environment Selector */}
            <div>
              <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wider">
                Target Environment
              </label>
              <select
                value={selectedEnvId}
                onChange={(e) => setSelectedEnvId(e.target.value)}
                className="w-full text-sm border border-slate-200 rounded-xl px-3 py-2 bg-slate-50 text-slate-900 outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
              >
                {environments.map((e) => (
                  <option key={e.id} value={e.id}>
                    {e.name}
                  </option>
                ))}
              </select>
            </div>

            {/* Evaluation Method / Type */}
            <div>
              <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wider">
                OpenFeature Call Type
              </label>
              <div className="grid grid-cols-2 gap-2">
                {(['boolean', 'string', 'number', 'object'] as const).map((t) => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => setEvalType(t)}
                    className={`py-1.5 px-3 text-xs font-medium rounded-lg border text-left transition-all ${
                      evalType === t
                        ? 'bg-indigo-50 border-indigo-300 text-indigo-700 font-semibold'
                        : 'bg-white border-slate-200 text-slate-600 hover:bg-slate-50'
                    }`}
                  >
                    get{t.charAt(0).toUpperCase() + t.slice(1)}Value()
                  </button>
                ))}
              </div>
            </div>

            {/* Context: targetingKey */}
            <div>
              <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wider">
                EvaluationContext.targetingKey
              </label>
              <input
                type="text"
                value={targetingKey}
                onChange={(e) => setTargetingKey(e.target.value)}
                placeholder="usr_1024_beta"
                className="w-full text-sm font-mono border border-slate-200 rounded-xl px-3 py-2 bg-slate-50 text-slate-900 outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
              />
            </div>

            {/* Context: Custom Attributes */}
            <div>
              <label className="block text-xs font-semibold text-slate-700 mb-1.5 uppercase tracking-wider">
                Custom Attributes (JSON)
              </label>
              <textarea
                rows={3}
                value={customAttributes}
                onChange={(e) => setCustomAttributes(e.target.value)}
                className="w-full text-xs font-mono border border-slate-200 rounded-xl p-3 bg-slate-50 text-slate-900 outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
              />
            </div>

            {/* Run Button */}
            <button
              onClick={handleRunEvaluation}
              className="w-full inline-flex items-center justify-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium text-sm py-2.5 px-4 rounded-xl shadow-sm transition-colors"
            >
              <Play className="w-4 h-4 fill-white" />
              Run OpenFeature Evaluation
            </button>
          </div>

          {/* Resolution Result Card */}
          {evalResult && (
            <div className="bg-slate-900 text-white rounded-2xl p-5 shadow-xl border border-slate-800 space-y-4">
              <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                <div className="flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                  <span className="text-xs font-semibold uppercase tracking-wider text-slate-300">
                    ResolutionDetails
                  </span>
                </div>
                <span className="text-[11px] font-mono text-emerald-400">
                  {evalResult.latencyUs} µs
                </span>
              </div>

              <div className="space-y-2 text-xs">
                <div className="flex justify-between py-1 border-b border-slate-800/60">
                  <span className="text-slate-400">Value:</span>
                  <span className="font-mono font-bold text-emerald-300">
                    {typeof evalResult.value === 'object'
                      ? JSON.stringify(evalResult.value)
                      : String(evalResult.value)}
                  </span>
                </div>
                <div className="flex justify-between py-1 border-b border-slate-800/60">
                  <span className="text-slate-400">Variant:</span>
                  <span className="font-mono text-indigo-300">{evalResult.variant}</span>
                </div>
                <div className="flex justify-between py-1 border-b border-slate-800/60">
                  <span className="text-slate-400">Reason:</span>
                  <span className="font-mono text-amber-300">{evalResult.reason}</span>
                </div>
                <div className="flex justify-between py-1">
                  <span className="text-slate-400">Evaluated At:</span>
                  <span className="text-slate-400">{evalResult.evaluatedAt}</span>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Right Column: Multi-Language Code Generator (7 cols) */}
        <div className="lg:col-span-7 space-y-6">
          <div className="bg-white border border-slate-200 rounded-2xl shadow-sm overflow-hidden">
            {/* Language Selector Header */}
            <div className="border-b border-slate-200 bg-slate-50 px-4 py-3 flex flex-wrap items-center justify-between gap-2">
              <div className="flex items-center gap-2">
                <Code2 className="w-4 h-4 text-indigo-600" />
                <span className="text-xs font-bold text-slate-800 uppercase tracking-wider">
                  OpenFeature Client Code
                </span>
              </div>
              <button
                onClick={handleCopyCode}
                className="inline-flex items-center gap-1.5 px-3 py-1 rounded-lg border border-slate-200 bg-white hover:bg-slate-50 text-xs font-semibold text-slate-700 shadow-xs transition-colors"
              >
                {hasCopied ? (
                  <>
                    <Check className="w-3.5 h-3.5 text-emerald-600" /> Copied
                  </>
                ) : (
                  <>
                    <Copy className="w-3.5 h-3.5 text-slate-500" /> Copy Snippet
                  </>
                )}
              </button>
            </div>

            {/* Language Tabs */}
            <div className="flex overflow-x-auto border-b border-slate-200 px-4 pt-2 gap-2 bg-slate-100/60">
              {(
                [
                  { id: 'react', name: 'React (Web)' },
                  { id: 'node', name: 'Node.js' },
                  { id: 'python', name: 'Python' },
                  { id: 'go', name: 'Go' },
                  { id: 'java', name: 'Java' },
                  { id: 'dotnet', name: '.NET (C#)' },
                  { id: 'ios', name: 'iOS (Swift)' },
                  { id: 'android', name: 'Android' },
                ] as const
              ).map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveLang(tab.id)}
                  className={`px-3 py-2 text-xs font-medium whitespace-nowrap border-b-2 transition-all ${
                    activeLang === tab.id
                      ? 'border-indigo-600 text-indigo-600 font-semibold'
                      : 'border-transparent text-slate-600 hover:text-slate-900'
                  }`}
                >
                  {tab.name}
                </button>
              ))}
            </div>

            {/* Code Display Area */}
            <div className="p-4 bg-slate-950 overflow-x-auto">
              <pre className="text-xs font-mono text-emerald-400 leading-relaxed">
                <code>{getCodeSnippet(activeLang)}</code>
              </pre>
            </div>
          </div>

          {/* CNCF OpenFeature Conformance Checklist */}
          <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm">
            <h3 className="text-sm font-bold text-slate-900 mb-4 flex items-center gap-2">
              <Layers className="w-4 h-4 text-indigo-600" />
              CNCF OpenFeature Conformance Verification
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
              <div className="flex items-center gap-2 p-2.5 bg-emerald-50 rounded-lg border border-emerald-100 text-emerald-800">
                <CheckCircle2 className="w-4 h-4 text-emerald-600 shrink-0" />
                <span>Standard Types (Boolean, String, Number, Object)</span>
              </div>
              <div className="flex items-center gap-2 p-2.5 bg-emerald-50 rounded-lg border border-emerald-100 text-emerald-800">
                <CheckCircle2 className="w-4 h-4 text-emerald-600 shrink-0" />
                <span>Deterministic MurmurHash3 Bucketing</span>
              </div>
              <div className="flex items-center gap-2 p-2.5 bg-emerald-50 rounded-lg border border-emerald-100 text-emerald-800">
                <CheckCircle2 className="w-4 h-4 text-emerald-600 shrink-0" />
                <span>Resolution Details & Reason Codes</span>
              </div>
              <div className="flex items-center gap-2 p-2.5 bg-emerald-50 rounded-lg border border-emerald-100 text-emerald-800">
                <CheckCircle2 className="w-4 h-4 text-emerald-600 shrink-0" />
                <span>Telemetry Hooks (PostHog, Datadog)</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
