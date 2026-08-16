import React from 'react';
import { Loader2 } from 'lucide-react';

interface SwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
  loading?: boolean;
  label?: string;
  size?: 'sm' | 'md' | 'lg';
}

export const Switch: React.FC<SwitchProps> = ({
  checked,
  onChange,
  disabled = false,
  loading = false,
  label,
  size = 'md',
}) => {
  const sizeClasses = {
    sm: 'w-8 h-4 outer-circle:w-3 outer-circle:h-3',
    md: 'w-11 h-6 outer-circle:w-5 outer-circle:h-5',
    lg: 'w-14 h-7 outer-circle:w-6 outer-circle:h-6',
  };

  const translateClasses = {
    sm: checked ? 'translate-x-4' : 'translate-x-0.5',
    md: checked ? 'translate-x-5' : 'translate-x-0.5',
    lg: checked ? 'translate-x-7' : 'translate-x-0.5',
  };

  const circleSizes = {
    sm: 'w-3 h-3',
    md: 'w-5 h-5',
    lg: 'w-6 h-6',
  };

  return (
    <label
      className={`inline-flex items-center gap-2 ${disabled || loading ? 'cursor-not-allowed opacity-60' : 'cursor-pointer'}`}
    >
      <button
        type="button"
        role="switch"
        aria-checked={checked}
        disabled={disabled || loading}
        onClick={() => !disabled && !loading && onChange(!checked)}
        className={`relative inline-flex flex-shrink-0 rounded-full transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 ${
          sizeClasses[size].split(' ')[0]
        } ${sizeClasses[size].split(' ')[1]} ${checked ? 'bg-emerald-500' : 'bg-slate-300'}`}
      >
        <span className="sr-only">{label || 'Toggle switch'}</span>
        <span
          className={`pointer-events-none inline-block transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out flex items-center justify-center ${
            circleSizes[size]
          } ${translateClasses[size]} my-auto top-0 bottom-0 relative`}
        >
          {loading && <Loader2 className="w-3 h-3 animate-spin text-slate-500" />}
        </span>
      </button>
      {label && <span className="text-sm font-medium text-slate-700">{label}</span>}
    </label>
  );
};
