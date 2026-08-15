import React, { useState, useEffect } from 'react';
import { Save, Mail } from 'lucide-react';
import { toast } from 'react-hot-toast';
import { useSMTPConfig, useUpdateSMTPConfig, useTestSMTP } from '../hooks/useConfig';

export const SystemSettings: React.FC = () => {
  const { data: smtpData } = useSMTPConfig();
  const updateSMTPMutation = useUpdateSMTPConfig();
  const testSMTPMutation = useTestSMTP();

  const [smtpForm, setSmtpForm] = useState({
    host: '',
    port: 1025,
    username: '',
    password: ''
  });
  const [testEmail, setTestEmail] = useState('');

  useEffect(() => {
    if (smtpData) {
      setSmtpForm({
        host: smtpData.host || '',
        port: smtpData.port || 1025,
        username: smtpData.username || '',
        password: '' // Keep empty
      });
    }
  }, [smtpData]);

  const handleSaveSMTP = (e: React.FormEvent) => {
    e.preventDefault();
    updateSMTPMutation.mutate(smtpForm, {
      onSuccess: () => toast.success('SMTP configuration saved!'),
      onError: (err: any) => toast.error(err?.response?.data || 'Failed to save SMTP configuration')
    });
  };

  const handleTestEmail = () => {
    if (!testEmail) {
      toast.error('Please enter a test email address');
      return;
    }
    testSMTPMutation.mutate(testEmail, {
      onSuccess: () => {
        toast.success('Test email sent successfully!');
        setTestEmail('');
      },
      onError: (err: any) => toast.error(err?.response?.data || 'Failed to send test email')
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">System Settings</h1>
          <p className="text-sm text-slate-500 mt-1">Configure global platform settings and integrations</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-1">
          <h2 className="text-lg font-semibold text-slate-900 flex items-center gap-2">
            <Mail className="w-5 h-5 text-indigo-500" />
            Email Server (SMTP)
          </h2>
          <p className="text-sm text-slate-500 mt-1">
            Configure the outbound mail server used for invitations and notifications.
          </p>
        </div>
        <div className="md:col-span-2 space-y-6">
          <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
            <form onSubmit={handleSaveSMTP} className="p-6 space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">SMTP Host</label>
                  <input
                    type="text"
                    required
                    value={smtpForm.host}
                    onChange={e => setSmtpForm({...smtpForm, host: e.target.value})}
                    placeholder="smtp.example.com"
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">SMTP Port</label>
                  <input
                    type="number"
                    required
                    value={smtpForm.port}
                    onChange={e => setSmtpForm({...smtpForm, port: parseInt(e.target.value)})}
                    placeholder="587"
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">Username (Optional)</label>
                  <input
                    type="text"
                    value={smtpForm.username}
                    onChange={e => setSmtpForm({...smtpForm, username: e.target.value})}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">Password (Optional)</label>
                  <input
                    type="password"
                    value={smtpForm.password}
                    onChange={e => setSmtpForm({...smtpForm, password: e.target.value})}
                    placeholder={smtpData?.username ? "••••••••" : ""}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                  />
                  {smtpData?.username && (
                    <p className="text-xs text-slate-500">Leave empty to keep existing password</p>
                  )}
                </div>
              </div>
              <div className="pt-4 flex justify-end border-t border-slate-100">
                <button
                  type="submit"
                  disabled={updateSMTPMutation.isPending}
                  className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50"
                >
                  <Save className="w-4 h-4" />
                  {updateSMTPMutation.isPending ? 'Saving...' : 'Save Settings'}
                </button>
              </div>
            </form>
          </div>

          <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
            <div className="p-6 space-y-4">
              <h3 className="text-md font-medium text-slate-900">Test Connection</h3>
              <p className="text-sm text-slate-500">Send a test email to verify your SMTP configuration.</p>
              <div className="flex gap-2">
                <input
                  type="email"
                  value={testEmail}
                  onChange={e => setTestEmail(e.target.value)}
                  placeholder="test@example.com"
                  className="flex-1 px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                />
                <button
                  type="button"
                  onClick={handleTestEmail}
                  disabled={testSMTPMutation.isPending || !testEmail}
                  className="px-4 py-2 text-sm font-medium text-indigo-600 bg-indigo-50 rounded-lg hover:bg-indigo-100 disabled:opacity-50"
                >
                  {testSMTPMutation.isPending ? 'Sending...' : 'Send Test'}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
