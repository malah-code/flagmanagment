import { useState } from 'react';
import { useProjects, useDeleteProject } from '../hooks/useProjects';
import { CreateProjectDialog } from '../components/projects/CreateProjectDialog';
import { Link } from 'react-router-dom';
import { Plus, Folder, Trash2, ArrowRight, Loader2, Search } from 'lucide-react';

export const ProjectsList = () => {
  const { data: projects = [], isLoading, isError, error } = useProjects();
  const deleteMutation = useDeleteProject();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const filteredProjects = projects.filter((p) =>
    p.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    p.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDelete = async (id: string, name: string) => {
    if (confirm(`Are you sure you want to delete project "${name}"? All associated environments and flags will be lost.`)) {
      await deleteMutation.mutateAsync(id);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-200 pb-5">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">Projects</h1>
          <p className="text-sm text-slate-500 mt-1">
            Manage your project workspaces, environments, and feature flag configurations.
          </p>
        </div>
        <button
          onClick={() => setIsCreateOpen(true)}
          className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium px-4 py-2.5 rounded-lg shadow-sm transition-colors self-start sm:self-auto"
        >
          <Plus className="w-4 h-4" />
          <span>New Project</span>
        </button>
      </div>

      <div className="flex items-center gap-3 bg-white px-3.5 py-2.5 rounded-lg border border-slate-200 shadow-sm max-w-md">
        <Search className="w-4 h-4 text-slate-400" />
        <input
          type="text"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          placeholder="Filter projects..."
          className="w-full text-sm outline-none bg-transparent text-slate-900 placeholder:text-slate-400"
        />
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center p-12">
          <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
        </div>
      ) : isError ? (
        <div className="bg-red-50 text-red-600 p-4 rounded-xl border border-red-200 text-sm">
          Failed to load projects: {(error as Error).message}
        </div>
      ) : filteredProjects.length === 0 ? (
        <div className="bg-white border border-slate-200 rounded-xl p-12 text-center space-y-3">
          <div className="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-full flex items-center justify-center mx-auto">
            <Folder className="w-6 h-6" />
          </div>
          <h3 className="text-base font-semibold text-slate-900">No projects found</h3>
          <p className="text-sm text-slate-500 max-w-sm mx-auto">
            {searchTerm ? 'No projects match your search query.' : 'Get started by creating your first project.'}
          </p>
          {!searchTerm && (
            <button
              onClick={() => setIsCreateOpen(true)}
              className="inline-flex items-center gap-2 bg-indigo-600 text-white text-sm font-medium px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors mt-2"
            >
              <Plus className="w-4 h-4" />
              <span>Create Project</span>
            </button>
          )}
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
          {filteredProjects.map((project) => (
            <div
              key={project.id}
              className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm hover:shadow-md transition-all flex flex-col justify-between group"
            >
              <div>
                <div className="flex items-start justify-between gap-3">
                  <div className="flex items-center gap-2.5 font-semibold text-slate-900 text-lg group-hover:text-indigo-600 transition-colors">
                    <Folder className="w-5 h-5 text-indigo-500 shrink-0" />
                    <span className="truncate">{project.name}</span>
                  </div>
                  <button
                    onClick={() => handleDelete(project.id, project.name)}
                    className="text-slate-400 hover:text-red-600 transition-colors p-1 rounded hover:bg-slate-50"
                    title="Delete Project"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
                <p className="text-sm text-slate-500 mt-2 line-clamp-2 min-h-[40px]">
                  {project.description || 'No description provided.'}
                </p>
              </div>

              <div className="pt-4 border-t border-slate-100 flex items-center justify-between text-xs text-slate-400 mt-4">
                <span>Created {new Date(project.createdAt).toLocaleDateString()}</span>
                <Link
                  to={`/projects/${project.id}`}
                  className="inline-flex items-center gap-1 font-medium text-indigo-600 hover:text-indigo-700 hover:underline text-sm"
                >
                  Manage <ArrowRight className="w-3.5 h-3.5" />
                </Link>
              </div>
            </div>
          ))}
        </div>
      )}

      <CreateProjectDialog isOpen={isCreateOpen} onClose={() => setIsCreateOpen(false)} />
    </div>
  );
};
