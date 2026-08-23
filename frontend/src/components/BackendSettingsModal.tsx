import React, { useState } from 'react';
import { api, getApiBaseUrl, setApiBaseUrl, DEFAULT_API_BASE_URL } from '../services/api';
import { useNotification } from '../context/NotificationContext';
import {
  X,
  Server,
  CheckCircle2,
  AlertCircle,
  RotateCcw,
  Activity,
  Terminal,
} from 'lucide-react';

interface BackendSettingsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const BackendSettingsModal: React.FC<BackendSettingsModalProps> = ({
  isOpen,
  onClose,
}) => {
  const [apiUrl, setApiUrl] = useState(getApiBaseUrl());
  const [isTesting, setIsTesting] = useState(false);
  const [testResult, setTestResult] = useState<{
    ok: boolean;
    latency: number;
    tested: boolean;
  }>({
    ok: false,
    latency: 0,
    tested: false,
  });

  const { addToast } = useNotification();

  if (!isOpen) return null;

  const handleTestConnection = async () => {
    setIsTesting(true);
    setApiBaseUrl(apiUrl);
    const res = await api.checkBackendHealth();
    setTestResult({ ok: res.ok, latency: res.latency, tested: true });
    setIsTesting(false);

    if (res.ok) {
      addToast('Backend Connected', `Go server responding in ${res.latency}ms`, 'success');
    } else {
      addToast('Backend Unreachable', 'Ensure your Go backend is running and reachable.', 'error');
    }
  };

  const handleSave = () => {
    setApiBaseUrl(apiUrl);
    addToast('Settings Saved', `API Base URL set to ${apiUrl || DEFAULT_API_BASE_URL}`, 'success');
    onClose();
  };

  const handleReset = () => {
    setApiBaseUrl('');
    setApiUrl(DEFAULT_API_BASE_URL);
    setTestResult({ ok: false, latency: 0, tested: false });
    addToast('Reset', `API Base URL reset to default (${DEFAULT_API_BASE_URL})`, 'info');
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80">
      <div className="relative w-full max-w-lg rounded-2xl bg-slate-900 border border-slate-800 p-6 sm:p-7 overflow-hidden space-y-5">
        <button
          onClick={onClose}
          className="absolute top-5 right-5 p-1.5 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-2xl bg-sky-500/20 border border-sky-500/30 flex items-center justify-center text-sky-400">
            <Server className="w-6 h-6" />
          </div>
          <div>
            <h3 className="text-xl font-bold text-white tracking-tight">
              Backend API Configuration
            </h3>
            <p className="text-xs text-slate-400 mt-0.5">
              Connect frontend to your Go backend instance (uses /api by default for same host).
            </p>
          </div>
        </div>

        <div className="space-y-3">
          <label className="block text-xs font-semibold text-slate-300">
            Go Backend Base URL
          </label>
          <div className="flex gap-2">
            <input
              type="text"
              value={apiUrl}
              onChange={(e) => {
                setApiUrl(e.target.value);
                setTestResult({ ok: false, latency: 0, tested: false });
              }}
              placeholder="/api (default) or http://localhost:8080"
              className="flex-1 px-4 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 placeholder-slate-500 text-sm font-mono focus:outline-none focus:border-sky-500 transition-all"
            />
            <button
              onClick={handleTestConnection}
              disabled={isTesting}
              className="px-4 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-100 text-xs font-bold transition-colors flex items-center gap-1.5 shrink-0"
            >
              {isTesting ? (
                <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              ) : (
                <>
                  <Activity className="w-4 h-4 text-sky-400" />
                  <span>Test Ping</span>
                </>
              )}
            </button>
          </div>

          {testResult.tested && (
            <div
              className={`p-3 rounded-xl border text-xs flex items-center justify-between ${
                testResult.ok
                  ? 'bg-emerald-950/40 border-emerald-500/30 text-emerald-200'
                  : 'bg-rose-950/40 border-rose-500/30 text-rose-200'
              }`}
            >
              <div className="flex items-center gap-2">
                {testResult.ok ? (
                  <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                ) : (
                  <AlertCircle className="w-4 h-4 text-rose-400 shrink-0" />
                )}
                <span>
                  {testResult.ok
                    ? `Go Backend is live and reachable (${testResult.latency}ms)`
                    : 'Backend server is not responding at this URL.'}
                </span>
              </div>
            </div>
          )}
        </div>

        <div className="rounded-2xl bg-slate-950/70 border border-slate-800/80 p-4 space-y-2 text-xs text-slate-400 leading-relaxed">
          <div className="flex items-center gap-1.5 text-slate-200 font-semibold">
            <Terminal className="w-4 h-4 text-brand-400" />
            <span>How to Run Backend:</span>
          </div>
          <p className="font-mono text-[11px] bg-slate-900 p-2 rounded-lg text-slate-300">
            go run main.go
          </p>
          <p>
            Make sure your <span className="font-mono text-slate-300">.env</span> file has valid PostgreSQL credentials
            and YouTube/Gemini API keys.
          </p>
        </div>

        <div className="flex items-center justify-between gap-3 pt-2">
          <button
            onClick={handleReset}
            className="flex items-center gap-1.5 px-3 py-2 rounded-xl text-xs font-semibold text-slate-400 hover:text-slate-200 transition-colors"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            <span>Reset to Default</span>
          </button>

          <div className="flex items-center gap-2">
            <button
              onClick={onClose}
              className="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSave}
              className="px-5 py-2 rounded-xl bg-sky-500 hover:bg-sky-400 text-slate-950 font-bold text-xs shadow-md shadow-sky-500/20 transition-all"
            >
              Save Configuration
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
