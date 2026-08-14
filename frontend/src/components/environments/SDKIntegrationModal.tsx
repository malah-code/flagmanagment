import React, { useState } from 'react';
import { X, Copy, Check, Code2, ShieldCheck, Terminal } from 'lucide-react';
import toast from 'react-hot-toast';

interface SDKIntegrationModalProps {
  isOpen: boolean;
  onClose: () => void;
  envName: string;
  envKey: string;
}

type SDKLanguage = 'react' | 'node' | 'python' | 'go';

export const SDKIntegrationModal: React.FC<SDKIntegrationModalProps> = ({
  isOpen,
  onClose,
  envName,
  envKey,
}) => {
  const [activeTab, setActiveTab] = useState<SDKLanguage>('react');
  const [copied, setCopied] = useState(false);

  if (!isOpen) return null;

  const displayKey = envKey || `env_${envName.toLowerCase().replace(/[^a-z0-9]/g, '_')}_token`;

  const snippets: Record<SDKLanguage, { title: string; install: string; code: string }> = {
    react: {
      title: 'React / Web SDK',
      install: 'npm install @flagmanagement/react-sdk',
      code: `import { FlagProvider, useFeatureFlag } from '@flagmanagement/react-sdk';

// 1. Wrap your application with FlagProvider
export const App = () => (
  <FlagProvider
    clientKey="${displayKey}"
    endpoint="http://localhost:8080"
    userContext={{ userId: "user_123" }}
  >
    <MyComponent />
  </FlagProvider>
);

// 2. Evaluate flags in any component
const MyComponent = () => {
  const isEnabled = useFeatureFlag('new-checkout-flow');
  return isEnabled ? <NewCheckout /> : <OldCheckout />;
};`,
    },
    node: {
      title: 'Node.js / Express SDK',
      install: 'npm install @flagmanagement/node-sdk',
      code: `const { FlagClient } = require('@flagmanagement/node-sdk');

const flagClient = new FlagClient({
  clientKey: "${displayKey}",
  endpoint: "http://localhost:8080",
});

async function main() {
  await flagClient.connect();
  
  const isEnabled = await flagClient.evaluate('new-checkout-flow', {
    userId: 'user_123'
  });

  console.log('Feature Enabled:', isEnabled);
}`,
    },
    python: {
      title: 'Python SDK',
      install: 'pip install flagmanagement-sdk',
      code: `from flagmanagement import FlagClient

client = FlagClient(
    client_key="${displayKey}",
    endpoint="http://localhost:8080"
)

# Evaluate flag for a user
is_enabled = client.evaluate("new-checkout-flow", user_context={
    "userId": "user_123"
})

print(f"Feature enabled: {is_enabled}")`,
    },
    go: {
      title: 'Go SDK',
      install: 'go get github.com/flagmanagment/go-sdk',
      code: `package main

import (
	"fmt"
	"github.com/flagmanagment/go-sdk"
)

func main() {
	client := flagclient.New(flagclient.Config{
		ClientKey: "${displayKey}",
		Endpoint:  "http://localhost:8080",
	})
	defer client.Close()

	ctx := flagclient.Context{"userId": "user_123"}
	isEnabled := client.Evaluate("new-checkout-flow", ctx)

	fmt.Println("Feature Enabled:", isEnabled)
}`,
    },
  };

  const handleCopyCode = () => {
    navigator.clipboard.writeText(snippets[activeTab].code);
    setCopied(true);
    toast.success('Code snippet copied to clipboard');
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl border border-slate-200 max-w-2xl w-full p-6 space-y-5 animate-in fade-in zoom-in-95 duration-200 max-h-[90vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-slate-100 pb-4">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-indigo-50 text-indigo-600 rounded-xl border border-indigo-100">
              <Code2 className="w-5 h-5" />
            </div>
            <div>
              <h2 className="text-lg font-bold text-slate-900 leading-tight">
                SDK Integration Guide
              </h2>
              <p className="text-xs text-slate-500 flex items-center gap-1.5 mt-0.5">
                <span>Environment:</span>
                <span className="font-semibold text-slate-800 bg-slate-100 px-2 py-0.5 rounded-md">
                  {envName}
                </span>
                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200 ml-1">
                  <ShieldCheck className="w-3 h-3" /> Client-Side / Public Key
                </span>
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-slate-600 transition-colors p-1.5 rounded-lg hover:bg-slate-100"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Tab Navigation */}
        <div className="flex border-b border-slate-200 gap-2">
          {(['react', 'node', 'python', 'go'] as SDKLanguage[]).map((lang) => (
            <button
              key={lang}
              onClick={() => setActiveTab(lang)}
              className={`px-4 py-2 text-xs font-bold transition-all border-b-2 capitalize ${
                activeTab === lang
                  ? 'border-indigo-600 text-indigo-600 bg-indigo-50/50 rounded-t-lg'
                  : 'border-transparent text-slate-500 hover:text-slate-800 hover:bg-slate-50 rounded-t-lg'
              }`}
            >
              {lang === 'react' ? 'React (Web)' : lang === 'node' ? 'Node.js' : lang}
            </button>
          ))}
        </div>

        {/* Snippet Area */}
        <div className="space-y-3 flex-1 overflow-y-auto pr-1">
          {/* Installation Command */}
          <div className="bg-slate-900 rounded-xl p-3 flex items-center justify-between font-mono text-xs text-slate-300 border border-slate-800">
            <div className="flex items-center gap-2">
              <Terminal className="w-4 h-4 text-emerald-400" />
              <span>{snippets[activeTab].install}</span>
            </div>
            <button
              onClick={() => {
                navigator.clipboard.writeText(snippets[activeTab].install);
                toast.success('Install command copied');
              }}
              className="text-slate-400 hover:text-white transition-colors p-1"
              title="Copy Install Command"
            >
              <Copy className="w-3.5 h-3.5" />
            </button>
          </div>

          {/* Code Block */}
          <div className="relative bg-slate-950 rounded-xl border border-slate-800 p-4 font-mono text-xs text-slate-200 overflow-x-auto">
            <button
              onClick={handleCopyCode}
              className="absolute top-3 right-3 bg-slate-800 hover:bg-slate-700 text-slate-200 hover:text-white px-2.5 py-1.5 rounded-lg text-[11px] font-semibold transition-all flex items-center gap-1.5 shadow-sm border border-slate-700"
            >
              {copied ? (
                <>
                  <Check className="w-3.5 h-3.5 text-emerald-400" />
                  <span>Copied!</span>
                </>
              ) : (
                <>
                  <Copy className="w-3.5 h-3.5" />
                  <span>Copy Code</span>
                </>
              )}
            </button>
            <pre className="whitespace-pre">{snippets[activeTab].code}</pre>
          </div>
        </div>

        {/* Footer */}
        <div className="pt-2 flex justify-end border-t border-slate-100">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-slate-900 hover:bg-slate-800 text-white text-xs font-semibold rounded-xl transition-colors shadow-sm"
          >
            Done
          </button>
        </div>
      </div>
    </div>
  );
};
