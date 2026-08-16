import { useState } from 'react';
import { Copy, Check, ShieldCheck, Code2 } from 'lucide-react';
import toast from 'react-hot-toast';
import { SDKIntegrationModal } from './SDKIntegrationModal';

interface ClientSideKeyCardProps {
  envName: string;
  apiKey?: string;
}

export const ClientSideKeyCard = ({ envName, apiKey }: ClientSideKeyCardProps) => {
  const [copied, setCopied] = useState(false);
  const [isGuideOpen, setIsGuideOpen] = useState(false);

  const displayKey = apiKey || `env_${envName.toLowerCase().replace(/[^a-z0-9]/g, '_')}_token`;

  const handleCopy = () => {
    navigator.clipboard.writeText(displayKey);
    setCopied(true);
    toast.success('Client-side Environment Key copied to clipboard!');
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="bg-white border border-slate-200 rounded-2xl p-6 shadow-sm space-y-4">
      <div className="flex items-start justify-between">
        <div>
          <div className="flex items-center gap-2">
            <h3 className="text-base font-bold text-slate-900">Client-side Environment Key</h3>
            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700 border border-emerald-200">
              <ShieldCheck className="w-3.5 h-3.5" /> Public Key
            </span>
          </div>
          <p className="text-sm text-slate-500 mt-1">
            Used by client-side applications (React, Mobile, Browser) for flag evaluations. Safe for
            public exposure.
          </p>
        </div>

        <button
          onClick={() => setIsGuideOpen(true)}
          className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold text-indigo-600 bg-indigo-50 hover:bg-indigo-100 rounded-xl border border-indigo-200/60 transition-colors"
        >
          <Code2 className="w-4 h-4" />
          <span>Integration Guide</span>
        </button>
      </div>

      <div className="flex items-center gap-3">
        <div className="flex-1 bg-slate-50 border border-slate-200 rounded-xl px-4 py-2.5 font-mono text-sm text-slate-800 select-all truncate">
          {displayKey}
        </div>
        <button
          onClick={handleCopy}
          className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold text-xs px-4 py-2.5 rounded-xl shadow-sm transition-colors shrink-0"
        >
          {copied ? (
            <>
              <Check className="w-4 h-4" />
              <span>Copied!</span>
            </>
          ) : (
            <>
              <Copy className="w-4 h-4" />
              <span>Copy Key</span>
            </>
          )}
        </button>
      </div>

      {isGuideOpen && (
        <SDKIntegrationModal
          isOpen={true}
          onClose={() => setIsGuideOpen(false)}
          envName={envName}
          envKey={displayKey}
        />
      )}
    </div>
  );
};
