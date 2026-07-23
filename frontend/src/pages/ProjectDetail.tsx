import { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useProject } from '../hooks/useProjects';
import { useEnvironments } from '../hooks/useEnvironments';
import { EnvironmentsList } from '../components/environments/EnvironmentsList';
import { FlagsList } from '../components/flags/FlagsList';
import { FlagStatesList } from '../components/flagStates/FlagStatesList';
import { ArrowLeft, Server, Flag, Sliders, Loader2, Folder } from 'lucide-react';

type TabType = 'environments' | 'flags' | 'targeting';

export const ProjectDetail = () => {
  const { projectId = '' } = useParams<{ projectId: string }>();
  const { data: project, isLoading: isLoadingProject } = useProject(projectId);
  const { data: environments = [] } = useEnvironments(projectId);

  const [activeTab, setActiveTab] = useState<TabType>('environments');
  const [selectedEnvId, setSelectedEnvId] = useState<string>('');

  // Default to first environment if available
  const currentEnvId = selectedEnvId || (environments[0]?.id ?? '');

  if (isLoadingProject) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
      </div>
    );
  }

  if (!project) {
    return (
      <div className="space-y-4">
        <Link to="/projects" className="inline-flex items-center gap-1 text-sm text-indigo-600 hover:underline">
          <ArrowLeft className="w-4 h-4" /> Back to Projects
        </Link>
        <div className="bg-red-50 text-red-600 p-4 rounded-xl border border-red-200">
          Project not found.
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumb & Header */}
      <div>
        <Link
          to="/projects"
          className="inline-flex items-center gap-1 text-sm font-medium text-slate-500 hover:text-slate-900 transition-colors mb-3"
        >
          <ArrowLeft className="w-4 h-4" /> Projects
        </Link>

        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-indigo-50 text-indigo-600 rounded-xl">
            <Folder className="w-6 h-6" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-slate-900">{project.name}</h1>
            <p className="text-sm text-slate-500">{project.description || 'No description provided.'}</p>
          </div>
        </div>
      </div>

      {/* Tabs Bar */}
      <div className="border-b border-slate-200 flex items-center justify-between gap-4">
        <nav className="flex gap-6 -mb-px">
          <button
            onClick={() => setActiveTab('environments')}
            className={`inline-flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-medium transition-colors ${
              activeTab === 'environments'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <Server className="w-4 h-4" />
            <span>Environments ({environments.length})</span>
          </button>

          <button
            onClick={() => setActiveTab('flags')}
            className={`inline-flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-medium transition-colors ${
              activeTab === 'flags'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <Flag className="w-4 h-4" />
            <span>Feature Flags</span>
          </button>

          <button
            onClick={() => setActiveTab('targeting')}
            className={`inline-flex items-center gap-2 py-3 px-1 border-b-2 text-sm font-medium transition-colors ${
              activeTab === 'targeting'
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:text-slate-700 hover:border-slate-300'
            }`}
          >
            <Sliders className="w-4 h-4" />
            <span>Flag Targeting & States</span>
          </button>
        </nav>

        {activeTab === 'targeting' && environments.length > 0 && (
          <div className="flex items-center gap-2 pb-2">
            <span className="text-xs font-semibold text-slate-500 uppercase">Environment:</span>
            <select
              value={currentEnvId}
              onChange={(e) => setSelectedEnvId(e.target.value)}
              className="text-sm border border-slate-300 rounded-lg px-3 py-1.5 bg-white text-slate-900 font-medium focus:ring-2 focus:ring-indigo-500 outline-none"
            >
              {environments.map((e) => (
                <option key={e.id} value={e.id}>
                  {e.name}
                </option>
              ))}
            </select>
          </div>
        )}
      </div>

      {/* Tab Content */}
      <div className="pt-2">
        {activeTab === 'environments' && <EnvironmentsList projectId={projectId} />}
        {activeTab === 'flags' && <FlagsList projectId={projectId} />}
        {activeTab === 'targeting' && (
          currentEnvId ? (
            <FlagStatesList projectId={projectId} environmentId={currentEnvId} />
          ) : (
            <div className="bg-white border border-slate-200 rounded-xl p-8 text-center text-sm text-slate-500">
              Please create an environment under the "Environments" tab first to configure flag targeting.
            </div>
          )
        )}
      </div>
    </div>
  );
};
