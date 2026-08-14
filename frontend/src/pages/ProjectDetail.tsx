import { useParams, Link, Routes, Route, useNavigate, Navigate } from 'react-router-dom';
import { useProject } from '../hooks/useProjects';
import { useEnvironments } from '../hooks/useEnvironments';
import { EnvironmentsList } from '../components/environments/EnvironmentsList';
import { FlagsList } from '../components/flags/FlagsList';
import { FlagStatesList } from '../components/flagStates/FlagStatesList';
import { ArrowLeft, Loader2 } from 'lucide-react';
import { EnvironmentSidebar } from '../components/layout/EnvironmentSidebar';

export const ProjectDetail = () => {
  const { projectId = '' } = useParams<{ projectId: string }>();
  const { data: project, isLoading: isLoadingProject } = useProject(projectId);
  const { data: environments = [] } = useEnvironments(projectId);
  const navigate = useNavigate();

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
    <div className="flex flex-col lg:flex-row gap-8 items-start">
      {/* Aesthetic Left Sidebar */}
      <EnvironmentSidebar
        projectId={projectId}
        projectName={project.name}
        environments={environments}
      />

      {/* Main Content Area */}
      <div className="flex-1 min-w-0">
        <Routes>
          <Route 
            index 
            element={
              environments.length > 0 
                ? <Navigate to={`/projects/${projectId}/env/${environments[0].id}`} replace />
                : <Navigate to={`/projects/${projectId}/flags`} replace />
            } 
          />
          <Route 
            path="flags" 
            element={
              <FlagsList 
                projectId={projectId} 
                onNavigateToTargeting={() => navigate(`/projects/${projectId}/env/${environments[0]?.id || ''}`)}
              />
            } 
          />
          <Route 
            path="env/:envId" 
            element={<EnvSpecificView projectId={projectId} />} 
          />
          <Route 
            path="settings" 
            element={<EnvironmentsList projectId={projectId} />} 
          />
        </Routes>
      </div>
    </div>
  );
};

const EnvSpecificView = ({ projectId }: { projectId: string }) => {
  const { envId } = useParams<{ envId: string }>();
  if (!envId) return null;
  return <FlagStatesList projectId={projectId} environmentId={envId} />;
};
