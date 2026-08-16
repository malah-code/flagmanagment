import { useState, useEffect } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useProject, useUpdateProject, useDeleteProject } from '../../hooks/useProjects';
import {
  Save,
  Trash2,
  AlertTriangle,
  Copy,
  Check,
  Folder,
  Calendar,
  Clock,
  Loader2,
} from 'lucide-react';
import toast from 'react-hot-toast';

interface ProjectSettingsProps {
  projectId: string;
}

export const ProjectSettings = ({ projectId }: ProjectSettingsProps) => {
  const navigate = useNavigate();
  const { data: project, isLoading } = useProject(projectId);
  const updateMutation = useUpdateProject();
  const deleteMutation = useDeleteProject();

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [copiedKey, setCopiedKey] = useState(false);

  useEffect(() => {
    if (project) {
      setName(project.name || '');
      setDescription(project.description || '');
    }
  }, [project]);

  if (isLoading) {
    return (
      <div className="flex justify-center p-12">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
      </div>
    );
  }

  if (!project) {
    return (
      <div className="bg-red-50 text-red-600 p-4 rounded-xl border border-red-200">
        Project not found.
      </div>
    );
  }

  const handleCopyKey = () => {
    navigator.clipboard.writeText(project.id);
    setCopiedKey(true);
    toast.success('Project ID copied to clipboard');
    setTimeout(() => setCopiedKey(false), 2000);
  };

  const handleUpdate = async (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    try {
      await updateMutation.mutateAsync({
        id: project.id,
        payload: {
          name: name.trim(),
          description: description.trim(),
        },
      });
      toast.success('Project settings updated successfully');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to update project';
      toast.error(message);
    }
  };

  const handleDelete = async () => {
    const confirmation = prompt(`To confirm deletion, type the project name "${project.name}":`);

    if (confirmation !== project.name) {
      if (confirmation !== null) {
        toast.error('Project name did not match. Deletion cancelled.');
      }
      return;
    }

    try {
      await deleteMutation.mutateAsync(project.id);
      toast.success(`Project "${project.name}" deleted successfully`);
      navigate('/projects');
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Failed to delete project';
      toast.error(message);
    }
  };

  const isDirty =
    name.trim() !== (project.name || '') || description.trim() !== (project.description || '');

  return (
    <div className="space-y-8 max-w-3xl">
      {/* Header */}
      <div>
        <h2 className="text-xl font-bold text-slate-900">Project Settings</h2>
        <p className="text-xs text-slate-500 mt-1">
          Manage general settings, metadata, and lifecycle options for this project.
        </p>
      </div>

      {/* Project Details Form */}
      <form
        onSubmit={handleUpdate}
        className="bg-white border border-slate-200 rounded-2xl shadow-sm overflow-hidden"
      >
        <div className="p-6 border-b border-slate-100 flex items-center justify-between">
          <div>
            <h3 className="text-base font-semibold text-slate-900">General Information</h3>
            <p className="text-sm text-slate-500 mt-0.5">
              Update the project's name and description.
            </p>
          </div>
          <div className="p-2.5 bg-indigo-50 text-indigo-600 rounded-xl">
            <Folder className="w-5 h-5" />
          </div>
        </div>

        <div className="p-6 space-y-6">
          {/* Project Name */}
          <div className="space-y-2">
            <label htmlFor="projectName" className="block text-sm font-medium text-slate-700">
              Project Name <span className="text-red-500">*</span>
            </label>
            <input
              id="projectName"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-900 focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              placeholder="e.g. Mobile Application"
              required
            />
          </div>

          {/* Project ID */}
          <div className="space-y-2">
            <label className="block text-sm font-medium text-slate-700">Project ID</label>
            <div className="flex items-center gap-2">
              <input
                type="text"
                value={project.id}
                readOnly
                className="w-full font-mono text-xs px-3.5 py-2.5 bg-slate-100 border border-slate-200 rounded-xl text-slate-600 select-all focus:outline-none"
              />
              <button
                type="button"
                onClick={handleCopyKey}
                className="inline-flex items-center gap-1.5 px-3 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-xl text-xs font-semibold transition-colors shrink-0"
                title="Copy Project ID"
              >
                {copiedKey ? (
                  <>
                    <Check className="w-3.5 h-3.5 text-emerald-600" />
                    <span>Copied</span>
                  </>
                ) : (
                  <>
                    <Copy className="w-3.5 h-3.5" />
                    <span>Copy</span>
                  </>
                )}
              </button>
            </div>
            <p className="text-xs text-slate-400">
              Unique identifier used in API requests and SDK configurations.
            </p>
          </div>

          {/* Project Description */}
          <div className="space-y-2">
            <label htmlFor="projectDesc" className="block text-sm font-medium text-slate-700">
              Description
            </label>
            <textarea
              id="projectDesc"
              rows={3}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-200 rounded-xl text-sm text-slate-900 focus:bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors resize-none"
              placeholder="Briefly describe what this project is for..."
            />
          </div>

          {/* Metadata info */}
          <div className="pt-2 border-t border-slate-100 flex flex-wrap gap-6 text-xs text-slate-500">
            {project.createdAt && (
              <div className="flex items-center gap-1.5">
                <Calendar className="w-3.5 h-3.5 text-slate-400" />
                <span>Created {new Date(project.createdAt).toLocaleDateString()}</span>
              </div>
            )}
            {project.updatedAt && (
              <div className="flex items-center gap-1.5">
                <Clock className="w-3.5 h-3.5 text-slate-400" />
                <span>Last updated {new Date(project.updatedAt).toLocaleDateString()}</span>
              </div>
            )}
          </div>
        </div>

        <div className="px-6 py-4 bg-slate-50 border-t border-slate-100 flex justify-end gap-3">
          <button
            type="submit"
            disabled={!isDirty || updateMutation.isPending}
            className="inline-flex items-center gap-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm font-medium px-4 py-2 rounded-xl shadow-sm transition-colors"
          >
            <Save className="w-4 h-4" />
            {updateMutation.isPending ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </form>

      {/* Danger Zone */}
      <div className="border border-red-200 rounded-2xl overflow-hidden shadow-sm">
        <div className="bg-red-50 p-6 border-b border-red-200">
          <div className="flex items-center gap-2 text-red-800">
            <AlertTriangle className="w-5 h-5" />
            <h3 className="text-base font-semibold">Danger Zone</h3>
          </div>
        </div>
        <div className="bg-white p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-6">
          <div className="space-y-1 max-w-lg">
            <h4 className="text-sm font-semibold text-slate-900">Delete Project</h4>
            <p className="text-sm text-slate-500">
              Permanently delete this project and all its environments, feature flags, targeting
              rules, server keys, and history. This action cannot be undone.
            </p>
          </div>
          <button
            type="button"
            onClick={handleDelete}
            disabled={deleteMutation.isPending}
            className="shrink-0 inline-flex items-center justify-center gap-2 px-4 py-2 bg-red-50 hover:bg-red-100 text-red-700 border border-red-200 hover:border-red-300 disabled:opacity-50 disabled:cursor-not-allowed rounded-xl text-sm font-semibold transition-colors"
          >
            <Trash2 className="w-4 h-4" />
            {deleteMutation.isPending ? 'Deleting...' : 'Delete Project'}
          </button>
        </div>
      </div>
    </div>
  );
};
