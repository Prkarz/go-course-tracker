import React, { useState } from 'react';
import { api } from '../services/api';
import { useNotification } from '../context/NotificationContext';
import {
  X,
  PlusCircle,
  Video,
  Sparkles,
  AlertCircle,
  CheckCircle2,
  HelpCircle,
  Bot,
} from 'lucide-react';

interface CreateCourseModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCourseCreated: (courseId?: number) => void;
}

export const CreateCourseModal: React.FC<CreateCourseModalProps> = ({
  isOpen,
  onClose,
  onCourseCreated,
}) => {
  const [title, setTitle] = useState('');
  const [url, setUrl] = useState('');
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisStep, setAnalysisStep] = useState(0);
  const [errorMsg, setErrorMsg] = useState('');

  const { addToast } = useNotification();

  if (!isOpen) return null;

  const validateUrl = (rawUrl: string): { isValid: boolean; message: string; listId?: string } => {
    if (!rawUrl.trim()) return { isValid: false, message: '' };
    try {
      const parsed = new URL(rawUrl.trim());
      const host = parsed.hostname.toLowerCase();
      const isYoutube = host === 'youtube.com' || host.endsWith('.youtube.com') || host === 'youtu.be';
      if (!isYoutube) {
        return { isValid: false, message: 'Must be a valid YouTube URL (youtube.com or youtu.be)' };
      }
      const listId = parsed.searchParams.get('list');
      if (!listId) {
        return {
          isValid: false,
          message: 'URL must contain a playlist identifier (?list=...). Single video URLs cannot be tracked as playlists.',
        };
      }
      return { isValid: true, message: `Valid playlist detected: ${listId}`, listId };
    } catch {
      return { isValid: false, message: 'Invalid URL format. Please include https://' };
    }
  };

  const urlValidation = validateUrl(url);

  const samplePlaylists = [
    {
      label: 'Golang Tutorial Series',
      title: 'Go Programming Crash Course',
      url: 'https://www.youtube.com/playlist?list=PL4cUxeGkcC9gC88BEo9czgyS21HJ46PbG',
    },
    {
      label: 'React & TypeScript Masterclass',
      title: 'React TypeScript Full Course',
      url: 'https://www.youtube.com/playlist?list=PLC3y8-rFHvwgu3GVPKtZ7ScW5Qh-fP25w',
    },
  ];

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMsg('');

    const cleanTitle = title.trim();
    const cleanUrl = url.trim();

    if (!cleanTitle || !cleanUrl) {
      setErrorMsg('Both Course Title and YouTube Playlist URL are required.');
      return;
    }

    if (!urlValidation.isValid) {
      setErrorMsg(urlValidation.message || 'Invalid YouTube playlist URL.');
      return;
    }

    setIsAnalyzing(true);
    setAnalysisStep(1);

    const ticker1 = setTimeout(() => setAnalysisStep(2), 1200);
    const ticker2 = setTimeout(() => setAnalysisStep(3), 2600);

    try {
      const res = await api.createCourse({
        title: cleanTitle,
        url: cleanUrl,
      });

      clearTimeout(ticker1);
      clearTimeout(ticker2);

      addToast(
        'Course Added with AI!',
        res.message || 'Course added to your courses successfully.',
        'success'
      );
      onCourseCreated(res.course_id);
      onClose();
      setTitle('');
      setUrl('');
    } catch (err: unknown) {
      clearTimeout(ticker1);
      clearTimeout(ticker2);
      setErrorMsg(err instanceof Error ? err.message : 'Failed to import course');
    } finally {
      setIsAnalyzing(false);
      setAnalysisStep(0);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80">
      <div className="relative w-full max-w-lg rounded-2xl bg-slate-900 border border-slate-800 p-6 sm:p-7 overflow-hidden">
        <button
          onClick={onClose}
          disabled={isAnalyzing}
          className="absolute top-5 right-5 p-1.5 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors disabled:opacity-40"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="flex items-center gap-3 mb-5">
          <div className="w-10 h-10 rounded-xl bg-brand-500 flex items-center justify-center">
            <Video className="w-5 h-5 text-slate-950" />
          </div>
          <div>
            <h3 className="text-xl font-extrabold text-white tracking-tight">
              Import YouTube Course
            </h3>
            <p className="text-xs text-slate-400 mt-0.5">
              Connect a playlist for automated video tracking & Gemini AI summary.
            </p>
          </div>
        </div>

        {errorMsg && (
          <div className="mb-4 p-3.5 rounded-xl bg-rose-950/40 border border-rose-500/30 text-rose-200 flex items-start gap-2.5 text-xs">
            <AlertCircle className="w-4 h-4 text-rose-400 shrink-0 mt-0.5" />
            <span className="leading-relaxed">{errorMsg}</span>
          </div>
        )}

        {isAnalyzing ? (
          <div className="py-8 text-center space-y-5">
            <div className="relative w-16 h-16 mx-auto">
              <div className="absolute inset-0 rounded-full border-4 border-slate-800" />
              <div className="absolute inset-0 rounded-full border-4 border-brand-400 border-t-transparent animate-spin" />
              <Bot className="w-8 h-8 text-brand-400 absolute inset-0 m-auto animate-pulse" />
            </div>

            <div className="space-y-2 max-w-sm mx-auto">
              <h4 className="text-base font-bold text-slate-100">
                {analysisStep === 1 && '1/3: Fetching Playlist from YouTube...'}
                {analysisStep === 2 && '2/3: Gemini AI analyzing topics & structure...'}
                {analysisStep === 3 && '3/3: Generating smart summary & curriculum tags...'}
              </h4>
              <p className="text-xs text-slate-400">
                This takes a few seconds while Google Gemini synthesizes insights for your dashboard.
              </p>
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-slate-300 mb-1.5">
                Course Title
              </label>
              <input
                type="text"
                required
                placeholder="e.g. Master Go Programming & Microservices"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full px-4 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 placeholder-slate-500 text-sm focus:outline-none focus:border-brand-500 transition-all"
              />
            </div>

            <div>
              <label className="block text-xs font-semibold text-slate-300 mb-1.5">
                YouTube Playlist URL
              </label>
              <input
                type="url"
                required
                placeholder="https://www.youtube.com/playlist?list=PL..."
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                className={`w-full px-4 py-2.5 rounded-xl bg-slate-950 border text-slate-100 placeholder-slate-500 text-sm focus:outline-none transition-all ${
                  url.trim() === ''
                    ? 'border-slate-800 focus:border-brand-500'
                    : urlValidation.isValid
                    ? 'border-emerald-500/70 focus:border-emerald-500'
                    : 'border-amber-500/70 focus:border-amber-500'
                }`}
              />

              {url.trim() !== '' && (
                <div className="mt-1.5 flex items-center gap-1.5 text-xs">
                  {urlValidation.isValid ? (
                    <span className="text-emerald-400 flex items-center gap-1 font-medium">
                      <CheckCircle2 className="w-3.5 h-3.5" />
                      <span>{urlValidation.message}</span>
                    </span>
                  ) : (
                    <span className="text-amber-400 flex items-center gap-1">
                      <HelpCircle className="w-3.5 h-3.5 shrink-0" />
                      <span>{urlValidation.message}</span>
                    </span>
                  )}
                </div>
              )}
            </div>

            <div className="pt-1">
              <span className="text-[11px] text-slate-400 font-semibold block mb-1.5">
                Try a sample learning playlist:
              </span>
              <div className="flex flex-wrap gap-2">
                {samplePlaylists.map((sample, idx) => (
                  <button
                    key={idx}
                    type="button"
                    onClick={() => {
                      setTitle(sample.title);
                      setUrl(sample.url);
                    }}
                    className="px-2.5 py-1 rounded-lg text-xs bg-slate-950 border border-slate-800 hover:border-slate-700 text-slate-300 transition-colors"
                  >
                    + {sample.label}
                  </button>
                ))}
              </div>
            </div>

            <button
              type="submit"
              disabled={!urlValidation.isValid || !title.trim()}
              className="w-full mt-2 py-3 rounded-xl bg-gradient-to-r from-brand-600 to-brand-500 hover:from-brand-500 hover:to-brand-400 text-slate-950 font-bold text-sm shadow-lg shadow-brand-500/20 flex items-center justify-center gap-2 transition-all transform active:scale-95 disabled:opacity-50"
            >
              <Sparkles className="w-4 h-4 fill-slate-950" />
              <span>Analyze & Import Course</span>
            </button>
          </form>
        )}
      </div>
    </div>
  );
};
