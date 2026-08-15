import React, { useState } from 'react';
import { Shield, UserPlus, Search, Edit2 } from 'lucide-react';
import { toast } from 'react-hot-toast';
import { useProjects } from '../hooks/useProjects';
import { useUsers, useInviteUser, useUpdateUserAccess } from '../hooks/useUsers';

export const UsersManagement: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const { data: projectsData = [] } = useProjects();
  const { data: usersData, isLoading } = useUsers();
  const inviteUserMutation = useInviteUser();
  const updateUserMutation = useUpdateUserAccess();
  
  const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
  const [inviteEmail, setInviteEmail] = useState('');
  const [inviteRole, setInviteRole] = useState('Read-Only Auditor');

  const allUsers = (usersData?.users || []).map((u: any) => ({
    id: u.id,
    email: u.email,
    name: u.email.split('@')[0],
    role: u.roles && u.roles.length > 0 ? u.roles[0] : 'No Role',
    projects: u.projects && u.projects.length > 0 ? u.projects : [],
    lastActive: new Date(u.updated_at).toLocaleDateString()
  }));

  const [editingUser, setEditingUser] = useState<{ id: string, email: string, name: string, role: string, projects: string[], lastActive: string } | null>(null);

  const handleSaveEdit = (e: React.FormEvent) => {
    e.preventDefault();
    if (editingUser) {
      const projects = editingUser.role === 'Global Administrator' ? [] : editingUser.projects;

      updateUserMutation.mutate(
        { id: editingUser.id, payload: { role: editingUser.role, project_ids: projects } },
        {
          onSuccess: () => {
            toast.success('User updated successfully!');
            setEditingUser(null);
          },
          onError: (error: any) => {
            toast.error(error?.response?.data || 'Failed to update user');
          }
        }
      );
    }
  };

  const handleInvite = (e: React.FormEvent) => {
    e.preventDefault();
    
    const projects = inviteRole === 'Global Administrator' ? [] : []; // Would need an array of UUIDs if they can select projects on invite, but MVP invite is fixed to Global Admin or Read-Only (no project selected initially)

    inviteUserMutation.mutate(
      { email: inviteEmail, role: inviteRole, project_ids: projects },
      {
        onSuccess: () => {
          toast.success(`Invitation sent to ${inviteEmail}!`);
          setIsInviteModalOpen(false);
          setInviteEmail('');
          setInviteRole('Read-Only Auditor');
        },
        onError: (error: any) => {
          toast.error(error?.response?.data || 'Failed to send invitation');
        }
      }
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Team Settings</h1>
          <p className="text-sm text-slate-500 mt-1">Manage user access and assign global or project-level roles.</p>
        </div>
        <button
          onClick={() => setIsInviteModalOpen(true)}
          className="inline-flex items-center justify-center gap-2 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors"
        >
          <UserPlus className="w-4 h-4" />
          <span>Invite User</span>
        </button>
      </div>

      <div className="bg-white border border-slate-200 rounded-xl shadow-sm overflow-hidden flex flex-col">
        <div className="p-4 border-b border-slate-200 bg-slate-50 flex items-center justify-between">
          <div className="relative w-full max-w-md">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input
              type="text"
              placeholder="Search by name or email..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-9 pr-4 py-2 text-sm border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent bg-white"
            />
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 text-xs uppercase font-semibold text-slate-500 border-b border-slate-200">
              <tr>
                <th className="px-6 py-4">User</th>
                <th className="px-6 py-4">Role</th>
                <th className="px-6 py-4">Projects</th>
                <th className="px-6 py-4">Last Active</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {isLoading && (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-slate-500">
                    Loading users...
                  </td>
                </tr>
              )}
              {!isLoading && allUsers.filter(u => u.email.includes(searchQuery) || u.name.toLowerCase().includes(searchQuery.toLowerCase())).map((user) => (
                <tr key={user.id} className="hover:bg-slate-50/50 transition-colors">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-indigo-100 text-indigo-700 flex items-center justify-center font-bold">
                        {user.name.charAt(0)}
                      </div>
                      <div>
                        <div className="font-medium text-slate-900">{user.name}</div>
                        <div className="text-xs text-slate-500">{user.email}</div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <span className="inline-flex items-center gap-1.5 rounded-full bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-700 border border-slate-200">
                      <Shield className="w-3 h-3 text-indigo-500" />
                      {user.role}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-slate-500 text-sm">
                    {user.role === 'Global Administrator' ? (
                      <span className="font-medium text-indigo-600">All Projects</span>
                    ) : (
                      user.projects.length > 0 
                        ? user.projects.map((pid: string) => projectsData.find((p: any) => p.id === pid)?.name || pid).join(', ')
                        : 'No Projects'
                    )}
                  </td>
                  <td className="px-6 py-4 text-slate-500">{user.lastActive}</td>
                  <td className="px-6 py-4 text-right">
                    <button 
                      onClick={() => setEditingUser(user)}
                      className="p-1.5 text-slate-400 hover:text-indigo-600 rounded-lg hover:bg-slate-100 transition-colors"
                      title="Edit User"
                    >
                      <Edit2 className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
              {!isLoading && allUsers.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-6 py-8 text-center text-slate-500">
                    No users found matching your search.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>

      {isInviteModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-100">
              <h3 className="text-lg font-semibold text-slate-900">Invite Team Member</h3>
            </div>
            <form onSubmit={handleInvite} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Email Address</label>
                <input
                  type="email"
                  required
                  value={inviteEmail}
                  onChange={e => setInviteEmail(e.target.value)}
                  placeholder="colleague@example.com"
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Role Assignment</label>
                <select 
                  value={inviteRole}
                  onChange={e => setInviteRole(e.target.value)}
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 bg-white"
                >
                  <option value="Read-Only Auditor">Read-Only Auditor</option>
                  <option value="Project Editor">Project Editor</option>
                  <option value="Global Administrator">Global Administrator</option>
                </select>
              </div>
              <div className="flex gap-3 pt-4 border-t border-slate-100 mt-6">
                <button
                  type="button"
                  onClick={() => setIsInviteModalOpen(false)}
                  className="flex-1 px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-lg hover:bg-slate-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={inviteUserMutation.isPending}
                  className="flex-1 px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50"
                >
                  {inviteUserMutation.isPending ? 'Sending...' : 'Send Invite'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {editingUser && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-xl w-full max-w-md overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-100">
              <h3 className="text-lg font-semibold text-slate-900">Edit Team Member</h3>
            </div>
            <form onSubmit={handleSaveEdit} className="p-6 space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Name</label>
                <input
                  type="text"
                  required
                  value={editingUser.name}
                  onChange={e => setEditingUser({...editingUser, name: e.target.value})}
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Email Address</label>
                <input
                  type="email"
                  required
                  value={editingUser.email}
                  onChange={e => setEditingUser({...editingUser, email: e.target.value})}
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-700 mb-1">Role Assignment</label>
                <select 
                  value={editingUser.role}
                  onChange={e => {
                    const role = e.target.value;
                    const defaultProjects = role === 'Global Administrator' ? ['All Projects'] : [];
                    setEditingUser({...editingUser, role, projects: defaultProjects});
                  }}
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 bg-white"
                >
                  <option value="Read-Only Auditor">Read-Only Auditor</option>
                  <option value="Project Editor">Project Editor</option>
                  <option value="Global Administrator">Global Administrator</option>
                </select>
              </div>
              {editingUser.role !== 'Global Administrator' && (
                <div>
                  <label className="block text-sm font-medium text-slate-700 mb-2">Assigned Projects</label>
                  <div className="space-y-2 max-h-40 overflow-y-auto p-2 border border-slate-200 rounded-lg bg-slate-50">
                    {projectsData.map((p: any) => (
                      <label key={p.id} className="flex items-center gap-2 text-sm text-slate-700 cursor-pointer">
                        <input
                          type="checkbox"
                          className="rounded text-indigo-600 focus:ring-indigo-500"
                          checked={editingUser.projects.includes(p.id)}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setEditingUser({ ...editingUser, projects: [...editingUser.projects, p.id] });
                            } else {
                              setEditingUser({ ...editingUser, projects: editingUser.projects.filter(id => id !== p.id) });
                            }
                          }}
                        />
                        {p.name}
                      </label>
                    ))}
                    {projectsData.length === 0 && <span className="text-xs text-slate-500">No projects found.</span>}
                  </div>
                </div>
              )}
              <div className="flex gap-3 pt-4 border-t border-slate-100 mt-6">
                <button
                  type="button"
                  onClick={() => setEditingUser(null)}
                  className="flex-1 px-4 py-2 text-sm font-medium text-slate-700 bg-white border border-slate-300 rounded-lg hover:bg-slate-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={updateUserMutation.isPending}
                  className="flex-1 px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50"
                >
                  {updateUserMutation.isPending ? 'Saving...' : 'Save Changes'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
