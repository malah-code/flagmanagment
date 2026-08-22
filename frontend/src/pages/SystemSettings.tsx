import React, { useState, useEffect } from 'react';
import { Save, Mail, ShieldCheck, KeyRound } from 'lucide-react';
import { toast } from 'react-hot-toast';
import {
  useSMTPConfig,
  useUpdateSMTPConfig,
  useTestSMTP,
  useSSOConfig,
  useUpdateSSOConfig,
} from '../hooks/useConfig';

export const SystemSettings: React.FC = () => {
  const { data: smtpData } = useSMTPConfig();
  const updateSMTPMutation = useUpdateSMTPConfig();
  const testSMTPMutation = useTestSMTP();

  const { data: ssoData } = useSSOConfig();
  const updateSSOMutation = useUpdateSSOConfig();

  const [smtpForm, setSmtpForm] = useState({
    host: '',
    port: 1025,
    username: '',
    password: '',
  });
  const [testEmail, setTestEmail] = useState('');

  const [ssoForm, setSsoForm] = useState({
    oidc: {
      enabled: false,
      issuer_url: '',
      client_id: '',
      client_secret: '',
    },
    saml: {
      enabled: false,
      idp_metadata_url: '',
      idp_entity_id: '',
      idp_sso_url: '',
      idp_cert: '',
      sp_entity_id: '',
    },
  });

  useEffect(() => {
    if (smtpData) {
      setSmtpForm({
        host: smtpData.host || '',
        port: smtpData.port || 1025,
        username: smtpData.username || '',
        password: '', // Keep empty
      });
    }
  }, [smtpData]);

  useEffect(() => {
    if (ssoData) {
      setSsoForm({
        oidc: {
          enabled: ssoData.oidc?.enabled ?? false,
          issuer_url: ssoData.oidc?.issuer_url || '',
          client_id: ssoData.oidc?.client_id || '',
          client_secret: '',
        },
        saml: {
          enabled: ssoData.saml?.enabled ?? false,
          idp_metadata_url: ssoData.saml?.idp_metadata_url || '',
          idp_entity_id: ssoData.saml?.idp_entity_id || '',
          idp_sso_url: ssoData.saml?.idp_sso_url || '',
          idp_cert: '',
          sp_entity_id: ssoData.saml?.sp_entity_id || '',
        },
      });
    }
  }, [ssoData]);

  const handleSaveSMTP = (e: React.FormEvent) => {
    e.preventDefault();
    updateSMTPMutation.mutate(smtpForm, {
      onSuccess: () => toast.success('SMTP configuration saved!'),
      onError: (err: any) =>
        toast.error(err?.response?.data || 'Failed to save SMTP configuration'),
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
      onError: (err: any) => toast.error(err?.response?.data || 'Failed to send test email'),
    });
  };

  const handleSaveSSO = (e: React.FormEvent) => {
    e.preventDefault();
    updateSSOMutation.mutate(ssoForm, {
      onSuccess: () => toast.success('SSO configuration saved!'),
      onError: (err: any) => toast.error(err?.response?.data || 'Failed to save SSO configuration'),
    });
  };

  return (
    <div className="space-y-10 pb-12">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">System Settings</h1>
          <p className="text-sm text-slate-500 mt-1">
            Configure global platform settings, identity providers, and mail integrations
          </p>
        </div>
      </div>

      {/* Section 1: Enterprise SSO */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 pt-2">
        <div className="md:col-span-1">
          <h2 className="text-lg font-semibold text-slate-900 flex items-center gap-2">
            <ShieldCheck className="w-5 h-5 text-indigo-600" />
            Enterprise SSO
          </h2>
          <p className="text-sm text-slate-500 mt-1">
            Configure Single Sign-On using OpenID Connect (OIDC) or SAML 2.0 identity providers
            (Okta, Azure AD, Google Workspace).
          </p>
        </div>
        <div className="md:col-span-2 space-y-6">
          <form onSubmit={handleSaveSSO} className="space-y-6">
            {/* OIDC Config Card */}
            <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
              <div className="px-6 py-4 bg-slate-50/70 border-b border-slate-200 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <KeyRound className="w-4 h-4 text-indigo-500" />
                  <span className="font-semibold text-slate-800 text-sm">
                    OpenID Connect (OIDC)
                  </span>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={ssoForm.oidc.enabled}
                    onChange={(e) =>
                      setSsoForm({
                        ...ssoForm,
                        oidc: { ...ssoForm.oidc, enabled: e.target.checked },
                      })
                    }
                    className="sr-only peer"
                  />
                  <div className="w-9 h-5 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-indigo-600"></div>
                  <span className="ml-2 text-xs font-medium text-slate-600">
                    {ssoForm.oidc.enabled ? 'Enabled' : 'Disabled'}
                  </span>
                </label>
              </div>

              <div className="p-6 space-y-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">
                    Issuer URL (Discovery Endpoint)
                  </label>
                  <input
                    type="url"
                    value={ssoForm.oidc.issuer_url}
                    onChange={(e) =>
                      setSsoForm({
                        ...ssoForm,
                        oidc: { ...ssoForm.oidc, issuer_url: e.target.value },
                      })
                    }
                    placeholder="https://accounts.google.com or https://dev-123.okta.com/oauth2/default"
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
                  />
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-700">Client ID</label>
                    <input
                      type="text"
                      value={ssoForm.oidc.client_id}
                      onChange={(e) =>
                        setSsoForm({
                          ...ssoForm,
                          oidc: { ...ssoForm.oidc, client_id: e.target.value },
                        })
                      }
                      placeholder="client-id-xyz"
                      className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-700">Client Secret</label>
                    <input
                      type="password"
                      value={ssoForm.oidc.client_secret}
                      onChange={(e) =>
                        setSsoForm({
                          ...ssoForm,
                          oidc: { ...ssoForm.oidc, client_secret: e.target.value },
                        })
                      }
                      placeholder={ssoData?.oidc?.client_id ? '••••••••' : ''}
                      className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
                    />
                  </div>
                </div>
              </div>
            </div>

            {/* SAML 2.0 Config Card */}
            <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
              <div className="px-6 py-4 bg-slate-50/70 border-b border-slate-200 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <ShieldCheck className="w-4 h-4 text-purple-500" />
                  <span className="font-semibold text-slate-800 text-sm">SAML 2.0</span>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={ssoForm.saml.enabled}
                    onChange={(e) =>
                      setSsoForm({
                        ...ssoForm,
                        saml: { ...ssoForm.saml, enabled: e.target.checked },
                      })
                    }
                    className="sr-only peer"
                  />
                  <div className="w-9 h-5 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-purple-600"></div>
                  <span className="ml-2 text-xs font-medium text-slate-600">
                    {ssoForm.saml.enabled ? 'Enabled' : 'Disabled'}
                  </span>
                </label>
              </div>

              <div className="p-6 space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-700">
                      IdP Single Sign-On URL
                    </label>
                    <input
                      type="url"
                      value={ssoForm.saml.idp_sso_url}
                      onChange={(e) =>
                        setSsoForm({
                          ...ssoForm,
                          saml: { ...ssoForm.saml, idp_sso_url: e.target.value },
                        })
                      }
                      placeholder="https://idp.example.com/app/saml/sso"
                      className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
                    />
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-700">
                      IdP Entity ID / Issuer
                    </label>
                    <input
                      type="text"
                      value={ssoForm.saml.idp_entity_id}
                      onChange={(e) =>
                        setSsoForm({
                          ...ssoForm,
                          saml: { ...ssoForm.saml, idp_entity_id: e.target.value },
                        })
                      }
                      placeholder="http://www.okta.com/exk123"
                      className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">
                    IdP X.509 Certificate (Optional)
                  </label>
                  <textarea
                    rows={2}
                    value={ssoForm.saml.idp_cert}
                    onChange={(e) =>
                      setSsoForm({
                        ...ssoForm,
                        saml: { ...ssoForm.saml, idp_cert: e.target.value },
                      })
                    }
                    placeholder="-----BEGIN CERTIFICATE----- ... -----END CERTIFICATE-----"
                    className="w-full px-3 py-2 font-mono text-xs border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="bg-slate-50 p-3 rounded-lg border border-slate-200 text-xs text-slate-600 space-y-1">
                  <p className="font-semibold text-slate-700">Service Provider (SP) Metadata:</p>
                  <p>
                    ACS URL:{' '}
                    <span className="font-mono text-indigo-600">
                      {window.location.origin}/api/v1/auth/saml/acs
                    </span>
                  </p>
                  <p>
                    Entity ID:{' '}
                    <span className="font-mono text-indigo-600">
                      {window.location.origin}/api/v1/auth/saml/metadata
                    </span>
                  </p>
                </div>
              </div>
            </div>

            <div className="flex justify-end">
              <button
                type="submit"
                disabled={updateSSOMutation.isPending}
                className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 disabled:opacity-50 transition-colors shadow-sm"
              >
                <Save className="w-4 h-4" />
                {updateSSOMutation.isPending ? 'Saving...' : 'Save SSO Settings'}
              </button>
            </div>
          </form>
        </div>
      </div>

      <hr className="border-slate-200" />

      {/* Section 2: SMTP Email Server */}
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
                    onChange={(e) => setSmtpForm({ ...smtpForm, host: e.target.value })}
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
                    onChange={(e) => setSmtpForm({ ...smtpForm, port: parseInt(e.target.value) })}
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
                    onChange={(e) => setSmtpForm({ ...smtpForm, username: e.target.value })}
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700">Password (Optional)</label>
                  <input
                    type="password"
                    value={smtpForm.password}
                    onChange={(e) => setSmtpForm({ ...smtpForm, password: e.target.value })}
                    placeholder={smtpData?.username ? '••••••••' : ''}
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
              <p className="text-sm text-slate-500">
                Send a test email to verify your SMTP configuration.
              </p>
              <div className="flex gap-2">
                <input
                  type="email"
                  value={testEmail}
                  onChange={(e) => setTestEmail(e.target.value)}
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
