import React from 'react';
import { Tag, X } from 'lucide-react';

interface TagFilterProps {
  allTags: string[];
  selectedTag: string;
  onChange: (tag: string) => void;
}

export const TagFilter: React.FC<TagFilterProps> = ({
  allTags,
  selectedTag,
  onChange,
}) => {
  if (allTags.length === 0) return null;

  return (
    <div className="flex items-center gap-2 flex-wrap">
      <span className="inline-flex items-center gap-1 text-xs font-semibold text-slate-500 uppercase tracking-wider">
        <Tag className="w-3.5 h-3.5 text-slate-400" /> Tags:
      </span>
      <button
        onClick={() => onChange('')}
        className={`px-2.5 py-1 rounded-lg text-xs font-medium transition-all ${
          !selectedTag
            ? 'bg-indigo-600 text-white shadow-sm'
            : 'bg-white text-slate-600 border border-slate-200 hover:bg-slate-50'
        }`}
      >
        All
      </button>
      {allTags.map((tag) => {
        const isSelected = selectedTag === tag;
        return (
          <button
            key={tag}
            onClick={() => onChange(isSelected ? '' : tag)}
            className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs font-medium transition-all ${
              isSelected
                ? 'bg-indigo-600 text-white shadow-sm'
                : 'bg-white text-slate-600 border border-slate-200 hover:bg-slate-50'
            }`}
          >
            <span>#{tag}</span>
            {isSelected && <X className="w-3 h-3 ml-0.5" />}
          </button>
        );
      })}
    </div>
  );
};
